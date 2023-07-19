package main

import (
	"log"
	"net/http"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

func analyzeInline(s *goquery.Selection) string {
	text := ""

	s.Contents().Each(func(_ int, cs *goquery.Selection) {
		n := goquery.NodeName(cs)
		if n == "#text" {
			text += strings.Trim(cs.Text(), "\n")
		} else {
			c := analyzeInline(cs)
			if n == "a" {
				h := cs.AttrOr("href", "")
				text += "[" + c + "](" + h + ")"
			} else if n == "strong" {
				c := analyzeInline(cs)
				text += "**" + c + "**"
			} else if n == "code" {
				text += "`" + c + "`"
			} else if n == "s" {
				text += "~~" + c + "~~"
			} else if cs.AttrOr("class", "") == "footnote-ref" {
				t := cs.Text()
				text += "[^" + t[1:len(t)-1] + "]"
			}
		}
	})
	return text
}

func analyzeBlock(s *goquery.Selection) *element {
	class := s.AttrOr("class", "")
	name := goquery.NodeName(s)

	if class == "code-block-container" {
		// code
		p := s.Find("pre")
		c := strings.TrimSpace(p.Text())
		class, exists := p.Attr("class")
		l := ""
		if exists && strings.HasPrefix(class, "language-") {
			l = strings.Replace(class, "language-", "", 1)
		}
		n := s.Find(".code-block-filename").Text()
		return &element{ttype: "code", content: c, language: l, src: n}
	} else if name == "ul" || name == "ol" {
		// list
		l := []string{}
		s.Find("li").Each(func(_ int, cs *goquery.Selection) {
			l = append(l, analyzeInline(cs))
		})
		return &element{ttype: name, list: l}
	} else if name == "table" {
		// table
		t := [][]string{}
		s.Find("tr").Each(func(_ int, cs *goquery.Selection) {
			row := []string{}
			cs.Find("th, td").Each(func(_ int, ds *goquery.Selection) {
				row = append(row, analyzeInline(ds))
			})
			t = append(t, row)
		})
		return &element{ttype: "table", table: t}
	} else if class == "footnotes" {
		// footnote
		l := []string{}
		s.Find("li p").Each(func(_ int, cs *goquery.Selection) {
			l = append(l, strings.TrimSpace(cs.Contents().First().Text()))
		})
		return &element{ttype: "footnote", list: l}
	} else if name == "details" {
		// accordion
		c := strings.TrimSpace(s.Find("summary").Text())
		l := []element{}
		s.Find(".details-content").Contents().Each(func(_ int, cs *goquery.Selection) {
			l = append(l, *analyzeBlock(cs))
		})
		return &element{ttype: "details", content: c, elements: l}
	} else {
		// others
		if name == "p" {
			children := s.Children()
			if goquery.NodeName(children.First()) == "img" {
				// img
				c := children.First().AttrOr("alt", "")
				src := children.First().AttrOr("src", "")
				caption := strings.TrimSpace(s.Find("em").Text())
				return &element{ttype: "img", content: c, src: src, caption: caption}
			} else {
				// paragraph
				return &element{ttype: "p", content: analyzeInline(s)}
			}
		} else {
			return &element{ttype: name, content: strings.TrimSpace(s.Text())}
		}
	}
}

func scrapeZenn(id string) article {
	res, err := http.Get("https://zenn.dev/inaniwaudon/articles/" + id)
	if err != nil {
		log.Fatal(err)
	}
	defer res.Body.Close()
	if res.StatusCode != 200 {
		log.Fatalf("status code error: %d %s", res.StatusCode, res.Status)
	}

	doc, err := goquery.NewDocumentFromReader(res.Body)
	if err != nil {
		log.Fatal(err)
	}
	title := doc.Find("h1").First().Text()
	tags := []string{}

	elements := []element{}
	body := doc.Find(".BodyContent_anchorToHeadings__Vl0_u")

	body.First().Children().Each(func(_ int, s *goquery.Selection) {
		elements = append(elements, *analyzeBlock(s))
	})
	return article{title, tags, elements}
}
