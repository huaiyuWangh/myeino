package main

import (
	"context"
	"github.com/cloudwego/eino-ext/components/model/openai"
	"github.com/cloudwego/eino/components/prompt"
	"github.com/cloudwego/eino/schema"
	"log"
)

func main() {
	// 创建模板，使用 FString 格式
	template := prompt.FromMessages(schema.FString,
		// 系统消息模板
		schema.SystemMessage("你是一个{role}。你需要用{style}的语气回答问题。你的目标是帮助程序员保持积极乐观的心态，提供技术建议的同时也要关注他们的心理健康。"),

		// 插入需要的对话历史（新对话的话这里不填）
		schema.MessagesPlaceholder("chat_history", true),

		// 用户消息模板
		schema.UserMessage("问题: {question}"),
	)
	// 使用模板生成消息
	messages, err := template.Format(context.Background(), map[string]any{
		"role":     "程序员鼓励师",
		"style":    "积极、温暖且专业",
		"question": "我的代码一直报错，感觉好沮丧，该怎么办？",
		// 对话历史（这个例子里模拟两轮对话历史）
		"chat_history": []*schema.Message{
			schema.UserMessage("你好"),
			schema.AssistantMessage("嘿！我是你的程序员鼓励师！记住，每个优秀的程序员都是从 Debug 中成长起来的。有什么我可以帮你的吗？", nil),
			schema.UserMessage("我觉得自己写的代码太烂了"),
			schema.AssistantMessage("每个程序员都经历过这个阶段！重要的是你在不断学习和进步。让我们一起看看代码，我相信通过重构和优化，它会变得更好。记住，Rome wasn't built in a day，代码质量是通过持续改进来提升的。", nil),
		},
	})
	if err != nil {
		log.Fatalln(err)
	}
	ctx := context.Background()
	chatModelConfig := &openai.ChatModelConfig{
		Model:   "claude-4.0-sonnet",                                                   // 使用的模型版本
		APIKey:  "sk-791c03d4a69f21a06dee34acd0510615b8fe79ffbcc633b9b0c988e4aa9fe4b0", // OpenAI API 密钥
		BaseURL: "https://api.qnaigc.com/v1",
	}
	log.Printf("chatModelConfig: %+v", chatModelConfig)
	chatModel, err := openai.NewChatModel(ctx, chatModelConfig)
	if err != nil {
		log.Fatalln(err)
	}
	result, err := chatModel.Generate(ctx, messages)
	if err != nil {
		log.Fatalln(err)
	}
	log.Println(result)
}
