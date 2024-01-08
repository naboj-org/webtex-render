package main

import (
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"github.com/naboj-org/webtex-render/webtex_api"
	"log"
	"net/url"
	"path"
	"strings"
)

func parseHtml(config Config) error {
	doc, err := goquery.NewDocumentFromReader(config.InputFile)
	if err != nil {
		log.Println("parsing input file failed:", err)
		return err
	}

	var toProcess []*goquery.Selection

	doc.Find("img").Each(func(i int, s *goquery.Selection) {
		toProcess = append(toProcess, s)
	})

	for _, s := range toProcess {
		src, _ := s.Attr("src")
		if !strings.HasPrefix(src, config.InputURL) {
			continue
		}

		src = strings.TrimPrefix(src, config.InputURL)
		src, err := url.QueryUnescape(src)
		if err != nil {
			return fmt.Errorf("failed to unescape query string: %w", err)
		}

		filepath, err := webtex_api.EquationSvg(src, config.EquationDirectory, config.Template, config.Engine)
		if err != nil {
			return fmt.Errorf("failed generation of '%s': %w", src, err)
		}
		filepath = path.Join(config.OutputURL, filepath)
		s.SetAttr("src", filepath)
	}

	var toExport *goquery.Selection
	if config.OnlyInnerHTML {
		toExport = doc.Find("body")
	} else {
		toExport = doc.Selection
	}
	html, err := toExport.Html()
	if err != nil {
		log.Println("failed generating output html:", err)
		return err
	}

	_, err = config.OutputFile.WriteString(html)
	return err
}
