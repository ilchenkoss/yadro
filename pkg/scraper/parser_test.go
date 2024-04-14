package scraper

import (
	"fmt"
	"reflect"
	"sync"
	"testing"
)

func TestParser(t *testing.T) {
	str := `{"month": "1", "num": 10, "link": "", "year": "2006", "news": "", "safe_title": "Pi Equals", "transcript": "Pi = 3.141592653589793helpimtrappedinauniversefactory7108914...", "alt": "My most famous drawing, and one of the first I did for the site", "img": "https://imgs.xkcd.com/comics/pi.jpg", "title": "Pi Equals", "day": "1"}`
	data := []byte(str)

	x := ParsedData{
		ID:       10,
		Keywords: []string{"famous", "draw", "one", "first", "site"},
		Url:      "https://imgs.xkcd.com/comics/pi.jpg",
	}

	y, _ := responseParser(data)
	fmt.Println(x, y)
	if !reflect.DeepEqual(x, y) {
		t.Error("\nResult was incorrect")
	}
}

func TestParserWorker(t *testing.T) {

	str1 := `{"month": "1", "num": 10, "link": "", "year": "2006", "news": "", "safe_title": "Pi Equals", "transcript": "Pi = 3.141592653589793helpimtrappedinauniversefactory7108914...", "alt": "My most famous drawing, and one of the first I did for the site", "img": "https://imgs.xkcd.com/comics/pi.jpg", "title": "Pi Equals", "day": "1"}`
	data1 := []byte(str1)
	str2 := `{"month": "1", "num": 3, "link": "", "year": "2006", "news": "", "safe_title": "Island (sketch)", "transcript": "[[A sketch of an Island]]\n{{Alt:Hello, island}}", "alt": "Hello, island", "img": "https://imgs.xkcd.com/comics/island_color.jpg", "title": "Island (sketch)", "day": "1"}`
	data2 := []byte(str2)

	wantResult := map[int]ScrapedData{
		10: {
			Keywords: []string{"famous", "draw", "one", "first", "site"},
			Url:      "https://imgs.xkcd.com/comics/pi.jpg",
		},
		3: {
			Keywords: []string{"hello", "island"},
			Url:      "https://imgs.xkcd.com/comics/island_color.jpg",
		},
	}

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

	result := make(map[int]ScrapedData)
	for r := range resultCh {
		result = r
		close(resultCh)
	}

	if !reflect.DeepEqual(result, wantResult) {
		t.Error("\nResult was incorrect")
	}
}
