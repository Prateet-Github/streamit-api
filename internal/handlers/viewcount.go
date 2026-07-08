package handlers

import (
	"net/http"
	"time"

	"github.com/Prateet-Github/streamit-api/internal/models"
	"github.com/Prateet-Github/streamit-api/internal/services/viewcount"
	"github.com/gin-gonic/gin"
)

type ViewCountHandler struct {
	producer *viewcount.Producer
}

func NewViewCountHandler(producer *viewcount.Producer) *ViewCountHandler {
	return &ViewCountHandler{
		producer: producer,
	}
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
	viewerID := c.GetString("userId")

	if viewerID == "" {
		viewerID = req.ViewerID
	}

	err := h.producer.Publish(
		c.Request.Context(),
		viewcount.ViewEvent{
			VideoID:   videoID,
			ViewerID:  viewerID,
			Elapsed:   req.Elapsed,
			Timestamp: time.Now(),
		},
	)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "failed to enqueue heartbeat",
		})
		return
	}

	c.Status(http.StatusAccepted)

}
