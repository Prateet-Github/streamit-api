package handlers

import (
	"log"
	"net/http"

	"github.com/Prateet-Github/streamit-api/internal/models"
	"github.com/gin-gonic/gin"
)

type ViewCountHandler struct {
}

func NewViewCountHandler() *ViewCountHandler {
	return &ViewCountHandler{}
}

func (h *ViewCountHandler) Heartbeat(c *gin.Context) {
	var req models.ViewHeartbeatRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid payload",
		})
		return
	}

	videoID := c.Param("id")

	log.Println(videoID, req.Elapsed)

	c.Status(http.StatusAccepted)
}
