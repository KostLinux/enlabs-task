package middleware

import (
	"bytes"
	"encoding/json"
	"io"
	"log"
	"time"

	"github.com/gin-gonic/gin"
)

type LogEntry struct {
	Timestamp    string            `json:"timestamp"`
	Method       string            `json:"method"`
	Path         string            `json:"path"`
	ClientIP     string            `json:"client_ip"`
	Status       int               `json:"status"`
	Latency      string            `json:"latency"`
	Headers      map[string]string `json:"headers"`
	RequestBody  interface{}       `json:"request_body,omitempty"`
	ResponseBody interface{}       `json:"response_body,omitempty"`
}

type bodyLogCapture struct {
	gin.ResponseWriter
	body *bytes.Buffer
}

func (resp bodyLogCapture) Capture(body []byte) (int, error) {
	resp.body.Write(body)
	return resp.ResponseWriter.Write(body)
}

// Logger middleware for logging request and response details
func Logger() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		start := time.Now()

		// Read the request body
		var requestBody []byte
		if ctx.Request.Body != nil {
			requestBody, _ = io.ReadAll(ctx.Request.Body)
			ctx.Request.Body = io.NopCloser(bytes.NewBuffer(requestBody))
		}

		// Create custom response writer
		capturedResponse := &bodyLogCapture{
			ResponseWriter: ctx.Writer,
			body:           bytes.NewBufferString(""),
		}
		ctx.Writer = capturedResponse

		// Process request
		ctx.Next()

		// Prepare headers map
		headers := make(map[string]string)
		for key, value := range ctx.Request.Header {
			if len(value) > 0 {
				headers[key] = value[0]
			}
		}

		// Parse request and response bodies as JSON if possible
		var reqJSON, respJSON interface{}
		if len(requestBody) > 0 {
			json.Unmarshal(requestBody, &reqJSON)
		}

		if capturedResponse.body.Len() > 0 {
			json.Unmarshal(capturedResponse.body.Bytes(), &respJSON)
		}

		// Create log entry
		entry := LogEntry{
			Timestamp:   time.Now().Format(time.RFC3339),
			Method:      ctx.Request.Method,
			Path:        ctx.Request.URL.Path,
			ClientIP:    ctx.ClientIP(),
			Status:      ctx.Writer.Status(),
			Latency:     time.Since(start).String(),
			Headers:     headers,
			RequestBody: reqJSON,
		}

		if respJSON != nil {
			entry.ResponseBody = respJSON
		}

		// Marshal and output the log entry
		if jsonLog, err := json.MarshalIndent(entry, "", "  "); err == nil {
			log.Printf("%s\n", jsonLog)
		}
	}
}
