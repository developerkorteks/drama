package middleware

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

// DynamicSwaggerHost middleware untuk mendeteksi host secara dinamis
func DynamicSwaggerHost() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Hanya untuk endpoint swagger
		if strings.HasPrefix(c.Request.URL.Path, "/swagger/") || strings.HasPrefix(c.Request.URL.Path, "/docs/") {
			// Deteksi host dari request
			host := c.Request.Host
			scheme := "http"

			// Deteksi HTTPS
			if c.Request.TLS != nil ||
				c.GetHeader("X-Forwarded-Proto") == "https" ||
				c.GetHeader("X-Forwarded-Ssl") == "on" {
				scheme = "https"
			}

			// Set dynamic host untuk swagger
			c.Set("swagger_host", host)
			c.Set("swagger_scheme", scheme)
		}

		c.Next()
	}
}

// SwaggerConfigHandler untuk menyediakan konfigurasi swagger yang dinamis
func SwaggerConfigHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		host := c.Request.Host
		scheme := "http"

		// Deteksi HTTPS
		if c.Request.TLS != nil ||
			c.GetHeader("X-Forwarded-Proto") == "https" ||
			c.GetHeader("X-Forwarded-Ssl") == "on" {
			scheme = "https"
		}

		config := map[string]interface{}{
			"swagger": "2.0",
			"info": map[string]interface{}{
				"title":       "DramaQu API",
				"description": "API untuk scraping data drama Korea dari dramaqu.ad",
				"version":     "1.0",
			},
			"host":     host,
			"basePath": "/",
			"schemes":  []string{scheme},
		}

		c.Header("Content-Type", "application/json")
		c.JSON(http.StatusOK, config)
	}
}
