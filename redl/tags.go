package redl

import (
	"strings"

	"github.com/gomarkdown/markdown"
	"github.com/gomarkdown/markdown/parser"
)

func ExtractTags(title string) ([]string, string) {
	tags := []string{}

	title = strings.TrimSpace(title)

	for strings.HasPrefix(title, "[") {
		end := strings.Index(title, "]")
		if end < 0 {
			break
		}
		tag := title[1:end]
		title = strings.TrimSpace(title[end+1:])

		if len(tag) > 0 {
			tags = append(tags, strings.TrimSpace(tag))
		}
	}
	for strings.HasSuffix(title, "]") {
		beginning := strings.LastIndex(title, "[")
		if beginning < 0 {
			break
		}
		tag := title[beginning+1 : len(title)-1]
		title = strings.TrimSpace(title[:beginning])

		if len(tag) > 0 {
			tags = append(tags, strings.TrimSpace(tag))
		}
	}

	return tags, title
}

func FormatPost(body string) string {
	extensions := parser.Tables | parser.Autolink | parser.Strikethrough | parser.SpaceHeadings | parser.HeadingIDs | parser.SuperSubscript
	parser := parser.NewWithExtensions(extensions)

	md := []byte(body)
	html := string(markdown.ToHTML(md, parser, nil))
	// After formatting, fix some common errors
	return strings.ReplaceAll(strings.ReplaceAll(html, "<hr>", "<hr/>"), "<br>", "<br/>")
}
