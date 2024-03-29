// Code generated by goagen v1.2.0-dirty, DO NOT EDIT.
//
// API "latex": Application Contexts
//
// Command:
// $ goagen
// --design=github.com/gianluca311/texsec/api/design
// --out=$(GOPATH)/src/github.com/gianluca311/texsec/api
// --version=v1.1.0-dirty

package app

import (
	"context"
	"github.com/goadesign/goa"
	"net/http"
)

// DownloadLatexContext provides the latex download action context.
type DownloadLatexContext struct {
	context.Context
	*goa.ResponseData
	*goa.RequestData
	UUID string
}

// NewDownloadLatexContext parses the incoming request URL and body, performs validations and creates the
// context used by the latex controller download action.
func NewDownloadLatexContext(ctx context.Context, r *http.Request, service *goa.Service) (*DownloadLatexContext, error) {
	var err error
	resp := goa.ContextResponse(ctx)
	resp.Service = service
	req := goa.ContextRequest(ctx)
	req.Request = r
	rctx := DownloadLatexContext{Context: ctx, ResponseData: resp, RequestData: req}
	paramUUID := req.Params["uuid"]
	if len(paramUUID) > 0 {
		rawUUID := paramUUID[0]
		rctx.UUID = rawUUID
	}
	return &rctx, err
}

// OK sends a HTTP response with status code 200.
func (ctx *DownloadLatexContext) OK(resp []byte) error {
	ctx.ResponseData.Header().Set("Content-Type", "text/plain")
	ctx.ResponseData.WriteHeader(200)
	_, err := ctx.ResponseData.Write(resp)
	return err
}

// TooMuchDownloads sends a HTTP response with status code 400.
func (ctx *DownloadLatexContext) TooMuchDownloads(resp []byte) error {
	ctx.ResponseData.Header().Set("Content-Type", "application/json")
	ctx.ResponseData.WriteHeader(400)
	_, err := ctx.ResponseData.Write(resp)
	return err
}

// NotFound sends a HTTP response with status code 404.
func (ctx *DownloadLatexContext) NotFound(resp []byte) error {
	ctx.ResponseData.Header().Set("Content-Type", "application/json")
	ctx.ResponseData.WriteHeader(404)
	_, err := ctx.ResponseData.Write(resp)
	return err
}

// StatusLatexContext provides the latex status action context.
type StatusLatexContext struct {
	context.Context
	*goa.ResponseData
	*goa.RequestData
	UUID string
}

// NewStatusLatexContext parses the incoming request URL and body, performs validations and creates the
// context used by the latex controller status action.
func NewStatusLatexContext(ctx context.Context, r *http.Request, service *goa.Service) (*StatusLatexContext, error) {
	var err error
	resp := goa.ContextResponse(ctx)
	resp.Service = service
	req := goa.ContextRequest(ctx)
	req.Request = r
	rctx := StatusLatexContext{Context: ctx, ResponseData: resp, RequestData: req}
	paramUUID := req.Params["uuid"]
	if len(paramUUID) > 0 {
		rawUUID := paramUUID[0]
		rctx.UUID = rawUUID
	}
	return &rctx, err
}

// OK sends a HTTP response with status code 200.
func (ctx *StatusLatexContext) OK(resp []byte) error {
	ctx.ResponseData.Header().Set("Content-Type", "application/json")
	ctx.ResponseData.WriteHeader(200)
	_, err := ctx.ResponseData.Write(resp)
	return err
}

// NotFound sends a HTTP response with status code 404.
func (ctx *StatusLatexContext) NotFound(resp []byte) error {
	ctx.ResponseData.Header().Set("Content-Type", "application/json")
	ctx.ResponseData.WriteHeader(404)
	_, err := ctx.ResponseData.Write(resp)
	return err
}

// UploadLatexContext provides the latex upload action context.
type UploadLatexContext struct {
	context.Context
	*goa.ResponseData
	*goa.RequestData
}

// NewUploadLatexContext parses the incoming request URL and body, performs validations and creates the
// context used by the latex controller upload action.
func NewUploadLatexContext(ctx context.Context, r *http.Request, service *goa.Service) (*UploadLatexContext, error) {
	var err error
	resp := goa.ContextResponse(ctx)
	resp.Service = service
	req := goa.ContextRequest(ctx)
	req.Request = r
	rctx := UploadLatexContext{Context: ctx, ResponseData: resp, RequestData: req}
	return &rctx, err
}

// OK sends a HTTP response with status code 200.
func (ctx *UploadLatexContext) OK(resp []byte) error {
	ctx.ResponseData.Header().Set("Content-Type", "application/json")
	ctx.ResponseData.WriteHeader(200)
	_, err := ctx.ResponseData.Write(resp)
	return err
}

// NotAcceptable sends a HTTP response with status code 406.
func (ctx *UploadLatexContext) NotAcceptable(resp []byte) error {
	ctx.ResponseData.Header().Set("Content-Type", "application/json")
	ctx.ResponseData.WriteHeader(406)
	_, err := ctx.ResponseData.Write(resp)
	return err
}
