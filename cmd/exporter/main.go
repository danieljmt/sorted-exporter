package main

import (
	"flag"
	"log"

	"github.com/danieljmt/sorted-exporter"
)

var (
	username, password string
	dest               string
)

func main() {
	flag.StringVar(&username, "username", "", "Your sorted.club username")
	flag.StringVar(&password, "password", "", "Your sorted.club password")
	flag.StringVar(&dest, "dest", ".", "Destination folder for your paprika recipe archive")

	flag.Parse()

	s, err := sorted.New(username, password, dest)
	if err != nil {
		log.Fatal(err)
	}
	err = s.Run()
	if err != nil {
		log.Fatal(err)
	}
	log.Println("Sorted Exporter exited successfully")
}
