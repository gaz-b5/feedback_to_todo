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
	prompt := `You are a highly accurate task extraction assistant. Given an email report from a customer, your goal is to extract all relevant tasks along with their detailed classifications and comprehensive descriptions, precisely formatted for the engineering team to act on.

	Instructions:

	Do not include any greetings, closings, explanations, or extra text.

	Follow the exact output format:
	<description>=<bool>;<description>=<bool>;...
	where each <description> is a clear, comprehensive, and detailed task description including all relevant context and specifics from the report, and <bool> is either true or false.

	The boolean value indicates whether the task is a bug (true) or a feature request (false).

	To determine if a task is a bug:

	Ask yourself: Is the customer reporting a problem, error, or malfunction in the product's current behavior?

	If yes, it is a bug (true).

	To determine if a task is a feature request:

	Ask yourself: Is the customer requesting new functionality, an enhancement, or something that does not currently exist but would improve convenience or capability?

	If yes, it is a feature (false).

	Only include tasks explicitly mentioned or clearly implied by the customer’s report.

	Avoid numbering, bullet points, or any additional formatting—stick precisely to task descriptions.

	Separate multiple tasks with semicolons (;) with no trailing semicolon at the end.

	Note: When constructing each task description, include all relevant details, examples, affected components, user scenarios, or error messages mentioned in the report to provide the engineering team with sufficient context to understand and prioritize the task properly.`

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
	prompt := `You are a string comparison assistant. Given two text strings, your goal is to determine whether both strings describe the same issue, problem, or task, even if phrased differently. Consider the meaning, context, and intent behind the strings, rather than just exact wording.

	If both strings point to the same core issue or task, respond with the string: "true".

	If the strings describe different issues, problems, or tasks, respond with the string: "false".

	Ignore minor phrasing differences, synonyms, or wording variations.

	Focus on the underlying problem or request conveyed.

	Do not provide any additional explanation or text, only respond "true" or "false" exactly.

	Example:

	String A: "The login button does not respond when clicked."

	String B: "Clicking the login button has no effect."

	Response: "true"

	String A: "Add an option to export reports as PDF."

	String B: "The export function fails to generate any output."

	Response: "false"`

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
