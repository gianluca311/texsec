package api

import (
	"github.com/gianluca311/texsec/api/app"
	"github.com/goadesign/goa"
)

// LatexController implements the latex resource.
type LatexController struct {
	*goa.Controller
}

// NewLatexController creates a latex controller.
func NewLatexController(service *goa.Service) *LatexController {
	return &LatexController{Controller: service.NewController("LatexController")}
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

	// LatexController_Upload: end_implement
	return nil
}
