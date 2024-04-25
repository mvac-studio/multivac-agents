package data

import (
	"context"
	"fmt"
	chromago "github.com/amikos-tech/chroma-go"
	"github.com/amikos-tech/chroma-go/collection"
	"github.com/amikos-tech/chroma-go/ollama"
	"github.com/amikos-tech/chroma-go/types"
	"log"
	"regexp"
	"strings"
)

type VectorStore struct {
	client     *chromago.Client
	secrets    *AgentDataStore
	userid     string
	collection *chromago.Collection
	embedFn    types.EmbeddingFunction
}

func NewVectorStore(userid string, index string) *VectorStore {
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
		userid:     userid,
		embedFn:    embeddingFunction,
		secrets:    NewAgentDataStore(),
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
	fmt.Printf("Memory wiped for %s\n", store.collection.Name)
	index := store.collection.Name
	_, err := store.client.DeleteCollection(context.Background(), index)
	if err != nil {
		log.Println(err)
	}
	col, err := store.client.NewCollection(context.Background(),
		collection.WithName(index),
		collection.WithEmbeddingFunction(store.embedFn),
		collection.WithCreateIfNotExist(true),
		collection.WithHNSWDistanceFunction(types.L2),
	)
	store.collection = col
}

func (store *VectorStore) Query(agentid string, query string, top int32, maxDistance float32) (memory string, containsSecret bool) {
	result, err := store.collection.QueryWithOptions(context.TODO(),
		types.WithNResults(top),
		types.WithQueryText(query),
	)
	if err != nil {
		log.Println(err)
	}
	containsSecret = false
	builder := strings.Builder{}
	secretMatch := regexp.MustCompile(`\[~SECRET] ref:(.*?)\[SECRET~]`)
	for _, document := range result.Documents {
		for _, v := range document {
			if secretMatch.MatchString(v) {
				secretRef := secretMatch.FindStringSubmatch(v)[1]
				result, _ := store.secrets.RetrieveSecret(secretRef, store.userid, agentid)
				v = secretMatch.ReplaceAllString(v, result)
				containsSecret = true
			}
			builder.WriteString(v)
		}
	}
	return builder.String(), containsSecret
}
