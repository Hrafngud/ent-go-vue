package httpapi

import (
	"bytes"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"io"
	"math"
	"mime"
	"net"
	"net/http"
	"net/netip"
	"strconv"
	"strings"
	"sync"
	"time"
	"unicode/utf8"
)

const maxBodyBytes = 16 * 1024

type SecurityOptions struct {
	// Only these socket peers may supply X-Forwarded-For. Empty trusts nobody.
	TrustedProxyCIDRs []netip.Prefix
}

type bucket struct {
	tokens  float64
	updated time.Time
}

// Bounded token buckets expire after ten idle minutes. Capacity exhaustion fails
// closed rather than evicting active buckets and resetting an attacker's quota.
type limiter struct {
	mu          sync.Mutex
	buckets     map[string]bucket
	rate, burst float64
	lastCleanup time.Time
	now         func() time.Time
}

func newLimiter(perMinute, burst int) *limiter {
	return &limiter{buckets: make(map[string]bucket), rate: float64(perMinute) / 60, burst: float64(burst), now: time.Now}
}

func (l *limiter) allow(key string) int {
	l.mu.Lock()
	defer l.mu.Unlock()
	now := l.now()
	if now.Sub(l.lastCleanup) >= time.Minute {
		for k, b := range l.buckets {
			if now.Sub(b.updated) >= 10*time.Minute {
				delete(l.buckets, k)
			}
		}
		l.lastCleanup = now
	}
	b, exists := l.buckets[key]
	if !exists {
		if len(l.buckets) >= 10000 {
			return 60
		}
		b = bucket{tokens: l.burst, updated: now}
	}
	b.tokens = math.Min(l.burst, b.tokens+now.Sub(b.updated).Seconds()*l.rate)
	b.updated = now
	retry := 0
	if b.tokens < 1 {
		retry = int(math.Ceil((1 - b.tokens) / l.rate))
	} else {
		b.tokens--
	}
	l.buckets[key] = b
	return retry
}

func secureHandler(next http.Handler, options SecurityOptions) http.Handler {
	general := newLimiter(120, 60)
	authIP := newLimiter(10, 10)
	account := newLimiter(10, 5)
	// Bound aggregate bcrypt work even if source IPs and account names rotate.
	authGlobal := newLimiter(60, 20)
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("X-Content-Type-Options", "nosniff")
		w.Header().Set("Referrer-Policy", "strict-origin-when-cross-origin")
		w.Header().Set("Cache-Control", "no-store")
		ip := clientIP(r, options.TrustedProxyCIDRs)
		// Internal probes remain available; public readiness requests are limited
		// because each one does database work.
		if internalProbe(r, ip) {
			next.ServeHTTP(w, r)
			return
		}
		if retry := general.allow(ip); retry > 0 {
			tooMany(w, retry)
			return
		}
		isAuth := authenticationRequest(r)
		if isAuth {
			if retry := authIP.allow(ip); retry > 0 {
				tooMany(w, retry)
				return
			}
			if retry := authGlobal.allow("all"); retry > 0 {
				tooMany(w, retry)
				return
			}
		}
		if r.Method == http.MethodPost || r.Method == http.MethodPut || r.Method == http.MethodPatch {
			body, ok := readJSON(w, r)
			if !ok {
				return
			}
			if isAuth && !allowAccount(w, body, account) {
				return
			}
			r.Body = io.NopCloser(bytes.NewReader(body))
			r.ContentLength = int64(len(body))
		}
		next.ServeHTTP(w, r)
	})
}

func readJSON(w http.ResponseWriter, r *http.Request) ([]byte, bool) {
	if r.ContentLength > maxBodyBytes {
		problem(w, http.StatusRequestEntityTooLarge, "request body exceeds 16 KiB")
		return nil, false
	}
	mediaType, params, err := mime.ParseMediaType(r.Header.Get("Content-Type"))
	if err != nil || mediaType != "application/json" || (params["charset"] != "" && !strings.EqualFold(params["charset"], "utf-8")) || r.Header.Get("Content-Encoding") != "" {
		problem(w, http.StatusUnsupportedMediaType, "send UTF-8 application/json without content encoding")
		return nil, false
	}
	r.Body = http.MaxBytesReader(w, r.Body, maxBodyBytes)
	body, err := io.ReadAll(r.Body)
	_ = r.Body.Close()
	if err != nil {
		var limitErr *http.MaxBytesError
		if errors.As(err, &limitErr) {
			problem(w, 413, "request body exceeds 16 KiB")
		} else {
			problem(w, 400, "could not read request body")
		}
		return nil, false
	}
	if !utf8.Valid(body) || !validJSONUnicode(body) || !validJSONObject(body) {
		problem(w, http.StatusBadRequest, "send one JSON object with unique keys and valid UTF-8")
		return nil, false
	}
	return body, true
}

