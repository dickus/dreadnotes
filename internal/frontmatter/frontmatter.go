package frontmatter

import (
	"time"

	"gopkg.in/yaml.v3"
)

const HumanTimeLayout = "2006-01-02 15:04"

type CustomTime struct {
	time.Time
}

func (ct *CustomTime) UnmarshalYAML(value *yaml.Node) error {
	t, err := time.Parse(HumanTimeLayout, value.Value)
	if err == nil {
		ct.Time = t

		return nil
	}

	t, err = time.Parse(time.RFC3339, value.Value)
	if err == nil {
		ct.Time = t

		return nil
	}

	t, err = time.Parse("2006-01-02", value.Value)
	if err == nil {
		ct.Time = t

		return nil
	}

	return err
}

// Frontmatter represents the metadata parsed from the YAML header of a note.
type Frontmatter struct {
	Title   string     `yaml:"title"`
	Created CustomTime `yaml:"created"`
	Updated CustomTime `yaml:"updated"`
	Tags    []string   `yaml:"tags"`
}

// Document represents a fully parsed Markdown file, including its metadata, body content, and file path.
type Document struct {
	Meta    Frontmatter
	Content []byte
	Path    string
}
