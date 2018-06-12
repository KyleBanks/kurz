package doc

import (
	"bytes"
	"errors"
	"io"
	"io/ioutil"
	"reflect"
	"testing"
)

type mockResolver struct {
	resolveFn func(string) (io.ReadCloser, error)
}

func (m mockResolver) Resolve(path string) (io.ReadCloser, error) {
	return m.resolveFn(path)
}

type mockParser struct {
	parseFn func(io.Reader) (Document, error)
}

func (m mockParser) Parse(r io.Reader) (Document, error) {
	return m.parseFn(r)
}

func TestNewDocument(t *testing.T) {
	expectPath := "/path/to/file"
	expectContent := "CONTENT"
	var expectDoc Document

	r := mockResolver{
		resolveFn: func(path string) (io.ReadCloser, error) {
			if path != expectPath {
				t.Fatalf("Unexpected path, expected=%v, got=%v", expectPath, path)
			}

			b := bytes.NewBufferString(expectContent)
			return ioutil.NopCloser(b), nil
		},
	}
	p := mockParser{
		parseFn: func(r io.Reader) (Document, error) {
			got, err := ioutil.ReadAll(r)
			if err != nil {
				t.Fatal(err)
			}

			if string(got) != expectContent {
				t.Fatalf("Unexpected content, expected=%v, got=%s", expectContent, got)
			}

			return expectDoc, nil
		},
	}

	d, err := NewDocument(expectPath, r, p)
	if err != nil {
		t.Fatal(err)
	}

	if !reflect.DeepEqual(d, expectDoc) {
		t.Fatalf("Unexpected document, expected=%v, got=%v", expectDoc, d)
	}
}

func TestNewDocument_errors(t *testing.T) {
	// Resolver Error
	{
		expectErr := errors.New("resolver error")
		r := mockResolver{
			resolveFn: func(path string) (io.ReadCloser, error) {
				return nil, expectErr
			},
		}

		if _, err := NewDocument("", r, nil); err != expectErr {
			t.Fatalf("Unexpected error, expected=%v, got=%v", expectErr, err)
		}
	}

	// Parser Error
	{
		expectErr := errors.New("parser error")
		r := mockResolver{
			resolveFn: func(path string) (io.ReadCloser, error) {
				return ioutil.NopCloser(&bytes.Buffer{}), nil
			},
		}
		p := mockParser{
			parseFn: func(r io.Reader) (Document, error) {
				return Document{}, expectErr
			},
		}

		if _, err := NewDocument("", r, p); err != expectErr {
			t.Fatalf("Unexpected error, expected=%v, got=%v", expectErr, err)
		}
	}
}

func TestNopStyler_Style(t *testing.T) {
	var ns NopStyler
	styles := []Style{Bold, Italic, Underline, Code}
	for _, s := range styles {
		got := ns.Style("string", s)
		if got != "string" {
			t.Errorf("Unexpected output for style %v, expected=string, got=%v", s, got)
		}
	}

}
