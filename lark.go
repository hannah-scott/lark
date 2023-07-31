package main

import (
	"html"
	"regexp"
	"strings"
	"unicode"
)

type Lark struct {
	Articles []Article
}

type Article struct {
	Sections []Section
}

type Section struct {
	Blocks []Block
}

type Block struct {
	Glyph    string
	Contents []string
}

func getGlyph(line string) string {
	// Handle empty lines
	if len(line) == 0 {
		return "paragraph"
	}

	// Handle dividers
	if len(line) > 2 {
		switch line[0:3] {
		case "***":
			return "article"
		case "---":
			return "section"
		}

	}

	char := line[0]
	glyphs := map[byte]string{
		'=':  "header",
		'-':  "subheader",
		'_':  "subsubheader",
		'+':  "date",
		'~':  "author",
		'@':  "link",
		'!':  "image",
		'>':  "blockquote",
		'*':  "ulist",
		':':  "olist",
		'\'': "pre",
		'`':  "code",
	}
	for k, v := range glyphs {
		if char == k {
			return v
		}
	}
	return "paragraph"
}

func parseLine(line string, glyph string) string {
	// Get the glyph and store it in a block
	var content string

	if glyph != "paragraph" {
		content = line[1:]
	} else {
		content = line
	}

	// Personal opinion
	content = strings.Replace(content, "---", "—", -1)
	content = strings.Replace(content, "--", "—", -1)

	return strings.TrimLeftFunc(content, unicode.IsSpace)
}

func isPre(line string, pre bool) bool {
	if pre && getGlyph(line) != "pre" {
		return true
	}
	if !pre && getGlyph(line) == "pre" {
		return true
	}
	return false
}

func isCode(line string, code bool) bool {
	if code && getGlyph(line) != "code" {
		return true
	}
	if !code && getGlyph(line) == "code" {
		return true
	}
	return false
}

func encodeLark(lines []string) Lark {
	lark := Lark{}
	article := Article{}
	section := Section{}
	block := Block{}
	context := getGlyph(lines[0])

	// Read in a line
	for _, line := range lines {
		// Check the glyph to get context
		glyph := getGlyph(line)
		// If context has changed then write current block out and start a new one
		if glyph != context {
			block.Glyph = context
			if block.Contents != nil {
				section.Blocks = append(section.Blocks, block)
				block = Block{}
			}
			context = glyph
		}
		if glyph == "section" {
			// Write out current block and section and reinitialize them
			block.Glyph = glyph

			if block.Contents != nil {
				section.Blocks = append(section.Blocks, block)
			}
			if section.Blocks != nil {
				article.Sections = append(article.Sections, section)
			}

			block = Block{}
			section = Section{}
		} else if glyph == "article" {
			// Write out current block, section, article and reinitialize them
			block.Glyph = glyph

			if block.Contents != nil {
				section.Blocks = append(section.Blocks, block)
			}
			if section.Blocks != nil {
				article.Sections = append(article.Sections, section)
			}
			if article.Sections != nil {
				lark.Articles = append(lark.Articles, article)
			}

			block = Block{}
			section = Section{}
			article = Article{}
		} else {
			if glyph == "pre" || glyph == "code" {
				block.Contents = append(block.Contents, line[1:])
			} else {
				block.Contents = append(block.Contents, parseLine(line, glyph))
			}
		}

	}
	// If it's the last line then write everything out
	if block.Contents != nil {
		block.Glyph = getGlyph(lines[len(lines)-1])
		section.Blocks = append(section.Blocks, block)
	}
	if section.Blocks != nil {
		article.Sections = append(article.Sections, section)
	}
	if article.Sections != nil {
		lark.Articles = append(lark.Articles, article)
	}

	return lark
}

func (b *Block) GetHTMLTags() string {
	tags := map[string]string{
		"header":       "h1",
		"subheader":    "h2",
		"subsubheader": "h3",
		"date":         "h3",
		"author":       "h3",
		"link":         "a",
		"image":        "a",
		"blockquote":   "blockquote",
		"ulist":        "ul",
		"olist":        "ol",
		"pre":          "pre",
		"code":         "code",
	}
	for k, v := range tags {
		if b.Glyph == k {
			return v
		}
	}
	return "p"
}

func (b *Block) GetTaggedContent() []string {
	var o []string
	for _, c := range b.Contents {
		if strings.TrimSpace(c) != "" {
			s := "<" + b.GetHTMLTags() + ">\n"
			s += c + "\n"
			s += "</" + b.GetHTMLTags() + ">\n"
			o = append(o, s)
		}
	}
	return o
}

func (b *Block) IsLink() bool {
	return b.Glyph == "link"
}

func (b *Block) IsImage() bool {
	return b.Glyph == "image"
}

func (b *Block) IsList() bool {
	return b.Glyph == "ulist" || b.Glyph == "olist"
}

func (b *Block) IsPre() bool {
	return b.Glyph == "pre"
}

func (b *Block) IsCode() bool {
	return b.Glyph == "code"
}

func (b *Block) EncodeImage() string {
	o := ""
	for _, c := range b.Contents {
		sep := strings.SplitN(c, " ", 2)
		link := sep[0]
		desc := sep[0]
		if len(sep) == 2 {
			desc = sep[1]
		}
		o += "<img src='" + link + "' alt='" + desc + "' loading='lazy' />\n"
	}

	return o
}

func (b *Block) EncodeLink() string {
	o := ""
	for _, c := range b.Contents {
		sep := strings.SplitN(c, " ", 2)
		link := sep[0]
		desc := sep[0]
		if len(sep) == 2 {
			desc = sep[1]
		}
		if strings.HasPrefix(link, "http://") || strings.HasPrefix(link, "https://") {
			o += "<li><a href='" + link + "' class='external-link'>" + desc + "</a></li>\n"
		} else if strings.HasPrefix(link, "gemini://") {
			o += "<li><a href='" + link + "' class='gemini-link'>" + desc + "</a></li>\n"
		} else {
			o += "<li><a href='" + link + "'>" + desc + "</a></li>\n"
		}
	}
	return o
}

func (b *Block) EncodeList() string {
	output := "<" + b.GetHTMLTags() + ">\n"
	for _, content := range b.Contents {
		output += "\t<li>" + content + "</li>\n"
	}
	output += "</" + b.GetHTMLTags() + ">\n"

	return output
}

func (b *Block) EncodePre() string {
	output := "<" + b.GetHTMLTags() + ">\n"
	for _, content := range b.Contents {
		output += content + "\n"
	}
	output += "</" + b.GetHTMLTags() + ">\n"

	return output
}

func (b *Block) EncodeCode() string {
	output := "<pre><code>"
	for _, content := range b.Contents {
		output += content + "\n"
	}
	output += "</code></pre>\n"

	return output
}

func (b *Block) EscapeString(s string) string {
	return html.EscapeString(s)
}

func (a *Article) GetID() string {
	for _, s := range a.Sections {
		for _, b := range s.Blocks {
			if b.Glyph == "header" {
				str := b.Contents[0]
				str = regexp.MustCompile(`[^a-zA-Z0-9 \-_]+`).ReplaceAllString(str, "")
				str = strings.ToLower(strings.Replace(str, " ", "-", -1))
				return str
			}
			if b.Glyph == "subheader" {
				str := b.Contents[0]
				str = regexp.MustCompile(`[^a-zA-Z0-9 \-_]+`).ReplaceAllString(str, "")
				str = strings.ToLower(strings.Replace(str, " ", "-", -1))
				return str
			}
		}
	}

	return ""
}
