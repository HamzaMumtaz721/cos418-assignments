package cos418_hw1_1

import (
	"fmt"
	"io/ioutil"
	"regexp"
	"sort"
	"strings"
)

func topWords(path string, numWords int, charThreshold int) []WordCount {
	data, err := ioutil.ReadFile(path)
	checkError(err)

	text := strings.ToLower(string(data))

	// Remove non-alphanumeric characters completely (don't replace with space)
	re := regexp.MustCompile("[^0-9a-zA-Z\\s]+")
	text = re.ReplaceAllString(text, "")

	tokens := strings.Fields(text)

	freq := make(map[string]int)
	for _, w := range tokens {
		if len(w) >= charThreshold {
			freq[w]++
		}
	}

	var wordCounts []WordCount
	for w, c := range freq {
		wordCounts = append(wordCounts, WordCount{w, c})
	}

	sortWordCounts(wordCounts)

	if numWords > len(wordCounts) {
		numWords = len(wordCounts)
	}
	return wordCounts[:numWords]
}

type WordCount struct {
	Word  string
	Count int
}

func (wc WordCount) String() string {
	return fmt.Sprintf("%v: %v", wc.Word, wc.Count)
}

func sortWordCounts(wordCounts []WordCount) {
	sort.Slice(wordCounts, func(i, j int) bool {
		wc1 := wordCounts[i]
		wc2 := wordCounts[j]
		if wc1.Count == wc2.Count {
			return wc1.Word < wc2.Word
		}
		return wc1.Count > wc2.Count
	})
}
