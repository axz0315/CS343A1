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

func read(path string) []string {
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

func single_threaded(files []string) {
	// initializes a map which will keep track of strings and the number of occurences of that string
	wordCounts := make(map[string]int)

	// Function to clean and split text into words
	cleanAndSplit := func(text string) []string {
		re := regexp.MustCompile(`[[:alnum:]]+`)
		words := re.FindAllString(text, -1)
		for i := range words {
			words[i] = strings.ToLower(words[i])
		}
		return words
	}

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
			for _, word := range words {
				wordCounts[word]++
			}
		}

		if err := scanner.Err(); err != nil {
			log.Printf("Error reading file %s: %v", filePath, err)
		}
	}

	// Sort the words alphabetically
	var words []string
	for word := range wordCounts {
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
		outputFile.WriteString(fmt.Sprintf("%s %d\n", word, wordCounts[word]))
	}
}

func multi_threaded(files []string) {
	// TODO: Your multi-threaded implementation
	// Get the total size of the files
	for _, filePath := range files {
		// Get the file size
		file, err := os.Stat(filePath)
		if err != nil {
			log.Printf("Error with determining file %s: %v", filePath, err)
		}
		// get the size
		size += int(file.Size())
	}
	// Determine the size of the string that each thread will handle
	size_per_thread := size / num_threads
	log.Printf("Size: %d", size_per_thread)
	bytes_per_thread := int64(1250000)

	var strings []string

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
		// myReader := bufio.NewReader(Myfile)

		for i := 0; i < int(file.Size()); i += int(bytes_per_thread) {
			Myfile.Seek(int64(i), 0)

			buffer := make([]byte, bytes_per_thread)
			n, err := Myfile.Read(buffer)
			if err != nil {
				log.Fatal(err)
			}

			strings = append(strings, string(buffer[:n]))

		}
	}
	fmt.Println(strings)
	fmt.Println(len(strings))

}

func main() {
	// TODO: add argument processing and run both single-threaded and multi-threaded functions

	files := read("/Users/bellasteedly/Library/Mobile Documents/com~apple~CloudDocs/Academics/Year4/Semester2/CS343/Assignment1/starter/input")
	// bella path: "/Users/bellasteedly/Library/Mobile Documents/com~apple~CloudDocs/Academics/Year4/Semester2/CS343/Assignment1/starter/input"
	// single_threaded(files)
	multi_threaded(files)
	// multi_threaded(files)
}
