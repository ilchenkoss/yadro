package scraper

import (
	"myapp/pkg/database"
	"reflect"
	"testing"
)

func TestDecoder(t *testing.T) {

	str := `{"month": "1", "num": 10, "link": "", "year": "2006", "news": "", "safe_title": "Pi Equals", "transcript": "Pi = 3.141592653589793helpimtrappedinauniversefactory7108914...", "alt": "My most famous drawing, and one of the first I did for the site", "img": "https://imgs.xkcd.com/comics/pi.jpg", "title": "Pi Equals", "day": "1"}`
	data := []byte(str)

	x := ResponseData{}
	x.Alt = "My most famous drawing, and one of the first I did for the site"
	x.Img = "https://imgs.xkcd.com/comics/pi.jpg"

	parserResult, _ := decodeResponse(data)

	if !reflect.DeepEqual(x, parserResult) {
		t.Error("\nResult was incorrect")
	}
}

func TestParser(t *testing.T) {
	str := `{"month": "1", "num": 10, "link": "", "year": "2006", "news": "", "safe_title": "Pi Equals", "transcript": "Pi = 3.141592653589793helpimtrappedinauniversefactory7108914...", "alt": "My most famous drawing, and one of the first I did for the site", "img": "https://imgs.xkcd.com/comics/pi.jpg", "title": "Pi Equals", "day": "1"}`
	data := []byte(str)

	x := database.ParsedData{}
	x.Url = "https://imgs.xkcd.com/comics/pi.jpg"
	x.Keywords = []string{"famous", "draw", "one", "first", "site"}

	y, _ := responseParser(data)

	if !reflect.DeepEqual(x, y) {
		t.Error("\nResult was incorrect")
	}
}
