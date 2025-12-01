package main

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/cloudwego/eino-examples/quickstart/eino_assistant/eino/einoagent"
	"github.com/cloudwego/eino-ext/components/model/openai"
	"github.com/cloudwego/eino-ext/components/tool/duckduckgo/v2"
	"github.com/cloudwego/eino/components/model"
	"github.com/cloudwego/eino/components/prompt"
	"github.com/cloudwego/eino/components/tool"
	"github.com/cloudwego/eino/compose"
	"github.com/cloudwego/eino/flow/agent/react"
	"github.com/cloudwego/eino/schema"
	"io"
	"log"
	"time"
)

type UserMessage struct {
	ID      string            `json:"id"`
	Query   string            `json:"query"`
	History []*schema.Message `json:"history"`
}

func main() {
	ctx := context.Background()
	runner, err := BuildEinoAgent(ctx)
	if err != nil {
		panic(err)
	}
	id := "1"
	sr, err := runner.Stream(ctx, &UserMessage{ID: "1", Query: "hello"})
	if err != nil {
		panic(err)
	}
outer:
	for {
		select {
		case <-ctx.Done():
			fmt.Printf("[Chat] Context done for chat ID: %s\n", id)
			return
		default:
			msg, err := sr.Recv()
			if errors.Is(err, io.EOF) {
				fmt.Printf("[Chat] EOF received for chat ID: %s\n", id)
				break outer
			}
			if err != nil {
				fmt.Printf("[Chat] Error receiving message: %v\n", err)
				break outer
			}
			fmt.Println(msg.Content)
		}
	}

}

func BuildEinoAgent(ctx context.Context) (r compose.Runnable[*UserMessage, *schema.Message], err error) {
	const (
		InputToQuery   = "InputToQuery"
		InputToHistory = "InputToHistory"
		Retriever      = "Retriever"
		ChatTemplate   = "ChatTemplate"
		ReactAgent     = "ReactAgent"
	)
	g := compose.NewGraph[*UserMessage, *schema.Message]()
	_ = g.AddLambdaNode(InputToQuery, compose.InvokableLambdaWithOption(newInputToQuery))
	//_ = g.AddLambdaNode(InputToHistory, compose.InvokableLambdaWithOption(newInputToHistory), compose.WithNodeName("UserMessageToVariables"))
	//retrieverKeyOfRetriever, err := newRetriever(ctx)
	//if err != nil {
	//	//	return nil, err
	//	//}
	//_ = g.AddRetrieverNode(Retriever, retrieverKeyOfRetriever, compose.WithOutputKey("documents"))
	chatTemplateKeyOfChatTemplate, err := newChatTemplate(ctx)
	if err != nil {
		return nil, err
	}
	_ = g.AddChatTemplateNode(ChatTemplate, chatTemplateKeyOfChatTemplate)
	reactAgentKeyOfLambda, err := newReact(ctx)
	if err != nil {
		return nil, err
	}
	_ = g.AddLambdaNode(ReactAgent, reactAgentKeyOfLambda, compose.WithNodeName("ReAct Agent"))
	_ = g.AddEdge(compose.START, InputToQuery)
	_ = g.AddEdge(InputToQuery, ChatTemplate)
	//_ = g.AddEdge(InputToQuery, Retriever)
	//_ = g.AddEdge(Retriever, ChatTemplate)
	//_ = g.AddEdge(compose.START, InputToHistory)
	//_ = g.AddEdge(InputToHistory, ChatTemplate)
	_ = g.AddEdge(ChatTemplate, ReactAgent)
	_ = g.AddEdge(ReactAgent, compose.END)
	r, err = g.Compile(ctx, compose.WithGraphName("EinoAgent"), compose.WithNodeTriggerMode(compose.AllPredecessor))
	if err != nil {
		return nil, err
	}
	return r, err
}

func newInputToQuery(ctx context.Context, input *UserMessage, opts ...any) (output map[string]any, err error) {
	// Log input
	inputJSON, _ := json.MarshalIndent(input, "", "  ")
	fmt.Printf("[InputToQuery] Input: %s", string(inputJSON))
	return map[string]any{
		"documents": nil,
		"content":   input.Query,
		"date":      time.Now().Format(time.DateTime),
	}, nil
}

