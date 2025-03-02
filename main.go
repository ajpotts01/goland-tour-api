package main

import (
	"github.com/ajpotts01/goland-tour-api/internal/db"
	"github.com/ajpotts01/goland-tour-api/internal/todo"
	"github.com/ajpotts01/goland-tour-api/internal/transport"
	"log"
)

func main() {
	d, err := db.New("postgres", "example", "localhost", "postgres", 5432)
	if err != nil {
		log.Fatal(err)
	}
	svc := todo.NewService(d)
	server := transport.NewServer(svc)

	if err := server.Serve(); err != nil {
		log.Fatal(err)
	}
}
