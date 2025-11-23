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

outer:
	for {
		select {
		case <-ctx.Done():
			log.Printf("[Chat] Context done for chat ID: %s\n", id)
			return
		default:
			msg, err := sr.Recv()
			if errors.Is(err, io.EOF) {
				log.Printf("[Chat] EOF received for chat ID: %s\n", id)
				break outer
			}
			if err != nil {
				log.Printf("[Chat] Error receiving message: %v\n", err)
				break outer
			}

			err = s.Publish(&sse.Event{
				Data: []byte(msg.Content),
			})
			if err != nil {
				log.Printf("[Chat] Error publishing message: %v\n", err)
				break outer
			}
		}
	}
}
