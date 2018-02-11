package handlers

import (
	"context"
	"github.com/gavincabbage/api.il2missionplanner.com/mock"
	"github.com/gorilla/mux"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestHealthHandler(t *testing.T) {
	testHandler("GET", "/health", http.StatusOK,
		`{"data":"OK","error":""}`, HealthHandler, t)
}

func TestNotFoundHandler(t *testing.T) {
	// TODO Improve API body comparison
	testHandler("GET", "/badurl", http.StatusNotFound,
		"", NotFoundHandler, t)
}

func TestConfigHandlerHandler(t *testing.T) {
	// TODO Improve API body comparison
	testHandler("GET", "/config", http.StatusOK, ``, ConfigHandler, t)
}

// testHandler and supporting methods below implement a reusable pattern for testing functions
// implementing the http.HandlerFunc interface.
func testHandler(method string, path string, expectedStatus int, expectedBody string,
	h http.HandlerFunc, t *testing.T) {

	mockResponse, mockRequest := buildMockResponseAndRequest(method, path, t)
	handler := http.HandlerFunc(h)
	handler.ServeHTTP(mockResponse, mockRequest)

	checkHandlerResponseCode(mockResponse, expectedStatus, t)
	if expectedBody != "" {
		checkHandlerResponseBody(mockResponse, expectedBody, t)
	}
}

func buildMockResponseAndRequest(method string, path string, t *testing.T) (*httptest.ResponseRecorder, *http.Request) {

	mockRequest, err := http.NewRequest(method, path, nil)
	if err != nil {
		t.Fatal(err)
	}
	mockRequest = mockRequest.WithContext(context.WithValue(mockRequest.Context(), "config", mock.Config()))
	mockRequest = mux.SetURLVars(mockRequest, map[string]string{"server": "testserver"})

	mockResponse := httptest.NewRecorder()
	return mockResponse, mockRequest
}

func checkHandlerResponseCode(mockResponse *httptest.ResponseRecorder, expectedStatus int, t *testing.T) {
	if status := mockResponse.Code; status != expectedStatus {
		t.Errorf("handler returned incorrect status code: got %v want %v",
			status, expectedStatus)
	}
}

func checkHandlerResponseBody(mockResponse *httptest.ResponseRecorder, expectedBody string, t *testing.T) {
	if body := mockResponse.Body.String(); body != expectedBody {
		t.Errorf("handler returned incorrect response body: got %v want %v",
			body, expectedBody)
	}
}
