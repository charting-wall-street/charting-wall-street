package main

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"sort"
	"strconv"
	"strings"
)

type ChapterEntry struct {
	Info FileInfo
	Body string
}

func removeComments(input string) string {
	regex := regexp.MustCompile(`<!--.*?-->`)
	return regex.ReplaceAllString(input, "")
}

func addHeader(title string, level int, content string) string {
	if level < 1 {
		level = 1
	} else if level > 6 {
		level = 6
	}
	header := strings.Repeat("#", level) + " " + title
	if level == 1 {
		header = "\\newpage\n" + header
	}
	return header + "\n\n" + content
}

type FileInfo struct {
	FileName   string
	Prefix     string
	Title      string
	Identifier string
	Header     bool
}

func parseHeaderPrefix(prefix string) int {
	parts := strings.Split(prefix, ".")
	lastPart, err := strconv.Atoi(parts[len(parts)-1])
	if err != nil {
		return 1
	}
	if lastPart == 0 {
		return len(parts) - 1
	} else {
		return len(parts)
	}
}

func parseFilePath(filePath string) FileInfo {
	fileName := filepath.Base(filePath)
	parts := strings.Split(fileName, "_")
	identifier := strings.TrimSuffix(parts[len(parts)-1], ".md")
	title := strings.Join(parts[1:len(parts)-1], " ")
	title = strings.ReplaceAll(title, "_", " ")
	return FileInfo{FileName: filePath, Prefix: parts[0], Title: title, Identifier: identifier, Header: parts[0][0] == '0'}
}

func countWords(content string) int {
	regex := regexp.MustCompile(`[^a-zA-Z0-9]+`)
	content = regex.ReplaceAllString(content, " ")
	words := strings.Fields(content)
	return len(words)
}

func addPlaceholder(content string) string {
	if strings.TrimSpace(content) == "" {
		return ""
	} else {
		return content
	}
}

func addSupplemental(name string, output *bytes.Buffer) {
	output.WriteString(readContents("./src/supplemental/" + name + ".md"))
}

func readContents(src string) string {
	contents, err := os.ReadFile(src)
	if err != nil {
		panic(err)
	}
	return string(contents)
}

func writeOutput(output *bytes.Buffer) {
	f, err := os.Create("./src/build/main.md")
	if err != nil {
		panic(err)
	}
	defer f.Close()
	_, err = f.Write(output.Bytes())
	if err != nil {
		panic(err)
	}
}

func main() {

	files, err := filepath.Glob("./src/sections/*.md")
	if err != nil {
		panic(err)
	}

	fmt.Println("Open and merge sections...")
	fileStructs := make([]ChapterEntry, 0)

	totalWords := 0
	for _, path := range files {
		info := parseFilePath(path)
		text := readContents(path)
		if !info.Header {
			text = removeComments(text)
			totalWords += countWords(text)
			text = addPlaceholder(text)
			headerLevel := parseHeaderPrefix(info.Prefix)
			text = addHeader(info.Title, headerLevel, text)
		}
		fileStruct := ChapterEntry{
			Info: info,
			Body: text,
		}
		fileStructs = append(fileStructs, fileStruct)
	}

	sort.Slice(fileStructs, func(i, j int) bool {
		return fileStructs[i].Info.FileName < fileStructs[j].Info.FileName
	})

	// Write main
	output := bytes.Buffer{}
	for _, file := range fileStructs {
		output.WriteString(file.Body)
		output.WriteString("\n\n\n")
	}
	addSupplemental("appendix", &output)
	addSupplemental("bibliography", &output)
	writeOutput(&output)

	fmt.Printf("Project contains %d words. %.2f%% complete.\n", totalWords, float64(totalWords)/24000*100)

	fmt.Println("Compiling abstract...")
	compileAbstract("en")
	compileAbstract("nl")

	fmt.Println("Compiling thesis...")
	compileThesis()

	fmt.Println("Done!")
}

func compileThesis() {
	cmd := exec.Command("pandoc", "--citeproc", "--csl=../vancouver.csl", "--bibliography=../references.bib", "../build/main.md", "-o", "../build/main.pdf", "--pdf-engine=pdflatex")
	cmd.Dir = "./src/sections"
	cmdout, err := cmd.CombinedOutput()
	if err != nil {
		fmt.Printf("Error: %s\n", err.Error())
		fmt.Printf("Command output: %s\n", cmdout)
		os.Exit(1)
	}
}

func compileAbstract(lang string) {
	cmd := exec.Command("pandoc", "--citeproc", "--csl=../vancouver.csl", "--bibliography=../references.bib", "./abstract_"+lang+".md", "-o", "../build/abstract_"+lang+".pdf", "--pdf-engine=pdflatex")
	cmd.Dir = "./src/abstract"
	cmdout, err := cmd.CombinedOutput()
	if err != nil {
		fmt.Printf("Error: %s\n", err.Error())
		fmt.Printf("Command output: %s\n", cmdout)
		os.Exit(1)
	}
}
