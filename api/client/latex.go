// API "latex": latex Resource Client
//
// Code generated by goagen v1.1.0-dirty, DO NOT EDIT.
//
// Command:
// $ goagen
// --design=github.com/gianluca311/texsec/api/design
// --out=$(GOPATH)/src/github.com/gianluca311/texsec/api
// --version=v1.1.0-dirty

package client

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
	"strconv"
)

// DownloadLatexPath computes a request path to the download action of latex.
func DownloadLatexPath(uuid int) string {
	param0 := strconv.Itoa(uuid)

	return fmt.Sprintf("/download/%s", param0)
}

// Download route for compilation
func (c *Client) DownloadLatex(ctx context.Context, path string) (*http.Response, error) {
	req, err := c.NewDownloadLatexRequest(ctx, path)
	if err != nil {
		return nil, err
	}
	return c.Client.Do(ctx, req)
}

// NewDownloadLatexRequest create the request corresponding to the download action endpoint of the latex resource.
func (c *Client) NewDownloadLatexRequest(ctx context.Context, path string) (*http.Request, error) {
	scheme := c.Scheme
	if scheme == "" {
		scheme = "http"
	}
	u := url.URL{Host: c.Host, Scheme: scheme, Path: path}
	req, err := http.NewRequest("GET", u.String(), nil)
	if err != nil {
		return nil, err
	}
	return req, nil
}

// StatusLatexPath computes a request path to the status action of latex.
func StatusLatexPath(uuid string) string {
	param0 := uuid

	return fmt.Sprintf("/status/%s", param0)
}

// Actual compilation status
func (c *Client) StatusLatex(ctx context.Context, path string) (*http.Response, error) {
	req, err := c.NewStatusLatexRequest(ctx, path)
	if err != nil {
		return nil, err
	}
	return c.Client.Do(ctx, req)
}

// NewStatusLatexRequest create the request corresponding to the status action endpoint of the latex resource.
func (c *Client) NewStatusLatexRequest(ctx context.Context, path string) (*http.Request, error) {
	scheme := c.Scheme
	if scheme == "" {
		scheme = "http"
	}
	u := url.URL{Host: c.Host, Scheme: scheme, Path: path}
	req, err := http.NewRequest("GET", u.String(), nil)
	if err != nil {
		return nil, err
	}
	return req, nil
}

// UploadLatexPath computes a request path to the upload action of latex.
func UploadLatexPath() string {

	return fmt.Sprintf("/upload")
}

// Route for uploading the Latex files
func (c *Client) UploadLatex(ctx context.Context, path string, debug *bool, maxDownloads *int, uuid *string) (*http.Response, error) {
	req, err := c.NewUploadLatexRequest(ctx, path, debug, maxDownloads, uuid)
	if err != nil {
		return nil, err
	}
	return c.Client.Do(ctx, req)
}

// NewUploadLatexRequest create the request corresponding to the upload action endpoint of the latex resource.
func (c *Client) NewUploadLatexRequest(ctx context.Context, path string, debug *bool, maxDownloads *int, uuid *string) (*http.Request, error) {
	scheme := c.Scheme
	if scheme == "" {
		scheme = "http"
	}
	u := url.URL{Host: c.Host, Scheme: scheme, Path: path}
	values := u.Query()
	if debug != nil {
		tmp5 := strconv.FormatBool(*debug)
		values.Set("debug", tmp5)
	}
	if maxDownloads != nil {
		tmp6 := strconv.Itoa(*maxDownloads)
		values.Set("max_downloads", tmp6)
	}
	if uuid != nil {
		values.Set("uuid", *uuid)
	}
	u.RawQuery = values.Encode()
	req, err := http.NewRequest("POST", u.String(), nil)
	if err != nil {
		return nil, err
	}
	return req, nil
}
