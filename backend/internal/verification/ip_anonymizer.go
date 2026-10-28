package verification

import (
	"crypto/sha256"
	"encoding/hex"
	"net"
	"net/http"
	"strings"
)

type IPAnonymizer struct {
	salt string
}

func NewIPAnonymizer(salt string) *IPAnonymizer {
	return &IPAnonymizer{salt: salt}
}

// ExtractClientIP extracts the real client IP from Cloudflare, reverse proxies, or RemoteAddr
func ExtractClientIP(r *http.Request) string {
	// 1. Cloudflare connecting IP
	if cfIP := r.Header.Get("CF-Connecting-IP"); cfIP != "" {
		return strings.TrimSpace(cfIP)
	}

	// 2. X-Forwarded-For (leftmost client IP)
	if xff := r.Header.Get("X-Forwarded-For"); xff != "" {
		parts := strings.Split(xff, ",")
		if len(parts) > 0 {
			ip := strings.TrimSpace(parts[0])
			if ip != "" {
				return ip
			}
		}
	}

	// 3. X-Real-IP
	if realIP := r.Header.Get("X-Real-IP"); realIP != "" {
		return strings.TrimSpace(realIP)
	}

	// 4. RemoteAddr fallback
	ip, _, err := net.SplitHostPort(r.RemoteAddr)
	if err == nil && ip != "" {
		return ip
	}

	return r.RemoteAddr
}

// HashIP creates a deterministic SHA-256 hash of the client IP with a server salt
func (a *IPAnonymizer) HashIP(rawIP string) string {
	cleanIP := strings.TrimSpace(rawIP)
	hasher := sha256.New()
	hasher.Write([]byte(a.salt + ":" + cleanIP))
	return hex.EncodeToString(hasher.Sum(nil))
}
