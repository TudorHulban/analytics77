package shared

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestGetClientIP(t *testing.T) {
	type testCase struct {
		description string
		setupReq    func() *http.Request
		expected    string
	}

	testCases := []testCase{
		{
			description: "X-Forwarded-For single IP",
			setupReq: func() *http.Request {
				req := httptest.NewRequest(http.MethodGet, "/", nil) //nolint:noctx
				req.Header.Set("X-Forwarded-For", "203.0.113.195")

				return req
			},
			expected: "203.0.113.195",
		},
		{
			description: "X-Forwarded-For multiple IPs (takes first and trims space)",
			setupReq: func() *http.Request {
				req := httptest.NewRequest(http.MethodGet, "/", nil) //nolint:noctx
				req.Header.Set("X-Forwarded-For", " 198.51.100.1, 203.0.113.195, 70.41.3.18 ")

				return req
			},
			expected: "198.51.100.1",
		},
		{
			description: "X-Real-IP fallback when X-Forwarded-For is missing",
			setupReq: func() *http.Request {
				req := httptest.NewRequest(http.MethodGet, "/", nil) //nolint:noctx
				req.Header.Set("X-Real-IP", "192.0.2.1")

				return req
			},
			expected: "192.0.2.1",
		},
		{
			description: "X-Forwarded-For takes precedence over X-Real-IP",
			setupReq: func() *http.Request {
				req := httptest.NewRequest(http.MethodGet, "/", nil) //nolint:noctx
				req.Header.Set("X-Forwarded-For", "203.0.113.195")
				req.Header.Set("X-Real-IP", "192.0.2.1")

				return req
			},
			expected: "203.0.113.195",
		},
		{
			description: "RemoteAddr fallback with port stripped",
			setupReq: func() *http.Request {
				req := httptest.NewRequest(http.MethodGet, "/", nil) //nolint:noctx
				req.RemoteAddr = "192.168.1.50:54321"

				return req
			},
			expected: "192.168.1.50",
		},
		{
			description: "RemoteAddr fallback without port (IPv4)",
			setupReq: func() *http.Request {
				req := httptest.NewRequest(http.MethodGet, "/", nil) //nolint:noctx
				req.RemoteAddr = "10.0.0.1"

				return req
			},
			expected: "10.0.0.1",
		},
		{
			description: "RemoteAddr fallback with IPv6 port parsing",
			setupReq: func() *http.Request {
				req := httptest.NewRequest(http.MethodGet, "/", nil) //nolint:noctx
				req.RemoteAddr = "[2001:db8::1]:443"

				return req
			},
			expected: "2001:db8::1",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.description,
			func(t *testing.T) {
				req := tc.setupReq()
				actual := GetClientIP(req)

				require.Equal(t,
					tc.expected,
					GetClientIP(req),

					"GetClientIP() = %q, want %q",
					actual,
					tc.expected,
				)
			},
		)
	}
}
