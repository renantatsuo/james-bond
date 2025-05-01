package env

import (
	"log"
	"os"
	"strconv"
)

type EnvVar[T any] struct {
	name     string
	required bool
}

func Get(name string) EnvVar[any] {
	return EnvVar[any]{
		name:     name,
		required: false,
	}
}

func (e EnvVar[any]) Required() EnvVar[any] {
	e.required = true
	return e
}

func (e EnvVar[any]) String() EnvVar[string] {
	return EnvVar[string](e)
}

func (e EnvVar[any]) Int() EnvVar[int] {
	return EnvVar[int](e)
}

func (e EnvVar[T]) Parse() T {
	value := os.Getenv(e.name)
	if value == "" && e.required {
		panic("required environment variable is missing: " + e.name)
	}
	var t T
	switch any(t).(type) {
	case string:
		return any(value).(T)
	case int:
		v, err := strconv.Atoi(value)
		if err != nil {
			log.Printf("failed to parse int value %q", value)
		}
		return any(v).(T)
	default:
		log.Printf("unexpected type for value %q", value)
		return any(value).(T)
	}
}
