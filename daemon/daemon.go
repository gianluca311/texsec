package daemon

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"net/rpc"
	"os"
	"path"
	"path/filepath"
	"strings"
	"time"

	"github.com/boltdb/bolt"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/client"
	"github.com/gianluca311/texsec/api"
	"github.com/spf13/viper"
)

var (
	bucketName = []byte("compilations")
)

const (
	compilationRunning = iota
	compilationDone    = iota
	compilationError   = iota
)

type compilation struct {
	UUID           string
	ContainerID    string
	InputfileName  string
	OutputFileName string
	UploadTime     time.Time
	MaxDownloads   int
	DownloadCount  int
	Status         int
	Logs           string
}

type Daemon struct {
	client *client.Client
	db     *bolt.DB
}

func (this *Daemon) Status(args *api.RPCRequest, reply *api.StatusMessage) error {

	this.db.View(func(tx *bolt.Tx) error {

		b := tx.Bucket(bucketName)

		data := b.Get([]byte(args.UUID))
		if data == nil {
			return fmt.Errorf("Compilation with UUID %s not found", args.UUID)
		}

		var cmpl compilation
		err := json.Unmarshal(data, &cmpl)
		if err != nil {
			return err
		}

		reply.OK = true
		reply.DownloadCount = cmpl.DownloadCount
		reply.Message = cmpl.Logs
		reply.UUID = cmpl.UUID
		reply.MaxDownloads = cmpl.MaxDownloads
		reply.UploadTime = cmpl.UploadTime

		return nil
	})

	return nil
}

// Download attempts to return the compiled Latex document from storage.
func (this *Daemon) Download(args *api.RPCRequest, reply *api.RPCDownloadResponse) error {
	log.Printf("Download request for %s", args.UUID)

	var err error
	err = this.db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket(bucketName)

		// Get DB data for requested UUID
		data := b.Get([]byte(args.UUID))
		if data == nil {
			log.Printf("Compilation with UUID %s not found", args.UUID)
			return fmt.Errorf("Compilation with UUID %s not found", args.UUID)
		}

		var cmpl compilation
		err = json.Unmarshal(data, &cmpl)
		if err != nil {
			return err
		}

		// Check for max download restriction. If there is a restriction (maxDownloads != 0)
		// and it is already reached, deny the download.
		if cmpl.MaxDownloads != 0 && cmpl.MaxDownloads <= cmpl.DownloadCount {
			reply.UUID = cmpl.UUID
			reply.OK = false
			reply.Message = "Download limit exceeded"
			return nil
		}

		cmpl.DownloadCount++

		fileData, err := ioutil.ReadFile(cmpl.OutputFileName)
		if err != nil {
			log.Printf("[%s] Unable to open output file %s", args.UUID, cmpl.OutputFileName)
			reply.UUID = cmpl.UUID
			reply.OK = false
			reply.Message = fmt.Sprintf("Unable to open output file %s", cmpl.OutputFileName)
			return err
		}

		reply.File = fileData
		reply.OK = true
		reply.UUID = args.UUID

		// Save modified data back into DB
		data, err = json.Marshal(&cmpl)
		if err != nil {
			return err
		}

		err = b.Put([]byte(args.UUID), data)
		return err
	})

	if err != nil {
		reply.Message = err.Error()
		reply.OK = false
		reply.UUID = args.UUID
		return err
	}

	return err
}

