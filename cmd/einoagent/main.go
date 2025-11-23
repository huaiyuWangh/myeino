package main

import (
	"context"
	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/app/server"
	"log"
	"myeino/cmd/einoagent/agent"
)

var port = "8090"

func main() {
	// 创建 Hertz 服务器
	h := server.Default(server.WithHostPorts(":" + port))

	// 注册 agent 路由组
	agentGroup := h.Group("/agent")
	if err := agent.BindRoutes(agentGroup); err != nil {
		log.Fatal("failed to bind agent routes:", err)
	}

	// Redirect root path to /agent
	h.GET("/", func(ctx context.Context, c *app.RequestContext) {
		c.Redirect(302, []byte("/agent"))
	})

	// 启动服务器
	h.Spin()
}
