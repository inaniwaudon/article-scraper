package main

import (
	"log"
	"net/http"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

func analyzeInline(s *goquery.Selection) string {
	n := goquery.NodeName(s)
	if n == "#text" {
		return strings.Trim(s.Text(), "\n")
	} else {
		c := analyzeLoopInline(s)
		if n == "a" {
			h := s.AttrOr("href", "")
			return "[" + c + "](" + h + ")"
		} else if n == "strong" {
			c := analyzeLoopInline(s)
			return "**" + c + "**"
		} else if n == "code" {
			return "`" + c + "`"
		} else if n == "s" {
			return "~~" + c + "~~"
		} else if s.AttrOr("class", "") == "footnote-ref" {
			t := s.Text()
			return "[^" + t[1:len(t)-1] + "]"
		}
	}
	return ""
}

func analyzeLoopInline(s *goquery.Selection) string {
	text := ""
	s.Contents().Each(func(_ int, cs *goquery.Selection) {
		text += analyzeInline(cs)
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
			l = append(l, analyzeLoopInline(cs))
		})
		return &element{ttype: name, list: l}
	} else if name == "table" {
		// table
		t := [][]string{}
		s.Find("tr").Each(func(_ int, cs *goquery.Selection) {
			row := []string{}
			cs.Find("th, td").Each(func(_ int, ds *goquery.Selection) {
				row = append(row, analyzeLoopInline(ds))
			})
			t = append(t, row)
		})
		return &element{ttype: "table", table: t}
	} else if class == "footnotes" {
		// footnote
		l := []string{}
		s.Find("li p").Each(func(_ int, cs *goquery.Selection) {
			t := ""
			cs.Contents().Each(func(i int, ds *goquery.Selection) {
				if ds.AttrOr("class", "") != "footnote-backref" {
					t += analyzeInline(ds)
				}
			})
			l = append(l, t)
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
				return &element{ttype: "p", content: analyzeLoopInline(s)}
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
