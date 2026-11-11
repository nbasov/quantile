package main

import "testing"

func TestJSONFieldNumericValueParser_Parse(t *testing.T) {
	tests := []struct {
		name      string
		fieldName string
		line      string
		wantVal   float32
		wantErr   string
		wantStat  Statistics
	}{
		{
			name:      "integer value",
			fieldName: "value",
			line:      `{"value": 42}`,
			wantVal:   42,
			wantStat:  Statistics{TotalLinesCount: 1},
		},
		{
			name:      "float value",
			fieldName: "value",
			line:      `{"value": 1.5}`,
			wantVal:   1.5,
			wantStat:  Statistics{TotalLinesCount: 1},
		},
		{
			name:      "negative value",
			fieldName: "value",
			line:      `{"value": -3.25}`,
			wantVal:   -3.25,
			wantStat:  Statistics{TotalLinesCount: 1},
		},
		{
			name:      "zero value",
			fieldName: "value",
			line:      `{"value": 0}`,
			wantVal:   0,
			wantStat:  Statistics{TotalLinesCount: 1},
		},
		{
			name:      "exponent notation",
			fieldName: "value",
			line:      `{"value": 1e3}`,
			wantVal:   1000,
			wantStat:  Statistics{TotalLinesCount: 1},
		},
		{
			name:      "field among other fields",
			fieldName: "value",
			line:      `{"id": "abc", "value": 7, "tags": [1, 2], "meta": {"value": 100}}`,
			wantVal:   7,
			wantStat:  Statistics{TotalLinesCount: 1},
		},
		{
			name:      "key in different case",
			fieldName: "value",
			line:      `{"VaLuE": 5}`,
			wantErr:   ErrMessageMissingField,
			wantStat:  Statistics{TotalLinesCount: 1, MissingFieldLinesCount: 1},
		},
		{
			name:      "field name in different case",
			fieldName: "VALUE",
			line:      `{"value": 5}`,
			wantErr:   ErrMessageMissingField,
			wantStat:  Statistics{TotalLinesCount: 1, MissingFieldLinesCount: 1},
		},
		{
			name:      "extra whitespace around value",
			fieldName: "value",
			line:      "{ \"value\" :\t 9 \n}",
			wantVal:   9,
			wantStat:  Statistics{TotalLinesCount: 1},
		},
		{
			name:      "empty line",
			fieldName: "value",
			line:      ``,
			wantErr:   ErrMessageJSONParse,
			wantStat:  Statistics{TotalLinesCount: 1, MalformedLinesCount: 1},
		},
		{
			name:      "broken json",
			fieldName: "value",
			line:      `{"value": 1`,
			wantErr:   ErrMessageJSONParse,
			wantStat:  Statistics{TotalLinesCount: 1, MalformedLinesCount: 1},
		},
		{
			name:      "plain text",
			fieldName: "value",
			line:      `not a json`,
			wantErr:   ErrMessageJSONParse,
			wantStat:  Statistics{TotalLinesCount: 1, MalformedLinesCount: 1},
		},
		{
			name:      "json array instead of object",
			fieldName: "value",
			line:      `[1, 2, 3]`,
			wantErr:   ErrMessageJSONParse,
			wantStat:  Statistics{TotalLinesCount: 1, MalformedLinesCount: 1},
		},
		{
			name:      "json number instead of object",
			fieldName: "value",
			line:      `42`,
			wantErr:   ErrMessageJSONParse,
			wantStat:  Statistics{TotalLinesCount: 1, MalformedLinesCount: 1},
		},
		{
			name:      "trailing garbage",
			fieldName: "value",
			line:      `{"value": 1} extra`,
			wantErr:   ErrMessageJSONParse,
			wantStat:  Statistics{TotalLinesCount: 1, MalformedLinesCount: 1},
		},
		{
			name:      "empty object",
			fieldName: "value",
			line:      `{}`,
			wantErr:   ErrMessageMissingField,
			wantStat:  Statistics{TotalLinesCount: 1, MissingFieldLinesCount: 1},
		},
		{
			name:      "field missing",
			fieldName: "value",
			line:      `{"other": 1}`,
			wantErr:   ErrMessageMissingField,
			wantStat:  Statistics{TotalLinesCount: 1, MissingFieldLinesCount: 1},
		},
		{
			name:      "field only in nested object",
			fieldName: "value",
			line:      `{"meta": {"value": 1}}`,
			wantErr:   ErrMessageMissingField,
			wantStat:  Statistics{TotalLinesCount: 1, MissingFieldLinesCount: 1},
		},
		{
			name:      "key is prefix of field name",
			fieldName: "value",
			line:      `{"val": 1}`,
			wantErr:   ErrMessageMissingField,
			wantStat:  Statistics{TotalLinesCount: 1, MissingFieldLinesCount: 1},
		},
		{
			name:      "numeric string",
			fieldName: "value",
			line:      `{"value": "1.5"}`,
			wantErr:   ErrMessageNonNumeric,
			wantStat:  Statistics{TotalLinesCount: 1, NonNumericFieldLines: 1},
		},
		{
			name:      "text string",
			fieldName: "value",
			line:      `{"value": "abc"}`,
			wantErr:   ErrMessageNonNumeric,
			wantStat:  Statistics{TotalLinesCount: 1, NonNumericFieldLines: 1},
		},
		{
			name:      "null",
			fieldName: "value",
			line:      `{"value": null}`,
			wantErr:   ErrMessageNonNumeric,
			wantStat:  Statistics{TotalLinesCount: 1, NonNumericFieldLines: 1},
		},
		{
			name:      "boolean",
			fieldName: "value",
			line:      `{"value": true}`,
			wantErr:   ErrMessageNonNumeric,
			wantStat:  Statistics{TotalLinesCount: 1, NonNumericFieldLines: 1},
		},
		{
			name:      "array",
			fieldName: "value",
			line:      `{"value": [1]}`,
			wantErr:   ErrMessageNonNumeric,
			wantStat:  Statistics{TotalLinesCount: 1, NonNumericFieldLines: 1},
		},
		{
			name:      "object",
			fieldName: "value",
			line:      `{"value": {"x": 1}}`,
			wantErr:   ErrMessageNonNumeric,
			wantStat:  Statistics{TotalLinesCount: 1, NonNumericFieldLines: 1},
		},
		{
			name:      "float32 overflow",
			fieldName: "value",
			line:      `{"value": 1e40}`,
			wantErr:   ErrMessageNonNumeric,
			wantStat:  Statistics{TotalLinesCount: 1, NonNumericFieldLines: 1},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := NewJSONFieldNumericValueParser(tt.fieldName)

			gotVal, gotErr := p.Parse([]byte(tt.line))

			if gotErr != tt.wantErr {
				t.Errorf("Parse(%q) error = %q, want %q", tt.line, gotErr, tt.wantErr)
			}
			if gotVal != tt.wantVal {
				t.Errorf("Parse(%q) value = %v, want %v", tt.line, gotVal, tt.wantVal)
			}
			if p.Stat != tt.wantStat {
				t.Errorf("Parse(%q) stat = %+v, want %+v", tt.line, p.Stat, tt.wantStat)
			}
		})
	}
}
