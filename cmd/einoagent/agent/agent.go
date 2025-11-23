package agent

import (
	"context"
	"errors"
	"fmt"
	"github.com/cloudwego/eino-examples/quickstart/eino_assistant/pkg/mem"
	"github.com/cloudwego/eino/compose"
	"github.com/cloudwego/eino/schema"
	"io"
	"myeino/agent"
)

var memory = mem.GetDefaultMemory()

func RunAgent(ctx context.Context, id string, msg string) (*schema.StreamReader[*schema.Message], error) {
	runner, err := agent.BuildEinoAgent(ctx)
	if err != nil {
		return nil, err
	}

	conversation := memory.GetConversation(id, true)

	userMessage := &agent.UserMessage{
		ID:      id,
		Query:   msg,
		History: conversation.GetMessages(),
	}

	sr, err := runner.Stream(ctx, userMessage, compose.WithCallbacks())
	if err != nil {
		return nil, err
	}

	srs := sr.Copy(2)

	go func() {
		fullMsgs := make([]*schema.Message, 0)

		defer func() {
			// close stream
			srs[1].Close()

			// add history
			conversation.Append(schema.UserMessage(msg))

			fullMsg, err := schema.ConcatMessages(fullMsgs)
			if err != nil {
				fmt.Println(err)
			}

			conversation.Append(fullMsg)
		}()

	outer:
		for {
			select {
			case <-ctx.Done():
				fmt.Println("context done", ctx.Err())
				return
			default:
				chunk, err := srs[1].Recv()
				if err != nil {
					if errors.Is(err, io.EOF) {
						break outer
					}
				}

				fullMsgs = append(fullMsgs, chunk)
			}
		}
	}()

	return sr, nil
}
