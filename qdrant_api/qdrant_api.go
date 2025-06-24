package qdrant_api

import (
	"context"
	"fmt"
	"log"
	"strconv"
	"strings"
	"time"

	"github.com/qdrant/go-client/qdrant"
	"github.com/tmc/langchaingo/llms"
	"github.com/tmc/langchaingo/llms/openai"
	dp "David/data_point"
	llm "David/llm_functions"
)

func UpdateAndCreateDataPoint(client *qdrant.Client, dataPoint dp.DataPoint, id int, collectionId string) {
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
			strconv.Itoa(id): CreatePayload(dataPoint),
		}),
	}

	ctx := context.Background()
	_, err := client.Upsert(ctx, &qdrant.UpsertPoints{
		CollectionName: collectionId,
		Points:         []*qdrant.PointStruct{point},
	})
	if err != nil {
		panic(err)
	}
}

func CreateQdrantPayload(data dp.DataPoint) map[string]any {
	return map[string]any{
		"Content":   data.Content,
		"IsBug":     strconv.FormatBool(data.IsBug),
		"RepCount":  strconv.Itoa(data.RepCount),
		"Priority":  fmt.Sprintf("%.2f", data.Priority),
		"Timestamp": data.Timestamp.Format(time.RFC3339),
		"Status":    data.Status.String(),
	}
}

func CreateCollection(client *qdrant.Client, collectionName string) error {
	_, err := client.CreateCollection(context.Background(), &qdrant.CreateCollection{
		CollectionName: collectionName,
		VectorsConfig: qdrant.NewVectorsConfig(&qdrant.VectorParams{
			Size:     384, // Adjust to your embedding size
			Distance: qdrant.Distance_Cosine,
		}),
	})
	return err
}