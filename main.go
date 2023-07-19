package main

import (
	"fmt"
	"strconv"
	"strings"
)

type element struct {
	ttype    string
	content  string
	language string
	name     string
	list     []string
}

type article struct {
	title string
	tags  []string
	body  []element
}

func elementsToMd(elements []element) string {
	lines := []string{}
	for _, e := range elements {
		var t string
		switch e.ttype {
		case "h1":
			t = "# " + e.content
		case "h2":
			t = "## " + e.content
		case "h3":
			t = "### " + e.content
		case "h4":
			t = "#### " + e.content
		case "h5":
			t = "##### " + e.content
		case "h6":
			t = "###### " + e.content
		case "code":
			t = "```"
			if e.language != "" {
				t += e.language
				if e.name != "" {
					t += ":" + e.name
				}
			}
			t += "\n" + e.content + "\n```"
		case "ol":
			l := []string{}
			for i, item := range e.list {
				l = append(l, strconv.Itoa(i+1)+". "+item)
			}
			t = strings.Join(l, "\n")
		case "ui":
			l := []string{}
			for _, item := range e.list {
				l = append(l, "- "+item)
			}
			t = strings.Join(l, "\n")
		default:
			t = e.content
		}
		lines = append(lines, t)
	}
	return strings.Join(lines, "\n\n")
}

func main() {
	article := scrapeZenn("7fa50a744cb67a")
	mdtext := elementsToMd(article.body)
	fmt.Println(mdtext)
}
