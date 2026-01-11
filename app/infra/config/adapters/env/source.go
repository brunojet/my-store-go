package env

import (
	"os"

	"github.com/brunojet/my-store-go/app/infra/config/contracts"
)

type source struct{}

// New creates a ports.Source that reads from environment variables.
func New() contracts.Source {
	return source{}
}

func (source) Lookup(key string) (string, bool) {
	return os.LookupEnv(key)
}
