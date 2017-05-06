package api

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/rpc"
	"os"
	"strconv"
	"sync"
	"time"

	"github.com/gianluca311/texsec/api/app"
	"github.com/goadesign/goa"
	"github.com/satori/go.uuid"
	"gopkg.in/h2non/filetype.v1"
)

// LatexController implements the latex resource.
type LatexController struct {
	*goa.Controller
	*sync.Mutex
}

type ResponseMessage struct {
	OK      bool   `json:"ok"`
	UUID    string `json:"uuid"`
	Message string `json:"message"`
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

	// Put your logic here

	// LatexController_Download: end_implement
	return nil
}

// Status runs the status action.
func (c *LatexController) Status(ctx *app.StatusLatexContext) error {
	// LatexController_Status: start_implement
	// Put your logic here

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

	//FIX ME
	file, handler, err := ctx.FormFile("file")
	if err != nil {
		return goa.ErrBadRequest("failed to load file: %s", err.Error())
	}
	defer file.Close()

	head := make([]byte, 261)
	_, err = file.Read(head)
	if err != nil {
		return goa.ErrBadRequest("unable to read file: %s", err.Error())
	}

	if !filetype.IsArchive(head) {
		res := &ResponseMessage{
			OK:      false,
			Message: "Uploaded file isn't an archive",
		}
		resJSON, _ := json.Marshal(res)
		return ctx.NotAcceptable(resJSON)
	}

	fileKind, _ := filetype.Match(head)

	var archive *app.LatexArchive

	fileName := "latex-archive-" + uuid + fileKind.Extension
	f, err := os.OpenFile("./archives/"+fileName, os.O_WRONLY|os.O_CREATE, 0666)
	if err != nil {
		return fmt.Errorf("failed to save file: %s", err)
	}

	defer f.Close()
	io.Copy(f, file)

	archive = &app.LatexArchive{UUID: uuid, Filename: fileName, UploadedAt: uploadedAt}

	res := &ResponseMessage{
		OK:      true,
		UUID:    archive.UUID,
		Message: "Upload of " + handler.Filename + " successfull. debug set to: " + strconv.FormatBool(debug) + ". max_downloads: " + strconv.Itoa(maxDownloads) + " Proccess UUID: " + archive.UUID,
	}

	client, err := rpc.DialHTTP("tcp", "localhost:1234")
	if err != nil {
		log.Fatal("dialing:", err)
	}

	var resp ResponseMessage
	client.Call("StartProcessing", archive, &resp)

	resJSON, _ := json.Marshal(res)
	ctx.OK(resJSON)
	// LatexController_Upload: end_implement
	return nil
}
