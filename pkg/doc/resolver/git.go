package resolver

import (
	"fmt"
	"io"
	"strings"
)

const (
	publicBitbucketTemplate = "https://bitbucket.org/%s/%s/raw/master/%s"
	publicGithubTemplate    = "https://raw.githubusercontent.com/%s/%s/master/%s"
	publicGitlabTemplate    = "https://gitlab.com/%s/%s/raw/master/%s"
)

var readmeFileNames = []string{
	"README.md", "readme.md",
	"Readme.md", "README",
}

// Git can be used to resolve a README file from its Git repository.
//
// Git supports the following repository types and formats:
//	Bitbucket: bitbucket.org/user/repo
//		ex. bitbucket.org/atlassian/aui
// 	Github: github.com/user/repo
//		ex. github.com/KyleBanks/kurz
//  Gitlab: gitlab.com/user/repo
//		ex. gitlab.com/openpowerlifting/opl-data
type Git struct{}

// Resolve attempts to find and load a remote README file from a git repository.
func (Git) Resolve(path string) (io.ReadCloser, error) {
	components := strings.Split(path, "/")
	if len(components) != 3 || len(components[1]) == 0 || len(components[2]) == 0 {
		return nil, ErrInvalidPath
	}

	var template, user, repo string
	switch strings.ToLower(components[0]) {
	case "bitbucket.org":
		template = publicBitbucketTemplate
		user = components[1]
		repo = components[2]
	case "github.com":
		template = publicGithubTemplate
		user = components[1]
		repo = components[2]
	case "gitlab.com":
		template = publicGitlabTemplate
		user = components[1]
		repo = components[2]
	default:
		return nil, ErrInvalidPath
	}

	var resolver URL
	for _, r := range readmeFileNames {
		content, err := resolver.Resolve(fmt.Sprintf(template, user, repo, r))
		if err != nil && err != ErrInvalidPath {
			return nil, err
		} else if content != nil {
			return content, nil
		}
	}

	return nil, ErrInvalidPath
}
