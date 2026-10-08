// Command seed creates an example clover database populated with a "movies"
// collection. It is meant to be used to manually test the clovercli commands.
//
// Usage:
//
//	go run ./examples/seed [-out examples/movies]
package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"

	cl "github.com/ostafen/clover/v2"
	d "github.com/ostafen/clover/v2/document"
	q "github.com/ostafen/clover/v2/query"
)

// movie is a sample document of the "movies" collection. Field names are
// lowercased through the "clover" struct tag, so they can be used as query
// criteria from the clovercli prompt.
type movie struct {
	Title   string   `clover:"title"`
	Year    int      `clover:"year"`
	Author  string   `clover:"author"`
	Cast    []string `clover:"cast"`
	Genres  []string `clover:"genres"`
	Rating  float64  `clover:"rating"`
	Runtime int      `clover:"runtime"`
}

func sampleMovies() []movie {
	return []movie{
		{
			Title:   "Inception",
			Year:    2010,
			Author:  "Christopher Nolan",
			Cast:    []string{"Leonardo DiCaprio", "Joseph Gordon-Levitt", "Elliot Page", "Tom Hardy"},
			Genres:  []string{"Action", "Sci-Fi", "Thriller"},
			Rating:  8.8,
			Runtime: 148,
		},
		{
			Title:   "The Dark Knight",
			Year:    2008,
			Author:  "Christopher Nolan",
			Cast:    []string{"Christian Bale", "Heath Ledger", "Aaron Eckhart", "Michael Caine"},
			Genres:  []string{"Action", "Crime", "Drama"},
			Rating:  9.0,
			Runtime: 152,
		},
		{
			Title:   "Interstellar",
			Year:    2014,
			Author:  "Christopher Nolan",
			Cast:    []string{"Matthew McConaughey", "Anne Hathaway", "Jessica Chastain"},
			Genres:  []string{"Adventure", "Drama", "Sci-Fi"},
			Rating:  8.7,
			Runtime: 169,
		},
		{
			Title:   "Pulp Fiction",
			Year:    1994,
			Author:  "Quentin Tarantino",
			Cast:    []string{"John Travolta", "Uma Thurman", "Samuel L. Jackson", "Bruce Willis"},
			Genres:  []string{"Crime", "Drama"},
			Rating:  8.9,
			Runtime: 154,
		},
		{
			Title:   "The Grand Budapest Hotel",
			Year:    2014,
			Author:  "Wes Anderson",
			Cast:    []string{"Ralph Fiennes", "Tony Revolori", "Saoirse Ronan", "Adrien Brody"},
			Genres:  []string{"Adventure", "Comedy", "Drama"},
			Rating:  8.1,
			Runtime: 99,
		},
		{
			Title:   "Parasite",
			Year:    2019,
			Author:  "Bong Joon-ho",
			Cast:    []string{"Song Kang-ho", "Lee Sun-kyun", "Cho Yeo-jeong", "Choi Woo-shik"},
			Genres:  []string{"Drama", "Thriller"},
			Rating:  8.5,
			Runtime: 132,
		},
		{
			Title:   "Spirited Away",
			Year:    2001,
			Author:  "Hayao Miyazaki",
			Cast:    []string{"Rumi Hiiragi", "Miyu Irino", "Mari Natsuki"},
			Genres:  []string{"Animation", "Adventure", "Family"},
			Rating:  8.6,
			Runtime: 125,
		},
		{
			Title:   "The Lord of the Rings: The Return of the King",
			Year:    2003,
			Author:  "Peter Jackson",
			Cast:    []string{"Elijah Wood", "Viggo Mortensen", "Ian McKellen", "Orlando Bloom"},
			Genres:  []string{"Action", "Adventure", "Drama", "Fantasy"},
			Rating:  9.0,
			Runtime: 201,
		},
		{
			Title:   "La La Land",
			Year:    2016,
			Author:  "Damien Chazelle",
			Cast:    []string{"Ryan Gosling", "Emma Stone", "John Legend"},
			Genres:  []string{"Comedy", "Drama", "Music", "Romance"},
			Rating:  8.0,
			Runtime: 128,
		},
		{
			Title:   "Jurassic Park",
			Year:    1993,
			Author:  "Steven Spielberg",
			Cast:    []string{"Sam Neill", "Laura Dern", "Jeff Goldblum", "Richard Attenborough"},
			Genres:  []string{"Action", "Adventure", "Sci-Fi"},
			Rating:  8.2,
			Runtime: 127,
		},
	}
}

func main() {
	log.SetFlags(0)

	out := flag.String("out", filepath.Join("examples", "movies"), "directory of the example clover database")
	flag.Parse()

	dir, err := filepath.Abs(*out)
	if err != nil {
		log.Fatal(err)
	}

	// start from a clean state, so the database can be regenerated at will
	if err := os.RemoveAll(dir); err != nil {
		log.Fatal(err)
	}
	if err := os.MkdirAll(dir, 0o755); err != nil {
		log.Fatal(err)
	}

	db, err := cl.Open(dir)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	const collection = "movies"
	if err := db.CreateCollection(collection); err != nil {
		log.Fatal(err)
	}

	for _, m := range sampleMovies() {
		docID, err := db.InsertOne(collection, d.NewDocumentOf(m))
		if err != nil {
			log.Fatal(err)
		}
		fmt.Printf("inserted %-50s _id: %s\n", m.Title, docID)
	}

	movies, err := db.FindAll(q.NewQuery(collection))
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("\ncreated %s with collection %q holding %d documents\n", filepath.Join(dir, "data.db"), collection, len(movies))
	fmt.Printf("run it with: go run . %s\n", *out)
}
