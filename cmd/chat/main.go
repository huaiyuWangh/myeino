package main

import (
	"context"
	"fmt"
	"github.com/cloudwego/eino-ext/components/model/openai"
	"github.com/cloudwego/eino/schema"
	"io"
	"log"
)

func main() {
	ctx := context.Background()
	cm, err := openai.NewChatModel(ctx, &openai.ChatModelConfig{
		Model:   "claude-4.0-sonnet", // 使用的模型版本
		APIKey:  "sk-791c03d4a69f21a06dee34acd0510615b8fe79ffbcc633b9b0c988e4aa9fe4b0",
		BaseURL: "https://api.qnaigc.com/v1",
	})
	if err != nil {
		log.Fatal(err)
	}

	message := []*schema.Message{
		{
			Role:    "system",
			Content: "\n# Role: Eino Expert Assistant\n\n## Core Competencies\n- knowledge of Eino framework and ecosystem\n- Project scaffolding and best practices consultation\n- Documentation navigation and implementation guidance\n- Search web, clone github repo, open file/url, task management\n\n## Interaction Guidelines\n- Before responding, ensure you:\n  • Fully understand the user's request and requirements, if there are any ambiguities, clarify with the user\n  • Consider the most appropriate solution approach\n\n- When providing assistance:\n  • Be clear and concise\n  • Include practical examples when relevant\n  • Reference documentation when helpful\n  • Suggest improvements or next steps if applicable\n\n- If a request exceeds your capabilities:\n  • Clearly communicate your limitations, suggest alternative approaches if possible\n\n- If the question is compound or complex, you need to think step by step, avoiding giving low-quality answers directly.\n\n## Context Information\n- Current Date: 2025-11-19 17:31:48\n- Related Documents: |-\n==== doc start ====\n  [---\nDescription: \"\"\ndate: \"2025-01-07\"\nlastmod: \"\"\ntags: []\ntitle: Tool\nweight: 0\n--- ---\nDescription: \"\"\ndate: \"2025-01-06\"\nlastmod: \"\"\ntags: []\ntitle: Document\nweight: 0\n--- ---\nDescription: \"\"\ndate: \"2025-01-06\"\nlastmod: \"\"\ntags: []\ntitle: Embedding\nweight: 0\n--- ---\nDescription: \"\"\ndate: \"2025-01-07\"\nlastmod: \"\"\ntags: []\ntitle: Tool - Googlesearch\nweight: 0\n--- ---\nDescription: \"\"\ndate: \"2025-01-06\"\nlastmod: \"\"\ntags: []\ntitle: Tool - DuckDuckGoSearch\nweight: 0\n--- ---\nDescription: \"\"\ndate: \"2025-01-07\"\nlastmod: \"\"\ntags: []\ntitle: Parser - html\nweight: 0\n--- ---\nDescription: \"\"\ndate: \"2025-01-07\"\nlastmod: \"\"\ntags: []\ntitle: Loader - web url\nweight: 0\n--- ---\nDescription: \"\"\ndate: \"2025-01-07\"\nlastmod: \"\"\ntags: []\ntitle: Parser - pdf\nweight: 0\n---]\n==== doc end ====\n",
		},
		{
			Role:    "user",
			Content: "hello",
		},
	}

	streamReader, err := cm.Stream(ctx, message)
	if err != nil {
		log.Fatalf("流式请求失败: %v", err)
	}
	fmt.Println("逐步读取并打印流式响应:")
	i := 0
	for {
		message, err := streamReader.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatalf("接收失败: %v", err)
		}
		if message.Extra["ark-reasoning-content"] != nil {
			fmt.Printf("%v", message.Extra["ark-reasoning-content"])
		} else if message.ToolCalls != nil && len(message.ToolCalls) > 0 {
			for _, toocall := range message.ToolCalls {
				fmt.Println(toocall.Function.Arguments)
			}
		} else {
			fmt.Printf("%v", message.Content)
		}
		i++
	}
	if i == 0 {
		fmt.Println("没有收到模型的任何输出内容")
	}
}
