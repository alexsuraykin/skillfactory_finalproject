package api

import (
	"context"
	"time"

	"math/rand"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
)

const (
	requestIDKey = "request_id"
	charset      = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
)

func requestIDMiddleware(c *gin.Context) {
	requestID := c.Query("request_id")
	if requestID == "" {
		requestID = generateRandomID(6)
	}

	ctx := context.WithValue(c.Request.Context(), requestIDKey, requestID)
	c.Request = c.Request.WithContext(ctx)

	c.Next()
}

func loggingMiddleware(c *gin.Context) {
	startTime := time.Now()

	c.Next()

	requestID, _ := c.Request.Context().Value(requestIDKey).(string)

	log.Printf("Time: %s, IP: %s, Status: %d, Request ID: %s",
		startTime.Format(time.RFC3339),
		c.ClientIP(),
		c.Writer.Status(),
		requestID,
	)
}

func generateRandomID(length int) string {
	seededRand := rand.New(rand.NewSource(time.Now().UnixNano()))
	b := make([]byte, length)
	for i := range b {
		b[i] = charset[seededRand.Intn(len(charset))]
	}
	return string(b)
}
