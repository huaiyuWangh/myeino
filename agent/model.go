package agent

import (
	"context"
	"github.com/cloudwego/eino-ext/components/model/openai"
	"github.com/cloudwego/eino/components/model"
	"github.com/cloudwego/eino/schema"
	"log"
)

// LoggingChatModel 是一个包装器，用于记录所有的模型调用参数
type LoggingChatModel struct {
	inner model.ChatModel
}

func (lcm *LoggingChatModel) Generate(ctx context.Context, in []*schema.Message, opts ...model.Option) (*schema.Message, error) {
	log.Printf("[LoggingChatModel] === GENERATE CALL ===")
	log.Printf("[LoggingChatModel] Input Messages: %d", len(in))
	for i, msg := range in {
		log.Printf("[LoggingChatModel] Message[%d]: Role=%s, Content(len=%d)", i, msg.Role, len(msg.Content))
		if len(msg.Content) > 500 {
			log.Printf("[LoggingChatModel] Message[%d] Content Preview: %s...", i, msg.Content[:500])
		} else {
			log.Printf("[LoggingChatModel] Message[%d] Content: %s", i, msg.Content)
		}
	}
	log.Printf("[LoggingChatModel] Options Count: %d", len(opts))
	return lcm.inner.Generate(ctx, in, opts...)
}

func (lcm *LoggingChatModel) Stream(ctx context.Context, in []*schema.Message, opts ...model.Option) (*schema.StreamReader[*schema.Message], error) {
	log.Printf("[LoggingChatModel] === STREAM CALL ===")
	log.Printf("[LoggingChatModel] Input Messages: %d", len(in))
	for i, msg := range in {
		log.Printf("[LoggingChatModel] Message[%d]: Role=%s, Content(len=%d)", i, msg.Role, len(msg.Content))
		if len(msg.Content) > 500 {
			log.Printf("[LoggingChatModel] Message[%d] Content Preview: %s...", i, msg.Content[:500])
		} else {
			log.Printf("[LoggingChatModel] Message[%d] Content: %s", i, msg.Content)
		}
		if len(msg.ToolCalls) > 0 {
			log.Printf("[LoggingChatModel] Message[%d] ToolCalls: %d", i, len(msg.ToolCalls))
		}
	}
	log.Printf("[LoggingChatModel] Options Count: %d", len(opts))
	for i, opt := range opts {
		log.Printf("[LoggingChatModel] Option[%d]: %T", i, opt)
	}

	log.Printf("[LoggingChatModel] === CALLING INNER STREAM ===")
	stream, err := lcm.inner.Stream(ctx, in, opts...)
	if err != nil {
		log.Printf("[LoggingChatModel] === STREAM CALL FAILED ===")
		log.Printf("[LoggingChatModel] Error Type: %T", err)
		log.Printf("[LoggingChatModel] Error: %v", err)
		return nil, err
	}

	log.Printf("[LoggingChatModel] === STREAM CALL SUCCESSFUL ===")
	return stream, nil
}

func newChatModel(ctx context.Context) (cm model.ChatModel, err error) {
	// TODO Modify component configuration here.
	maxTokens := 4096
	config := &openai.ChatModelConfig{Model: "claude-4.0-sonnet", // 使用的模型版本
		APIKey:    "sk-791c03d4a69f21a06dee34acd0510615b8fe79ffbcc633b9b0c988e4aa9fe4b0",
		BaseURL:   "https://api.qnaigc.com/v1",
		MaxTokens: &maxTokens, // 添加 max_tokens 参数，使用指针
	}

	// 详细打印模型配置参数
	log.Printf("[ChatModel] === DETAILED MODEL CONFIGURATION ===")
	log.Printf("[ChatModel] Model Name: %s", config.Model)
	log.Printf("[ChatModel] Base URL: %s", config.BaseURL)
	log.Printf("[ChatModel] API Key (first 10 chars): %s...", config.APIKey[:10])
	if config.MaxTokens != nil {
		log.Printf("[ChatModel] Max Tokens: %d", *config.MaxTokens)
	} else {
		log.Printf("[ChatModel] Max Tokens: Not set")
	}
	log.Printf("[ChatModel] Config Type: %T", config)

	innerCm, err := openai.NewChatModel(ctx, config)
	if err != nil {
		log.Printf("[ChatModel] Model creation failed: %v", err)
		return nil, err
	}

	log.Printf("[ChatModel] Model created successfully, Type: %T", innerCm)

	// 返回包装的日志模型
	cm = &LoggingChatModel{inner: innerCm}
	log.Printf("[ChatModel] Wrapped with LoggingChatModel")
	return cm, nil
}
