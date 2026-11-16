package adapters

import (
	"fmt"
	"sync/atomic"
)

type SequentialIDGenerator struct {
	counter atomic.Uint64 // гарантирует безопасное увеличение счётчика
}

func NewSequentialIDGenerator() *SequentialIDGenerator {
	return &SequentialIDGenerator{}
}

// принимает префикс STU/CRT/GRD
func (g *SequentialIDGenerator) NewID(prefix string) string {
	id := g.counter.Add(1)
	return fmt.Sprintf("%s-%04d", prefix, id)
}
