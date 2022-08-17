package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"multiNews/pkg/api"
	"multiNews/pkg/rss"
	"multiNews/pkg/storage"
	"net/http"
	"time"
)

type config struct {
	URLs          []string `json:"rss"`
	RequestPeriod int      `json:"req_period"`
}

func main() {
	var conf config
	// init db
	db, err := storage.New()
	if err != nil {
		log.Fatal(err)
	}
	// init api
	api := api.New(db)
	// parse config
	cFile, err := ioutil.ReadFile("./config.json")
	if err != nil {
		log.Fatal(err)
	}
	err = json.Unmarshal(cFile, &conf)
	if err != nil {
		log.Fatal(err)
	}

	// создаем каналы для новостей и ошибок, согласно задания
	chNews := make(chan []storage.NewsRecord)
	chErrs := make(chan error)

	// получаем данные
	for _, url := range conf.URLs {
		go urlParser(url, db, chNews, chErrs, conf.RequestPeriod)
	}

	// write data
	go func() {
		for news := range chNews {
			db.InsertNews(news)
		}
	}()

	go func() {
		for err := range chErrs {
			log.Println("обнаружена ошибка: ", err)
		}
	}()

	// server start
	err = http.ListenAndServe(":80", api.Router())
	if err != nil {
		log.Fatal(err)
	}
}

func urlParser(url string, db *storage.DB, chNews chan<- []storage.NewsRecord, chErrs chan<- error, reqPeriod int) {
	// запускаем бесконечный цикл со слипом как задание по расписанию
	for {
		news, err := rss.Parser(url)
		if err != nil {
			chErrs <- err
			continue
		}
		chNews <- news
		time.Sleep(time.Minute * time.Duration(reqPeriod))
	}
}
