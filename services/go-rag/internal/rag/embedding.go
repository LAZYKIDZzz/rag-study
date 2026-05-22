package rag

import (
	"hash/fnv"
	"math"
	"strings"
)

const embeddingDimensions = 64

type embeddingProvider interface {
	Embed(text string) []float64
	Name() string
}

type deterministicEmbeddingProvider struct{}

func newEmbeddingProvider() embeddingProvider {
	return deterministicEmbeddingProvider{}
}

func (deterministicEmbeddingProvider) Name() string {
	return "deterministic-hash"
}

func (deterministicEmbeddingProvider) Embed(text string) []float64 {
	vector := make([]float64, embeddingDimensions)
	for _, token := range tokenize(text) {
		hash := fnv.New64a()
		_, _ = hash.Write([]byte(token))
		sum := hash.Sum64()
		index := int(sum % embeddingDimensions)
		sign := 1.0
		if sum&1 == 0 {
			sign = -1.0
		}
		vector[index] += sign
	}
	return normalizeVector(vector)
}

func tokenize(text string) []string {
	fields := strings.Fields(strings.ToLower(text))
	tokens := make([]string, 0, len(fields))
	for _, field := range fields {
		token := strings.Trim(field, ".,!?;:\"'()[]{}<>")
		if token != "" {
			tokens = append(tokens, token)
		}
	}
	return tokens
}

func normalizeVector(vector []float64) []float64 {
	var sum float64
	for _, value := range vector {
		sum += value * value
	}
	if sum == 0 {
		return vector
	}

	magnitude := math.Sqrt(sum)
	for i := range vector {
		vector[i] = vector[i] / magnitude
	}
	return vector
}
