package tmpl

import (
	"bytes"
	"io"
	"os"
	"text/template/parse"
)

type Template struct {
	File      string
	Trees     map[string]*parse.Tree
	Inspector *Inspector

	r io.ReadSeeker
}

func Parse(filepath string) (*Template, error) {
	src, errF := os.ReadFile(filepath)
	if errF != nil {
		return nil, errF
	}

	return ParseBytes(filepath, src)
}

func ParseString(name, content string) (*Template, error) {
	return ParseBytes(name, []byte(content))
}

func ParseBytes(name string, src []byte) (*Template, error) {
	t := &Template{
		File:      name,
		r:         nil,
		Trees:     make(map[string]*parse.Tree),
		Inspector: nil,
	}

	t.r = bytes.NewReader(src)
	t.Trees = make(map[string]*parse.Tree)

	tree := &parse.Tree{
		Name: name,
		Mode: parse.ParseComments | parse.SkipFuncCheck,
	}

	_, err := tree.Parse(string(src), "{{", "}}", t.Trees, map[string]any{})
	if err != nil {
		return nil, err
	}

	roodNotes := make([]parse.Node, 0, len(t.Trees))
	for _, tree := range t.Trees {
		roodNotes = append(roodNotes, tree.Root)
	}

	t.Inspector = newInspector(roodNotes)

	return t, nil
}
