package rag

import (
	"math"
	"sort"
	"sync"
)

type vectorStore struct {
	mu     sync.RWMutex
	chunks map[string]Chunk
}

func newVectorStore() *vectorStore {
	return &vectorStore{
		chunks: make(map[string]Chunk),
	}
}

func (s *vectorStore) replaceDocumentChunks(documentID string, chunks []Chunk) {
	s.mu.Lock()
	defer s.mu.Unlock()

	for id, chunk := range s.chunks {
		if chunk.DocumentID == documentID {
			delete(s.chunks, id)
		}
	}
	for _, chunk := range chunks {
		chunk.Metadata = cloneStringMap(chunk.Metadata)
		chunk.Embedding = cloneFloatSlice(chunk.Embedding)
		s.chunks[chunk.ID] = chunk
	}
}

func (s *vectorStore) search(queryEmbedding []float64, topK int, filters map[string]string) []RetrievedChunk {
	s.mu.RLock()
	defer s.mu.RUnlock()

	if topK <= 0 {
		topK = 5
	}

	results := make([]RetrievedChunk, 0, len(s.chunks))
	for _, chunk := range s.chunks {
		if !metadataMatches(chunk.Metadata, filters) {
			continue
		}

		results = append(results, RetrievedChunk{
			ChunkID:       chunk.ID,
			DocumentID:    chunk.DocumentID,
			DocumentTitle: chunk.Metadata["document_title"],
			Ordinal:       chunk.Ordinal,
			Content:       chunk.Content,
			Metadata:      cloneStringMap(chunk.Metadata),
			Score:         cosine(queryEmbedding, chunk.Embedding),
		})
	}

	sort.Slice(results, func(i int, j int) bool {
		if results[i].Score == results[j].Score {
			return results[i].ChunkID < results[j].ChunkID
		}
		return results[i].Score > results[j].Score
	})

	if len(results) > topK {
		results = results[:topK]
	}
	return results
}

func cosine(left []float64, right []float64) float64 {
	limit := len(left)
	if len(right) < limit {
		limit = len(right)
	}

	var dot float64
	var leftMagnitude float64
	var rightMagnitude float64
	for i := 0; i < limit; i++ {
		dot += left[i] * right[i]
		leftMagnitude += left[i] * left[i]
		rightMagnitude += right[i] * right[i]
	}
	if leftMagnitude == 0 || rightMagnitude == 0 {
		return 0
	}
	return dot / (math.Sqrt(leftMagnitude) * math.Sqrt(rightMagnitude))
}

func metadataMatches(metadata map[string]string, filters map[string]string) bool {
	for key, expected := range filters {
		if metadata[key] != expected {
			return false
		}
	}
	return true
}
