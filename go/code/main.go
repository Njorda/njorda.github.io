package main

import (
	"bufio"
	"context"
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"google.golang.org/genai"
)

func main() {
	// This starts the main loop of the agent.
	fmt.Println("Hello, World!")
	agent, err := NewAgent(context.Background(), os.Getenv("GEMINI_API_KEY"))
	if err != nil {
		panic(err)
	}

	if err := agent.agenLoop(); err != nil {
		panic(err)
	}
}

// Struct to keep the tool definitions
// The json tags are kind of important to keep the tool descriptions complete to the LLM.
type tool struct {
	Name     string                      `json:"name"`
	Function func(args ...string) string `json:"-"`
	Args     []string                    `json:"args"`
}

// The struct for the LLM
type agent struct {
	client *genai.Client
}

func NewAgent(ctx context.Context, apiKey string) (*agent, error) {

	client, err := genai.NewClient(ctx, &genai.ClientConfig{
		APIKey:  apiKey,
		Backend: genai.BackendGeminiAPI,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create client: %v", err)
	}
	return &agent{
		client: client,
	}, nil
}

var tools = map[string]tool{
	"read": {
		Name:     "read",
		Function: read,
		Args:     []string{"fileName"},
	},
	"edit": {
		Name:     "edit",
		Function: edit,
		Args:     []string{"fileName", "oldCode", "newCode"},
	},
	"list": {
		Name:     "list",
		Function: list,
		Args:     []string{"fileName"},
	},
}

func (a *agent) getToolNames() (string, error) {
	toolNames, err := json.Marshal(tools)
	if err != nil {
		return "", fmt.Errorf("failed to marshal tools: %v", err)
	}
	return string(toolNames), nil
}

func (a *agent) getSystemPrompt() (string, error) {
	toolNames, err := a.getToolNames()
	if err != nil {
		return "", fmt.Errorf("failed to get tool names: %v", err)
	}
	return fmt.Sprintf(`
	You are a coding assistant whose goal it is to help us solve coding tasks. 
	You have access to a series of tools you can execute. Here are the tools you can execute:

	%v

	When you want to use a tool, reply with exactly one line in the format: 'tool: TOOL_NAME({"arg1":"value1"})' and nothing else.
	Use the EXACT argument names from the tool definitions above as JSON keys.
	Use compact single-line JSON with double quotes. After receiving a tool_result(...) message, continue the task.
	If no tool is needed, respond normally.
	`, toolNames), nil
}

func read(args ...string) string {
	if len(args) == 0 {
		return "error: fileName argument is required"
	}
	content, err := os.ReadFile(args[0])
	if err != nil {
		return fmt.Errorf("failed to read file: %v", err).Error()
	}
	return string(content)
}

func edit(args ...string) string {
	if len(args) == 3 {
		return fmt.Errorf("error: fileName, old code and new code are required in order to do string replacement").Error()
	}
	content, err := os.ReadFile(args[0])
	if err != nil {
		return fmt.Errorf("failed to read file: %v", err).Error()
	}

	newContent := strings.Replace(string(content), args[1], args[2], 1)

	return string(newContent)
}

func list(args ...string) string {
	fmt.Println("Listing...")
	files, err := os.ReadDir(".")
	if err != nil {
		return fmt.Errorf("failed to read directory: %v", err).Error()
	}
	var fileNames []string
	for _, file := range files {
		fileNames = append(fileNames, file.Name())
	}
	return strings.Join(fileNames, "\n")
}

func toolArgs(toolArgsJSON string, expectedArgs []string) ([]string, error) {
	var args []string
	if toolArgsJSON == "" || toolArgsJSON == "{}" {
		return args, nil
	}

	// Parse the JSON into a map
	var argsMap map[string]string
	if err := json.Unmarshal([]byte(toolArgsJSON), &argsMap); err != nil {
		return nil, fmt.Errorf("failed to parse tool args: %v", err)
	}

	// Extract args in the order defined by expectedArgs
	for _, argName := range expectedArgs {
		if val, ok := argsMap[argName]; ok {
			args = append(args, val)
		}
	}
	return args, nil
}

func (a *agent) agenLoop() error {
	systemPrompt, err := a.getSystemPrompt()
	if err != nil {
		return fmt.Errorf("failed to get system prompt: %v", err)
	}
	previousResponse := []*genai.Content{}
	for {
		reader := bufio.NewReader(os.Stdin)
		fmt.Print("Enter a command: ")
		input, err := reader.ReadString('\n')
		if err != nil {
			return fmt.Errorf("failed to read command: %v", err)
		}
		userInput := genai.NewContentFromText(input, genai.RoleUser)
		previousResponse = append(previousResponse, userInput)
		for {
			systemContent := genai.NewContentFromText(systemPrompt, genai.RoleUser)
			contents := append([]*genai.Content{systemContent}, previousResponse...)

			response, err := a.client.Models.GenerateContent(
				context.Background(),
				"gemini-2.5-flash",
				contents,
				nil,
			)
			if err != nil {
				return fmt.Errorf("failed to generate content: %v", err)
			}

			responseText := response.Candidates[0].Content.Parts[0].Text
			fmt.Println(responseText)
			previousResponse = append(previousResponse, response.Candidates[0].Content)

			if strings.HasPrefix(responseText, "tool:") {
				// Extract tool name
				toolName := strings.Split(responseText, "tool:")[1]
				toolName = strings.Split(toolName, "(")[0]
				toolName = strings.TrimSpace(toolName)

				// Extract JSON args from between ( and )
				startIdx := strings.Index(responseText, "(")
				endIdx := strings.LastIndex(responseText, ")")
				argsJSON := ""
				if startIdx != -1 && endIdx != -1 && endIdx > startIdx {
					argsJSON = responseText[startIdx+1 : endIdx]
				}

				tool, ok := tools[toolName]
				if !ok {
					return fmt.Errorf("tool not found: %v", toolName)
				}

				// Parse args
				parsedArgs, err := toolArgs(argsJSON, tool.Args)
				if err != nil {
					return fmt.Errorf("failed to parse tool args: %v", err)
				}

				result := tool.Function(parsedArgs...)
				fmt.Println(result)
				previousResponse = append(previousResponse, genai.NewContentFromText(result, genai.RoleUser))
			} else {
				break
			}
		}
	}
}
