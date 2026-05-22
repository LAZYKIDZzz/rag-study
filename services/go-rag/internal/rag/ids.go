package rag

import (
	"fmt"
	"sync/atomic"
)

type idGenerator struct {
	value atomic.Uint64
}

func (g *idGenerator) next(prefix string) string {
	id := g.value.Add(1)
	return fmt.Sprintf("%s_%06d", prefix, id)
}
