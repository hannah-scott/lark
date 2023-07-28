package main

import (
	"bufio"
	"fmt"
	"html"
	"os"
	"strings"
	"unicode"
)

type Lark struct {
	articles []Article
}

type Article struct {
	sections []Section
}

type Section struct {
	blocks []Block
}

type Block struct {
	glyph    string
	contents []string
}

func getGlyph(line string) string {
	// Handle empty lines
	if len(line) == 0 {
		return "paragraph"
	}

	// Handle dividers
	if len(line) > 2 {
		switch line[0:3] {
		case "```":
			return "pre"
		case "'''":
			return "pre"
		case "***":
			return "article"
		case "---":
			return "section"
		}

	}

	char := line[0]
	glyphs := map[byte]string{
		'=': "header",
		'-': "subheader",
		'_': "subsubheader",
		'+': "date",
		'~': "author",
		'@': "link",
		'!': "image",
		'>': "blockquote",
		'*': "ulist",
		':': "olist",
	}
	for k, v := range glyphs {
		if char == k {
			return v
		}
	}
	return "paragraph"
}

func parseLine(line string, glyph string) Block {
	// Get the glyph and store it in a block
	var content string
	block := Block{glyph: glyph}

	if glyph != "paragraph" {
		content = line[1:]
	} else {
		content = line
	}

	block.contents = []string{strings.TrimLeftFunc(content, unicode.IsSpace)}
	return block
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

func encodeLark(lines []string) Lark {
	lark := Lark{}
	article := Article{}
	section := Section{}
	block := Block{}
	preblock := Block{glyph: "pre"}
	ublock := Block{glyph: "ulist"}
	oblock := Block{glyph: "olist"}

	pre := false

	for _, line := range lines {
		glyph := getGlyph(line)

		pre = isPre(line, pre)
		if pre {
			if glyph != "pre" {
				preblock.contents = append(preblock.contents, line)
			}
		} else {
			if glyph == "article" {
				if section.blocks != nil {
					article.sections = append(article.sections, section)
					section.blocks = nil
				}
				if article.sections != nil {
					lark.articles = append(lark.articles, article)
					article.sections = nil
				}
			} else if glyph == "section" {
				if section.blocks != nil {
					article.sections = append(article.sections, section)
					section.blocks = nil
				}
			} else {
				// Handle pre block printing now
				if preblock.contents != nil {
					section.blocks = append(section.blocks, preblock)
					preblock.contents = nil
				}
				if glyph != "pre" {
					if glyph == "ulist" {
						ublock.contents = append(ublock.contents, parseLine(line, glyph).contents[0])
					} else if glyph == "olist" {
						oblock.contents = append(oblock.contents, parseLine(line, glyph).contents[0])
					} else {
						if ublock.contents != nil {
							section.blocks = append(section.blocks, ublock)
							ublock.contents = nil
						}

						if oblock.contents != nil {
							section.blocks = append(section.blocks, oblock)
							oblock.contents = nil
						}

						block = parseLine(line, glyph)
						section.blocks = append(section.blocks, block)
					}
				}
			}
		}
	}

	if section.blocks != nil {
		article.sections = append(article.sections, section)
	}
	if article.sections != nil {
		lark.articles = append(lark.articles, article)
	}

	return lark
}

func getHTMLTags(glyph string) string {
	tags := map[string]string{
		"header":       "h1",
		"subheader":    "h2",
		"subsubheader": "h3",
		"date":         "h2 class='date'",
		"author":       "h2 class='author'",
		"link":         "a",
		"image":        "a",
		"blockquote":   "blockquote",
		"ulist":        "ul",
		"olist":        "ol",
		"pre":          "pre",
	}
	for k, v := range tags {
		if glyph == k {
			return v
		}
	}
	return "p"
}

func parseLarkToHTML(lark Lark) string {
	output := ""

	for _, article := range lark.articles {
		output += "<article>\n"

		for _, section := range article.sections {
			if len(article.sections) > 1 {
				output += "<section>\n"
			}

			for _, block := range section.blocks {
				// Handle images and links
				if block.glyph == "link" {
					sep := strings.SplitN(block.contents[0], " ", 2)
					link := sep[0]
					desc := sep[0]
					if len(sep) == 2 {
						desc = sep[1]
					}
					output += "<p><a href='" + link + "'>" + desc + "</a></p>\n"
				} else if block.glyph == "image" {
					sep := strings.SplitN(block.contents[0], " ", 2)
					link := sep[0]
					desc := sep[0]
					if len(sep) == 2 {
						desc = sep[1]
					}

					output += "<img src='" + link + "' alt='" + desc + "' />\n"
				} else {
					output += "<" + getHTMLTags(block.glyph) + ">"

					for _, content := range block.contents {
						if block.glyph == "pre" {
							output += html.EscapeString(content) + "\n"
						} else if block.glyph == "ulist" || block.glyph == "olist" {
							output += "\n<li>" + content + "</li>\n"
						} else {
							if content != "" {
								output += "\n" + content + "\n"
							}
						}
					}
					output += "</" + strings.Split(getHTMLTags(block.glyph), " ")[0] + ">\n"
				}

			}

			if len(article.sections) > 1 {
				output += "</section>\n"
			}
		}

		output += "</article>\n"
	}

	return output
}

func main() {
	// Open file and read in line by line
	var lines []string

	file, err := os.Open("tgt.lark")
	if err != nil {
		panic(err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}

	if err := scanner.Err(); err != nil {
		panic(err)
	}

	html := parseLarkToHTML(encodeLark(lines))

	fmt.Println("<html><head><link rel='stylesheet' href='style.css'></head><body>")
	fmt.Println(html)
	fmt.Println("</body></html>")
}
