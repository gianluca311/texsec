package api

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"strconv"
	"sync"
	"time"

	"github.com/gianluca311/texsec/api/app"
	"github.com/goadesign/goa"
	"github.com/satori/go.uuid"
)

// LatexController implements the latex resource.
type LatexController struct {
	*goa.Controller
	*sync.Mutex
	archive *LatexUploadArchiveData
}

type LatexUploadArchiveData struct {
	UUID       string
	FileName   string
	UploadedAt time.Time
}

type ResponseMessage struct {
	OK      bool
	UUID    string
	Message string
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

	debug, _ := strconv.ParseBool(ctx.PostFormValue("debug"))
	maxDownloads, _ := strconv.Atoi(ctx.PostFormValue("max_downloads"))
	//FIX ME
	file, handler, err := ctx.FormFile("file")
	if err != nil {
		return goa.ErrBadRequest("failed to load file: %s", err.Error())
	}
	defer file.Close()
	var archive *app.LatexArchive

	fileName := "latex-archive-" + time.Now().String() + handler.Filename
	f, err := os.OpenFile("./archives/"+fileName, os.O_WRONLY|os.O_CREATE, 0666)
	if err != nil {
		return fmt.Errorf("failed to save file: %s", err)
	}

	defer f.Close()
	io.Copy(f, file)

	data := c.saveArchive(fileName)
	archive = &app.LatexArchive{UUID: data.UUID, Filename: data.FileName, UploadedAt: data.UploadedAt}

	res := &ResponseMessage{
		OK:      true,
		UUID:    archive.UUID,
		Message: "Upload success full. debug set to: " + strconv.FormatBool(debug) + " max_downloads: " + strconv.Itoa(maxDownloads) + " Proccess UUID: " + archive.UUID,
	}

	resJSON, _ := json.Marshal(res)
	ctx.OK(resJSON)
	// LatexController_Upload: end_implement
	return nil
}

func (c *LatexController) saveArchive(filename string) *LatexUploadArchiveData {
	c.Lock()
	defer c.Unlock()
	data := &LatexUploadArchiveData{
		UUID:       uuid.NewV4().String(),
		FileName:   filename,
		UploadedAt: time.Now(),
	}
	c.archive = data
	return data
}
