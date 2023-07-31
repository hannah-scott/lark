package main

import (
	"io/ioutil"
	"os"
	"strings"
	"text/template"
)

func main() {
	// Read in reference file
	content, err := ioutil.ReadFile("docs/index.lark")
	if err != nil {
		panic(err)
	}
	text := string(content)
	lines := strings.Split(text, "\n")
	lark := encodeLark(lines)

	f, err := os.Create("docs/index.html")
	if err != nil {
		panic(err)
	}
	defer f.Close()

	// Execute the template
	tmpl, err := template.New("lark.tmpl").ParseFiles("lark.tmpl")
	if err != nil {
		panic(err)
	}
	err = tmpl.Execute(f, lark)
	if err != nil {
		panic(err)
	}
	f.Sync()
}
