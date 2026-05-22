package rag

import (
	"fmt"
	"strings"
)

func cloneStringMap(input map[string]string) map[string]string {
	output := make(map[string]string)
	for key, value := range input {
		output[key] = value
	}
	return output
}

func cloneFloatSlice(input []float64) []float64 {
	return append([]float64(nil), input...)
}

func normalizeWhitespace(input string) string {
	return strings.Join(strings.Fields(input), " ")
}

func intToFixed(value int) string {
	return fmt.Sprintf("%04d", value)
}

func trimForAnswer(input string, limit int) string {
	input = normalizeWhitespace(input)
	if len(input) <= limit {
		return input
	}
	if limit <= 3 {
		return input[:limit]
	}
	return input[:limit-3] + "..."
}

func chunksForResponse(chunks []Chunk) []Chunk {
	output := make([]Chunk, 0, len(chunks))
	for _, chunk := range chunks {
		chunk.Metadata = cloneStringMap(chunk.Metadata)
		chunk.Embedding = nil
		output = append(output, chunk)
	}
	return output
}
