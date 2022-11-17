package csv

import (
	"bytes"
	"encoding/csv"
	"io"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCsv(t *testing.T) {
	data := bytes.NewBuffer([]byte{})
	records := [][]string{
		{"first_name", "last_name", "username"},
		{"Rob", "Pike", "rob"},
		{"Ken", "Thompson", "ken"},
		{"Robert", "Griesemer", "gri"},
	}

	/*
	   Write
	*/
	w := csv.NewWriter(data)
	for _, record := range records {
		err := w.Write(record)
		assert.NoError(t, err)
	}
	// Write any buffered data to the underlying writer (standard output).
	w.Flush()
	assert.NoError(t, w.Error())

	/*
	   Read
	*/
	r := csv.NewReader(data)

	var result = make([][]string, 0, len(records))
	for {
		record, err := r.Read()
		if err == io.EOF {
			break
		}
		assert.NoError(t, err)
		result = append(result, record)
	}

	assert.Equal(t, records, result)
}
