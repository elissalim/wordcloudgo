package textprocessing

import (
	"log"
	"strings"
	"regexp"
	"io/ioutil"
	"sort"

	//external library
	"github.com/aaaton/golem"
	"github.com/elissalim/wordcloudgo/textmining"
)

func stopWordsList() []string {
	content, err := ioutil.ReadFile("stopwords.txt")
	if err != nil {
		log.Println(err)
	}
	words := string(content)
	stopWords := strings.Split(words, "\n")
	return stopWords
}

func reviewsList(str string) []string {
	reviews := strings.Split(str, " ")
	return reviews
}

func lemmatizeWords(reviews []string) []string {
	lemmatizer, err := golem.New("english")
	if err != nil {
		log.Println(err)
	}
	reviewsLength := len(reviews)
	var result []string
	for i := 0; i < reviewsLength; i++ {
		word, err := lemmatizer.Lemma(reviews[i])
		if err != nil {
			result = append(result, reviews[i])
		}
		result = append(result, word)
	}
	return result
}

func wordCount(str []string) map[string]int {
	wordCountMap := make(map[string]int)
	for _, v := range str {
		wordCountMap[v] += 1
	}
	return wordCountMap
}

func removeStopWords(stopWords []string, reviews map[string]int) map[string]int {
	for _, v := range stopWords {
		if reviews[v] > 0 {
		   delete(reviews, v)
		}
	}
	return reviews
}

func sortReviewsList(reviews map[string]int) PairList {
	reviewsList := make(PairList, len(reviews))
	i := 0
	for key, value := range reviews {
		reviewsList[i] = Pair{key, value}
		i++
	}
	sort.Sort(sort.Reverse(reviewsList))
	return reviewsList
}

type Pair struct {
	Key string
	Value int
}

type PairList []Pair

func (p PairList) Swap(i, j int) { p[i], p[j] = p[j], p[i] }

func (p PairList) Len() int { return len(p) }

func (p PairList) Less(i, j int) bool { return p[i].Value < p[j].Value }

func SortedResult() PairList {
	stopWords := stopWordsList()
	content := textmining.TextMining()
	reg, err := regexp.Compile("[^a-zA-Z0-9]+")
	if err != nil {
		log.Println(err)
	}
	processedString := reg.ReplaceAllString(content, " ")
	lowerString := strings.ToLower(processedString)
	reviewsList := reviewsList(lowerString)
	lemmatizeWordsList := lemmatizeWords(reviewsList)
	reviewsListCount := wordCount(lemmatizeWordsList)
	removeStopWords := removeStopWords(stopWords, reviewsListCount)
	return sortReviewsList(removeStopWords)
}
