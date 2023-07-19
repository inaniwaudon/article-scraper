package main

import (
	"log"
	"net/http"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

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
		class, exists := s.Attr("class")
		name := goquery.NodeName(s)
		if exists && class == "code-block-container" {
			// code
			p := s.Find("pre")
			c := strings.TrimSpace(p.Text())
			class, exists := p.Attr("class")
			l := ""
			if exists && strings.HasPrefix(class, "language-") {
				l = strings.Replace(class, "language-", "", 1)
			}
			n := s.Find(".code-block-filename").Text()
			e := element{ttype: "code", content: c, language: l, name: n}
			elements = append(elements, e)
		} else if name == "ul" || name == "ol" {
			// list
			l := []string{}
			s.Find("li").Each(func(_ int, cs *goquery.Selection) {
				l = append(l, cs.Text())
			})
			e := element{ttype: name, list: l}
			elements = append(elements, e)
		} else {
			// others
			c := strings.TrimSpace(s.Text())
			e := element{ttype: name, content: c}
			elements = append(elements, e)
		}
	})
	return article{title, tags, elements}
}
