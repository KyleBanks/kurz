package resolver

import (
	"bytes"
	"io/ioutil"
	"net/http"
	"testing"
)

func TestGit_Resolve(t *testing.T) {
	oldDefaultHttp := DefaultHttpGetter
	defer func() {
		DefaultHttpGetter = oldDefaultHttp
	}()

	tests := []struct {
		repo      string
		expectURL string
	}{
		{"github.com/KyleBanks/kurz", "https://raw.githubusercontent.com/KyleBanks/kurz/master/README.md"},
		{"bitbucket.org/atlassian/aui", "https://bitbucket.org/atlassian/aui/raw/master/README.md"},
		{"gitlab.com/openpowerlifting/opl-data", "https://gitlab.com/openpowerlifting/opl-data/raw/master/README.md"},
	}

	for idx, tt := range tests {
		expectContent := "GITHUB MARKDOWN"

		DefaultHttpGetter = &mockHttpGetter{
			getFn: func(url string) (*http.Response, error) {
				if url != tt.expectURL {
					t.Errorf("[%d] Unexpected url, expected=%v, got=%v", idx, tt.expectURL, url)
				}

				return &http.Response{
					StatusCode: 200,
					Body:       ioutil.NopCloser(bytes.NewBufferString(expectContent)),
				}, nil
			},
		}

		var g Git
		res, err := g.Resolve(tt.repo)
		if err != nil {
			t.Fatal(err)
		}

		content, _ := ioutil.ReadAll(res)
		if string(content) != expectContent {
			t.Errorf("[%d] Unexpected content, expected=%v, got=%s", idx, expectContent, content)
		}
	}
}

func TestGit_Resolve_invalidPath(t *testing.T) {
	tests := []string{
		"google.com",
		"google.com/path/to/file",
		"",
		"/path/to/file",
		"file.md",
		"github.com",
		"github.com/user",
		"github.com/user/",
		"github.com//repo",
		"http://github.com/user/repo",
		"bitbucket.org",
		"bitbucket.org/user",
		"bitbucket.org/user/",
		"bitbucket.org//repo",
		"http://bitbucket.org/user/repo",
		"gitlab.com",
		"gitlab.com/user",
		"gitlab.com/user/",
		"gitlab.com//repo",
		"http://gitlab.com/user/repo",
	}

	for idx, host := range tests {
		var g Git
		if _, err := g.Resolve(host); err != ErrInvalidPath {
			t.Errorf("[%d] Unexpected err, expected=%v, got=%v", idx, ErrInvalidPath, err)
		}
	}
}
