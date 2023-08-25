package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"sort"
	"strings"
)

var apiKey = os.Getenv("API_KEY")
var apiUrl = "https://api.openai.com/v1/completions"

type Response struct {
	ID      string   `json:"id"`
	Object  string   `json:"object"`
	Created int64    `json:"created"`
	Model   string   `json:"model"`
	Choices []Choice `json:"choices"`
	Usage   struct {
		PromptTokens     int `json:"prompt_tokens"`
		CompletionTokens int `json:"completion_tokens"`
		TotalTokens      int `json:"total_tokens"`
	} `json:"usage"`
}

type Choice struct {
	Text         string  `json:"text"`
	Index        int     `json:"index"`
	Logprobs     *string `json:"logprobs"`
	FinishReason string  `json:"finish_reason"`
}

func Summarize(body string) string {
	payload := map[string]interface{}{
		"model":             "text-davinci-003",
		"prompt":            "Can you briefly summarize the following section of a research paper, in a couple of sentences:\n\n" + body + "\n\n\nSummary:\n\n",
		"temperature":       0,
		"max_tokens":        100,
		"top_p":             1.0,
		"frequency_penalty": 0.0,
		"presence_penalty":  0.0,
	}

	payloadBytes, err := json.Marshal(payload)
	if err != nil {
		log.Fatalln("Error marshaling payload:", err)
	}

	req, err := http.NewRequest("POST", apiUrl, bytes.NewBuffer(payloadBytes))
	if err != nil {
		log.Fatalln("Error creating request:", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+apiKey)

	client := http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Fatalln("Error sending request:", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		log.Fatalf("Request failed, got status code: %d\n", resp.StatusCode)
	}

	var response Response
	err = json.NewDecoder(resp.Body).Decode(&response)
	if err != nil {
		log.Fatalln("Error decoding response:", err)
	}

	if len(response.Choices) == 0 {
		fmt.Println("Prompt Tokens:", response.Usage.PromptTokens)
		fmt.Println("Completion Tokens:", response.Usage.CompletionTokens)
		fmt.Println("Total Tokens:", response.Usage.TotalTokens)
		log.Fatalln("No response!")
	}

	choice := response.Choices[0]
	return choice.Text
}

func fileExists(filename string) bool {
	_, err := os.Stat(filename)
	return err == nil
}

func concatFiles() {

	outputDir := "./src/summaries/"

	fileInfos, err := os.ReadDir(outputDir)
	if err != nil {
		panic(err)
	}

	sort.Slice(fileInfos, func(i, j int) bool {
		return fileInfos[i].Name() < fileInfos[j].Name()
	})

	var concatenatedSummary strings.Builder

	for _, fileInfo := range fileInfos {
		if !fileInfo.IsDir() && filepath.Ext(fileInfo.Name()) == ".txt" {
			filePath := filepath.Join(outputDir, fileInfo.Name())

			fileContent, err := os.ReadFile(filePath)
			if err != nil {
				fmt.Printf("Error reading file %s: %v\n", fileInfo.Name(), err)
				continue
			}

			concatenatedSummary.WriteString(fileInfo.Name())
			concatenatedSummary.WriteString("\n")
			concatenatedSummary.Write(fileContent)
			concatenatedSummary.WriteString("\n\n")
		}
	}

	outputFilePath := filepath.Join("./src/build", "summary.txt")
	err = os.WriteFile(outputFilePath, []byte(concatenatedSummary.String()), os.ModePerm)
	if err != nil {
		panic(err)
	}
}

func main() {

	files, err := filepath.Glob("./src/sections/*.md")
	if err != nil {
		panic(err)
	}

	for _, path := range files {
		fileName := filepath.Base(path)
		fileNameWithoutExt := fileName[:len(fileName)-len(filepath.Ext(fileName))]

		fmt.Print(fileNameWithoutExt)

		if fileExists("./src/summaries/" + fileNameWithoutExt + ".txt") {
			fmt.Println(": skipped (file already exists)")
			continue
		}

		contents, err := os.ReadFile(path)
		if err != nil {
			panic(err)
		}

		if len(contents) < 100 {
			fmt.Println(": skipped (no content)")
			continue
		}

		// Generate summaries to verify of the core of each section is clearly conveyed
		summary := Summarize(string(contents))

		// Trim the summary to the last dot
		lastDotIndex := strings.LastIndex(summary, ".")
		if lastDotIndex != -1 {
			summary = summary[:lastDotIndex+1]
		}

		out, err := os.Create("./src/summaries/" + fileNameWithoutExt + ".txt")
		if err != nil {
			panic(err)
		}
		out.Write([]byte(summary))
		out.Close()

		fmt.Println(": done")
	}

	concatFiles()
}
