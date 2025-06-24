package main

import (
	// "context"
	// "fmt"
	"David/routes"
	"log"
	"os"

	// "strconv"
	"strings"
	// "time"

	// "github.com/anush008/fastembed-go"
	// "github.com/joho/godotenv"
	"github.com/pocketbase/pocketbase"
	// "github.com/pocketbase/pocketbase/apis"
	"github.com/pocketbase/pocketbase/core"
	"github.com/pocketbase/pocketbase/plugins/migratecmd"
	// "github.com/qdrant/go-client/qdrant"
	// "David/qdrant_api"
	// "github.com/tmc/langchaingo/llms/openai"
)


func main() {
	app := pocketbase.New()

	 // loosely check if it was executed using "go run"
    isGoRun := strings.HasPrefix(os.Args[0], os.TempDir())

    migratecmd.MustRegister(app, app.RootCmd, migratecmd.Config{
        // enable auto creation of migration files when making collection changes in the Dashboard
        // (the isGoRun check is to enable it only during development)
        Automigrate: isGoRun,
    })

	app.OnServe().BindFunc(func(se *core.ServeEvent) error {
		// serves static files from the provided public dir (if exists)
		routes.RegisterUserRoutes(se)
		return se.Next()
	})
	if err := app.Start(); err != nil {
		log.Fatal(err)
	}
	// err := godotenv.Load()
	// if err != nil {
	// 	log.Fatalf("Error loading .env file")
	// }

	// dbKey := os.Getenv("QDRANT_DB_KEY")

	// // Create a new Qdrant client with API key authentication and TLS enabled
	// client, err := qdrant.NewClient(&qdrant.Config{
	// 	Host:   "ec9d9f59-6f8c-4cdc-ae05-fa7bc0a465e7.us-west-2-0.aws.cloud.qdrant.io",
	// 	Port:   6334,
	// 	APIKey: dbKey,
	// 	UseTLS: true,
	// })
	// if err != nil {
	// 	panic(err)
	// }

	// collections, err := client.ListCollections(context.Background())
	// if err != nil {
	// 	panic(err)
	// }

	// fmt.Println(collections)
	// defer client.Close()

	// apiKey := os.Getenv("GROQ_API_KEY")

	// llm, err := openai.New(
	// 	openai.WithModel("llama3-8b-8192"),
	// 	openai.WithBaseURL("https://api.groq.com/openai/v1"),
	// 	openai.WithToken(apiKey),
	// )
	// if err != nil {
	// 	log.Fatal(err)
	// }

	// query := "Hi Team, I’ve been facing a couple of issues with the app recently. First, the app crashes every time I try to log in, especially when my internet connection is slow. Also, push notifications seem to have stopped working entirely after the last update, even though they’re enabled in the settings. On another note, it would be great if you could add a dark mode option to make the app easier to use at night. Please let me know if you need more details to look into these issues.Thanks,[Customer Name]"

	// tasks := qdrant_api.GetTasks(query, llm)

	// options := fastembed.InitOptions{
	// 	Model:     fastembed.BGESmallEN, // Correct identifier
	// 	CacheDir:  "models",
	// 	MaxLength: 512,
	// }

	// embedder, err := fastembed.NewFlagEmbedding(&options)
	// if err != nil {
	// 	panic(fmt.Sprintf("Model initialization failed: %v", err))
	// }
	// defer embedder.Destroy()

	// for _, task := range tasks {
	// 	isBug := false
	// 	taskQuery := strings.Split(task, "=")

	// 	if taskQuery[1] == "true" {
	// 		isBug = true
	// 	}

	// 	dataPoint := qdrant_api.DataPoint{
	// 		Content:   taskQuery[0],
	// 		IsBug:     isBug,
	// 		RepCount:  1,
	// 		Priority:  0.5,
	// 		Timestamp: time.Now(),
	// 		Status:    qdrant_api.Pending,
	// 	}
	// 	//payload := createPayload(dataPoint)

	// 	embedding, err := embedder.Embed([]string{taskQuery[0]}, 25)
	// 	if err != nil {
	// 		panic(fmt.Sprintf("Embedding failed: %v", err))
	// 	}
	// 	dataPoint.Embedding = embedding[0]

	// 	limit := uint64(1)

	// 	results, err := client.Query(context.Background(), &qdrant.QueryPoints{
	// 		CollectionName: "task_data",
	// 		Query:          qdrant.NewQuery(embedding[0]...),
	// 		Limit:          &limit,
	// 		WithPayload:    qdrant.NewWithPayload(true),
	// 		WithVectors:    qdrant.NewWithVectors(true),
	// 	})
	// 	if err != nil {
	// 		panic(err)
	// 	}

	// 	countResponse, err := client.Count(context.Background(), &qdrant.CountPoints{
	// 		CollectionName: "task_data",
	// 	})
	// 	if err != nil {
	// 		panic(err)
	// 	}
	// 	var result map[string]*qdrant.Value
	// 	if len(results) > 0 {
	// 		for _, value := range results[0].GetPayload() {
	// 		result = value.GetStructValue().GetFields()
	// 	}
	// 	}
	

	// 	fmt.Println("Result: ", result["Content"].GetStringValue())

	// 	if len(results)>0 && qdrant_api.CompareStrings(taskQuery[0], result["Content"].GetStringValue(), llm) {
	// 		fmt.Println("Similar task found")
	// 		dataPoint.Content = result["Content"].GetStringValue()
	// 		dataPoint.IsBug, _ = strconv.ParseBool(result["IsBug"].GetStringValue())
	// 		dataPoint.RepCount, _ = strconv.Atoi(result["RepCount"].GetStringValue())
	// 		dataPoint.RepCount++
	// 		dataPoint.Timestamp, _ = time.Parse(time.RFC3339, result["Timestamp"].GetStringValue())
	// 		dataPoint.Embedding = results[0].GetVectors().GetVector().Data
	// 		qdrant_api.UpdateAndCreateDataPoint(client, dataPoint, int(results[0].GetId().GetNum()))
	// 	} else {
	// 		dataPoint.Content = qdrant_api.ExpandTask(taskQuery[0], llm)
	// 		qdrant_api.UpdateAndCreateDataPoint(client, dataPoint, int(countResponse)+1)
	// 	}

	// 	fmt.Printf("Content: %s\n", dataPoint.Content)
	// 	fmt.Printf("IsBug: %t\n", dataPoint.IsBug)
	// 	fmt.Printf("RepCount: %d\n", dataPoint.RepCount)
	// 	fmt.Printf("Priority: %.2f\n", dataPoint.Priority)
	// 	fmt.Printf("Timestamp: %s\n", dataPoint.Timestamp.Format(time.RFC3339))
	// 	fmt.Printf("Status: %s\n", dataPoint.Status.String())
	// 	// fmt.Print("Embedding: \n", dataPoint.Embedding)
	// 	fmt.Println()

	// }

}



