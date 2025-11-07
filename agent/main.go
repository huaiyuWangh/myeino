package main

import (
	"context"
	"data-analyze/util"
	"github.com/cloudwego/eino-ext/components/model/openai"
	"github.com/cloudwego/eino/compose"
	"log"
	"os"

	"github.com/cloudwego/eino-ext/components/tool/duckduckgo"
	"github.com/cloudwego/eino/components/tool"
	"github.com/cloudwego/eino/components/tool/utils"
	"github.com/cloudwego/eino/schema"
)

func main() {
	// 初始化 tools
	todoTools := []tool.BaseTool{
		getAddTodoTool(),    // NewTool 构建
		getUpdateTodoTool(), // InferTool 构建
		getListTodoTool(),   // 实现Tool接口
		//getSearchTool(),     // 官方封装的工具
	}

	// 创建并配置 ChatModel
	chatModel, err := openai.NewChatModel(context.Background(), &openai.ChatModelConfig{
		Model:   os.Getenv("ANTHROPIC_DEFAULT_SONNET_MODEL"), // 使用的模型版本
		APIKey:  os.Getenv("ANTHROPIC_API_KEY"),              // OpenAI API 密钥
		BaseURL: "https://api.qnaigc.com/v1",
	})
	if err != nil {
		util.Fatal(err)
	}
	ctx := context.Background()
	// 获取工具信息并绑定到 ChatModel
	toolInfos := make([]*schema.ToolInfo, 0, len(todoTools))
	for _, tool := range todoTools {
		info, err := tool.Info(ctx)
		if err != nil {
			util.Fatal(err)
		}
		toolInfos = append(toolInfos, info)
	}
	err = chatModel.BindTools(toolInfos)
	if err != nil {
		util.Fatal(err)
	}

	// 创建 tools 节点
	todoToolsNode, err := compose.NewToolNode(context.Background(), &compose.ToolsNodeConfig{
		Tools: todoTools,
	})
	if err != nil {
		util.Fatal(err)
	}

	// 构建完整的处理链
	chain := compose.NewChain[[]*schema.Message, []*schema.Message]()
	chain.
		AppendChatModel(chatModel, compose.WithNodeName("chat_model")).
		AppendToolsNode(todoToolsNode, compose.WithNodeName("tools"))

	// 编译并运行 chain
	agent, err := chain.Compile(ctx)
	if err != nil {
		util.Fatal(err)
	}

	// 运行示例
	resp, err := agent.Invoke(ctx, []*schema.Message{
		{
			Role:    schema.User,
			Content: "添加一个学习 Eino 的 TODO，展示TODO list",
		},
	})
	if err != nil {
		util.Fatal(err)
	}

	// 输出结果
	util.Info(resp)
	for _, msg := range resp {
		util.Info(msg.Content)
	}
}

func getAddTodoTool() tool.InvokableTool {
	// 工具信息
	info := &schema.ToolInfo{
		Name: "add_todo",
		Desc: "Add a todo item",
		ParamsOneOf: schema.NewParamsOneOfByParams(map[string]*schema.ParameterInfo{
			"content": {
				Desc:     "The content of the todo item",
				Type:     schema.String,
				Required: true,
			},
			"started_at": {
				Desc: "The started time of the todo item, in unix timestamp",
				Type: schema.Integer,
			},
			"deadline": {
				Desc: "The deadline of the todo item, in unix timestamp",
				Type: schema.Integer,
			},
		}),
	}

	// 使用NewTool创建工具
	return utils.NewTool(info, AddTodoFunc)
}

type TodoAddParams struct{}

// 处理函数
func AddTodoFunc(_ context.Context, params *TodoAddParams) (string, error) {
	// Mock处理逻辑
	return `{"msg": "add todo success"}`, nil
}

// 参数结构体
type TodoUpdateParams struct {
	ID        string  `json:"id" jsonschema:"description=id of the todo"`
	Content   *string `json:"content,omitempty" jsonschema:"description=content of the todo"`
	StartedAt *int64  `json:"started_at,omitempty" jsonschema:"description=start time in unix timestamp"`
	Deadline  *int64  `json:"deadline,omitempty" jsonschema:"description=deadline of the todo in unix timestamp"`
	Done      *bool   `json:"done,omitempty" jsonschema:"description=done status"`
}

// 处理函数
func UpdateTodoFunc(_ context.Context, params *TodoUpdateParams) (string, error) {
	// Mock处理逻辑
	return `{"msg": "update todo success"}`, nil
}

func getUpdateTodoTool() tool.InvokableTool {
	// 使用 InferTool 创建工具
	updateTool, _ := utils.InferTool(
		"update_todo", // tool name
		"Update a todo item, eg: content,deadline...", // tool description
		UpdateTodoFunc)
	return updateTool
}

type ListTodoTool struct{}

func getListTodoTool() tool.InvokableTool {
	return &ListTodoTool{}
}

func (lt *ListTodoTool) Info(ctx context.Context) (*schema.ToolInfo, error) {
	return &schema.ToolInfo{
		Name: "list_todo",
		Desc: "List all todo items",
		ParamsOneOf: schema.NewParamsOneOfByParams(map[string]*schema.ParameterInfo{
			"finished": {
				Desc:     "filter todo items if finished",
				Type:     schema.Boolean,
				Required: false,
			},
		}),
	}, nil
}

func (lt *ListTodoTool) InvokableRun(ctx context.Context, argumentsInJSON string, opts ...tool.Option) (string, error) {
	// Mock调用逻辑
	return `{"todos": [{"id": "1", "content": "在2024年12月10日之前完成Eino项目演示文稿的准备工作", "started_at": 1717401600, "deadline": 1717488000, "done": false}]}`, nil
}

func getSearchTool() tool.InvokableTool {
	// 创建 duckduckgo Search 工具
	searchTool, err := duckduckgo.NewTool(context.Background(), &duckduckgo.Config{})
	if err != nil {
		log.Fatal(err)
	}
	return searchTool
}
