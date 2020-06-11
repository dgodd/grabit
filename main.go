package main

import (
	"fmt"
	"github.com/PuerkitoBio/goquery"
	mailjet "github.com/mailjet/mailjet-apiv3-go"
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
	mailjetClient := mailjet.NewMailjetClient(os.Getenv("MJ_APIKEY_PUBLIC"), os.Getenv("MJ_APIKEY_PRIVATE"))
	messagesInfo := []mailjet.InfoMessagesV31{
		mailjet.InfoMessagesV31{
			From: &mailjet.RecipientV31{
				Email: "dave@godard.id.au",
				Name:  "Dave Goddard",
			},
			To: &mailjet.RecipientsV31{
				mailjet.RecipientV31{
					Email: "dave@godard.id.au",
					Name:  "Dave Goddard",
				},
				mailjet.RecipientV31{
					Email: "catherine@lypc.com.au",
					Name:  "Catherine Garro",
				},
			},
			Subject:  "Grabit - Worldmarksp.com",
			TextPart: text,
		},
	}
	messages := mailjet.MessagesV31{Info: messagesInfo}
	res, err := mailjetClient.SendMailV31(&messages)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Data: %+v\n", res)
}
