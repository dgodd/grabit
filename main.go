package main

import (
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"log"
	"net/http"
	"regexp"
	"strings"
)

func main() {
	resp, err := http.Get("https://grabit.clubwyndhamsp.com/page/7/Packages")
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()

	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		log.Fatal(err)
	}

	doc.Find(".listing").Each(func(i int, s *goquery.Selection) {
		body := ""
		s.Find(".listing-header,h4,p").Each(func(_ int, s *goquery.Selection) {
			// TODO: Ignore style="display:none;"
			if style, exists := s.Attr("style"); exists {
				if style == "display:none;" {
					return
				}
			}
			body += s.Text() + "\n"
		})
		body = strings.TrimSpace(body)
		body = regexp.MustCompile(`[\n\r\s]*\n+[\n\r\s]*`).ReplaceAllString(body, "\n")

		// body := s.Text()
		fmt.Println(body)
		fmt.Println("========")
	})
}
