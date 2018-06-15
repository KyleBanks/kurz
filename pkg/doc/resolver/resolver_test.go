package resolver

import (
	"bytes"
	"errors"
	"io"
	"io/ioutil"
	"net/http"
	"testing"

	"github.com/KyleBanks/kurz/pkg/doc"
)

type mockResolver struct {
	resolveFn func(string) (io.ReadCloser, error)
}

func (m *mockResolver) Resolve(p string) (io.ReadCloser, error) {
	return m.resolveFn(p)
}

type mockHttpGetter struct {
	getFn func(string) (*http.Response, error)
}

func (m *mockHttpGetter) Get(url string) (*http.Response, error) {
	return m.getFn(url)
}

func TestChain_Resolve(t *testing.T) {
	// Happy path, middle resolver resolves successfully
	expectRes := "MARKDOWN\nFILE"
	expectPath := "/path/to/file"

	c := Chain{
		Resolvers: []doc.Resolver{
			&mockResolver{
				resolveFn: func(p string) (io.ReadCloser, error) {
					if p != expectPath {
						t.Errorf("Unexpected path, expected=%v, got=%v", expectPath, p)
					}
					return nil, ErrInvalidPath
				},
			},
			&mockResolver{
				resolveFn: func(p string) (io.ReadCloser, error) {
					if p != expectPath {
						t.Errorf("Unexpected path, expected=%v, got=%v", expectPath, p)
					}
					return ioutil.NopCloser(bytes.NewBufferString(expectRes)), nil
				},
			},
			&mockResolver{
				resolveFn: func(p string) (io.ReadCloser, error) {
					t.Fatal("Final resolver should not have been invoked.")
					return nil, nil
				},
			},
		},
	}

	rc, err := c.Resolve(expectPath)
	if err != nil {
		t.Fatal(err)
	}

	res, _ := ioutil.ReadAll(rc)
	if string(res) != expectRes {
		t.Errorf("Unexpected response, expected=%v, got=%s", expectRes, res)
	}
}

func TestChain_Resolve_wholeChain(t *testing.T) {
	// Runs through the whole chain, never resolving
	expectErr := ErrInvalidPath
	expectSequence := []int{1, 2, 3}

	var sequence []int
	c := Chain{
		Resolvers: []doc.Resolver{
			&mockResolver{
				resolveFn: func(p string) (io.ReadCloser, error) {
					sequence = append(sequence, 1)
					return nil, ErrInvalidPath
				},
			},
			&mockResolver{
				resolveFn: func(p string) (io.ReadCloser, error) {
					sequence = append(sequence, 2)
					return nil, ErrInvalidPath
				},
			},
			&mockResolver{
				resolveFn: func(p string) (io.ReadCloser, error) {
					sequence = append(sequence, 3)
					return nil, ErrInvalidPath
				},
			},
		},
	}

	_, err := c.Resolve("path")
	if err != expectErr {
		t.Errorf("Unexpected error, expected=%v, got=%v", expectErr, err)
	}

	if len(sequence) != len(expectSequence) {
		t.Errorf("Unexpected number of resolvers invoked, expected=%v, got=%v", len(expectSequence), len(sequence))
	}

	for i := 0; i < len(sequence); i++ {
		if sequence[i] != expectSequence[i] {
			t.Errorf("Unexpected invocation at index %d, expected=%v, got=%v", i, expectSequence[i], sequence[i])
		}
	}
}

func TestURL_Resolve(t *testing.T) {
	var m mockHttpGetter
	u := URL{
		HttpGetter: &m,
	}

	// Happy Path
	{
		expectRes := "MARKDOWN\nFILE\nCONTENTS"
		expectUrl := "http://example.com/FILE.md"
		m.getFn = func(url string) (*http.Response, error) {
			if url != expectUrl {
				t.Errorf("Unexpected URL, expected=%v, got=%v", expectUrl, url)
			}

			return &http.Response{
				StatusCode: 200,
				Body:       ioutil.NopCloser(bytes.NewBufferString(expectRes)),
			}, nil
		}

		rc, err := u.Resolve(expectUrl)
		if err != nil {
			t.Fatal(err)
		}

		res, _ := ioutil.ReadAll(rc)
		if string(res) != expectRes {
			t.Errorf("Unexpected response, expected=%v, got=%s", expectRes, res)
		}
	}

	// Bad status code
	{
		expectErr := ErrInvalidPath
		m.getFn = func(url string) (*http.Response, error) {
			return &http.Response{
				StatusCode: 400,
			}, nil
		}

		_, err := u.Resolve("URL")
		if err != expectErr {
			t.Errorf("Unexpected error for bad status code, expected=%v, got=%v", expectErr, err)
		}
	}

	// HTTP error
	{
		expectErr := errors.New("sample error")
		m.getFn = func(url string) (*http.Response, error) {
			return nil, expectErr
		}

		_, err := u.Resolve("URL")
		if err != expectErr {
			t.Errorf("Unexpected error for http error, expected=%v, got=%v", expectErr, err)
		}
	}
}
