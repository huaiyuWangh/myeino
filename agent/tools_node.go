package agent

import (
	"context"
	"github.com/cloudwego/eino-examples/quickstart/eino_assistant/eino/einoagent"

	"github.com/cloudwego/eino-ext/components/tool/duckduckgo/v2"
	"github.com/cloudwego/eino/components/tool"
)

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

	toolDDGSearch, err := NewDDGSearch(ctx)
	if err != nil {
		return nil, err
	}
	return []tool.BaseTool{
		einoAssistantTool,
		toolTask,
		toolOpen,
		toolGitClone,
		toolDDGSearch,
	}, nil
}

func NewDDGSearch(ctx context.Context) (bt tool.BaseTool, err error) {
	// TODO Modify component configuration here.
	config := &duckduckgo.Config{}
	bt, err = duckduckgo.NewTextSearchTool(ctx, config)
	if err != nil {
		return nil, err
	}
	return bt, nil
}
