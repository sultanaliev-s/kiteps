package health

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"
)

func TestPing(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/health/ping", nil)
	w := httptest.NewRecorder()

	var (
		serviceName  = "test-service"
		expectedResp = fmt.Sprintf("pong from %s", serviceName)
	)

	handlePing([]byte(expectedResp))(w, req)
	resp := w.Result()
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Errorf("expected status code %d, got %d", http.StatusOK, resp.StatusCode)
	}

	if resp.Header.Get("Content-Type") != "text/plain" {
		t.Errorf(
			"expected content type %s, got %s", "text/plain",
			resp.Header.Get("Content-Type"),
		)
	}

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		t.Errorf("expected no error, got %s", err.Error())
	}

	if string(data) != expectedResp {
		t.Errorf("expected body %q, got %q", expectedResp, string(data))
	}
}

func TestReady(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/health/ready", nil)
	w := httptest.NewRecorder()

	var (
		expectedRespCheck = Check{
			Name:     "databaseConnection",
			Status:   "UP",
			Critical: true,
			Message:  "database connection is up",
		}

		checkers = []Checker{
			func(ctx context.Context) Check {
				return expectedRespCheck
			},
		}

		expectedResp = response{
			Status: "UP",
			Checks: []Check{expectedRespCheck},
		}
	)

	handleReady(checkers)(w, req)
	resp := w.Result()
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Errorf("expected status code %d, got %d", http.StatusOK, resp.StatusCode)
	}

	if resp.Header.Get("Content-Type") != "application/json" {
		t.Errorf(
			"expected content type %s, got %s", "application/json",
			resp.Header.Get("Content-Type"),
		)
	}

	actualResp := response{}
	if err := json.NewDecoder(resp.Body).Decode(&actualResp); err != nil {
		t.Errorf("decoding response, expected no error, got %s", err.Error())
	}

	if actualResp.Status != expectedResp.Status {
		t.Errorf(
			"incorrect status in response, expected status %s, got %s",
			expectedResp.Status,
			actualResp.Status,
		)
	}

	if len(actualResp.Checks) != len(expectedResp.Checks) {
		t.Fatalf(
			"incorrect number of checks in response, expected %d checks, got %d, actual checks: %v",
			len(expectedResp.Checks),
			len(actualResp.Checks),
			actualResp.Checks,
		)
	}

	for i := 0; i < len(actualResp.Checks)-1; i++ {
		if !reflect.DeepEqual(actualResp.Checks[i], expectedResp.Checks[i]) {
			t.Errorf(
				"expected check %v, got %v",
				expectedResp.Checks[i],
				actualResp.Checks[i],
			)
		}
	}
}

func TestLive(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/health/live", nil)
	w := httptest.NewRecorder()

	handleLive(w, req)
	resp := w.Result()
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Errorf("expected status code %d, got %d", http.StatusOK, resp.StatusCode)
	}

	if resp.Header.Get("Content-Type") != "application/json" {
		t.Errorf(
			"expected content type %s, got %s", "application/json",
			resp.Header.Get("Content-Type"),
		)
	}

	actualResp := response{}
	if err := json.NewDecoder(resp.Body).Decode(&actualResp); err != nil {
		t.Errorf("decoding response, expected no error, got %s", err.Error())
	}

	if actualResp.Status != "UP" {
		t.Errorf("expected status %s, got %s", "UP", actualResp.Status)
	}

	if len(actualResp.Checks) != 0 {
		t.Errorf("expected no checks, got %d", len(actualResp.Checks))
	}

	if actualResp.Data == nil {
		t.Fatal("expected data to be present, got nil")
	}

	data, ok := actualResp.Data.(map[string]any)

	if !ok {
		t.Errorf("expected data to be of type map[string]any, got %T", actualResp.Data)
	}

	m, ok := data["memory"]
	if !ok {
		t.Errorf("expected data to have field memory, got %v", data)
	}

	mem, ok := m.(map[string]any)
	if !ok {
		t.Errorf("expected data.memory to be of type memory, got %T", m)
	}

	if _, ok := mem["used"]; !ok {
		t.Errorf("expected data.memory to have field used, got %v", mem)
	}

	if _, ok := mem["free"]; !ok {
		t.Errorf("expected data.memory to have field free, got %v", mem)
	}
}
