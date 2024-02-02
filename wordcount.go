package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strings"
)

var num_threads = 3
var size = 0
var wordCount = make(map[string]int)
var threadInput []string

func readFilesFromFolder(path string) []string {
	var files []string

	// Open the directory
	dir, err := os.Open(path)
	if err != nil {
		log.Fatal(err)
	}
	defer dir.Close()

	// Read the contents of the directory
	fileInfos, err := dir.Readdir(-1)
	if err != nil {
		log.Fatal(err)
	}

	// Construct the full path for each file and add it to the slice
	for _, fileInfo := range fileInfos {
		if fileInfo.IsDir() {
			// Skip directories, you can include them if needed
			continue
		}

		filePath := filepath.Join(path, fileInfo.Name())
		files = append(files, filePath)
	}

	return files
}

// Function to clean and split text into words
// can happen at same time
func cleanAndSplit(text string) []string {
	re := regexp.MustCompile(`[[:alnum:]]+`)
	words := re.FindAllString(text, -1)
	for i := range words {
		words[i] = strings.ToLower(words[i])
	}
	return words
}

// critical
func fillHashMap(words []string) {
	// Update word frequency count
	for _, word := range words {
		wordCount[word]++
	}
}

func generateOutputFile() {
	// Sort the words alphabetically
	var words []string
	for word := range wordCount {
		words = append(words, word)
	}
	sort.Strings(words)

	// Write the word count to the output file
	outputFile, err := os.Create("output/single.txt")
	if err != nil {
		log.Fatal(err)
	}
	defer outputFile.Close()

	for _, word := range words {
		outputFile.WriteString(fmt.Sprintf("%s %d\n", word, wordCount[word]))
	}
}

func single_threaded(files []string) {
	// initializes a map which will keep track of strings and the number of occurences of that string
	// Process each file
	for _, filePath := range files {
		// Open the file
		file, err := os.Open(filePath)
		if err != nil {
			log.Printf("Error opening file %s: %v", filePath, err)
			continue
		}
		defer file.Close()

		// Read file line by line
		scanner := bufio.NewScanner(file)
		for scanner.Scan() {
			// Clean and split each line into words
			words := cleanAndSplit(scanner.Text())

			// Update word frequency count
			fillHashMap(words)
		}

		if err := scanner.Err(); err != nil {
			log.Printf("Error reading file %s: %v", filePath, err)
		}
	}

	generateOutputFile()
}

func multi_threaded(files []string) {
	// Set the size of the string that each thread will handle
	bytes_per_thread := int64(1250000)

	// Split the file into strings with size_per_thread bytes in each
	for _, filePath := range files {
		Myfile, err := os.Open(filePath)
		if err != nil {
			fmt.Println(err)
		}
		file, err := os.Stat(filePath)
		if err != nil {
			log.Printf("Error with determining file %s: %v", filePath, err)
		}

		for i := 0; i < int(file.Size()); i += int(bytes_per_thread) {
			Myfile.Seek(int64(i), 0)

			buffer := make([]byte, bytes_per_thread)
			n, err := Myfile.Read(buffer)
			if err != nil {
				log.Fatal(err)
			}

			threadInput = append(threadInput, string(buffer[:n]))

		}
	}
}

func main() {
	// TODO: add argument processing and run both single-threaded and multi-threaded functions

	files := readFilesFromFolder("/Users/bellasteedly/Library/Mobile Documents/com~apple~CloudDocs/Academics/Year4/Semester2/CS343/Assignment1/starter/input")
	// bella path: "/Users/bellasteedly/Library/Mobile Documents/com~apple~CloudDocs/Academics/Year4/Semester2/CS343/Assignment1/starter/input"
	single_threaded(files)
	// multi_threaded(files)
	// multi_threaded(files)
}
