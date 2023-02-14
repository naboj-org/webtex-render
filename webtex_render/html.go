package main

import (
	"github.com/PuerkitoBio/goquery"
	"log"
	"net/url"
	"strings"
)

func parseHtml(config Config) error {
	doc, err := goquery.NewDocumentFromReader(config.InputFile)
	if err != nil {
		log.Println("parsing input file failed:", err)
		return err
	}

	doc.Find("img").Each(func(i int, s *goquery.Selection) {
		src, _ := s.Attr("src")
		if !strings.HasPrefix(src, config.InputURL) {
			return
		}

		src = strings.TrimPrefix(src, config.InputURL)
		src, err := url.QueryUnescape(src)
		if err != nil {
			log.Println("failed query unescaping:", err)
			return
		}

		sha, err := EquationSvg(config, src)
		s.SetAttr("src", sha)
	})

	html, err := doc.Html()
	if err != nil {
		log.Println("failed generating output html:", err)
		return err
	}

	_, err = config.OutputFile.WriteString(html)
	return err
}
