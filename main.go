package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/anush008/fastembed-go"
	"github.com/joho/godotenv"
	"github.com/pocketbase/pocketbase"
	"github.com/pocketbase/pocketbase/apis"
	"github.com/pocketbase/pocketbase/core"
	"github.com/qdrant/go-client/qdrant"

	"github.com/tmc/langchaingo/llms"

	// "David/qdrant_api"
	"github.com/tmc/langchaingo/llms/openai"
	"google.golang.org/grpc"
)

type Status int

const (
	Pending Status = iota
	InProgress
	Completed
)

func (s Status) String() string {
	return [...]string{"Pending", "In Progress", "Completed"}[s]
}

type DataPoint struct {
	Content   string
	IsBug     bool
	RepCount  int
	Priority  float32
	Timestamp time.Time
	Status    Status
	Embedding []float32
}

func main() {
	app := pocketbase.New()
	app.OnServe().BindFunc(func(se *core.ServeEvent) error {
		// serves static files from the provided public dir (if exists)
		se.Router.GET("/{path...}", apis.Static(os.DirFS("./pb_public"), false))
		return se.Next()
	})

	if err := app.Start(); err != nil {
		log.Fatal(err)
	}
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Error loading .env file")
	}

	dbKey := os.Getenv("QDRANT_DB_KEY")

	// Create a new Qdrant client with API key authentication and TLS enabled
	client, err := qdrant.NewClient(&qdrant.Config{
		Host: "ec9d9f59-6f8c-4cdc-ae05-fa7bc0a465e7.us-west-2-0.aws.cloud.qdrant.io",
		//Port:   6334,
		APIKey:                 dbKey,
		UseTLS:                 true,
		SkipCompatibilityCheck: true,
		GrpcOptions: []grpc.DialOption{
			grpc.WithAuthority("ec9d9f59-6f8c-4cdc-ae05-fa7bc0a465e7.us-west-2-0.aws.cloud.qdrant.io:6334"), // Explicitly set authority with port
		},
	})
	if err != nil {
		panic(err)
	}

	collections, err := client.ListCollections(context.Background())
	if err != nil {
		panic(err)
	}

	fmt.Println(collections)
	defer client.Close()

	apiKey := os.Getenv("GROQ_API_KEY")

	llm, err := openai.New(
		openai.WithModel("llama3-8b-8192"),
		openai.WithBaseURL("https://api.groq.com/openai/v1"),
		openai.WithToken(apiKey),
	)
	if err != nil {
		log.Fatal(err)
	}

	query := "Hi Team, I’ve been facing a couple of issues with the app recently. First, the app crashes every time I try to log in, especially when my internet connection is slow. Also, push notifications seem to have stopped working entirely after the last update, even though they’re enabled in the settings. On another note, it would be great if you could add a dark mode option to make the app easier to use at night. Please let me know if you need more details to look into these issues.Thanks,[Customer Name]"

	tasks := getTasks(query, llm)

	options := fastembed.InitOptions{
		Model:     fastembed.BGESmallEN, // Correct identifier
		CacheDir:  "models",
		MaxLength: 512,
	}

	embedder, err := fastembed.NewFlagEmbedding(&options)
	if err != nil {
		panic(fmt.Sprintf("Model initialization failed: %v", err))
	}
	defer embedder.Destroy()

	fmt.Println("Model initialised successfully")

	for _, task := range tasks {
		isBug := false
		taskQuery := strings.Split(task, "=")

		if taskQuery[1] == "true" {
			isBug = true
		}

		dataPoint := DataPoint{
			Content:   taskQuery[0],
			IsBug:     isBug,
			RepCount:  1,
			Priority:  0.5,
			Timestamp: time.Now(),
			Status:    Pending,
		}
		//payload := createPayload(dataPoint)

		embedding, err := embedder.Embed([]string{taskQuery[0]}, 25)
		if err != nil {
			panic(fmt.Sprintf("Embedding failed: %v", err))
		}
		dataPoint.Embedding = embedding[0]

		limit := uint64(1)

		results, err := client.Query(context.Background(), &qdrant.QueryPoints{
			CollectionName: "task_data",
			Query:          qdrant.NewQuery(embedding[0]...),
			Limit:          &limit,
			WithPayload:    qdrant.NewWithPayload(true),
			WithVectors:    qdrant.NewWithVectors(true),
		})
		if err != nil {
			panic(err)
		}

		countResponse, err := client.Count(context.Background(), &qdrant.CountPoints{
			CollectionName: "task_data",
		})
		if err != nil {
			panic(err)
		}

		if countResponse == 0 {
			dataPoint.Content = expandTask(taskQuery[0], llm)
			updateAndCreateDataPoint(client, dataPoint, int(countResponse)+1)
		} else {
			var result map[string]*qdrant.Value

			for _, value := range results[0].GetPayload() {
				result = value.GetStructValue().GetFields()
			}

			fmt.Println("Result: ", result["Content"].GetStringValue())

			if compareStrings(taskQuery[0], result["Content"].GetStringValue(), llm) {
				fmt.Println("Similar task found")
				dataPoint.Content = result["Content"].GetStringValue()
				dataPoint.IsBug, _ = strconv.ParseBool(result["IsBug"].GetStringValue())
				dataPoint.RepCount, _ = strconv.Atoi(result["RepCount"].GetStringValue())
				dataPoint.RepCount++
				dataPoint.Timestamp, _ = time.Parse(time.RFC3339, result["Timestamp"].GetStringValue())
				dataPoint.Embedding = results[0].GetVectors().GetVector().Data
				updateAndCreateDataPoint(client, dataPoint, int(results[0].GetId().GetNum()))
			} else {
				dataPoint.Content = expandTask(taskQuery[0], llm)
				updateAndCreateDataPoint(client, dataPoint, int(countResponse)+1)
			}

		}

		fmt.Printf("Content: %s\n", dataPoint.Content)
		fmt.Printf("IsBug: %t\n", dataPoint.IsBug)
		fmt.Printf("RepCount: %d\n", dataPoint.RepCount)
		fmt.Printf("Priority: %.2f\n", dataPoint.Priority)
		fmt.Printf("Timestamp: %s\n", dataPoint.Timestamp.Format(time.RFC3339))
		fmt.Printf("Status: %s\n", dataPoint.Status.String())
		// fmt.Print("Embedding: \n", dataPoint.Embedding)
		fmt.Println()

	}

}
func updateAndCreateDataPoint(client *qdrant.Client, dataPoint DataPoint, id int) {
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
			strconv.Itoa(id): createPayload(dataPoint),
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

func getTasks(query string, llm *openai.LLM) []string {
	// Define the prompt for the LLM
	prompt := `You are a task extraction assistant. Given the email report from a customer, extract the tasks and their details to be sent to the engineering team. Do not add any greeting or ending sentence, stick to the format given, do not add any index like "Task1=", etc. The tasks should be in the following format,where the <bool> implies a bool value which is true if the task is reported bug, or if it is a requested feature:
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

func compareStrings(a string, b string, llm *openai.LLM) bool {
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

func createPayload(dp DataPoint) map[string]any {
	return map[string]any{
		"Content":   dp.Content,
		"IsBug":     strconv.FormatBool(dp.IsBug),
		"RepCount":  strconv.Itoa(dp.RepCount),
		"Priority":  fmt.Sprintf("%.2f", dp.Priority),
		"Timestamp": dp.Timestamp.Format(time.RFC3339),
		"Status":    dp.Status.String(),
	}
}

func expandTask(task string, llm *openai.LLM) string {
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
