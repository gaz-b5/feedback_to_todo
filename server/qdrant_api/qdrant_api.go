package qdrant_api

import (
	"context"

	// "crypto/rand"
	// "encoding/binary"
	"fmt"

	// "log"
	// "strconv"
	// "strings"
	// "time"

	"github.com/google/uuid"
	"github.com/qdrant/go-client/qdrant"

	// "github.com/tmc/langchaingo/llms"
	// "github.com/tmc/langchaingo/llms/openai"
	dp "David/data_point"
	// llm "David/llm_functions"
)

var CLIENT *qdrant.Client

func UpdateAndCreateDataPoint(dataPoint dp.DataPoint, collectionId string) string {

	ctx := context.Background()

	qdrantDBId := uuid.New().String()

	point := &qdrant.PointStruct{
		Id: qdrant.NewIDUUID(qdrantDBId),
		Vectors: &qdrant.Vectors{
			VectorsOptions: &qdrant.Vectors_Vector{
				Vector: &qdrant.Vector{
					Data: dataPoint.Embedding,
				},
			},
		},
		Payload: qdrant.NewValueMap(map[string]any{
			"TaskID": dataPoint.TaskID,
		}),
	}

	_, err := CLIENT.Upsert(ctx, &qdrant.UpsertPoints{
		CollectionName: collectionId,
		Points:         []*qdrant.PointStruct{point},
	})
	if err != nil {
		panic(err)
	}

	return qdrantDBId
}

func CreateQdrantPayload(data dp.DataPoint) map[string]any {
	return map[string]any{
		"TaskID": data.TaskID,
		// "Content":   data.Content,
		// "IsBug":     strconv.FormatBool(data.IsBug),
		// "RepCount":  strconv.Itoa(data.RepCount),
		// "Priority":  fmt.Sprintf("%.2f", data.Priority),
		// "Timestamp": data.Timestamp.Format(time.RFC3339),
		// "Status":    data.Status.String(),
	}
}

// qdrant collection name will be the same as project name
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

// DeleteCollection deletes a Qdrant collection by its name.
func DeleteCollection(collectionName string) error {
	ctx := context.Background()
	err := CLIENT.DeleteCollection(ctx, collectionName)
	return err
}

func ReturnClosestTaskID(collectionName string, limit uint64, embedding []float32) (string, int) {

	results, err := CLIENT.Query(context.Background(), &qdrant.QueryPoints{
		CollectionName: collectionName,
		Query:          qdrant.NewQuery(embedding...),
		Limit:          &limit,
		WithPayload:    qdrant.NewWithPayload(true),
		WithVectors:    qdrant.NewWithVectors(true),
	})
	if err != nil {
		panic(err)
	}

	// var result map[string]*qdrant.Value
	// if len(results) > 0 {
	// 	for _, value := range results[0].GetPayload() {
	// 		result = value.GetStructValue().GetFields()
	// 	}
	// }

	taskId := ""
	if len(results) > 0 {
		payload := results[0].GetPayload()
		if val, ok := payload["TaskID"]; ok {
			taskId = val.GetStringValue()
		} else {
			panic("Closest task found but could not get payload")
		}
	}
	fmt.Println(taskId)

	numResults := len(results)
	return taskId, numResults
}

func DeleteTaskQdrant(collectionName string, qdrantDBId string) error {
	_, err := CLIENT.Delete(context.Background(), &qdrant.DeletePoints{
		CollectionName: collectionName,
		Points: qdrant.NewPointsSelector(
			qdrant.NewIDUUID(qdrantDBId), // Replace with your point ID (int64)
		),
	})
	if err != nil {
		return err
	}
	return nil

}

func DeleteTasksQdrantBulk(collectionName string, qdrantDBIds []string) error {
	// Create a slice of PointIDs
	var pointIDs []*qdrant.PointId
	for _, id := range qdrantDBIds {
		pointIDs = append(pointIDs, qdrant.NewIDUUID(id))
	}

	_, err := CLIENT.Delete(context.Background(), &qdrant.DeletePoints{
		CollectionName: collectionName,
		Points:         qdrant.NewPointsSelector(pointIDs...),
	})
	if err != nil {
		return err
	}
	return nil
}
