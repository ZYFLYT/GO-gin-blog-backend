package utils

import (
	"net"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

// GetClientIP 从 http.Request 获取客户端真实 IP
func GetClientIP(r *http.Request) string {
	xForwardedFor := r.Header.Get("X-Forwarded-For")
	if xForwardedFor != "" {
		ips := strings.Split(xForwardedFor, ",")
		for _, ip := range ips {
			ip = strings.TrimSpace(ip)
			if ip != "" {
				return ip
			}
		}
	}

	xRealIP := r.Header.Get("X-Real-IP")
	if xRealIP != "" {
		return xRealIP
	}

	ip, _, err := net.SplitHostPort(r.RemoteAddr)
	if err != nil {
		return r.RemoteAddr
	}
	return ip
}

// GetClientIPFromGin 从 gin.Context 获取客户端真实 IP
func GetClientIPFromGin(c *gin.Context) string {
	return GetClientIP(c.Request)
}
