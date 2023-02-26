package operations

import (
	"fmt"
	"github.com/go-lean/bevaluate/config"
	"github.com/go-lean/bevaluate/storage"
	"gopkg.in/yaml.v3"
	"io"
)

type (
	InitOperation struct {
		store InitStore
	}

	InitStore interface {
		TryAccessing(path string) error
		OpenCreate(path string) (io.WriteCloser, error)
	}
)

func NewInitOperation(store InitStore) InitOperation {
	return InitOperation{store: store}
}

func (o InitOperation) Run(path string) error {
	errExist := o.store.TryAccessing(path)
	if errExist == nil {
		return fmt.Errorf("config file aleady exists")
	}

	if errExist != storage.ErrNotExisting {
		return fmt.Errorf("could not determine if config file aleady exists")
	}

	cfg := config.Default()
	data, errMarshal := yaml.Marshal(cfg)
	if errMarshal != nil {
		return fmt.Errorf("could not marshal config: %w", errMarshal)
	}

	file, errCreate := o.store.OpenCreate(path)
	if errCreate != nil {
		return fmt.Errorf("could not create config file: %w", errCreate)
	}
	defer func() {
		_ = file.Close()
	}()

	if _, errWrite := file.Write(data); errWrite != nil {
		return fmt.Errorf("could not write to config file: %w", errWrite)
	}

	return nil
}
