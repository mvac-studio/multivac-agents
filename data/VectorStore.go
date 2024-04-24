package data

import (
	"context"
	chromago "github.com/amikos-tech/chroma-go"
	"github.com/amikos-tech/chroma-go/collection"
	"github.com/amikos-tech/chroma-go/ollama"
	"github.com/amikos-tech/chroma-go/types"
	"log"
	"strings"
)

type VectorStore struct {
	client     *chromago.Client
	collection *chromago.Collection
	embedFn    types.EmbeddingFunction
}

func NewVectorStore(index string) *VectorStore {
	client, err := chromago.NewClient("http://chromadb-service.default.svc.cluster.local:8000")
	if err != nil {
		log.Println(err)
	}
	embeddingFunction, err := ollama.NewOllamaEmbeddingFunction(ollama.WithBaseURL("http://ollama-service.default.svc.cluster.local:11434"), ollama.WithModel("all-minilm:l6-v2"))
	col, err := client.NewCollection(context.Background(),
		collection.WithName(index),
		collection.WithEmbeddingFunction(embeddingFunction),
		collection.WithCreateIfNotExist(true),
		collection.WithHNSWDistanceFunction(types.L2),
	)
	if err != nil {
		panic(err)
	}

	return &VectorStore{
		client:     client,
		embedFn:    embeddingFunction,
		collection: col,
	}
}

func (store *VectorStore) Commit(chunk string) error {

	rs, err := types.NewRecordSet(types.WithEmbeddingFunction(store.embedFn), types.WithIDGenerator(types.NewULIDGenerator()))
	rs.WithRecord(types.WithDocument(chunk))
	if err != nil {
		log.Println(err)
	}

	_, err = rs.BuildAndValidate(context.TODO())
	_, err = store.collection.AddRecords(context.Background(), rs)
	if err != nil {
		return err
	}
	return nil
}

func (store *VectorStore) Clear() {
	_, err := store.client.DeleteCollection(context.Background(), store.collection.Name)
	if err != nil {
		log.Println(err)
	}
}

func (store *VectorStore) Query(query string, top int32, maxDistance float32) string {
	result, err := store.collection.QueryWithOptions(context.TODO(),
		types.WithNResults(top),
		types.WithQueryText(query),
	)
	if err != nil {
		log.Println(err)
	}
	builder := strings.Builder{}
	for _, document := range result.Documents {
		for _, v := range document {
			builder.WriteString(v)
		}
	}
	return builder.String()
}
