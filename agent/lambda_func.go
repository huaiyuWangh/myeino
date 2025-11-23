package agent

import (
	"context"
	"encoding/json"
	"log"
	"time"
)

// newLambda component initialization function of node 'InputToQuery' in graph 'EinoAgent'
func newInputToQuery(ctx context.Context, input *UserMessage, opts ...any) (output string, err error) {
	// Log input
	inputJSON, _ := json.MarshalIndent(input, "", "  ")
	log.Printf("[InputToQuery] Input: %s", string(inputJSON))

	output = input.Query

	// Log output
	log.Printf("[InputToQuery] Output: %s", output)

	return output, nil
}

// newLambda1 component initialization function of node 'InputToHistory' in graph 'EinoAgent'
func newInputToHistory(ctx context.Context, input *UserMessage, opts ...any) (output map[string]any, err error) {
	// Log input
	inputJSON, _ := json.MarshalIndent(input, "", "  ")
	log.Printf("[InputToHistory] Input: %s", string(inputJSON))

	output = map[string]any{
		"content": input.Query,
		"history": input.History,
		"date":    time.Now().Format(time.DateTime),
	}

	// Log output
	outputJSON, _ := json.MarshalIndent(output, "", "  ")
	log.Printf("[InputToHistory] Output: %s", string(outputJSON))

	return output, nil
}
