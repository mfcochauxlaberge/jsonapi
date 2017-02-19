package jsonapi

import "time"

var (
	users Collection
	books Collection
	urls  []*URL
)

func init() {
	loc, _ := time.LoadLocation("")

	// Resources
	users = &WrapperCollection{}
	users.Add(
		Wrap(&user{
			ID:           "1",
			Name:         "Bob",
			Age:          36,
			CreatedAt:    time.Date(2017, 1, 2, 3, 4, 5, 6, loc),
			BestFriend:   "2",
			Contacts:     []string{"2", "3"},
			FavoriteBook: "",
			Readings:     []string{},
		}),
	)
	users.Add(
		Wrap(&user{
			ID:           "2",
			Name:         "Noam Chomsky",
			Age:          100,
			CreatedAt:    time.Date(2017, 1, 2, 3, 4, 5, 6, loc),
			BestFriend:   "1",
			Contacts:     []string{"3"},
			FavoriteBook: "1",
			Readings:     []string{"1", "2"},
		}),
	)
	users.Add(
		Wrap(&user{
			ID:           "3",
			Name:         "The Dude",
			Age:          45,
			CreatedAt:    time.Date(2017, 1, 2, 3, 4, 5, 6, loc),
			BestFriend:   "",
			Contacts:     []string{"1", "2"},
			FavoriteBook: "1",
			Readings:     []string{"1", "2"},
		}),
	)

	books = &WrapperCollection{}
	books.Add(
		Wrap(&book{
			ID:        "1",
			Title:     "Understanding Power",
			CreatedAt: time.Now(),
			Author:    "3",
		}),
	)
	books.Add(
		Wrap(&book{
			ID:        "2",
			Title:     "The Title of a Book",
			CreatedAt: time.Now(),
			Author:    "1",
		}),
	)

	urls = []*URL{
		&URL{
			Params: &Params{
				RelData: map[string][]string{
					"users": []string{"best-friend", "contacts"},
				},
			},
		},
		&URL{
			Params: &Params{
				Fields: map[string][]string{
					"users": []string{"name", "readings"},
				},
				RelData: map[string][]string{
					"users": []string{"contacts", "readings"},
				},
			},
		},
	}
}

// User
type user struct {
	ID string `json:"id" api:"users"`

	// Attributes
	Name      string    `json:"name" api:"attr"`
	Age       uint8     `json:"age" api:"attr"`
	CreatedAt time.Time `json:"created-at" api:"attr"`

	// Relationships
	BestFriend   string   `json:"best-friend" api:"rel,users"`
	Contacts     []string `json:"contacts" api:"rel,users"`
	FavoriteBook string   `json:"favorite-book" api:"rel,books"`
	Readings     []string `json:"readings" api:"rel,books"`
}

// Book
type book struct {
	ID string `json:"id" api:"books"`

	// Attributes
	Title     string    `json:"title" api:"attr"`
	CreatedAt time.Time `json:"written-at" api:"attr"`

	// Relationships
	Author string `json:"author" api:"rel,users"`
}
