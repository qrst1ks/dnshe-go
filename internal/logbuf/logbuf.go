package logbuf

import (
	"fmt"
	"sync"
	"time"
)

type Entry struct {
	Time    string `json:"time"`
	Level   string `json:"level"`
	Message string `json:"message"`
}

type Buffer struct {
	mu      sync.RWMutex
	entries []Entry
	limit   int
}

func New(limit int) *Buffer {
	if limit <= 0 {
		limit = 200
	}
	return &Buffer{limit: limit}
}

func (b *Buffer) Addf(level, format string, args ...any) {
	message := fmt.Sprintf(format, args...)
	entry := Entry{
		Time:    time.Now().Format("2006-01-02 15:04:05"),
		Level:   level,
		Message: message,
	}

	b.mu.Lock()
	b.entries = append(b.entries, entry)
	if len(b.entries) > b.limit {
		b.entries = b.entries[len(b.entries)-b.limit:]
	}
	b.mu.Unlock()

	fmt.Printf("[%s] [%s] %s\n", entry.Time, entry.Level, entry.Message)
}

func (b *Buffer) List() []Entry {
	b.mu.RLock()
	defer b.mu.RUnlock()
	out := make([]Entry, len(b.entries))
	copy(out, b.entries)
	return out
}

func (b *Buffer) Clear() {
	b.mu.Lock()
	b.entries = nil
	b.mu.Unlock()
}
