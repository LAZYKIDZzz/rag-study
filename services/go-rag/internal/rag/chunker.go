package rag

import "strings"

type chunker struct {
	maxWords int
	overlap  int
}

func newChunker() chunker {
	return chunker{
		maxWords: 120,
		overlap:  20,
	}
}

func (c chunker) split(doc Document) []Chunk {
	words := strings.Fields(normalizeWhitespace(doc.Content))
	if len(words) == 0 {
		return nil
	}

	step := c.maxWords - c.overlap
	if step <= 0 {
		step = c.maxWords
	}

	var chunks []Chunk
	for start, ordinal := 0, 0; start < len(words); start, ordinal = start+step, ordinal+1 {
		end := start + c.maxWords
		if end > len(words) {
			end = len(words)
		}

		metadata := cloneStringMap(doc.Metadata)
		metadata["document_title"] = doc.Title

		chunks = append(chunks, Chunk{
			ID:         doc.ID + "_chunk_" + intToFixed(ordinal),
			DocumentID: doc.ID,
			Ordinal:    ordinal,
			Content:    strings.Join(words[start:end], " "),
			Metadata:   metadata,
		})

		if end == len(words) {
			break
		}
	}

	return chunks
}
