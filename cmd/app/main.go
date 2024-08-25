package main

import (
	"github.com/ananikitina/notes-rest/internal/app"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	_ "github.com/lib/pq"
)

func main() {
	app.Start()
}
