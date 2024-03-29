package api

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/rpc"
	"strconv"
	"sync"
	"time"

	"bytes"
	"io"

	"github.com/gianluca311/texsec/api/app"
	"github.com/goadesign/goa"
	"github.com/satori/go.uuid"
	"github.com/spf13/viper"
	"gopkg.in/h2non/filetype.v1"
)

// LatexController implements the latex resource.
type LatexController struct {
	*goa.Controller
	*sync.Mutex
}

type RPCCompileRequest struct {
	ArchiveInfo  *app.LatexArchive
	File         []byte
	MaxDownloads int
	Debug        bool
}

type RPCDownloadResponse struct {
	ResponseMessage
	File []byte
}

type RPCRequest struct {
	UUID string
}

type ResponseMessage struct {
	OK      bool   `json:"ok"`
	UUID    string `json:"uuid"`
	Message string `json:"message"`
}

type StatusMessage struct {
	OK            bool      `json:"ok"`
	UUID          string    `json:"uuid"`
	MaxDownloads  int       `json:"max_downloads"`
	DownloadCount int       `json:"download_count"`
	UploadTime    time.Time `json:"upload_time"`
	Message       string    `json:"message"`
}

// NewLatexController creates a latex controller.
func NewLatexController(service *goa.Service) *LatexController {
	return &LatexController{
		Controller: service.NewController("LatexController"),
		Mutex:      &sync.Mutex{},
	}
}

// Download runs the download action.
func (c *LatexController) Download(ctx *app.DownloadLatexContext) error {
	// LatexController_Download: start_implement

	client, err := rpc.DialHTTP("tcp", viper.GetString("daemonEndpoint"))
	if err != nil {
		return goa.ErrInternal(err.Error())
	}

	uuidParam := ctx.Params.Get("uuid")
	var resp RPCDownloadResponse
	err = client.Call("Daemon.Download", RPCRequest{uuidParam}, &resp)
	if err != nil {
		return goa.ErrInternal(err.Error())
	}
	fileReader := bytes.NewReader(resp.File)

	ctx.ResponseData.Header().Add("Content-Disposition", "attachment; filename="+uuidParam+".pdf")
	ctx.ResponseData.Header().Set("Content-Type", ctx.ResponseData.Header().Get("Content-Type"))
	io.Copy(ctx.ResponseWriter, fileReader)
	// LatexController_Download: end_implement
	return nil
}

// Status runs the status action.
func (c *LatexController) Status(ctx *app.StatusLatexContext) error {
	// LatexController_Status: start_implement
	// Put your logic here
	client, err := rpc.DialHTTP("tcp", viper.GetString("daemonEndpoint"))
	if err != nil {
		return goa.ErrInternal(err.Error())
	}

	uuidParam := ctx.Params.Get("uuid")
	var resp StatusMessage
	err = client.Call("Daemon.Status", RPCRequest{uuidParam}, &resp)
	if err != nil {
		return goa.ErrInternal(err.Error())
	}

	if resp.OK == false {
		jsonResp, _ := json.Marshal(resp)
		return ctx.NotFound(jsonResp)
	}

	jsonResp, _ := json.Marshal(resp)
	ctx.OK(jsonResp)

	// LatexController_Status: end_implement
	return nil
}

// Upload runs the upload action.
func (c *LatexController) Upload(ctx *app.UploadLatexContext) error {
	// LatexController_Upload: start_implement

	// Put your logic here
	uuid := uuid.NewV4().String()
	uploadedAt := time.Now()

	debug, _ := strconv.ParseBool(ctx.PostFormValue("debug"))
	maxDownloads, _ := strconv.Atoi(ctx.PostFormValue("max_downloads"))
	if maxDownloads < 0 {
		maxDownloads = 0
	}

	//FIX ME
	file, handler, err := ctx.FormFile("file")
	if err != nil {
		return goa.ErrBadRequest("failed to load file: %s", err.Error())
	}
	defer file.Close()

	fileContent, err := ioutil.ReadAll(file)
	if err != nil {
		return goa.ErrBadRequest("unable to read file: %s", err.Error())
	}

	fileHead := fileContent[:261]

	if !filetype.IsArchive(fileHead) {
		res := &ResponseMessage{
			OK:      false,
			Message: "Uploaded file isn't an archive",
		}
		resJSON, _ := json.Marshal(res)
		return ctx.NotAcceptable(resJSON)
	}

	fileKind, _ := filetype.Match(fileHead)

	var archive *app.LatexArchive

	fileName := "latex-archive-" + uuid + fileKind.Extension
	/*f, err := os.OpenFile("./archives/"+fileName, os.O_WRONLY|os.O_CREATE, 0666)
	if err != nil {
		return fmt.Errorf("failed to save file: %s", err)
	}

	defer f.Close()
	io.Copy(f, file)*/

	archive = &app.LatexArchive{UUID: uuid, Filename: fileName, UploadedAt: uploadedAt}

	res := &ResponseMessage{
		OK:      true,
		UUID:    archive.UUID,
		Message: "Upload of " + handler.Filename + " successfull. debug set to: " + strconv.FormatBool(debug) + ". max_downloads: " + strconv.Itoa(maxDownloads) + " Proccess UUID: " + archive.UUID,
	}

	client, err := rpc.DialHTTP("tcp", viper.GetString("daemonEndpoint"))
	if err != nil {
		log.Fatal("dialing:", err)
	}

	var resp ResponseMessage
	uploadRequest := RPCCompileRequest{ArchiveInfo: archive, File: fileContent, MaxDownloads: maxDownloads}
	client.Call("Daemon.Compile", uploadRequest, &resp)

	resJSON, _ := json.Marshal(res)
	ctx.OK(resJSON)
	// LatexController_Upload: end_implement
	return nil
}
