package main

import (
	"io"
	"io/ioutil"
	"os"
	"strings"
	"text/template"
)

type TemplateParser struct{}

func (p *TemplateParser) Parse(reader io.Reader, writer io.Writer) error {
	ctx := make(map[string]string)
	t, err := ioutil.ReadAll(os.Stdin)
	if err != nil {
		return err
	}
	tmpl, err := template.New("parser").Parse(string(t))
	if err != nil {
		return err
	}
	for _, e := range os.Environ() {
		chunks := strings.Split(e, "=")
		ctx[chunks[0]] = chunks[1]
	}
	if err := tmpl.Execute(os.Stdout, ctx); err != nil {
		return err
	}
	return nil
}

func main() {
	new(TemplateParser).Parse(os.Stdin, os.Stdout)
}
