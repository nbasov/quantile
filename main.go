package main

import (
	"bufio"
	"errors"
	"flag"
	"fmt"
	"os"
	"path/filepath"
)

func main() {
	p, fieldName, filePath, verbose, err := prepareParams()
	if err != nil {
		fmt.Println("Error:", err.Error())
		os.Exit(1)
	}

	if verbose {
		fmt.Printf("Arguments: p=%g, field=%s, file=%s\n\n", p, fieldName, filePath)
	}

	parser := NewJSONFieldNumericValueParser(fieldName)
	//computer := NewQuantileSimple(p)
	computer := NewQuantileTwoHeaps(p)

	// Pipeline definition
	linesCh := readLinesChan(filePath)
	valueCh := parser.ParseChan(linesCh)
	quantileCh := computeChan(valueCh, &computer)

	// Results printing
	w := bufio.NewWriter(os.Stdout)
	for q := range quantileCh {
		fmt.Fprintln(w, q.String())
	}

	if verbose {
		fmt.Fprintln(w, "\n"+parser.Stat.String())
	}

	w.Flush()
}

func prepareParams() (p float32, fieldName string, filePath string, verbose bool, err error) {
	p, fieldName, verbose, err = parseAndValidateArgs()
	if err != nil {
		return 0, "", "", false, err
	}

	fmt.Println("Enter file path and press Enter:")
	filePath, err = waitForFilePath()
	if err != nil {
		return 0, "", "", false, err
	}

	return p, fieldName, filePath, verbose, nil
}

func parseAndValidateArgs() (float32, string, bool, error) {
	q := flag.Float64("q", 0, "Quantile in range [0.0, 1.0] - REQUIRED")
	fieldName := flag.String("f", "", "JSON field name for quantile computation, case sensitive - REQUIRED")
	verbose := flag.Bool("v", false, "Verbose output")
	flag.Parse()

	if flag.NFlag() < 2 {
		flag.CommandLine.Usage()
		return 0, "", false, errors.New("missing required command line arguments")
	}

	if *q < 0 || *q > 1 {
		return 0, "", false, errors.New("quantile should be in range [0.0, 1.0]")
	}

	if len(*fieldName) == 0 {
		return 0, "", false, errors.New("JSON field name should not be empty")
	}

	return float32(*q), *fieldName, *verbose, nil
}

func waitForFilePath() (string, error) {
	scanner := bufio.NewScanner(os.Stdin)
	scanner.Scan()
	filePath := scanner.Text()

	if !filepath.IsAbs(filePath) {
		return "", errors.New("absolute file path required")
	}

	_, err := os.Stat(filePath)
	if os.IsNotExist(err) {
		return "", errors.New("file does not exist")
	} else if err != nil {
		return "", err
	}

	return filePath, nil
}

// ChannelValue represents a value along with an optional error message.
type ChannelValue struct {
	Value        float32
	ErrorMessage string
}

func (cv ChannelValue) String() string {
	if cv.ErrorMessage != "" {
		return cv.ErrorMessage
	}

	return fmt.Sprintf("%g", cv.Value)
}
