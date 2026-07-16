package shared

import (
	"net"
	"net/http"
	"strings"
)

func GetClientIP(r *http.Request) string {
	// 1. Check standard proxy headers
	if xff := r.Header.Get("X-Forwarded-For"); xff != "" {
		// X-Forwarded-For can be a comma-separated list;
		// strings.Cut returns everything before the first comma in 'before'
		before, _, _ := strings.Cut(xff, ",")

		return strings.TrimSpace(before)
	}

	if realIP := r.Header.Get("X-Real-IP"); realIP != "" {
		return realIP
	}

	// 2. Fallback to RemoteAddr (strip the port number)
	result, _, errSplit := net.SplitHostPort(r.RemoteAddr)
	if errSplit != nil {
		return r.RemoteAddr
	}

	return result
}
