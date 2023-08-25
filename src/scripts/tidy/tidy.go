package main

import (
	"bufio"
	"fmt"
	"log"
	"math/rand"
	"os"
	"path"
	"strconv"
	"strings"
)

func randIdentifier() string {
	const charset = "ABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	identifier := make([]byte, 4)
	for i := range identifier {
		identifier[i] = charset[rand.Intn(len(charset))]
	}
	return string(identifier)
}

type Section struct {
	Id     string
	Hidden bool
	Hashes string
	Prefix string
	Title  string
	Header bool
}

func main() {

	file, err := os.Open("src/sections.txt")
	if err != nil {
		fmt.Println("Error opening file:", err)
		return
	}
	defer file.Close()

	sections := make([]Section, 0)
	scanner := bufio.NewScanner(file)

	reformattedLines := []string{}
	identifierMask := make(map[string]bool)
	depthMarker := make(map[string]bool)

	// Read the file line by line
	chapter := []int{0}
	lastPrefix := ""
	for scanner.Scan() {
		line := scanner.Text()
		titleStart := strings.IndexByte(line, ' ')
		if titleStart == -1 {
			continue
		}
		identifierStart := strings.IndexByte(line, '[')
		if identifierStart == -1 {
			identifierStart = len(line)
		}
		heading := strings.Trim(line[:titleStart], " ")
		newDepth := false
		if len(heading) > len(chapter) {
			if len(heading) != len(chapter)+1 {
				log.Fatalln("invalid section: ", line)
			}
			chapter = append(chapter, 0)
			newDepth = true
		} else if len(heading) < len(chapter) {
			chapter = chapter[:len(heading)]
		}
		chapter[len(chapter)-1]++

		indexS := make([]string, 0)
		for _, v := range chapter {
			indexS = append(indexS, strconv.Itoa(v))
		}
		prefix := strings.Join(indexS, ".")
		if newDepth {
			depthMarker[lastPrefix] = true
		}
		lastPrefix = prefix

		title := strings.Trim(line[titleStart:identifierStart], " ")

		identifier := strings.Trim(line[identifierStart:], "[] ")
		identifierMask[identifier] = true

		output := fmt.Sprintf("%s %s %s", heading, title, identifier)
		reformattedLines = append(reformattedLines, output)

		strings.Split(line, " ")

		sections = append(sections, Section{
			Id:     identifier,
			Hashes: heading,
			Prefix: prefix,
			Title:  title,
			Hidden: strings.IndexByte(title, '(') != -1,
		})
	}

	if err := scanner.Err(); err != nil {
		fmt.Println("Error scanning file:", err)
		return
	}

	sections = append(sections, Section{
		Id:     "H0",
		Hidden: false,
		Hashes: "",
		Prefix: "0.0",
		Title:  "METADATA",
		Header: true,
	})

	sections = append(sections, Section{
		Id:     "SUMM",
		Hidden: false,
		Hashes: "",
		Prefix: "0.10",
		Title:  "SUMMARY",
		Header: true,
	})

	sections = append(sections, Section{
		Id:     "AAAA",
		Hidden: false,
		Hashes: "",
		Prefix: "0.11",
		Title:  "ABSTRACT",
		Header: true,
	})

	sections = append(sections, Section{
		Id:     "PREF",
		Hidden: false,
		Hashes: "",
		Prefix: "0.12",
		Title:  "PREFACE",
		Header: true,
	})

	sections = append(sections, Section{
		Id:     "TOEL",
		Hidden: false,
		Hashes: "",
		Prefix: "0.13",
		Title:  "TTB",
		Header: true,
	})

	sections = append(sections, Section{
		Id:     "TABL",
		Hidden: false,
		Hashes: "",
		Prefix: "0.98",
		Title:  "TABLE",
		Header: true,
	})

	// Assign unique id to each section
	for i := range sections {
		if sections[i].Id == "" {
			for {
				newId := randIdentifier()
				if _, ok := identifierMask[newId]; ok {
					continue
				}
				identifierMask[newId] = true
				sections[i].Id = newId
				break
			}
		}
	}

	// Write the lines to a temporary file
	out, err := os.Create("src/sections.txt")
	if err != nil {
		fmt.Println("Error creating temporary file:", err)
		return
	}
	defer out.Close()

	for _, section := range sections {
		var output string
		if section.Header {
			continue
		} else if !section.Hidden {
			output = fmt.Sprintf("%s %s [%s]\n", section.Hashes, section.Title, section.Id)
		} else {
			output = fmt.Sprintf("%s %s\n", section.Hashes, section.Title)
		}
		_, _ = out.Write([]byte(output))
	}

	sectionsDir := path.Join("src", "sections")

	for _, section := range sections {

		if section.Hidden {
			continue
		}

		files, err := os.ReadDir(sectionsDir)
		if err != nil {
			fmt.Println("Error reading directory:", err)
			return
		}

		suffix := fmt.Sprintf("%s.md", section.Id)

		exists := ""
		for _, f := range files {
			if !f.IsDir() && strings.HasSuffix(f.Name(), suffix) {
				exists = f.Name()
				break
			}
		}
		prefix := section.Prefix
		if depthMarker[prefix] {
			prefix = prefix + ".0_"
		} else {
			prefix = prefix + "_"
		}

		title := strings.ReplaceAll(section.Title, " ", "_")
		name := fmt.Sprintf("%s%s_%s.md", prefix, title, section.Id)

		// Create new empty file when no file exists
		if exists == "" {
			f, err := os.Create(path.Join(sectionsDir, name))
			if err != nil {
				log.Fatalln(err)
			}
			f.Close()
			fmt.Println("Created ", name)
			continue
		}

		// Do not do anything when the file is as expected
		if name == exists {
			continue
		}

		// Move the file if the title has changed
		err = os.Rename(path.Join(sectionsDir, exists), path.Join(sectionsDir, name))
		if err != nil {
			log.Fatalln("Failed to rename file: ", err)
		}
		fmt.Println("Renamed ", exists, " -> ", name)
	}
}
