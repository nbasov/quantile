package main

import (
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
)

const (
	ErrMessageJSONParse    = "Malformed line"
	ErrMessageNonNumeric   = "Non-numeric field value"
	ErrMessageMissingField = "Missing field"
)

func parseChan(lineCh chan []byte, parser *JSONFieldNumericValueParser) chan ChannelValue {
	resCh := make(chan ChannelValue)

	go func() {
		for line := range lineCh {
			val, errMessage := parser.Parse(line)
			resCh <- ChannelValue{
				Value:        val,
				ErrorMessage: errMessage,
			}
		}

		close(resCh)
	}()

	return resCh
}

// JSONFieldNumericValueParser is a parser that extracts a numeric value from a specified field in a JSON object.
type JSONFieldNumericValueParser struct {
	Stat Statistics

	fieldName string
	fieldsMap map[string]json.RawMessage
}

func NewJSONFieldNumericValueParser(fieldName string) JSONFieldNumericValueParser {
	return JSONFieldNumericValueParser{
		Stat:      Statistics{},
		fieldName: fieldName,
		fieldsMap: make(map[string]json.RawMessage),
	}
}

// Parse parses a single line of JSON, extracting the numeric value from the specified field.
func (p *JSONFieldNumericValueParser) Parse(line []byte) (float32, string) {
	p.Stat.TotalLinesCount++
	return p.parseLineUnmarshal(line)
}

func (p *JSONFieldNumericValueParser) parseLineUnmarshal(line []byte) (val float32, errMessage string) {
	clear(p.fieldsMap)

	err := json.Unmarshal(line, &p.fieldsMap)
	if err != nil {
		p.Stat.MalformedLinesCount++
		return 0, ErrMessageJSONParse
	}

	raw, ok := p.fieldsMap[p.fieldName]
	if !ok {
		p.Stat.MissingFieldLinesCount++
		return 0, ErrMessageMissingField
	}

	valParsed, err := strconv.ParseFloat(string(raw), 32)
	if err != nil {
		p.Stat.NonNumericFieldLines++
		return 0, ErrMessageNonNumeric
	}

	return float32(valParsed), ""
}

// Statistics holds the statistics of the parsing process, including counts of total lines, malformed lines,
// missing field lines, and non-numeric field value lines.
type Statistics struct {
	TotalLinesCount        int64
	MalformedLinesCount    int64
	MissingFieldLinesCount int64
	NonNumericFieldLines   int64
}

func (stat Statistics) String() string {
	sb := strings.Builder{}

	if stat.MalformedLinesCount > 0 {
		fmt.Fprintf(&sb, "Malformed lines: %d\n", stat.MalformedLinesCount)
	}
	if stat.MissingFieldLinesCount > 0 {
		fmt.Fprintf(&sb, "Missing field lines: %d\n", stat.MissingFieldLinesCount)
	}
	if stat.NonNumericFieldLines > 0 {
		fmt.Fprintf(&sb, "Non-numeric field value lines: %d\n", stat.NonNumericFieldLines)
	}

	fmt.Fprintf(&sb, "Total lines: %d\n", stat.TotalLinesCount)

	return strings.TrimSuffix(sb.String(), "\n")
}
