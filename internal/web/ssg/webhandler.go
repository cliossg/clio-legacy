package ssg

import (
	hm "github.com/hermesgen/hm"
)

const (
	// WIP: This will be obtained from configuration.
	defaultAPIBaseURL = "http://localhost:8081/api/v1"
)

const (
	ssgFeat = "ssg"
	ssgPath = "/ssg"
)

type WebHandler struct {
	*hm.WebHandler
	apiClient *hm.APIClient
}

func NewWebHandler(tm *hm.TemplateManager, flash *hm.FlashManager, opts ...hm.Option) *WebHandler {
	handler := hm.NewWebHandler(tm, flash, opts...)
	apiClient := hm.NewAPIClient("web-api-client", func() string { return "" }, defaultAPIBaseURL, opts...)
	return &WebHandler{
		WebHandler: handler,
		apiClient:  apiClient,
	}
}
