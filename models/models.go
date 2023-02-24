package models

type (
	DirEntry interface {
		Name() string
		IsDir() bool
	}
)
