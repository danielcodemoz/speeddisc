//go:build linux

package main

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"testing/fstest"
	"time"

	"github.com/danielcodemoz/speeddisc/internal/report"
)

const testHost = "127.0.0.1:41234"
const testPrefix = "/secret/"

func newTestDesk(demo bool) *desk {
	return newDesk(testHost, "secret", fstest.MapFS{"index.html": &fstest.MapFile{Data: []byte("<!doctype html><title>SpeedDisc</title>")}}, func() {}, demo)
}
func makeRequest(t *testing.T, d *desk, method, url, body, host, origin, remote string) *httptest.ResponseRecorder {
	t.Helper()
	r := httptest.NewRequest(method, "http://"+testHost+url, strings.NewReader(body))
	r.Host = host
	r.RemoteAddr = remote
	if origin != "" {
		r.Header.Set("Origin", origin)
	}
	if body != "" {
		r.Header.Set("Content-Type", "application/json")
	}
	w := httptest.NewRecorder()
	d.ServeHTTP(w, r)
	return w
}
func TestAccessBoundToLoopbackHostTokenAndOrigin(t *testing.T) {
	d := newTestDesk(true)
	for _, v := range []struct {
		host, path, origin, remote string
		code                       int
	}{
		{testHost, testPrefix, "", "127.0.0.1:5555", 200},
		{"foreign.example", testPrefix, "", "127.0.0.1:5555", 404},
		{testHost, "/other/", "", "127.0.0.1:5555", 404},
		{testHost, testPrefix, "", "192.0.2.20:5555", 404},
		{testHost, testPrefix, "https://foreign.example", "127.0.0.1:5555", 403},
	} {
		w := makeRequest(t, d, "GET", v.path, "", v.host, v.origin, v.remote)
		if w.Code != v.code {
			t.Fatalf("%+v got %d", v, w.Code)
		}
		if w.Header().Get("Cache-Control") != "no-store" || !strings.Contains(w.Header().Get("Content-Security-Policy"), "frame-ancestors 'none'") {
			t.Fatal("missing security headers")
		}
	}
}
func TestRunAndPrivateReportExports(t *testing.T) {
	d := newTestDesk(true)
	req := func(method, path, body string) *httptest.ResponseRecorder {
		return makeRequest(t, d, method, testPrefix+path, body, testHost, "", "127.0.0.1:5555")
	}
	for _, bad := range []string{`{"id":"disk","lang":"en","command":"rm"}`, `{"id":"disk","lang":"en"} {}`} {
		if w := req("POST", "api/run", bad); w.Code != 400 {
			t.Fatalf("bad request accepted: %s", w.Body.String())
		}
	}
	if w := req("POST", "api/run", `{"id":"other","lang":"en"}`); w.Code != 404 {
		t.Fatalf("unknown task status %d", w.Code)
	}
	var out resultView
	w := req("POST", "api/run", `{"id":"disk","lang":"en"}`)
	if w.Code != 200 || json.Unmarshal(w.Body.Bytes(), &out) != nil || out.ID == "" || out.Status != "Demo" || len(out.Rows) == 0 {
		t.Fatalf("bad demo result: %d %s", w.Code, w.Body.String())
	}
	for _, f := range []string{"txt", "pdf"} {
		x := req("GET", "api/report?id="+out.ID+"&format="+f, "")
		if x.Code != 200 || !strings.Contains(x.Header().Get("Content-Disposition"), "."+f) {
			t.Fatalf("bad %s export: %d", f, x.Code)
		}
		if f == "pdf" && !strings.HasPrefix(x.Body.String(), "%PDF") {
			t.Fatal("bad PDF signature")
		}
		if f == "txt" && !strings.Contains(x.Body.String(), "Fictional values") {
			t.Fatal("bad text export")
		}
	}
	if w := req("GET", "api/report?id=other&format=pdf", ""); w.Code != 404 {
		t.Fatalf("unknown report status %d", w.Code)
	}
}
func TestLivePathCallsOnlyCatalogueTask(t *testing.T) {
	d := newTestDesk(false)
	calls := 0
	d.runTask = func(ctx context.Context, id, lang string) report.Result {
		calls++
		if id != "memory" || lang != "pt" {
			t.Fatalf("unexpected route %s %s", id, lang)
		}
		if _, ok := ctx.Deadline(); !ok {
			t.Fatal("missing task deadline")
		}
		return report.Result{Title: "Memória", Status: "Concluído", Lang: lang, Started: time.Now()}
	}
	w := makeRequest(t, d, http.MethodPost, testPrefix+"api/run", `{"id":"memory","lang":"pt"}`, testHost, "", "127.0.0.1:5555")
	if calls != 1 || w.Code != 200 {
		t.Fatalf("live path called %d times, status %d", calls, w.Code)
	}
}
