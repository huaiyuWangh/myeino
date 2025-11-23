package agent

import (
	"context"
	"encoding/json"
	"log"

	"github.com/cloudwego/eino/compose"
	"github.com/cloudwego/eino/flow/agent"
	"github.com/cloudwego/eino/flow/agent/react"
	"github.com/cloudwego/eino/schema"
)

// loggedGenerate wraps the react agent Generate function with logging
func loggedGenerate(originalGenerate func(context.Context, []*schema.Message, ...agent.AgentOption) (*schema.Message, error)) func(context.Context, []*schema.Message, ...agent.AgentOption) (*schema.Message, error) {
	return func(ctx context.Context, input []*schema.Message, opts ...agent.AgentOption) (*schema.Message, error) {
		// Log input
		inputJSON, _ := json.MarshalIndent(input, "", "  ")
		log.Printf("[ReactAgent] Input: %d messages: %s", len(input), string(inputJSON))

		// Call original function
		output, err := originalGenerate(ctx, input, opts...)

		// Log output
		if err != nil {
			log.Printf("[ReactAgent] Error: %v", err)
			return nil, err
		}

		outputJSON, _ := json.MarshalIndent(output, "", "  ")
		log.Printf("[ReactAgent] Output: %s", string(outputJSON))

		return output, nil
	}
}

// loggedStream wraps the react agent Stream function with logging
func loggedStream(originalStream func(context.Context, []*schema.Message, ...agent.AgentOption) (*schema.StreamReader[*schema.Message], error)) func(context.Context, []*schema.Message, ...agent.AgentOption) (*schema.StreamReader[*schema.Message], error) {
	return func(ctx context.Context, input []*schema.Message, opts ...agent.AgentOption) (*schema.StreamReader[*schema.Message], error) {
		// Log input messages
		inputJSON, _ := json.MarshalIndent(input, "", "  ")
		log.Printf("[ReactAgent-Stream] === DETAILED MODEL CALL PARAMETERS ===")
		log.Printf("[ReactAgent-Stream] Input Messages Count: %d", len(input))
		log.Printf("[ReactAgent-Stream] Input Messages: %s", string(inputJSON))

		// Log agent options
		log.Printf("[ReactAgent-Stream] Agent Options Count: %d", len(opts))
		for i, opt := range opts {
			log.Printf("[ReactAgent-Stream] Agent Option[%d]: %T", i, opt)
		}

		// Log context details
		if deadline, ok := ctx.Deadline(); ok {
			log.Printf("[ReactAgent-Stream] Context Deadline: %v", deadline)
		} else {
			log.Printf("[ReactAgent-Stream] Context Deadline: No deadline set")
		}

		// Log context values (be careful not to log sensitive data)
		log.Printf("[ReactAgent-Stream] Context Type: %T", ctx)

		// Call original function
		log.Printf("[ReactAgent-Stream] === CALLING ORIGINAL STREAM FUNCTION ===")
		reader, err := originalStream(ctx, input, opts...)

		// Log output
		if err != nil {
			log.Printf("[ReactAgent-Stream] === STREAM CALL FAILED ===")
			log.Printf("[ReactAgent-Stream] Error Type: %T", err)
			log.Printf("[ReactAgent-Stream] Error: %v", err)
			return nil, err
		}

		log.Printf("[ReactAgent-Stream] === STREAM CALL SUCCESSFUL ===")
		log.Printf("[ReactAgent-Stream] Stream Reader Type: %T", reader)
		log.Printf("[ReactAgent-Stream] Output: Stream reader created successfully")

		return reader, nil
	}
}

// newLambda2 component initialization function of node 'ReactAgent' in graph 'EinoAgent'
func newLambda2(ctx context.Context) (lba *compose.Lambda, err error) {
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

	// Wrap the generate and stream functions with logging
	loggedGen := loggedGenerate(ins.Generate)
	loggedStr := loggedStream(ins.Stream)

	lba, err = compose.AnyLambda(loggedGen, loggedStr, nil, nil)
	if err != nil {
		return nil, err
	}
	return lba, nil
}