func GetTools(ctx context.Context) ([]tool.BaseTool, error) {
	einoAssistantTool, err := einoagent.NewEinoAssistantTool(ctx)
	if err != nil {
		return nil, err
	}

	toolTask, err := einoagent.NewTaskTool(ctx)
	if err != nil {
		return nil, err
	}

	toolOpen, err := einoagent.NewOpenFileTool(ctx)
	if err != nil {
		return nil, err
	}

	toolGitClone, err := einoagent.NewGitCloneFile(ctx)
	if err != nil {
		return nil, err
	}

	_, err = duckduckgo.NewTextSearchTool(ctx, &duckduckgo.Config{})
	return []tool.BaseTool{
		einoAssistantTool,
		toolTask,
		toolOpen,
		toolGitClone,
		//toolDDGSearch,
	}, nil
}

func newReact(ctx context.Context) (lba *compose.Lambda, err error) {
	// TODO Modify component configuration here.
	config := &react.AgentConfig{}
	chatModelIns11, err := newChatModel(ctx)
	if err != nil {
		return nil, err
	}
	config.Model = chatModelIns11
	tools, err := GetTools(ctx)
	if err != nil {
		return nil, err
	}
	config.ToolsConfig.Tools = tools
	ins, err := react.NewAgent(ctx, config)
	if err != nil {
		return nil, err
	}
	lba, err = compose.AnyLambda(ins.Generate, ins.Stream, nil, nil)
	if err != nil {
		return nil, err
	}
	return lba, nil
}

func newChatModel(ctx context.Context) (cm model.ChatModel, err error) {
	// TODO Modify component configuration here.
	maxTokens := 4096
	config := &openai.ChatModelConfig{Model: "claude-4.0-sonnet", // 使用的模型版本
		APIKey:    "sk-791c03d4a69f21a06dee34acd0510615b8fe79ffbcc633b9b0c988e4aa9fe4b0",
		BaseURL:   "https://api.qnaigc.com/v1",
		MaxTokens: &maxTokens, // 添加 max_tokens 参数，使用指针
	}
	cm, err = openai.NewChatModel(ctx, config)
	if err != nil {
		fmt.Printf("[ChatModel] Model creation failed: %v", err)
		return nil, err
	}

	return cm, nil
}

var systemPrompt = `
# Role: Eino Expert Assistant

## Core Competencies
- knowledge of Eino framework and ecosystem
- Project scaffolding and best practices consultation
- Documentation navigation and implementation guidance
- Search web, clone github repo, open file/url, task management

## Interaction Guidelines
- Before responding, ensure you:
  • Fully understand the user's request and requirements, if there are any ambiguities, clarify with the user
  • Consider the most appropriate solution approach

- When providing assistance:
  • Be clear and concise
  • Include practical examples when relevant
  • Reference documentation when helpful
  • Suggest improvements or next steps if applicable

- If a request exceeds your capabilities:
  • Clearly communicate your limitations, suggest alternative approaches if possible

- If the question is compound or complex, you need to think step by step, avoiding giving low-quality answers directly.

## Context Information
- Current Date: {date}
- Related Documents: |-
==== doc start ====
  {documents}
==== doc end ====
`

type ChatTemplateConfig struct {
	FormatType schema.FormatType
	Templates  []schema.MessagesTemplate
}

// LoggedChatTemplate wraps a chat template to add logging
type LoggedChatTemplate struct {
	inner prompt.ChatTemplate
}

func (lct *LoggedChatTemplate) Format(ctx context.Context, input map[string]any, opts ...prompt.Option) ([]*schema.Message, error) {
	// Log input
	inputJSON, _ := json.MarshalIndent(input, "", "  ")
	log.Printf("[ChatTemplate] Input: %s", string(inputJSON))

	// Call inner chat template
	messages, err := lct.inner.Format(ctx, input, opts...)

	// Log output
	if err != nil {
		log.Printf("[ChatTemplate] Error: %v", err)
		return nil, err
	}

	outputJSON, _ := json.MarshalIndent(messages, "", "  ")
	log.Printf("[ChatTemplate] Output: %d messages formatted: %s", len(messages), string(outputJSON))

	return messages, nil
}

// newChatTemplate component initialization function of node 'ChatTemplate' in graph 'EinoAgent'
func newChatTemplate(ctx context.Context) (ctp prompt.ChatTemplate, err error) {
	// TODO Modify component configuration here.
	config := &ChatTemplateConfig{
		FormatType: schema.FString,
		Templates: []schema.MessagesTemplate{
			schema.SystemMessage(systemPrompt),
			schema.MessagesPlaceholder("history", true),
			schema.UserMessage("{content}"),
		},
	}
	baseChatTemplate := prompt.FromMessages(config.FormatType, config.Templates...)

	// Wrap with logging
	ctp = &LoggedChatTemplate{inner: baseChatTemplate}
	return ctp, nil
}
