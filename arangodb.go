package arangom

import "fmt"

var (
	// ErrNoCollection is returned when no collection is provided.
	ErrNoCollection = fmt.Errorf("no collection provided")
	// ErrNoDatabase is returned when no database is provided.
	ErrNoDatabase = fmt.Errorf("no database provided")
)