func (this *Daemon) Compile(args *api.RPCCompileRequest, reply *api.ResponseMessage) error {
	log.Printf("Compilation request %s received", args.ArchiveInfo.UUID)

	volumePath := path.Join("compilations/", args.ArchiveInfo.UUID)
	volumePathAbsolute, err := filepath.Abs(volumePath)
	if err != nil {
		reply.Message = err.Error()
		reply.OK = false
		reply.UUID = args.ArchiveInfo.UUID
		return err
	}

	zipFilePath := path.Join(volumePathAbsolute, args.ArchiveInfo.Filename)

	err = os.MkdirAll(volumePathAbsolute, 0700)
	if err != nil {
		reply.Message = err.Error()
		reply.OK = false
		reply.UUID = args.ArchiveInfo.UUID
		log.Printf("[%s] Failed to create volume directory: %s", args.ArchiveInfo.UUID, err.Error())
		return err
	}

	err = ioutil.WriteFile(zipFilePath, args.File, 0700)
	if err != nil {
		reply.Message = err.Error()
		reply.OK = false
		reply.UUID = args.ArchiveInfo.UUID
		log.Printf("[%s] Failed to write zip file to directory: %s", args.ArchiveInfo.UUID, err.Error())
		return err
	}

	err = unzip(zipFilePath, volumePathAbsolute)
	if err != nil {
		reply.Message = err.Error()
		reply.OK = false
		reply.UUID = args.ArchiveInfo.UUID
		log.Printf("[%s] Failed to unzip: %s", args.ArchiveInfo.UUID, err.Error())
		return err
	}

	texFiles := checkExt(".tex", volumePathAbsolute)

	if len(texFiles) > 1 {
		reply.Message = fmt.Sprintf("Multiple tex files in archive: %s", strings.Join(texFiles, ","))
		reply.OK = false
		reply.UUID = args.ArchiveInfo.UUID
		log.Printf("[%s] Multiple tex files in archive: %s", args.ArchiveInfo.UUID, strings.Join(texFiles, ","))
		return err
	}

	if len(texFiles) == 0 {
		reply.Message = "No tex files found in archive."
		reply.OK = false
		reply.UUID = args.ArchiveInfo.UUID
		log.Printf("[%s] No tex files in archive", args.ArchiveInfo.UUID)
		return err
	}

	_, err = this.client.ImagePull(context.Background(), viper.GetString("dockerImage"), types.ImagePullOptions{})
	if err != nil {
		reply.Message = err.Error()
		reply.OK = false
		reply.UUID = args.ArchiveInfo.UUID
		log.Printf("[%s] Failed to pull image: %s", args.ArchiveInfo.UUID, err.Error())
		return err
	}

	cfg := &container.Config{
		Image:           viper.GetString("dockerImage"),
		Cmd:             []string{viper.GetString("latexCommand"), "-output-directory=/texdata", viper.GetString("latexCommandParam"), fmt.Sprintf("/texdata/%s", path.Base(texFiles[0]))},
		NetworkDisabled: true,
		WorkingDir:      "/texdata",
	}

	hostCfg := &container.HostConfig{
		Binds: []string{volumePathAbsolute + "://texdata"},
	}

	resp, err := this.client.ContainerCreate(context.Background(), cfg, hostCfg, nil, args.ArchiveInfo.UUID)
	if err != nil {
		reply.Message = err.Error()
		reply.OK = false
		reply.UUID = args.ArchiveInfo.UUID
		log.Printf("[%s] Failed to create container: %s", args.ArchiveInfo.UUID, err.Error())
		return err
	}

	err = this.client.ContainerStart(context.Background(), resp.ID, types.ContainerStartOptions{})
	if err != nil {
		reply.Message = err.Error()
		reply.OK = false
		reply.UUID = args.ArchiveInfo.UUID
		log.Printf("[%s] Failed to start container: %s", args.ArchiveInfo.UUID, err.Error())
		return err
	}

	ext := path.Ext(texFiles[0])
	outfile := texFiles[0][0:len(texFiles[0])-len(ext)] + ".pdf"
	cpl := compilation{
		args.ArchiveInfo.UUID,
		resp.ID,
		args.ArchiveInfo.Filename,
		path.Join(volumePathAbsolute, outfile),
		args.ArchiveInfo.UploadedAt,
		args.MaxDownloads,
		0,
		compilationRunning,
		"",
	}

	data, err := json.Marshal(&cpl)
	if err != nil {
		reply.Message = err.Error()
		reply.OK = false
		reply.UUID = args.ArchiveInfo.UUID
		log.Printf("[%s] Failed to persist compilation: %s", args.ArchiveInfo.UUID, err.Error())
		return err
	}

	this.db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket(bucketName)
		err = b.Put([]byte(args.ArchiveInfo.UUID), data)
		return err
	})

	reply.OK = true
	reply.UUID = args.ArchiveInfo.UUID
	return nil
}

func (this *Daemon) containerWatcher() {

	ticker := time.NewTicker(10 * time.Second)

	for t := range ticker.C {

		log.Printf("Starting cleanup round at %s", t.String())

		this.db.Update(func(tx *bolt.Tx) error {
			b := tx.Bucket(bucketName)
			b.ForEach(func(key, value []byte) error {

				var cmpl compilation
				err := json.Unmarshal(value, &cmpl)
				if err != nil {
					return err
				}

				if cmpl.Status == compilationRunning {

					resp, err := this.client.ContainerInspect(context.Background(), cmpl.ContainerID)
					if err != nil {
						return err
					}

					if !resp.State.Running {

						log.Printf("[%s] Container is no longer running", cmpl.UUID)
						if resp.State.ExitCode != 0 {

							cmpl.Status = compilationError

						} else {

							cmpl.Status = compilationDone

						}

						err = this.getCompilationLogs(&cmpl)
						if err != nil {
							log.Printf("[%s] Failed to collect logs: %s", cmpl.UUID, err.Error())
							return err
						}

						log.Printf("[%s] Removing container", cmpl.UUID)
						this.client.ContainerRemove(context.Background(), cmpl.ContainerID, types.ContainerRemoveOptions{})
					}

					data, err := json.Marshal(&cmpl)
					err = b.Put([]byte(cmpl.UUID), data)
					if err != nil {
						return err
					}

				} else {
					if time.Since(cmpl.UploadTime) > (7 * 24 * time.Hour) {
						b.Delete([]byte(cmpl.UUID))
					}
				}

				return nil
			})

			return nil
		})
	}
}

func (this *Daemon) getCompilationLogs(cmpl *compilation) error {

	opt := types.ContainerLogsOptions{ShowStdout: true, ShowStderr: true}
	logsRdr, err := this.client.ContainerLogs(context.Background(), cmpl.ContainerID, opt)
	if err != nil {
		return err
	}

	logsData, err := ioutil.ReadAll(logsRdr)
	cmpl.Logs = string(logsData)

	return nil
}

func StartDaemon() {
	cli, err := client.NewEnvClient()
	if err != nil {
		log.Fatalf("Failed to create docker client: %s", err.Error())
	}

	db, err := bolt.Open("daemon.db", 0600, nil)
	if err != nil {
		log.Fatalf("Failed to create database client: %s", err.Error())
	}

	db.Update(func(tx *bolt.Tx) error {
		_, err := tx.CreateBucketIfNotExists(bucketName)
		return err
	})

	daemon := &Daemon{cli, db}
	go daemon.containerWatcher()

	rpc.Register(daemon)

	rpc.HandleHTTP()
	l, err := net.Listen("tcp", viper.GetString("daemonEndpoint"))
	if err != nil {
		log.Fatal("listen error:", err)
	}
	http.Serve(l, nil)
}
