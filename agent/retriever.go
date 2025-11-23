package agent

import (
	"context"
	"encoding/json"
	"fmt"
	redispkg "github.com/cloudwego/eino-examples/quickstart/eino_assistant/pkg/redis"
	"github.com/cloudwego/eino/schema"
	rds "github.com/redis/go-redis/v9"
	"log"
	"strconv"

	"github.com/cloudwego/eino-ext/components/retriever/redis"
	"github.com/cloudwego/eino/components/retriever"
)

// LoggedRetriever wraps a retriever to add logging
type LoggedRetriever struct {
	inner retriever.Retriever
}

func (lr *LoggedRetriever) Retrieve(ctx context.Context, query string, opts ...retriever.Option) ([]*schema.Document, error) {
	// Log input
	log.Printf("[Retriever] Input: query=%s", query)

	// Call inner retriever
	docs, err := lr.inner.Retrieve(ctx, query, opts...)

	// Log output
	if err != nil {
		log.Printf("[Retriever] Error: %v", err)
		return nil, err
	}

	// Format documents for better logging
	formattedDocs := make([]map[string]interface{}, len(docs))
	for i, doc := range docs {
		contentPreview := doc.Content
		if len(doc.Content) > 100 {
			contentPreview = doc.Content[:100] + "..."
		}
		formattedDocs[i] = map[string]interface{}{
			"id":      doc.ID,
			"content": contentPreview,
			"score":   doc.Score(),
		}
	}

	outputJSON, _ := json.MarshalIndent(formattedDocs, "", "  ")
	log.Printf("[Retriever] Output: %d documents retrieved: %s", len(docs), string(outputJSON))

	return docs, nil
}

// newRetriever component initialization function of node 'Retriever' in graph 'EinoAgent'
func newRetriever(ctx context.Context) (rtr retriever.Retriever, err error) {
	// TODO Modify component configuration here.
	config := &redis.RetrieverConfig{
		Client: rds.NewClient(&rds.Options{
			Addr:     "localhost:6479",
			Protocol: 2,
		}),
		Index:        fmt.Sprintf("%s%s", redispkg.RedisPrefix, redispkg.IndexName),
		Dialect:      2,
		ReturnFields: []string{redispkg.ContentField, redispkg.MetadataField, redispkg.DistanceField},
		TopK:         8,
		VectorField:  redispkg.VectorField,
		DocumentConverter: func(ctx context.Context, doc rds.Document) (*schema.Document, error) {
			resp := &schema.Document{
				ID:       doc.ID,
				Content:  "",
				MetaData: map[string]any{},
			}
			for field, val := range doc.Fields {
				if field == redispkg.ContentField {
					resp.Content = val
				} else if field == redispkg.MetadataField {
					resp.MetaData[field] = val
				} else if field == redispkg.DistanceField {
					distance, err := strconv.ParseFloat(val, 64)
					if err != nil {
						continue
					}
					resp.WithScore(1 - distance)
				}
			}

			return resp, nil
		},
	}
	embeddingIns11, err := newEmbedding(ctx)
	if err != nil {
		return nil, err
	}
	config.Embedding = embeddingIns11
	baseRetriever, err := redis.NewRetriever(ctx, config)
	if err != nil {
		return nil, err
	}

	// Wrap with logging
	rtr = &LoggedRetriever{inner: baseRetriever}
	return rtr, nil
}
