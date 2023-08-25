package main

import (
	"fmt"
	"os"
	"strings"
	"text/template"
)

func (b *Block) EncodeMarkdown() string {
	glyphs := map[string]string{
		"header":       "#",
		"subheader":    "##",
		"subsubheader": "###",
		"date":         "",
		"author":       "",
		"blockquote":   ">",
		"ulist":        "-",
		"pre":          "```",
		"code":         "```",
	}

	if b.Glyph == "link" || b.Glyph == "image" {
		o := ""
		for _, c := range b.Contents {
			sep := strings.SplitN(c, " ", 2)
			link := sep[0]
			desc := sep[0]
			if len(sep) == 2 {
				desc = sep[1]
			}
			if b.Glyph == "image" {
				o += "!"
			}
			o += "[" + desc + "](" + link + ")\n"
		}
		return o
	}
	if b.Glyph == "pre" || b.Glyph == "code" {
		o := glyphs[b.Glyph] + "\n"
		for _, c := range b.Contents {
			o += c + "\n"
		}
		o += "```"
		return o
	}
	if b.Glyph == "olist" || b.Glyph == "ulist" {
		o := ""
		for i, c := range b.Contents {
			if b.Glyph == "olist" {
				o += fmt.Sprint(i+1, ". ")
			} else {
				o += glyphs[b.Glyph] + " "
			}
			o += c + "\n"
		}
		return o
	}

	o := glyphs[b.Glyph] + " "
	for _, c := range b.Contents {
		o += c
	}
	return o
}

func lark2md(lark Lark) error {
	f, err := os.Create("docs/index.md")
	if err != nil {
		return err
	}
	defer f.Close()

	// Execute the template
	tmpl, err := template.New("lark2md.tmpl").ParseFiles("lark2md.tmpl")
	if err != nil {
		panic(err)
	}
	err = tmpl.Execute(f, lark)
	if err != nil {
		return err
	}
	f.Sync()
	return nil
}
