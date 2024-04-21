package scraper

import (
	"reflect"
	"sync"
	"testing"
)

func TestParserWorker(t *testing.T) {

	str1 := `{"month": "1", "num": 10, "link": "", "year": "2006", "news": "", "safe_title": "Pi Equals", "transcript": "Pi = 3.141592653589793helpimtrappedinauniversefactory7108914...", "alt": "My most famous drawing", "img": "https://imgs.xkcd.com/comics/pi.jpg", "title": "Pi Equals", "day": "1"}`
	data1 := []byte(str1)
	str2 := `{"month": "1", "num": 3, "link": "", "year": "2006", "news": "", "safe_title": "Island (sketch)", "transcript": "[[A sketch of an Island]]", "alt": "Hello, island", "img": "https://imgs.xkcd.com/comics/island_color.jpg", "title": "Island (sketch)", "day": "1"}`
	data2 := []byte(str2)

	wantResult := map[int]ScrapedData{
		3: {
			Keywords: map[string][]KeywordInfo{
				"hello":  {KeywordInfo{1, 0, "alt"}},
				"island": {KeywordInfo{1, 0, "title"}, KeywordInfo{1, 1, "transcript"}, KeywordInfo{1, 1, "alt"}},
				"sketch": {KeywordInfo{1, 1, "title"}, KeywordInfo{1, 0, "transcript"}}},
			Url: "https://imgs.xkcd.com/comics/island_color.jpg",
		},
		10: {
			Keywords: map[string][]KeywordInfo{
				"draw":                            {KeywordInfo{1, 1, "alt"}},
				"equal":                           {KeywordInfo{1, 0, "title"}},
				"famous":                          {KeywordInfo{1, 0, "alt"}},
				"helpimtrappedinauniversefactori": {KeywordInfo{1, 0, "transcript"}}},
			Url: "https://imgs.xkcd.com/comics/pi.jpg",
		}}

	goodScrapesCh := make(chan []byte, 1)
	resultCh := make(chan map[int]ScrapedData, 1)
	var pwg sync.WaitGroup
	pwg.Add(2)

	go func() {
		goodScrapesCh <- data1
		goodScrapesCh <- data2
	}()

	go parserWorker(map[int]ScrapedData{}, goodScrapesCh, &pwg, resultCh)

	pwg.Wait()
	close(goodScrapesCh)

	result := <-resultCh
	close(resultCh)

	if !reflect.DeepEqual(result, wantResult) {
		t.Error("\nResult was incorrect")
	}
}
