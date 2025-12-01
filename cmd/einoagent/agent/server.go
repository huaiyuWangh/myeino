package agent

import (
	"context"
	"errors"
	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/protocol/consts"
	"github.com/cloudwego/hertz/pkg/route"
	"github.com/hertz-contrib/sse"
	"io"
	"log"
	"strings"
	"time"
)

func BindRoutes(r *route.RouterGroup) error {
	// API 路由
	r.GET("/api/chat", HandleChat)
	return nil
}

func HandleChat(ctx context.Context, c *app.RequestContext) {
	id := c.Query("id")
	message := c.Query("message")
	if id == "" || message == "" {
		c.JSON(consts.StatusBadRequest, map[string]string{
			"status": "error",
			"error":  "missing id or message parameter",
		})
		return
	}

	log.Printf("[Chat] Starting chat with ID: %s, Message: %s\n", id, message)

	sr, err := RunAgent(ctx, id, message)
	if err != nil {
		log.Printf("[Chat] Error running agent: %v\n", err)
		log.Printf("[Chat] Error type: %T\n", err)

		// 尝试获取更详细的错误信息
		if errString := err.Error(); len(errString) > 0 {
			log.Printf("[Chat] Full error string: %s\n", errString)
		}

		// 如果是复合错误，尝试展开
		if unwrapper, ok := err.(interface{ Unwrap() error }); ok {
			if innerErr := unwrapper.Unwrap(); innerErr != nil {
				log.Printf("[Chat] Inner error: %v (type: %T)\n", innerErr, innerErr)
			}
		}

		c.JSON(consts.StatusInternalServerError, map[string]string{
			"status": "error",
			"error":  err.Error(),
		})
		return
	}

	s := sse.NewStream(c)
	defer func() {
		sr.Close()
		c.Flush()

		log.Printf("[Chat] Finished chat with ID: %s\n", id)
	}()

	// 消息缓冲区
	var buffer strings.Builder
	ticker := time.NewTicker(100 * time.Millisecond) // 每100ms刷新一次缓冲区
	defer ticker.Stop()

	flushBuffer := func() error {
		if buffer.Len() > 0 {
			content := buffer.String()
			buffer.Reset()
			return s.Publish(&sse.Event{
				Data: []byte(content),
			})
		}
		return nil
	}

outer:
	for {
		select {
		case <-ctx.Done():
			log.Printf("[Chat] Context done for chat ID: %s\n", id)
			// 发送剩余缓冲内容
			if err := flushBuffer(); err != nil {
				log.Printf("[Chat] Error flushing final buffer: %v\n", err)
			}
			return
		case <-ticker.C:
			// 定时刷新缓冲区
			if err := flushBuffer(); err != nil {
				log.Printf("[Chat] Error flushing buffer: %v\n", err)
				break outer
			}
		default:
			msg, err := sr.Recv()
			if errors.Is(err, io.EOF) {
				log.Printf("[Chat] EOF received for chat ID: %s\n", id)
				// 发送剩余缓冲内容
				if err := flushBuffer(); err != nil {
					log.Printf("[Chat] Error flushing final buffer on EOF: %v\n", err)
				}
				break outer
			}
			if err != nil {
				log.Printf("[Chat] Error receiving message: %v\n", err)
				break outer
			}

			// 将消息添加到缓冲区
			buffer.WriteString(msg.Content)

			// 如果缓冲区过大或遇到句子结束标记，立即刷新
			if buffer.Len() > 200 || strings.HasSuffix(msg.Content, ".") ||
				strings.HasSuffix(msg.Content, "!") || strings.HasSuffix(msg.Content, "?") ||
				strings.HasSuffix(msg.Content, "\n") {
				if err := flushBuffer(); err != nil {
					log.Printf("[Chat] Error publishing message: %v\n", err)
					break outer
				}
			}
		}
	}
}
