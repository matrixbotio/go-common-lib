package config

import "fmt"

type EmptyFieldError struct {
	Field string
}

func (e EmptyFieldError) Error() string {
	return fmt.Sprintf("Field '%s' is empty", e.Field)
}

type LogLevelIncorrectError struct {
	Level string
	Field string
}

func (e LogLevelIncorrectError) Error() string {
	return fmt.Sprintf("Level '%s' for field '%s' is incorrect", e.Level, e.Field)
}
