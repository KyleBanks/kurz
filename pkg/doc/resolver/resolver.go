package resolver

import (
	"errors"
	"io"
	"os"
	"path/filepath"

	"github.com/KyleBanks/kurz/pkg/doc"
)

var (
	// ErrInvalidPath indicates that the path provided to a Resolver
	// is not the appropriate type.
	//
	// For example, if you pass a remote URL to a resolver.File,
	// you'll receive this error.
	ErrInvalidPath = errors.New("cannot resolve the path provided")
)

// Chain runs a sequence of doc.Resolver types, returning the first
// content to be successfully resolved.
//
// If any of the resolvers returns an error other than ErrInvalidPath,
// the error will be returned. If none of the resolvers are able to
// resolve a content body, an ErrInvalidPath error is returned.
type Chain struct {
	Resolvers []doc.Resolver
}

func (c Chain) Resolve(path string) (io.ReadCloser, error) {
	for _, r := range c.Resolvers {
		content, err := r.Resolve(path)
		if err != nil && err != ErrInvalidPath {
			return nil, err
		} else if content != nil {
			return content, nil
		}
	}

	return nil, ErrInvalidPath
}

type File struct{}

func (File) Resolve(path string) (io.ReadCloser, error) {
	path, err := filepath.Abs(path)
	if err != nil {
		return nil, err
	}

	f, err := os.Open(path)
	if err != nil {
		if os.IsNotExist(err) {
			err = ErrInvalidPath
		}
		return nil, err
	}

	return f, nil
}

type URL struct{}

func (URL) Resolve(url string) (io.ReadCloser, error) {
	return nil, ErrInvalidPath
}

type Git struct{}

func (Git) Resolve(repo string) (io.ReadCloser, error) {
	return nil, ErrInvalidPath
}
