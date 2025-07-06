package qdrant_api

import (
	"context"
	"fmt"
	// "log"
	"strconv"
	// "strings"
	"time"

	"github.com/qdrant/go-client/qdrant"
	// "github.com/tmc/langchaingo/llms"
	// "github.com/tmc/langchaingo/llms/openai"
	dp "David/data_point"
	// llm "David/llm_functions"
)

var CLIENT *qdrant.Client

func UpdateAndCreateDataPoint(dataPoint dp.DataPoint, id int, collectionId string) {
	point := &qdrant.PointStruct{
		Id: qdrant.NewIDNum(uint64(id)),
		Vectors: &qdrant.Vectors{
			VectorsOptions: &qdrant.Vectors_Vector{
				Vector: &qdrant.Vector{
					Data: dataPoint.Embedding,
				},
			},
		},
		Payload: qdrant.NewValueMap(map[string]any{
			strconv.Itoa(id): CreateQdrantPayload(dataPoint),
		}),
	}

	ctx := context.Background()
	_, err := CLIENT.Upsert(ctx, &qdrant.UpsertPoints{
		CollectionName: collectionId,
		Points:         []*qdrant.PointStruct{point},
	})
	if err != nil {
		panic(err)
	}
}

func CreateQdrantPayload(data dp.DataPoint) map[string]any {
	return map[string]any{
		// ""
		"Content":   data.Content,
		"IsBug":     strconv.FormatBool(data.IsBug),
		"RepCount":  strconv.Itoa(data.RepCount),
		"Priority":  fmt.Sprintf("%.2f", data.Priority),
		"Timestamp": data.Timestamp.Format(time.RFC3339),
		"Status":    data.Status.String(),
	}
}

//qdrant collection name will be the same as project name
func CreateCollection(collectionName string) error {
	 	err := CLIENT.CreateCollection(context.Background(), &qdrant.CreateCollection{
		CollectionName: collectionName,
		VectorsConfig: qdrant.NewVectorsConfig(&qdrant.VectorParams{
			Size:     384, // Adjust to your embedding size
			Distance: qdrant.Distance_Cosine,
		}),
	})
	return err
}