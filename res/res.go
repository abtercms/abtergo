package res

import (
	"embed"
)

//go:embed content/* templates/*
var content embed.FS

func Read(path string) ([]byte, error) {
	return content.ReadFile(path)
}
