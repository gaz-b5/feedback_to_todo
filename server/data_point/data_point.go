package data_point

import (
	// "context"
	// "fmt"
	// "log"
	// "strconv"
	// "strings"
	"time"
	// "github.com/qdrant/go-client/qdrant"
)

type Status int

const (
	Pending Status = iota
	InProgress
	Completed
)

type DataPoint struct {
	Content   string
	IsBug     bool
	RepCount  int
	Priority  int
	Timestamp time.Time
	Status    Status
	Embedding []float32
	TaskID    string
}

func (s Status) String() string {
	return [...]string{"Pending", "In Progress", "Completed"}[s]
}
