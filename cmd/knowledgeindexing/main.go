package main

import (
	"context"
	"fmt"
	"github.com/cloudwego/eino/components/document"
	"log"
	"myeino/examples"
)

// 1775577c-277f-4d16-bffa-b275ab38880c

func main() {
	ctx := context.Background()
	runner, err := examples.Buildmyeino(ctx)
	if err != nil {
		log.Fatal(err)
	}
	path := "/Users/wanghuaiyu/projects/myeino/cmd/knowledgeindexing/eino.md"
	ids, err := runner.Invoke(ctx, document.Source{URI: path})
	fmt.Println("ids:", ids)
	fmt.Println("err:", err)
}
