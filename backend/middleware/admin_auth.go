package middleware

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/Seeyko/casamia-api/services"
	"github.com/Seeyko/casamia-api/services/ratelimit"
)

type AdminAuthMiddleware struct {
	jwt        *services.JWTService
	allowedIPs []string
	limiter    *ratelimit.RateLimiter
}

func NewAdminAuthMiddleware(jwt *services.JWTService, allowedIPs []string) *AdminAuthMiddleware {
	limiter := ratelimit.NewRateLimiter(
		10,
		20,
		15*time.Minute,
		5*time.Minute,
	)
	return &AdminAuthMiddleware{
		jwt:        jwt,
		allowedIPs: allowedIPs,
		limiter:    limiter,
	}
}

func setSecurityHeaders(w http.ResponseWriter) {
	w.Header().Set("Cache-Control", "no-store")
	w.Header().Set("X-Content-Type-Options", "nosniff")
	w.Header().Set("Referrer-Policy", "no-referrer")
	w.Header().Set("X-Frame-Options", "DENY")
}

func writeUnauthorized(w http.ResponseWriter) {
	setSecurityHeaders(w)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusUnauthorized)
	json.NewEncoder(w).Encode(map[string]string{"error": "unauthorized"})
}

func (m *AdminAuthMiddleware) Authenticate(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ip := ExtractIP(r)
		setSecurityHeaders(w)

		if len(m.allowedIPs) > 0 {
			allowed := false
			for _, allowedIP := range m.allowedIPs {
				if ip == allowedIP {
					allowed = true
					break
				}
			}
			if !allowed {
				log.Printf("[ADMIN-AUTH] IP not in allowlist: %s", ip)
				time.Sleep(200 * time.Millisecond)
				writeUnauthorized(w)
				return
			}
		}

		allowed, _, _ := m.limiter.CheckRequest(ip, "admin")
		if !allowed {
			time.Sleep(200 * time.Millisecond)
			writeUnauthorized(w)
			return
		}

		authHeader := r.Header.Get("Authorization")
		if authHeader == "" || !strings.HasPrefix(authHeader, "Bearer ") {
			m.limiter.RecordFailure(ip, "admin")
			time.Sleep(200 * time.Millisecond)
			writeUnauthorized(w)
			return
		}

		tokenStr := strings.TrimPrefix(authHeader, "Bearer ")

		claims, err := m.jwt.ValidateToken(tokenStr)
		if err != nil {
			m.limiter.RecordFailure(ip, "admin")
			log.Printf("[ADMIN-AUTH] Invalid JWT: IP=%s, err=%v", ip, err)
			time.Sleep(200 * time.Millisecond)
			writeUnauthorized(w)
			return
		}

		m.limiter.RecordSuccess(ip, "admin")

		// Store claims in context for handlers
		ctx := context.WithValue(r.Context(), "claims", claims)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func ExtractIP(r *http.Request) string {
	if xff := r.Header.Get("X-Forwarded-For"); xff != "" {
		ips := strings.Split(xff, ",")
		return strings.TrimSpace(ips[0])
	}
	if xri := r.Header.Get("X-Real-IP"); xri != "" {
		return xri
	}
	ip := r.RemoteAddr
	if colon := strings.LastIndex(ip, ":"); colon != -1 {
		ip = ip[:colon]
	}
	return ip
}
