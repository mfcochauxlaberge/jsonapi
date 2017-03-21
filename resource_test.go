package jsonapi

import (
	"testing"
	"time"

	"github.com/kkaribu/tchek"
)

func TestResource(t *testing.T) {
	loc, _ := time.LoadLocation("")

	p1 := &painting{
		ID:        "persistence-memory",
		Title:     "The Persistence of Memory",
		PaintedIn: time.Date(1931, 0, 0, 0, 0, 0, 0, loc),
		Author:    "some-artist",
	}

	res := Wrap(p1)

	// Get
	tchek.AreEqual(t, 0, p1.Title, res.Get("title"))
	tchek.AreEqual(t, 1, "some-artist", res.GetToOne("author"))

	// Set
	res.Set("title", "New Title")
	tchek.AreEqual(t, 2, "New Title", p1.Title)
	tchek.AreEqual(t, 3, "New Title", res.Get("title"))

	p1.PaintedIn = time.Date(1932, 0, 0, 0, 0, 0, 0, loc)
	tchek.AreEqual(t, 4, p1.PaintedIn, res.Get("painted-in"))

	res.SetToOne("author", "another-artist")
	tchek.AreEqual(t, 5, "another-artist", p1.Author)
	tchek.AreEqual(t, 6, "another-artist", res.GetToOne("author"))
}

type painting struct {
	ID string `json:"id" api:"paintings"`

	Title     string    `json:"title" api:"attr"`
	Value     uint      `json:"value" api:"attr"`
	PaintedIn time.Time `json:"painted-in" api:"attr"`

	Author string `json:"author" api:"rel,artists,paintings"`
}

type artist struct {
	ID string `json:"id" api:"artists"`

	Name   string    `json:"name" api:"attr"`
	BornAt time.Time `json:"born-at" api:"attr"`

	Paintings string `json:"paintings" api:"rel,paintings,author"`
}