func authenticationRequest(r *http.Request) bool {
	return r.Method == http.MethodPost && (r.URL.Path == "/api/auth/login" || r.URL.Path == "/api/auth/register")
}

// encoding/json replaces unpaired UTF-16 escape surrogates with U+FFFD.
// Reject them rather than silently changing input before validation/hashing.
func validJSONUnicode(body []byte) bool {
	for i := 0; i < len(body); i++ {
		if body[i] != '\\' {
			continue
		}
		i++
		if i+4 >= len(body) || body[i] != 'u' {
			continue
		}
		code, err := strconv.ParseUint(string(body[i+1:i+5]), 16, 16)
		if err != nil {
			return false
		}
		i += 4
		if code >= 0xdc00 && code <= 0xdfff {
			return false
		}
		if code >= 0xd800 && code <= 0xdbff {
			if !validLowSurrogate(body, i) {
				return false
			}
			i += 6
		}
	}
	return true
}

func validLowSurrogate(body []byte, i int) bool {
	if i+6 >= len(body) || body[i+1] != '\\' || body[i+2] != 'u' {
		return false
	}
	low, err := strconv.ParseUint(string(body[i+3:i+7]), 16, 16)
	return err == nil && low >= 0xdc00 && low <= 0xdfff
}

func internalProbe(r *http.Request, ip string) bool {
	address, err := netip.ParseAddr(ip)
	return err == nil && address.IsLoopback() && r.Method == http.MethodGet && (r.URL.Path == "/api/health" || r.URL.Path == "/api/ready")
}

// Token parsing rejects duplicate keys (including escaped equivalents), trailing
// data, non-object roots, and excessive nesting before Huma schema validation.
func validJSONObject(body []byte) bool {
	decoder := json.NewDecoder(bytes.NewReader(body))
	decoder.UseNumber()
	first, err := decoder.Token()
	if err != nil || first != json.Delim('{') {
		return false
	}
	if !validJSONContainer(decoder, json.Delim('{'), 1) {
		return false
	}
	_, err = decoder.Token()
	return errors.Is(err, io.EOF)
}

func validJSONContainer(decoder *json.Decoder, open json.Delim, depth int) bool {
	if depth > 16 {
		return false
	}
	keys := make(map[string]bool)
	for decoder.More() {
		if open == '{' {
			key, err := decoder.Token()
			if err != nil {
				return false
			}
			name, ok := key.(string)
			if !ok || keys[name] {
				return false
			}
			keys[name] = true
		}
		value, err := decoder.Token()
		if err != nil {
			return false
		}
		if delim, ok := value.(json.Delim); ok && !validJSONContainer(decoder, delim, depth+1) {
			return false
		}
	}
	close, err := decoder.Token()
	return err == nil && ((open == '{' && close == json.Delim('}')) || (open == '[' && close == json.Delim(']')))
}

func clientIP(r *http.Request, trusted []netip.Prefix) string {
	host, _, err := net.SplitHostPort(r.RemoteAddr)
	if err != nil {
		host = r.RemoteAddr
	}
	peer, err := netip.ParseAddr(host)
	if err != nil {
		return "unknown"
	}
	peer = peer.Unmap()
	// Walk from the socket peer toward the client, stopping at the first
	// untrusted hop. Never accept a header just because it exists.
	hops := strings.Split(r.Header.Get("X-Forwarded-For"), ",")
	for i := len(hops) - 1; i >= 0 && trustedIP(peer, trusted); i-- {
		forwarded, err := netip.ParseAddr(strings.TrimSpace(hops[i]))
		if err != nil {
			return host
		}
		peer = forwarded.Unmap()
	}
	return peer.String()
}

func trustedIP(ip netip.Addr, trusted []netip.Prefix) bool {
	for _, prefix := range trusted {
		if prefix.Contains(ip) {
			return true
		}
	}
	return false
}

func tooMany(w http.ResponseWriter, retry int) {
	w.Header().Set("Retry-After", strconv.Itoa(retry))
	problem(w, http.StatusTooManyRequests, "too many requests; retry after the indicated delay")
}

func problem(w http.ResponseWriter, status int, detail string) {
	w.Header().Set("Content-Type", "application/problem+json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(map[string]any{"title": http.StatusText(status), "status": status, "detail": detail})
}

func allowAccount(w http.ResponseWriter, body []byte, account *limiter) bool {
	var fields map[string]json.RawMessage
	_ = json.Unmarshal(body, &fields)
	var email string
	_ = json.Unmarshal(fields["email"], &email)
	if email == "" {
		return true
	}
	key := sha256.Sum256([]byte(strings.ToLower(strings.TrimSpace(email))))
	if retry := account.allow(hex.EncodeToString(key[:])); retry > 0 {
		tooMany(w, retry)
		return false
	}
	return true
}
