package util

import (
	"bufio"
	"fmt"
	"math/rand"
	"os"
)

func DeduplicateLinesRandomly(inputFilepath string, outputFilepath string) error {
	file, err := os.Open(inputFilepath)
	if err != nil {
		return err
	}
	defer file.Close()

	outputFile, err := os.Create(outputFilepath)
	if err != nil {
		return err
	}
	defer outputFile.Close()

	scanner := bufio.NewScanner(file)
	uniqueLines := make(map[string]bool)

	for scanner.Scan() {
		line := scanner.Text()
		if !uniqueLines[line] {
			uniqueLines[line] = true
		}
	}

	if err := scanner.Err(); err != nil {
		return err
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

	return nil
}
