package handler

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/MGT06/EventHub_Backend.git/internal/dto"
	"github.com/MGT06/EventHub_Backend.git/internal/service"
	"github.com/gin-gonic/gin"
)

type StreamHandler struct {
	hub    *service.Hub
	notify *service.NotificationService
}

func NewStreamHandler(hub *service.Hub, notify *service.NotificationService) *StreamHandler {
	return &StreamHandler{hub: hub, notify: notify}
}

func writeEvent(w io.Writer, n dto.Notification) error {
	b, err := json.Marshal(n)
	if err != nil {
		return err
	}
	_, err = fmt.Fprintf(w, "id: %d\nevent: notification\ndata: %s\n\n", n.Id, b)
	return err
}

func (h *StreamHandler) Stream(c *gin.Context) {
	userId := c.GetInt("idUser")

	flusher, ok := c.Writer.(http.Flusher)
	if !ok {
		c.JSON(http.StatusInternalServerError, dto.ErrorResponse{
			Success: false,
			Message: "streamer not supported",
		})
		return
	}

	hdr := c.Writer.Header()
	hdr.Set("Content-Type", "text/event-stream")
	hdr.Set("Cache-Control", "no-cache")
	hdr.Set("Connection", "keep-alive")
	hdr.Set("X-Accel-Buffering", "no") 
	c.Status(http.StatusOK)
	flusher.Flush()

	client := h.hub.Register(userId)
	defer h.hub.Unregister(userId, client)

	var lastSent int
	if v := c.GetHeader("Last-Event-ID"); v != "" {
		if id, err := strconv.Atoi(v); err == nil && id > 0 {
			lastSent = id
			missed, err := h.notify.ListAfterID(c.Request.Context(), userId, id) 
			if err != nil {
				log.Printf("sse catch-up user %d: %v", userId, err)
			}
			for _, n := range missed {
				if err := writeEvent(c.Writer, n); err != nil {
					return
				}
				lastSent = n.Id
			}
			flusher.Flush()
		}
	}

	heartbeat := time.NewTicker(20 * time.Second)
	defer heartbeat.Stop()
	ctx := c.Request.Context()

	for {
		select {
		case <-ctx.Done():
			return

		case n, ok := <-client.Ch:
			if !ok {
				return
			}
			if n.Id <= lastSent { 
				continue
			}
			if err := writeEvent(c.Writer, n); err != nil {
				return
			}
			lastSent = n.Id
			flusher.Flush()

		case <-heartbeat.C:
			if _, err := fmt.Fprint(c.Writer, ": ping\n\n"); err != nil {
				return
			}
			flusher.Flush()
		}
	}
}
