package main

import (
	"fmt"
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
	if match, _ := regexp.MatchString(`^[a-zA-Z0-9-]{3,39}$`, username); !match {
		return
	}
	fmt.Println(username)
}
