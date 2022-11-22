package main

import (
	"bufio"
	"flag"
	"fmt"
	"log"
	"os"
	"sort"
	"strings"
	"time"
	"strconv"
)

func main() {
	inputFileName := flag.String("i", "", "Use a file with the name file-name as an input.")
	outputFileName := flag.String("o", "", "Use a file with the name file-name as an output.")
	sortingFilesIndex := flag.Int("f", 0, "Sort input lines by value number N.")
	isNotIgnoreHeader := flag.Bool("h", false, "The first line is a header that must be ignored during sorting but included in the output.")
	isReversedOrder := flag.Bool("r", false, "Sort input lines in reverse order.")
	flag.Parse()

	var content string
	if *inputFileName == "" {
		content = ReadFromConsole(*sortingFilesIndex, *isReversedOrder, *isNotIgnoreHeader)
	} else {
		content = ReadFromFile(*sortingFilesIndex, *isReversedOrder, *isNotIgnoreHeader, *inputFileName)
	}

	if content != "" {
		fmt.Println("Sorted data:\n" + content)

		if *outputFileName == "" {
			dateTimeNow := time.Now()
			*outputFileName = strconv.Itoa(dateTimeNow.Year()) + "-" + dateTimeNow.Month().String() + "-" + strconv.Itoa(dateTimeNow.Day()) + "_" + strconv.Itoa(dateTimeNow.Hour()) + "-" + strconv.Itoa(dateTimeNow.Minute()) + "-" + strconv.Itoa(dateTimeNow.Second()) + ".csv"
		}

		WriteToFile(content, *outputFileName)
	}
}

func ReadFromConsole(sortingFieldIndex int, isReversedOrder, isNotIgnoreHeader bool) string {
	scanner := bufio.NewScanner(os.Stdin)
	return StartProcessing(sortingFieldIndex, isReversedOrder, isNotIgnoreHeader, scanner)
}

func ReadFromFile(sortingFieldIndex int, isReversedOrder, isNotIgnoreHeader bool, inputFile string) string {
	file, err := os.Open(inputFile)

	if err != nil {
		log.Fatal(err)
	}

	defer file.Close()
	fileScanner := bufio.NewScanner(file)
	fileScanner.Split(bufio.ScanLines)

	content := StartProcessing(sortingFieldIndex, isReversedOrder, isNotIgnoreHeader, fileScanner)
	return content
}

func StartProcessing(sortingFieldIndex int, isReversedOrder, isNotIgnoreHeader bool, scanner *bufio.Scanner) string {
	var header string
	n := 0
	table := [][]string{}

	for scanner.Scan() {
		line := scanner.Text()
		row := strings.Split(line, ",")

		if n == 0 {
			n = len(row)
			if isNotIgnoreHeader {
				header = line
				continue
			}
		}

		if line == "" {
			break
		}

		if n != len(row) {
			log.Fatalf("Row has %d columns, but must have %d\n", len(row), n)
		}

		table = append(table, row)
	}

	if scanner.Err() != nil {
		log.Fatal(scanner.Err())
	}

	sort.Slice(table, func(i, j int) bool {
		return Compare(table[i][sortingFieldIndex], table[j][sortingFieldIndex], isReversedOrder)
	})

	var result strings.Builder

	if header != "" {
		result.WriteString(header)
		result.WriteString("\n")
	}

	for _, row := range table {
		result.WriteString(strings.Join(row, ","))
		result.WriteString("\n")
	}

	return result.String()
}

func WriteToFile(content, fileName string) {
	if fileName != "" {
		file, err := os.Create(fileName)

		if err != nil {
			log.Fatal(err)
		}

		defer file.Close()
		_, err = file.WriteString(content)

		if err != nil {
			log.Fatal(err)
		}
	}
}

func Compare(first, next string, isReversed bool) bool {
	return first < next != isReversed
}