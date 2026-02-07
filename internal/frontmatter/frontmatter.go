package frontmatter

type Frontmatter struct {
	Title string
	Created string
	Updated string
	Tags []string
}

var Front Frontmatter

type Document struct {
	Meta Frontmatter
	Content []byte
	Path string
}

var Doc Document

