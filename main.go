package main

import (
	"fmt"
	"log"

	"github.com/jub0bs/namecheck/bluesky"
	"github.com/jub0bs/namecheck/github"
)

func main() {
	username := "jub0bs"
	if !github.IsValid(username) {
		fmt.Printf("%q is invalid on GitHub.\n", username)
		return
	}
	avail, err := github.IsAvailable(username)
	if err != nil {
		log.Fatal(err)
	}
	if !avail {
		fmt.Printf("%q is valid but unavailable on GitHub.\n", username)
		return
	}
	fmt.Printf("%q is valid and available on GitHub.\n", username)

	if !bluesky.IsValid(username) {
		fmt.Printf("%q is invalid on Bluesky.\n", username)
		return
	}
	avail, err = bluesky.IsAvailable(username)
	if err != nil {
		log.Fatal(err)
	}
	if !avail {
		fmt.Printf("%q is valid but unavailable on Bluesky.\n", username)
		return
	}
	fmt.Printf("%q is valid and available on Bluesky.\n", username)
}
