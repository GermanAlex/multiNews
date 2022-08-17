package rss

import (
	"testing"
)

func TestParser(t *testing.T) {
	tfeed, err := Parser("https://habr.com/ru/rss/best/daily/?fl=ru")
	if err != nil {
		t.Fatal(err)
	}
	if len(tfeed) == 0 {
		t.Fatal("данные не раскодированы")
	}
}
