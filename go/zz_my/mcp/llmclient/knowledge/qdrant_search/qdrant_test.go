package qdrant_search

import (
	"context"
	"testing"

	"github.com/qdrant/go-client/qdrant"
)

func TestQdrantCreateCollection(t *testing.T) {
	var ctx = context.Background()
	qdClient, err := qdrant.NewClient(&qdrant.Config{
		Host:                   "localhost",
		Port:                   6334,
		SkipCompatibilityCheck: true,
	})
	if err != nil {
		return
	}

	var collectionName = "qdrant"
	exists, err := qdClient.CollectionExists(ctx, collectionName)
	if err != nil {
		t.Fatal(err)
	}
	if !exists {
		if err := qdClient.CreateCollection(ctx, &qdrant.CreateCollection{
			CollectionName: collectionName,
			VectorsConfig: qdrant.NewVectorsConfig(&qdrant.VectorParams{
				// 与embedding的Dimension对其，不然插入数据的时候不会成功并且不会抛出错误
				Size:     1024,
				Distance: qdrant.Distance_Cosine,
			}),
		}); err != nil {
			t.Fatal(err)
		}
	}

	collections, err := qdClient.ListCollections(ctx)
	if err != nil {
		t.Fatal(err)
	}
	t.Log(collections)
}
