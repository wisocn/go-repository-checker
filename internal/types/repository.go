package types

import "github.com/google/go-github/github"

type Repository struct {
	url   string
	name  string
	files []*github.RepositoryContent
}
