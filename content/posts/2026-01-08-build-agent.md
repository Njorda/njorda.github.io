---
layout: post
title: "Building a coding agent to learn how it works"
subtitle: "Yet another coding agent implemented in go"
date: 2026-01-12
author: "Niklas Hansson"
URL: "/2025/01/12/coding-agent"
---


In this blog post I will do what is so often(at least for the stuff I read) recomended these days and implement a coding agent from "scratch". The perception should be that it it is kind of like [The emporors New Clothes](https://en.wikipedia.org/wiki/The_Emperor%27s_New_Clothes). Where the coding agent is actually relativly small and dont contain that much. We will take it in steps and start with hard coding some tools before we also add some skills that we allow the agent to use.

Please be aware that you need to be careful with this and executing LLM generade code. Thus we will containerize it but it will note be ANYTHING production read so be aware of your own computer and your stuff, now you have been warned so lets go!

For the sake of simplicity I will assume you have docker installed. We will just copy the code in to the a go container and link the content, we will give the container read permissions to one folder and thus "ANYTHING CAN BE REMOVED IN THAT FOLDER", so be carefull, now you have been warned twice and I will not bring it up again. 


```bash
docker run -it --rm \
  -v $(pwd)/go/code:/app \
  -e GEMINI_API_KEY=$GEMINI_API_KEY \
  -w /app \
  golang:1.21 \
  go run main.go
```

As you can see above, we will in this case use [Gemini](https://gemini.google.com) for the models and specifically we will use 


The first step of our go code is the following gemini-2.5-flash, which is great bang for the buck and more than enough what we need for doing this. 

The first step is to set up the boiler plate go stuf with go init and add the go.main file. 

```bash
go init github.
```

The first thing to understand it that an agent sounds way more fancy that it is, in our case it is a "infinite loop" with an LLM call where we parse the return to figure our if we want to make a tool call(call a go function) or return to the user for new input. Thats it, there is no more magic. However the tricky part is to get this to work by having good promgs, descriptions and good tools in order for this loop to be able to validate the work and jump between the tools to make progress, similar to how a human works. 


```go
func main() {
	// This starts the main loop of the agent.
	agent, err := NewAgent(context.Background(), os.Getenv("GEMINI_API_KEY"))
	if err != nil {
		panic(err)
	}

	if err := agent.agenLoop(); err != nil {
		panic(err)
	}
}
```

The main function pretty much only starts the agent and let it run, i set it up a struct with some methods in order to just make the code more clean. The main thing to take awaya and understand is the `agentLoop` where non magic is happening. 


```go

func (a *agent) agenLoop() error {
	systemPrompt, err := a.getSystemPrompt()
	if err != nil {
		return fmt.Errorf("failed to get system prompt: %v", err)
	}
	previousResponse := []*genai.Content{}
	for {
		reader := bufio.NewReader(os.Stdin)
		fmt.Print("Question: ")
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
				fmt.Println("Calling tool...")
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
					return fmt.Errorf("failed to parse tool args: %v argsJSON: %v", err, argsJSON)
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
```

So lets take it in steps: 

1) The user inputs some kind of question or action they want the agent to do. 
2) The agent calls the LLM with the following content: The prompt, the tool descriptions and the user query. 
3) The response is parsed, if a tool is detected we parse the arguments and call the tool. If a tool is called we call the LLM again, as long as a tool is called we keep on taking the response and calling the LLM. For each response we add it to the history. 

As you migth notice we have cut some corners for example:
- We dont do anything about context length
- We dont limit the nbr of LLM calls in any way

The tools are injected dynamically for each call, in our case here there are allways one set of tools, but it could very much be dynamicallty injected or selected per call do keep the context resonable. There is some interestig work currently on moving more to [dynamic context](https://cursor.com/blog/dynamic-context-discovery)(due to the model improvements)

The tools follow this struct schema in order to try to give the LLM as much context about the tools as possible while keeping it resonable short as well: 

```go
type tool struct {
	Name        string                      `json:"name"`
	Function    func(args ...string) string `json:"-"`
	Args        []string                    `json:"args"`
	Description string                      `json:"description"`
}
```

I implemented three different tools that we use in the coding agent: 
- list, list files in folder
- read, reads file
- edit, edits files. The way we do it is by simply doing a string replace and overwriting the file. 

Lets take a look at one of the tools as example: 

```go
func edit(args ...string) string {
	if len(args) != 3 {
		return fmt.Errorf("error: fileName, old code and new code are required in order to do string replacement").Error()
	}
	content, err := os.ReadFile(args[0])
	if err != nil {
		return fmt.Errorf("failed to read file: %v", err).Error()
	}

	newContent := strings.Replace(string(content), args[1], args[2], 1)
	if err := os.WriteFile(args[0], []byte(newContent), 0644); err != nil {
		return fmt.Errorf("failed to write file: %v", err).Error()
	}

	return string(newContent)
}
```

This tool could also be extendedn in order to handle when we want to create a new file, so when the input file dont exsist we could just allow it to create the file and input the code. 

That said, the coolest part of this is that even with just theese three functions we can make our little coding agent to improve upon it. 

Thanks for reading