package textmining

import (
	"log"
	"strings"
	"io/ioutil"

	//external library
	"github.com/PuerkitoBio/goquery"
)

func websitesList() []string {
	content, err := ioutil.ReadFile("websites.txt")
	if err != nil {
		log.Println(err)
	}
	websites := string(content)
	websitesList := strings.Split(websites, "\n")
	return websitesList
}

func TextMining() string {
	result := ""
	websites := websitesList()
	for _, v := range websites {
		doc, err := goquery.NewDocument(v)
		if err != nil {
			log.Println(err)
		}
		text := doc.Find(".ui_qtext_expanded").Text()
		result = result + " " + text
	}
	return result
}
