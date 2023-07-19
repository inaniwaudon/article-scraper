package main

import (
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
)

type article struct {
	title string
	tags  []string
	body  []element
}

type element struct {
	ttype    string
	content  string
	language string
	src      string
	caption  string
	list     []string
	elements []element
	table    [][]string
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
				if e.src != "" {
					t += ":" + e.src
				}
			}
			t += "\n" + e.content + "\n```"
		case "ol":
			l := []string{}
			for i, item := range e.list {
				l = append(l, strconv.Itoa(i+1)+". "+item)
			}
			t = strings.Join(l, "\n")
		case "ul":
			l := []string{}
			for _, item := range e.list {
				l = append(l, "- "+item)
			}
			t = strings.Join(l, "\n")
		case "table":
			l := []string{}
			for i, row := range e.table {
				l = append(l, "| "+strings.Join(row, " | ")+" |")
				if i == 0 {
					d := []string{}
					for _ = range e.table[0] {
						d = append(d, "---")
					}
					l = append(l, "| "+strings.Join(d, " | ")+" |")
				}
			}
			t = strings.Join(l, "\n")
		case "img":
			t = "![" + e.content + "](" + e.src + ")"
			if e.caption != "" {
				t += "\n*" + e.caption + "*"
			}
		case "details":
			t = ":::details " + e.content + "\n" + elementsToMd(e.elements) + "\n:::"
			break
		case "footnote":
			l := []string{}
			for i, item := range e.list {
				l = append(l, "[^"+strconv.Itoa(i+1)+"]: "+item)
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
	if len(os.Args) < 3 {
		log.Fatal("Not enough arguments.")
	}
	t := os.Args[1]
	id := os.Args[2]

	if t == "zenn" {
		article := scrapeZenn(id)
		mdtext := elementsToMd(article.body)
		fmt.Println(mdtext)
	}
}
