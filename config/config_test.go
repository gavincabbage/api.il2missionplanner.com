package config

import (
	"context"
	"net/http"
	"testing"
)

func TestFromFile(t *testing.T) {
	// TODO Mock ioutil file read
}

func TestFromRequestContext(t *testing.T) {
	mockConfig := &Config{}
	mockConfig.Env = "test"
	mockRequest := &http.Request{}
	mockRequest = mockRequest.WithContext(context.WithValue(mockRequest.Context(), "config", mockConfig))
	resultConfig := FromRequestContext(mockRequest)
	if resultConfig.Env != mockConfig.Env {
		t.Fail()
	}
}
