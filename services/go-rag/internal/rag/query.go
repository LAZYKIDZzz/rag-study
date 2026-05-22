package rag

import "strings"

type queryProcessor struct{}

func newQueryProcessor() queryProcessor {
	return queryProcessor{}
}

func (queryProcessor) rewrite(query string, memory UserMemory) string {
	clean := strings.TrimSpace(query)
	if clean == "" {
		return clean
	}

	if len(memory.Facts) == 0 {
		return clean
	}

	return clean + " Context from user memory: " + strings.Join(memory.Facts, "; ")
}
