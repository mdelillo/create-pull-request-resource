package utils

type Source struct {
	GithubToken      string `json:"github_token"`
	RemoteRepository string `json:"remote_repository"`
	ForkedRepository string `json:"forked_repository"`
}
type Params struct {
	RepoLocation string `json:"repo_location,omitempty"`
	Base         string `json:"base,omitempty"`
	Description  string `json:"description,omitempty"`
	Title        string `json:"title,omitempty"`
	BranchPrefix string `json:"branch_prefix,omitempty"`
	AutoMerge    bool   `json:"auto_merge,omitempty"`
}

type OutRequest struct {
	Source Source `json:"source"`
	Params Params `json:"params"`
}

type InRequest struct {
	Version Version `json:"version"`
	Source  Source  `json:"source"`
	Params  Params  `json:"params"`
}
type Version struct {
	Ref string `json:"ref"`
}

type PrRespeonse struct {
	Number int  `json:"number"`
	Head   Head `json:"head"`
}
type Head struct {
	SHA string `json:"sha"`
}
