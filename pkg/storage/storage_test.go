package storage

import (
	"math/rand"
	"testing"
	"time"
)

func TestNew(t *testing.T) {
	_, err := New()
	if err != nil {
		t.Fatal(err)
	}
}

func TestDB_ListNews(t *testing.T) {
	rand.Seed(time.Now().UnixNano())
	newsrecord := []NewsRecord{
		{
			Title:       "Hello World!",
			Description: "My Meta World in Golang Develop",
			Link:        "https://skillfactory.ru",
		},
	}
	db, err := New()
	if err != nil {
		t.Fatal(err)
	}
	err = db.InsertNews(newsrecord)
	if err != nil {
		t.Fatal(err)
	}
	listnews, err := db.ListNews(1)
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("%+v", listnews)
}
