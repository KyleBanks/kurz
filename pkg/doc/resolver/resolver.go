package resolver

import (
	"errors"
	"io"
	"net/http"
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

	readmeFileNames = []string{
		"README.md", "readme.md", "Readme.md",
		"README",
	}
)

// HttpGetter defines a type that can send GET requests
// over HTTP.
type HttpGetter interface {
	Get(string) (*http.Response, error)
}

// DefaultHttpGetter is the default HTTP implementation
// used when performing HTTP requests.
var DefaultHttpGetter HttpGetter = http.DefaultClient

// Chain runs a sequence of doc.Resolver types, returning the first
// content to be successfully resolved.
type Chain struct {
	Resolvers []doc.Resolver
}

// Resolve attempts to resolve the provided path using all of the
// underlying resolvers in the chain.
//
// If any of the resolvers returns an error other than ErrInvalidPath,
// the error will be returned. If none of the resolvers are able to
// resolve a content body, an ErrInvalidPath error is returned.
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

// File can be used to resolve a local file by its path.
type File struct{}

// Resolve finds and loads a local file by its path.
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

// URL can be used to resolve a remote file by its URL.
type URL struct {
	// HttpGetter allows for a custom HTTP client implementation
	// to be used by the URL resolver. If this property is not
	// set, the DefaultHttpGetter will be used.
	HttpGetter HttpGetter
}

// Resolve finds and loads a remote file by its URL.
func (u URL) Resolve(url string) (io.ReadCloser, error) {
	var h HttpGetter = u.HttpGetter
	if h == nil {
		h = DefaultHttpGetter
	}

	resp, err := h.Get(url)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode < 200 || resp.StatusCode > 299 {
		return nil, ErrInvalidPath
	}

	return resp.Body, nil
}

// Git can be used to resolve a README file from its Git repository.
type Git struct{}

// Resolve attempts to find a load a remote README file from a git repository.
func (Git) Resolve(repo string) (io.ReadCloser, error) {
	return nil, ErrInvalidPath
}
