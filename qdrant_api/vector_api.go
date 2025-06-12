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
)

type Status int

const (
	Pending Status = iota
	InProgress
	Completed
)

type DataPoint struct {
	Content   string
	IsBug     bool
	RepCount  int
	Priority  float32
	Timestamp time.Time
	Status    Status
	Embedding []float32
}


func (s Status) String() string {
	return [...]string{"Pending", "In Progress", "Completed"}[s]
}

func UpdateAndCreateDataPoint(client *qdrant.Client, dataPoint DataPoint, id int) {
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
		CollectionName: "task_data",
		Points:         []*qdrant.PointStruct{point},
	})
	if err != nil {
		panic(err)
	}
}

func GetTasks(query string, llm *openai.LLM) []string {
	// Define the prompt for the LLM
	prompt := `You are a task extraction assistant. Given the email report from a customer, extract the tasks and their details to be sent to the engineering team. Do not add any greeting or ending sentence, stick to the format given, do not add any index like "Task1=", etc. The tasks should be in the following format:
	<task>=<bool>;<task>=<bool>;...`

	ctx := context.Background()

	fullPrompt := prompt + "\n\nCustomer Message:\n" + query

	// Use the correct Call method
	completion, err := llm.Call(ctx, fullPrompt,
		llms.WithTemperature(0.1),
	)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("LLM Response:\n", completion)

	return strings.Split(completion, ";")
}

func CompareStrings(a string, b string, llm *openai.LLM) bool {
	prompt := `You are a string comparison assistant. Given two string, determine if they are similar or not, that is are both the strings pointing out the same issue or not. Respond with "true" if they are similar and "false" if they are not.`

	ctx := context.Background()

	fullPrompt := prompt + "\n\nString 1:\n" + a + "\n\nString 2:\n" + b

	// Use the correct Call method
	completion, err := llm.Call(ctx, fullPrompt,
		llms.WithTemperature(0.1),
	)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("LLM Response:\n", completion)

	return strings.ToLower(completion) == "true"
}

func CreatePayload(dp DataPoint) map[string]any {
	return map[string]any{
		"Content":   dp.Content,
		"IsBug":     strconv.FormatBool(dp.IsBug),
		"RepCount":  strconv.Itoa(dp.RepCount),
		"Priority":  fmt.Sprintf("%.2f", dp.Priority),
		"Timestamp": dp.Timestamp.Format(time.RFC3339),
		"Status":    dp.Status,
	}
}

func ExpandTask(task string, llm *openai.LLM) string {
	// Define the prompt for the LLM
	// TODO: Want to make it so that the llm can get the actual values from all task history.
	prompt := `You are a task expansion assistant. Given a task, expand it to include more details and context. Make it 30 to 40 words long. Just simply respond with the expanded task. Avoid any greeting or ending sentence. Do not use placeholder variables for unknown values.`

	ctx := context.Background()

	fullPrompt := prompt + "\n\nTask:\n" + task

	// Use the correct Call method
	completion, err := llm.Call(ctx, fullPrompt,
		llms.WithTemperature(0.1),
	)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("LLM Response:\n", completion)

	return completion

}
