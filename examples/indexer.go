package examples

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/cloudwego/eino-ext/components/indexer/redis"

	"github.com/cloudwego/eino/components/indexer"
	"github.com/cloudwego/eino/schema"
	"github.com/google/uuid"
	rds "github.com/redis/go-redis/v9"
)

const (
	RedisPrefix = "eino:doc:"
	IndexName   = "vector_index"

	ContentField  = "content"
	MetadataField = "metadata"
	VectorField   = "content_vector"
	DistanceField = "distance"
)

// customDocumentToFields generates IDs for documents that don't have them
func customDocumentToFields(ctx context.Context, doc *schema.Document) (*redis.Hashes, error) {
	if doc.ID == "" {
		doc.ID = uuid.New().String()
	}
	key := doc.ID

	metadataBytes, err := json.Marshal(doc.MetaData)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal metadata: %w", err)
	}

	return &redis.Hashes{
		Key: key,
		Field2Value: map[string]redis.FieldValue{
			ContentField:  {Value: doc.Content, EmbedKey: VectorField},
			MetadataField: {Value: metadataBytes},
		},
	}, nil
}

// newIndexer component initialization function of node 'RedisIndexer' in graph 'myeino'
func newIndexer(ctx context.Context) (idr indexer.Indexer, err error) {
	// TODO Modify component configuration here.
	config := &redis.IndexerConfig{
		KeyPrefix: "eino:doc:",
		Client: rds.NewClient(&rds.Options{
			Addr: "localhost:6479",
		}),
		// Use custom document to fields mapping that handles missing IDs
		DocumentToHashes: customDocumentToFields,
	}
	embeddingIns11, err := newEmbedding(ctx)
	if err != nil {
		return nil, err
	}
	config.Embedding = embeddingIns11
	idr, err = redis.NewIndexer(ctx, config)
	if err != nil {
		return nil, err
	}
	return idr, nil
}
