//go:generate goagen bootstrap -d github.com/gianluca311/texsec/api/design

package api

import (
	"github.com/gianluca311/texsec/api/app"
	"github.com/goadesign/goa"
	"github.com/goadesign/goa/middleware"
)

func StartAPI() {
	// Create service
	service := goa.New("latex")

	// Mount middleware
	service.Use(middleware.RequestID())
	service.Use(middleware.LogRequest(true))
	service.Use(middleware.ErrorHandler(service, true))
	service.Use(middleware.Recover())

	// Mount "latex" controller
	c := NewLatexController(service)
	app.MountLatexController(service, c)

	// Start service
	if err := service.ListenAndServe(":8080"); err != nil {
		service.LogError("startup", "err", err)
	}

}
