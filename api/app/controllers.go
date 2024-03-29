// Code generated by goagen v1.2.0-dirty, DO NOT EDIT.
//
// API "latex": Application Controllers
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
	"github.com/goadesign/goa/cors"
	"net/http"
)

// initService sets up the service encoders, decoders and mux.
func initService(service *goa.Service) {
	// Setup encoders and decoders
	service.Encoder.Register(goa.NewJSONEncoder, "application/json")
	service.Decoder.Register(goa.NewJSONDecoder, "application/json")
	service.Decoder.Register(goa.NewGobDecoder, "application/gob", "application/x-gob")
	service.Decoder.Register(goa.NewXMLDecoder, "application/xml")

	// Setup default encoder and decoder
	service.Encoder.Register(goa.NewJSONEncoder, "*/*")
	service.Decoder.Register(goa.NewJSONDecoder, "*/*")
}

// LatexController is the controller interface for the Latex actions.
type LatexController interface {
	goa.Muxer
	Download(*DownloadLatexContext) error
	Status(*StatusLatexContext) error
	Upload(*UploadLatexContext) error
}

// MountLatexController "mounts" a Latex resource controller on the given service.
func MountLatexController(service *goa.Service, ctrl LatexController) {
	initService(service)
	var h goa.Handler
	service.Mux.Handle("OPTIONS", "/download/:uuid", ctrl.MuxHandler("preflight", handleLatexOrigin(cors.HandlePreflight()), nil))
	service.Mux.Handle("OPTIONS", "/status/:uuid", ctrl.MuxHandler("preflight", handleLatexOrigin(cors.HandlePreflight()), nil))
	service.Mux.Handle("OPTIONS", "/upload", ctrl.MuxHandler("preflight", handleLatexOrigin(cors.HandlePreflight()), nil))

	h = func(ctx context.Context, rw http.ResponseWriter, req *http.Request) error {
		// Check if there was an error loading the request
		if err := goa.ContextError(ctx); err != nil {
			return err
		}
		// Build the context
		rctx, err := NewDownloadLatexContext(ctx, req, service)
		if err != nil {
			return err
		}
		return ctrl.Download(rctx)
	}
	h = handleLatexOrigin(h)
	service.Mux.Handle("GET", "/download/:uuid", ctrl.MuxHandler("download", h, nil))
	service.LogInfo("mount", "ctrl", "Latex", "action", "Download", "route", "GET /download/:uuid")

	h = func(ctx context.Context, rw http.ResponseWriter, req *http.Request) error {
		// Check if there was an error loading the request
		if err := goa.ContextError(ctx); err != nil {
			return err
		}
		// Build the context
		rctx, err := NewStatusLatexContext(ctx, req, service)
		if err != nil {
			return err
		}
		return ctrl.Status(rctx)
	}
	h = handleLatexOrigin(h)
	service.Mux.Handle("GET", "/status/:uuid", ctrl.MuxHandler("status", h, nil))
	service.LogInfo("mount", "ctrl", "Latex", "action", "Status", "route", "GET /status/:uuid")

	h = func(ctx context.Context, rw http.ResponseWriter, req *http.Request) error {
		// Check if there was an error loading the request
		if err := goa.ContextError(ctx); err != nil {
			return err
		}
		// Build the context
		rctx, err := NewUploadLatexContext(ctx, req, service)
		if err != nil {
			return err
		}
		return ctrl.Upload(rctx)
	}
	h = handleLatexOrigin(h)
	service.Mux.Handle("POST", "/upload", ctrl.MuxHandler("upload", h, nil))
	service.LogInfo("mount", "ctrl", "Latex", "action", "Upload", "route", "POST /upload")
}

// handleLatexOrigin applies the CORS response headers corresponding to the origin.
func handleLatexOrigin(h goa.Handler) goa.Handler {

	return func(ctx context.Context, rw http.ResponseWriter, req *http.Request) error {
		origin := req.Header.Get("Origin")
		if origin == "" {
			// Not a CORS request
			return h(ctx, rw, req)
		}
		if cors.MatchOrigin(origin, "*") {
			ctx = goa.WithLogContext(ctx, "origin", origin)
			rw.Header().Set("Access-Control-Allow-Origin", origin)
			rw.Header().Set("Access-Control-Expose-Headers", "Cache-Control, Content-Length, Content-Type, Date, Expires, Host, Keep-Alive, Last-Modified, Location, Server, Status, Strict-Transport-Security, X-Requested-With, Accept, Origin, X-File-Name")
			rw.Header().Set("Access-Control-Max-Age", "600")
			rw.Header().Set("Access-Control-Allow-Credentials", "true")
			if acrm := req.Header.Get("Access-Control-Request-Method"); acrm != "" {
				// We are handling a preflight request
				rw.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
			}
			return h(ctx, rw, req)
		}

		return h(ctx, rw, req)
	}
}
