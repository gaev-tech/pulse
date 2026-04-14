package tests

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	v1 "github.com/gaevivan/pulse/internal/handler/v1"
)

func newTestRouter() http.Handler {
	return v1.NewRouter(v1.Deps{})
}

func TestSendMagicLink_InvalidEmail(t *testing.T) {
	tests := []struct {
		name string
		body string
	}{
		{"invalid email", `{"email":"not-an-email"}`},
		{"empty email", `{"email":""}`},
		{"missing email", `{}`},
		{"empty body", ``},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			router := newTestRouter()
			request := httptest.NewRequest(http.MethodPost, "/v1/auth/magic-link",
				bytes.NewBufferString(test.body))
			request.Header.Set("Content-Type", "application/json")
			recorder := httptest.NewRecorder()

			router.ServeHTTP(recorder, request)

			if recorder.Code != http.StatusBadRequest {
				t.Fatalf("expected 400, got %d", recorder.Code)
			}
			assertErrorCode(t, recorder.Body.Bytes(), "VALIDATION_ERROR")
		})
	}
}

func TestVerifyMagicLink_MissingToken(t *testing.T) {
	router := newTestRouter()
	request := httptest.NewRequest(http.MethodPost, "/v1/auth/magic-link/verify",
		bytes.NewBufferString(`{"token":""}`))
	request.Header.Set("Content-Type", "application/json")
	recorder := httptest.NewRecorder()

	router.ServeHTTP(recorder, request)

	if recorder.Code != http.StatusBadRequest {
		t.Fatalf("expected 400, got %d", recorder.Code)
	}
	assertErrorCode(t, recorder.Body.Bytes(), "VALIDATION_ERROR")
}

func TestRefresh_MissingToken(t *testing.T) {
	router := newTestRouter()
	request := httptest.NewRequest(http.MethodPost, "/v1/auth/refresh",
		bytes.NewBufferString(`{}`))
	request.Header.Set("Content-Type", "application/json")
	recorder := httptest.NewRecorder()

	router.ServeHTTP(recorder, request)

	if recorder.Code != http.StatusBadRequest {
		t.Fatalf("expected 400, got %d", recorder.Code)
	}
	assertErrorCode(t, recorder.Body.Bytes(), "VALIDATION_ERROR")
}

func TestLogout_MissingToken(t *testing.T) {
	router := newTestRouter()
	request := httptest.NewRequest(http.MethodPost, "/v1/auth/logout",
		bytes.NewBufferString(`{}`))
	request.Header.Set("Content-Type", "application/json")
	recorder := httptest.NewRecorder()

	router.ServeHTTP(recorder, request)

	if recorder.Code != http.StatusBadRequest {
		t.Fatalf("expected 400, got %d", recorder.Code)
	}
	assertErrorCode(t, recorder.Body.Bytes(), "VALIDATION_ERROR")
}

func TestMe_Unauthorized(t *testing.T) {
	router := newTestRouter()
	request := httptest.NewRequest(http.MethodGet, "/v1/auth/me", nil)
	recorder := httptest.NewRecorder()

	router.ServeHTTP(recorder, request)

	if recorder.Code != http.StatusUnauthorized {
		t.Fatalf("expected 401, got %d", recorder.Code)
	}
	assertErrorCode(t, recorder.Body.Bytes(), "UNAUTHORIZED")
}

func assertErrorCode(t *testing.T, body []byte, expectedCode string) {
	t.Helper()
	var response struct {
		Error struct {
			Code string `json:"code"`
		} `json:"error"`
	}
	if err := json.Unmarshal(body, &response); err != nil {
		t.Fatalf("failed to decode error response: %v", err)
	}
	if response.Error.Code != expectedCode {
		t.Fatalf("expected error code %q, got %q", expectedCode, response.Error.Code)
	}
}
