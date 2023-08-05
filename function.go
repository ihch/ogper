package function

import (
	"encoding/json"
	"io"
	"log"
	"net/http"

	"github.com/GoogleCloudPlatform/functions-framework-go/functions"
	"github.com/PuerkitoBio/goquery"
)

// struct定義
type OGP struct {
	Icon        string `json:"icon"`
	Url         string `json:"url"`
	Title       string `json:"title"`
	Image       string `json:"image"`
	Description string `json:"description"`
	SiteName    string `json:"siteName"`
}

func GetOGP(url string) (og OGP, err error) {
	// URLからHTMLを取得
	response, err := http.Get(url)
	if err != nil {
		return og, err
	}

	defer response.Body.Close()

	// html内からogを取得
	doc, err := goquery.NewDocumentFromReader(response.Body)

	icon := doc.Find("link[rel='icon']").First().AttrOr("href", "No icon")
	ogurl := doc.Find("meta[property='og:url']").First().AttrOr("content", url)
	pageTitle := doc.Find("title").First().AttrOr("content", "No title")
	title := doc.Find("meta[property='og:title']").First().AttrOr("content", pageTitle)
	image := doc.Find("meta[property='og:image']").First().AttrOr("content", "No image")
	description := doc.Find("meta[property='og:description']").First().AttrOr("content", "No description")
	siteName := doc.Find("meta[property='og:site_name']").First().AttrOr("content", "No siteName")

	// ogをjsonで返す
	og = OGP{
		icon,
		ogurl,
		title,
		image,
		description,
		siteName,
	}

	return og, err
}

func httpFunction(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Headers", "*")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "GET,OPTIONS")
	log.Println("requested")

	// URLからOGPのデータを取得
	urls := r.URL.Query()["url"]
	ogps := make([]OGP, 0)

	for _, url := range urls {
		ogp, err := GetOGP(url)
		if err != nil {
			log.Fatalln(err)
			continue
		}
		ogps = append(ogps, ogp)
	}

	responseJson, err := json.Marshal(struct {
		Ogps []OGP `json:"ogps"`
	}{ogps})

	if err != nil {
		log.Fatalln(err)
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Internal Server Error"))
	}

	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Cache-Control", "public, max-age=86400")
	io.WriteString(w, string(responseJson))
}

func init() {
	functions.HTTP("httpFunction", httpFunction)
}
