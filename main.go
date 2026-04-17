package main

import (
	"log"

	climain "github.com/mehmetalidsy/madget-cli/apps/cli"
)

func main() {
	if err := climain.Run(); err != nil {
		log.Fatal(err)
	}
}
