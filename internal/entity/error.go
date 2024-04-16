package entity

import (
	"fmt"
	"strings"
)

var (
	ErrorConflict = NewErrConflict("object")
	ErrorNotFound = NewErrNotFound("object")
)

// error not found
type ErrNotFound struct {
	name string
}

func (e *ErrNotFound) Error() string {
	return e.name + " not found"
}

func NewErrNotFound(text string) *ErrNotFound {
	return &ErrNotFound{text}
}

// error conflict
type ErrConflict struct {
	name string
}

func (e *ErrConflict) Error() string {
	return e.name + " already exist"
}

func NewErrConflict(text string) *ErrConflict {
	return &ErrConflict{text}
}

// error validation
type ErrValidation struct {
	Err    error
	Errors map[string]string
}

func (e ErrValidation) Error() string {
	return e.Err.Error()
}

func NewErrValidation() *ErrValidation {
	return &ErrValidation{Errors: make(map[string]string)}
}

type ErrNoRequiredParameter struct {
	parameters []string
}

func NewErrNoRequiredParameter(parameters ...string) *ErrNoRequiredParameter {
	return &ErrNoRequiredParameter{parameters: parameters}
}

func (e ErrNoRequiredParameter) Error() string {
	var str strings.Builder
	for _, param := range e.parameters {
		if str.Len() != 0 {
			str.WriteString(", ")
		}
		str.WriteString(param)
	}

	return fmt.Sprintf("could not find required parameter(s): %s", str.String())
}
