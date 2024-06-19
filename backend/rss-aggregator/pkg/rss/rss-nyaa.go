package rss

import (
	"encoding/xml"
	"net/http"
	"time"
)

type RSS struct {
	XMLName xml.Name `xml:"rss"`
	Version string   `xml:"version,attr"`
	Channel Channel  `xml:"channel"`
}

type Channel struct {
	Title       string `xml:"title"`
	Description string `xml:"description"`
	Link        string `xml:"link"`
	AtomLink    struct {
		Href string `xml:"href,attr"`
	} `xml:"atom\\:link"`
	Items []Item `xml:"item"`
}

type Item struct {
	Title       string `xml:"title"`
	Link        string `xml:"link"`
	GUID        string `xml:"guid"`
	PubDate     string `xml:"pubDate"`
	Seeders     int    `xml:"nyaa\\:seeders"`
	Leechers    int    `xml:"nyaa\\:leechers"`
	Downloads   int    `xml:"nyaa\\:downloads"`
	InfoHash    string `xml:"nyaa\\:infoHash"`
	CategoryID  string `xml:"nyaa\\:categoryId"`
	Category    string `xml:"nyaa\\:category"`
	Size        string `xml:"nyaa\\:size"`
	Comments    int    `xml:"nyaa\\:comments"`
	Trusted     string `xml:"nyaa\\:trusted"`
	Remake      string `xml:"nyaa\\:remake"`
	Description string `xml:"description"`
}

type Resolution string

const (
	Resolution480p  Resolution = "480p"
	Resolution720p  Resolution = "720p"
	Resolution1080p Resolution = "1080p"
	Resolution4K    Resolution = "4K"
)

type Language string

const (
	ENG   Language = "ENG"
	PORBR Language = "POR-BR"
	SPALA Language = "SPA-LA"
	SPA   Language = "SPA"
	ARA   Language = "ARA"
	FRE   Language = "FRE"
	GER   Language = "GER"
	ITA   Language = "ITA"
	RUS   Language = "RUS"
)

// ParseRSS parses the RSS feed from the given byte slice.
func FetchAndParseRSS(language Language, resolution Resolution) (*RSS, error) {

	url := "https://nyaa.si/?page=rss&c=1_0&f=0&q=" + string(language) + "+" + string(resolution)

	rss := &RSS{}
	httpClient := &http.Client{
		Timeout: time.Second * 30,
	}

	resp, err := httpClient.Get(url)

	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	err = xml.NewDecoder(resp.Body).Decode(rss)

	if err != nil {
		return nil, err
	}

	return rss, nil
}
