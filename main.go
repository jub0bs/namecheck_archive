package main

import (
	"fmt"
	"log"
	"regexp"
)

func main() {
	username := "jub0bs"
	ok, err := regexp.MatchString("^([[:alnum:]]-?){2,38}[[:alnum:]]$", username)
	if err != nil {
		log.Fatal(err)
	}
	if !ok {
		return
	}
	fmt.Println(username)
}
