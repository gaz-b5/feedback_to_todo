package llm_functions

import (
	"context"
	"fmt"
	"log"

	// "strconv"
	"strings"
	// "time"

	"github.com/anush008/fastembed-go"
	"github.com/tmc/langchaingo/llms"
	"github.com/tmc/langchaingo/llms/openai"
)

var LLM *openai.LLM

var MODEL *fastembed.FlagEmbedding

func GetTasks(query string) []string {
	// Define the prompt for the LLM
	prompt := `You are a task extraction assistant. Given the email report from a customer, extract the tasks and their details to be sent to the engineering team. Do not add any greeting or ending sentence, stick to the format given, do not add any index like "Task1=", etc. The tasks should be in the following format,where the <bool> implies a bool value which is true if the task is reported bug, or if it is a requested feature, where <description> implies the description of the task. To differentiate a Bug from a feature, ask yourself, if whatever the user is reporting, is it an mistake in the current functionality of the product, if it is then it is a bug, if the user is suggesting adding something for convinience, or add a functionality, then it is a feature.:
	<description>=<bool>;<description>=<bool>;...`

	ctx := context.Background()

	fullPrompt := prompt + "\n\nCustomer Message:\n" + query

	// Use the correct Call method
	completion, err := LLM.Call(ctx, fullPrompt,
		llms.WithTemperature(0.1),
	)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("LLM Response:\n", completion)

	return strings.Split(completion, ";")
}

func CompareStrings(a string, b string) bool {
	prompt := `You are a string comparison assistant. Given two strings, determine if they are similar or not, that is are both the strings pointing out the same issue or not. Respond with "true" if they are similar and "false" if they are not.`

	ctx := context.Background()

	fullPrompt := prompt + "\n\nString 1:\n" + a + "\n\nString 2:\n" + b

	// Use the correct Call method
	completion, err := LLM.Call(ctx, fullPrompt,
		llms.WithTemperature(0.1),
	)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("LLM Response:\n", completion)

	return strings.ToLower(completion) == "true"
}

func ExpandTask(task string) string {
	// Define the prompt for the LLM
	// TODO: Want to make it so that the llm can get the actual values from all task history.
	prompt := `You are a task expansion assistant. Given a task, expand it to include more details and context, so that a proper description of the task is created. Make it ATLEAST 2 to 3 sentences long. Just simply respond with the expanded task. Avoid any greeting or ending sentence. Do not use placeholder variables for unknown values.`

	ctx := context.Background()

	fullPrompt := prompt + "\n\nTask:\n" + task

	// Use the correct Call method
	completion, err := LLM.Call(ctx, fullPrompt,
		llms.WithTemperature(0.1),
	)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("LLM Response(description):\n", completion)

	return completion
}

func GenerateTitle(task string) string {
	// Define the prompt for the LLM
	prompt := `You are a title generation assistant. Given a task description, generate a short title for the task. Make it ATMOST 7-8 words long. Just simply respond with the title for the task. Avoid any greeting or ending sentence. Do not use placeholder variables for unknown values. Do not use quotes.`

	ctx := context.Background()

	fullPrompt := prompt + "\n\nTask:\n" + task

	// Use the correct Call method
	completion, err := LLM.Call(ctx, fullPrompt,
		llms.WithTemperature(0.1),
	)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("LLM Response(title):\n", completion)

	return completion
}

func GenerateEmbedding(query string) []float32 {
	embedding, err := MODEL.Embed([]string{query}, 25)
	if err != nil {
		panic(fmt.Sprintf("Embedding failed: %v", err))
	}
	return embedding[0]
}
