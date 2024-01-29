package api

// TODO problems with imports otherwise
import (
	_ "github.com/ogen-go/ogen"
	_ "github.com/ogen-go/ogen/gen"
)

//go:generate go run github.com/ogen-go/ogen/cmd/ogen --target . --package api --clean openapi.yaml
