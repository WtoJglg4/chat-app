package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func (h *Handler) signUp(c *gin.Context) {
	c.JSON(http.StatusOK, map[string]interface{}{
		"sign": "up",
	})
}

func (h *Handler) signIn(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"sign": "In"})
}
