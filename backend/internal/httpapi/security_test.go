package httpapi

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"net/netip"
	"strings"
	"sync"
	"testing"
	"time"
)

func TestRequestBoundary(t *testing.T) {
	for _, tc := range []struct {
		name, contentType, body string
		status                  int
		chunked                 bool
	}{
		{"valid", "application/json; charset=UTF-8", `{"name":"O'Connor"}`, 204, false},
		{"wrong type", "text/plain", `{}`, 415, false},
		{"wrong charset", "application/json; charset=latin1", `{}`, 415, false},
		{"malformed", "application/json", `{"name":`, 400, false},
		{"duplicate", "application/json", `{"name":"one","na\u006de":"two"}`, 400, false},
		{"nested duplicate", "application/json", `{"nested":{"key":1,"key":2}}`, 400, false},
		{"trailing data", "application/json", `{} {}`, 400, false},
		{"array root", "application/json", `[]`, 400, false},
		{"null root", "application/json", `null`, 400, false},
		{"invalid UTF-8", "application/json", "{\"name\":\"\xff\"}", 400, false},
		{"unpaired high surrogate", "application/json", `{"name":"\ud800"}`, 400, false},
		{"unpaired low surrogate", "application/json", `{"name":"\udfff"}`, 400, false},
		{"valid surrogate pair", "application/json", `{"name":"\ud83d\ude00"}`, 204, false},
		{"escaped backslash", "application/json", `{"name":"\\ud800"}`, 204, false},
		{"too deep", "application/json", `{"nested":` + strings.Repeat("[", 17) + "0" + strings.Repeat("]", 17) + "}", 400, false},
		{"too large", "application/json", `{"name":"` + strings.Repeat("x", maxBodyBytes) + `"}`, 413, false},
		{"chunked too large", "application/json", `{"name":"` + strings.Repeat("x", maxBodyBytes) + `"}`, 413, true},
	} {
		t.Run(tc.name, func(t *testing.T) {
			calls := 0
			handler := secureHandler(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) { calls++; w.WriteHeader(204) }), SecurityOptions{})
			req := httptest.NewRequest(http.MethodPost, "/api/admin/users", strings.NewReader(tc.body))
			req.Header.Set("Content-Type", tc.contentType)
			if tc.chunked {
				req.ContentLength = -1
			}
			rec := httptest.NewRecorder()
			handler.ServeHTTP(rec, req)
			if rec.Code != tc.status {
				t.Fatalf("got %d: %s", rec.Code, rec.Body.String())
			}
			if (calls == 1) != (tc.status == 204) {
				t.Fatal("invalid input reached handler")
			}
		})
	}
}

func TestClientIP(t *testing.T) {
	trusted := []netip.Prefix{netip.MustParsePrefix("10.0.0.2/32")}
	for _, tc := range []struct{ peer, forwarded, want string }{
		{"192.0.2.1:1000", "203.0.113.1", "192.0.2.1"},
		{"10.0.0.2:1000", "203.0.113.1", "203.0.113.1"},
		{"10.0.0.2:1000", "198.51.100.1, 192.0.2.1", "192.0.2.1"},
		{"10.0.0.2:1000", "invalid", "10.0.0.2"},
		{"[::ffff:192.0.2.1]:1000", "203.0.113.1", "192.0.2.1"},
	} {
		req := httptest.NewRequest(http.MethodGet, "/", nil)
		req.RemoteAddr = tc.peer
		req.Header.Set("X-Forwarded-For", tc.forwarded)
		if got := clientIP(req, trusted); got != tc.want {
			t.Errorf("%s / %s: got %s, want %s", tc.peer, tc.forwarded, got, tc.want)
		}
	}
}

func TestLimiterRefillAndBound(t *testing.T) {
	now := time.Now()
	l := newLimiter(10, 2)
	l.now = func() time.Time { return now }
	first, second, third := l.allow("client"), l.allow("client"), l.allow("client")
	if first != 0 || second != 0 || third != 6 {
		t.Fatal("burst or Retry-After incorrect")
	}
	now = now.Add(6 * time.Second)
	if l.allow("client") != 0 {
		t.Fatal("bucket did not refill")
	}
	for i := range 9999 {
		l.allow(fmt.Sprint(i))
	}
	if l.allow("new-client") != 60 || len(l.buckets) != 10000 {
		t.Fatal("bucket capacity is unbounded or fails open")
	}
	now = now.Add(11 * time.Minute)
	if l.allow("new-client") != 0 || len(l.buckets) != 1 {
		t.Fatal("idle buckets did not expire")
	}
}

func TestLimiterConcurrent(t *testing.T) {
	l := newLimiter(1, 5)
	var wg sync.WaitGroup
	for range 50 {
		wg.Go(func() { l.allow("same-client") })
	}
	wg.Wait()
	if l.allow("same-client") == 0 {
		t.Fatal("concurrent attempts exceeded burst")
	}
}

func TestAuthenticationAccountQuota(t *testing.T) {
	calls := 0
	handler := secureHandler(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) { calls++; w.WriteHeader(204) }), SecurityOptions{})
	for i := range 6 {
		// Changing source IP and casing must not reset the account quota.
		email := "person@example.com"
		if i%2 == 0 {
			email = " Person@EXAMPLE.com "
		}
		req := httptest.NewRequest(http.MethodPost, "/api/auth/login", strings.NewReader(fmt.Sprintf(`{"email":%q,"password":"password"}`, email)))
		req.RemoteAddr = fmt.Sprintf("192.0.2.%d:1000", i+1)
		req.Header.Set("Content-Type", "application/json")
		rec := httptest.NewRecorder()
		handler.ServeHTTP(rec, req)
		if i < 5 && rec.Code != 204 {
			t.Fatal("valid attempt throttled prematurely")
		}
		if i == 5 && (rec.Code != 429 || rec.Header().Get("Retry-After") == "") {
			t.Fatal("account quota not enforced")
		}
	}
	if calls != 5 {
		t.Fatal("blocked attempt reached expensive work")
	}
}

func TestIPQuotaAndHealthExemption(t *testing.T) {
	handler := secureHandler(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) { w.WriteHeader(204) }), SecurityOptions{})
	for i := range 61 {
		req := httptest.NewRequest(http.MethodGet, "/api/admin/users", nil)
		req.Header.Set("X-Forwarded-For", fmt.Sprintf("192.0.2.%d", i))
		rec := httptest.NewRecorder()
		handler.ServeHTTP(rec, req)
		if i == 60 && rec.Code != 429 {
			t.Fatal("IP quota bypassed with forged header")
		}
	}
	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, httptest.NewRequest(http.MethodGet, "/api/ready", nil))
	if rec.Code != 429 {
		t.Fatal("public readiness bypassed rate limit")
	}
	rec = httptest.NewRecorder()
	probe := httptest.NewRequest(http.MethodGet, "/api/health", nil)
	probe.RemoteAddr = "127.0.0.1:1000"
	handler.ServeHTTP(rec, probe)
	if rec.Code != 204 || rec.Header().Get("Cache-Control") != "no-store" {
		t.Fatal("health blocked or API cache protection missing")
	}
}
