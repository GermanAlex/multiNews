package rss

import (
	"encoding/xml"
	"io/ioutil"
	"multiNews/pkg/storage"
	"net/http"
	"strings"
	"time"
)

/*Определяем стандартные RSS-структуры*/

type Feed struct {
	XMLName xml.Name `xml:"rss"`
	Channel Channel  `xml:"channel"`
}

type Channel struct {
	Title       string `xml:"title"`
	Description string `xml:"description"`
	Link        string `xml:"link"`
	Items       []Item `xml:"item"`
}

type Item struct {
	Title       string `xml:"title"`
	Description string `xml:"description"`
	PubDate     string `xml:"pubDate"`
	Link        string `xml:"link"`
}

func Parser(url string) ([]storage.NewsRecord, error) {
	var f Feed
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	err = xml.Unmarshal(body, &f)
	if err != nil {
		return nil, err
	}
	var news []storage.NewsRecord
	for _, item := range f.Channel.Items {
		var n storage.NewsRecord
		n.Title = item.Title
		n.Description = item.Description
		n.Link = item.Link
		// обработка даты из формата RSS и стринга в Unix Time
		item.PubDate = strings.ReplaceAll(item.PubDate, ",", "")
		// подгоняем по формату RFC1123Z
		timeResult, err := time.Parse("Mon 2 Jan 2006 15:04:05 -0700", item.PubDate)
		if err != nil {
			// если не получилось гоним в RFC850
			timeResult, err = time.Parse("Mon 2 Jan 2006 15:04:05 GMT", item.PubDate)
		}
		if err == nil {
			// если ошибки конвертации пустые, то переводим в unix и заполняем
			n.PublicTime = timeResult.Unix()
		}

		news = append(news, n)
	}
	return news, nil
}
