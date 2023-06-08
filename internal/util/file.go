package util

import (
	"bufio"
	"fmt"
	"math/rand"
	"os"
)

func DeduplicateLinesRandomly(inputFilepath string, outputFilepath string) (int, error) {
	var outputFile *os.File
	var err error
	if outputFilepath == "-" {
		outputFile = os.Stdout
	} else {
		outputFile, err = os.Create(outputFilepath)
		if err != nil {
			return -1, err
		}
		defer outputFile.Close()
	}

	file, err := os.Open(inputFilepath)
	if err != nil {
		return -1, err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	uniqueLines := make(map[string]bool)

	for scanner.Scan() {
		line := scanner.Text()
		if !uniqueLines[line] {
			uniqueLines[line] = true
		}
	}

	if err := scanner.Err(); err != nil {
		return -1, err
	}

	lines := make([]string, 0, len(uniqueLines))
	for line := range uniqueLines {
		lines = append(lines, line)
	}

	rand.Shuffle(len(lines), func(i, j int) {
		lines[i], lines[j] = lines[j], lines[i]
	})

	for _, line := range lines {
		fmt.Fprintln(outputFile, line)
	}

	return len(lines), nil
}
