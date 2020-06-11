package main

import (
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"github.com/keighl/postmark"
	"log"
	"net/http"
	"os"
	"regexp"
	"strings"
)

func main() {
	reClean := regexp.MustCompile(`[\n\r\s]*\n+[\n\r\s]*`)

	resp, err := http.Get("https://grabit.clubwyndhamsp.com/page/7/Packages")
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()

	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		log.Fatal(err)
	}

	fullText := ""
	doc.Find(".listing").Each(func(i int, s *goquery.Selection) {
		body := ""
		s.Find(".listing-header,h4,p").Each(func(_ int, s *goquery.Selection) {
			if style, exists := s.Attr("style"); exists {
				if style == "display:none;" {
					return
				}
			}
			body += s.Text() + "\n"
		})
		body = strings.TrimSpace(body)
		body = reClean.ReplaceAllString(body, "\n")

		fullText += body + "\n============\n"
	})

	fmt.Println(fullText)
	sendEmail(fullText)
}

func sendEmail(text string) {
	client := postmark.NewClient(os.Getenv("POSTMARK_SERVER_TOKEN"), os.Getenv("POSTMARK_ACCOUNT_TOKEN"))

	for _, emailAddress := range []string{"dave@goddard.id.au", "catherine@lypc.com.au"} {
		email := postmark.Email{
			From:       emailAddress,
			To:         emailAddress,
			Subject:    "Grabit - Worldmarksp.com",
			HtmlBody:   "<pre>" + text + "</pre>",
			TextBody:   text,
			Tag:        "grabit-changed",
			TrackOpens: true,
		}
		_, err := client.SendEmail(email)
		if err != nil {
			panic(err)
		}
	}
}
