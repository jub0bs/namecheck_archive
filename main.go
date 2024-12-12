package main

import (
	"fmt"
	"log"
	"regexp"
	"strings"
)

func main() {
	username := "jub0bs"
	if strings.HasPrefix(username, "-") ||
		strings.HasSuffix(username, "-") ||
		strings.Contains(username, "--") {
		return
	}
	ok, err := regexp.MatchString("^[a-zA-Z0-9-]{3,39}$", username)
	if err != nil {
		log.Fatal(err)
	}
	if !ok {
		return
	}
	fmt.Println(username)
}
