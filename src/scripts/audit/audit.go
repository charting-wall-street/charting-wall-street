package main

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

func countAll(content string) int {
	regex := regexp.MustCompile(`[^a-zA-Z0-9]+`)
	content = regex.ReplaceAllString(content, " ")
	words := strings.Fields(content)
	return len(words)
}

func contains(s string, token string) bool {
	lowercase := strings.ToLower(s)
	return strings.Contains(lowercase, token)
}

func removeComments(input string) string {
	regex := regexp.MustCompile(`<!--.*?-->`)
	return regex.ReplaceAllString(input, "")
}

func countRefs(content string) int {
	regex := regexp.MustCompile(`\\ref\{[^\}]+\}`)
	matches := regex.FindAllString(content, -1)
	return len(matches)
}

func countCits(content string) int {
	regex := regexp.MustCompile(`\[@[^]]+\]`)
	matches := regex.FindAllString(content, -1)
	return len(matches)
}

func countTODORefs(content string) int {
	regex := regexp.MustCompile(`\\ref\{[^\}]*todo[^\}]*\}`)
	matches := regex.FindAllString(strings.ToLower(content), -1)
	return len(matches)
}

func countTODOCitations(content string) int {
	regex := regexp.MustCompile(`\[@ref\]`)
	matches := regex.FindAllString(strings.ToLower(content), -1)
	return len(matches)
}

type SectionEntry struct {
	FileName  string
	WordCount int
	Body      string
}

var minimalWords = 24000
var minimalRefToWord = 0.1

func main() {

	files, err := filepath.Glob("./src/sections/*.md")
	if err != nil {
		panic(err)
	}

	entries := make([]*SectionEntry, 0)

	totalWords := 0
	allCitations := 0
	for _, path := range files {
		fileName := filepath.Base(path)
		fileNameWithoutExt := fileName[:len(fileName)-len(filepath.Ext(fileName))]

		contents, err := os.ReadFile(path)
		if err != nil {
			panic(err)
		}

		if len(contents) < 100 {
			continue
		}

		cleaned := removeComments(string(contents))
		allCitations += countCits(cleaned)
		count := countAll(cleaned)
		totalWords += count
		entries = append(entries, &SectionEntry{
			FileName:  fileNameWithoutExt,
			WordCount: count,
			Body:      string(contents),
		})
	}

	todoCitations := 0
	fmt.Printf("%-60s %10s %10s\n", "Filename", "Count", "Percentage")
	for _, entry := range entries {
		if entry.FileName[0] == '0' {
			continue
		}
		percentage := float64(entry.WordCount) / float64(totalWords) * 100
		fmt.Printf("%-60s %10d %10.2f%%\n", entry.FileName, entry.WordCount, percentage)
		if contains(entry.Body, "todo") {
			fmt.Println(" - TODO")
		}
		if contains(entry.Body, "note") {
			fmt.Println(" - NOTE")
		}
		citCount := countTODOCitations(entry.Body)
		if citCount > 0 {
			fmt.Println(" - CITATION")
		}
		todoCitations += citCount
	}

	fmt.Println()
	wordCompletion := float64(totalWords) / float64(minimalWords) * 100
	refCompletion := float64(allCitations-todoCitations) / (float64(todoCitations)) * 100

	fmt.Printf("Words: %-8d %.2f%% completion\n", totalWords, wordCompletion)
	fmt.Printf("Refs:  %d/%-8d %.2f%% completion\n", allCitations-todoCitations, allCitations, refCompletion)
}
