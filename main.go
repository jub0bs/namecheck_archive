package main

import (
	"fmt"
	"net/http"
	"os"

	"github.com/jub0bs/namecheck/github"
	"github.com/jub0bs/namecheck/reddit"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Fprintln(os.Stderr, "usage: namecheck <username>")
		os.Exit(1)
	}
	username := os.Args[1]
	gh := github.GitHub{
		Client: http.DefaultClient,
	}
	valid := gh.IsValid(username)
	fmt.Printf("validity of %q on GitHub: %t\n", username, valid)
	if valid {
		avail, err := gh.IsAvailable(username)
		if err != nil {
			fmt.Println(err)
			return
		}
		fmt.Printf("availability of %q on GitHub: %t\n", username, avail)
	}
	rd := reddit.Reddit{
		Client: http.DefaultClient,
	}
	valid = rd.IsValid(username)
	fmt.Printf("validity of %q on Reddit: %t\n", username, valid)
	if valid {
		avail, err := rd.IsAvailable(username)
		if err != nil {
			fmt.Println(err)
			return
		}
		fmt.Printf("availability of %q on Reddit: %t\n", username, avail)
	}
}
