package api

import (
	"net/http"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

// corsMiddleware habilita CORS para el dev server del frontend (Vite en :5173).
// En fase 2 (go:embed, mismo origen) esto ya no hará falta.
func corsMiddleware(allowedOrigins []string) gin.HandlerFunc {
	cfg := cors.DefaultConfig()
	cfg.AllowOrigins = allowedOrigins
	cfg.AllowMethods = []string{"GET", "POST", "DELETE", "OPTIONS"}
	cfg.AllowHeaders = []string{"Origin", "Content-Type", "Accept"}
	return cors.New(cfg)
}

// respondError centraliza el formato de error JSON.
func respondError(c *gin.Context, status int, msg string) {
	c.JSON(status, gin.H{"error": msg})
}

// notFound es un helper para 404.
func notFound(c *gin.Context, msg string) {
	respondError(c, http.StatusNotFound, msg)
}
