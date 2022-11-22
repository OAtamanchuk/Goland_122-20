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
	"path/filepath"
)

func main() {
	inputFileName := flag.String("i", "", "Use a file with the name file-name as an input.")
	outputFileName := flag.String("o", "", "Use a file with the name file-name as an output.")
	sortingFilesIndex := flag.Int("f", 0, "Sort input lines by value number N.")
	isNotIgnoreHeader := flag.Bool("h", false, "The first line is a header that must be ignored during sorting but included in the output.")
	isReversedOrder := flag.Bool("r", false, "Sort input lines in reverse order.")
	inputDirectory := flag.String("d", "", "Specify a directory where the application must read input files from.")
	flag.Parse()

	if *inputFileName != "" && *inputDirectory != "" {
		fmt.Println("You can use only one of these flags: -i/-d")
		return
	}

	var content string

	content = Sort(ReadFile(ReadDirectory(content, *isNotIgnoreHeader, *inputDirectory), *isNotIgnoreHeader, *inputFileName), *sortingFilesIndex, *isReversedOrder)

	if content != "" {
		fmt.Println("Sorted data:\n" + content)

		if *outputFileName == "" {
			dateTimeNow := time.Now()
			*outputFileName = strconv.Itoa(dateTimeNow.Year()) + "-" + dateTimeNow.Month().String() + "-" + strconv.Itoa(dateTimeNow.Day()) + "_" + strconv.Itoa(dateTimeNow.Hour()) + "-" + strconv.Itoa(dateTimeNow.Minute()) + "-" + strconv.Itoa(dateTimeNow.Second()) + ".csv"
		}

		WriteToFile(content, *outputFileName)
	}
}

func ReadDirectory(content string, isNotIgnoreHeader bool, inputDirectory string) string {
	if inputDirectory != "" {
		err := filepath.Walk(inputDirectory,
			func(path string, info os.FileInfo, err error) error {
				if err != nil {
					return err
				}

				if strings.HasSuffix(path, ".csv") {
					content = content + ReadFromFile(isNotIgnoreHeader, path)
				}
				
				return nil
			})

		if err != nil {
			fmt.Println(err)
		}
	}

	return content
}

func ReadFile(content string, isNotIgnoreHeader bool, inputFileName string) string {
	if inputFileName == "" {
		content = content + ReadFromConsole(isNotIgnoreHeader)
	} else {
		content = content + ReadFromFile(isNotIgnoreHeader, inputFileName)
	}

	return content
}

func ReadFromConsole(isNotIgnoreHeader bool) string {
	scanner := bufio.NewScanner(os.Stdin)
	return StartProcessing(isNotIgnoreHeader, scanner)
}

func ReadFromFile(isNotIgnoreHeader bool, inputFile string) string {
	file, err := os.Open(inputFile)

	if err != nil {
		log.Fatal(err)
	}

	defer file.Close()
	fileScanner := bufio.NewScanner(file)
	fileScanner.Split(bufio.ScanLines)

	content := StartProcessing(isNotIgnoreHeader, fileScanner)
	return content
}

func StartProcessing(isNotIgnoreHeader bool, scanner *bufio.Scanner) string {
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

func Sort(content string, sortingFieldIndex int, isReversedOrder bool) string {
	table := [][]string{}

	rows := strings.Split(content, "\n")

	for i := 0; i < len(rows); i++ {
		cols := strings.Split(rows[i], ",")
		table = append(table, cols)
	}

	sort.Slice(table, func(i, j int) bool {
		return Compare(table[i][sortingFieldIndex], table[j][sortingFieldIndex], isReversedOrder)
	})

	var result strings.Builder
	
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