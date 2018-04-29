package main

import (
	"github.com/elissalim/wordcloudgo/textmining"
	"github.com/elissalim/wordcloudgo/textprocessing"
	"github.com/elissalim/wordcloudgo/wordcloud"
)

func main() {
	//extract web content from websites list
	content := textmining.TextMining()

	//use stop words list to do text processing
	processedText := textprocessing.SortedResult(content)

	//create word cloud with specified criteria
	wordcloud.WordCloud(processedText)
}
