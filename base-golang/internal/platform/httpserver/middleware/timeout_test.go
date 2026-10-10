package middleware

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/gofiber/fiber/v3"

	"base-golang/internal/platform/httpserver"
)

func TestTimeoutPreservesCORSAndRequestID(t *testing.T) {
	for _, id := range []string{"timeout-test", ""} {
		router := httpserver.New(nil)
		router.Use(httpserver.ResponseHeaders(), RequestID(), Logging(), CORS([]string{"https://app.example.com"}), httpserver.RequestTimeout(10*time.Millisecond), Recovery())
		finished := make(chan struct{})
		router.Get("/api/test", func(c fiber.Ctx) error {
			defer close(finished)
			<-c.Context().Done()
			return c.Context().Err()
		})
		request := httptest.NewRequest("GET", "/api/test", nil)
		request.Header.Set("Origin", "https://app.example.com")
		request.Header.Set("X-Request-ID", id)
		response, err := router.Test(request)
		if err != nil {
			t.Fatal(err)
		}
		var body struct {
			Error struct {
				Code string `json:"code"`
			} `json:"error"`
		}
		err = json.NewDecoder(response.Body).Decode(&body)
		response.Body.Close()
		if err != nil || response.StatusCode != http.StatusServiceUnavailable || body.Error.Code != "request_timeout" {
			t.Fatalf("status=%d body=%+v error=%v", response.StatusCode, body, err)
		}
		if response.Header.Get("Access-Control-Allow-Origin") != "https://app.example.com" || !strings.Contains(strings.ToLower(response.Header.Get("Access-Control-Expose-Headers")), "x-request-id") {
			t.Fatalf("CORS headers=%v", response.Header)
		}
		actual := response.Header.Get("X-Request-ID")
		if actual == "" || (id != "" && id != actual) {
			t.Fatalf("request ID=%q", actual)
		}
		select {
		case <-finished:
		case <-time.After(time.Second):
			t.Fatal("handler did not cancel")
		}
	}
}
