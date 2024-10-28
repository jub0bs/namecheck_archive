package main

import (
	"fmt"
	"regexp"
	"strings"
)

func main() {
	fmt.Println(IsValid("jub0bs"))
	fmt.Println(IsValid("khaled"))
	fmt.Println(IsValid("alaa"))
	fmt.Println(IsValid("khaled-"))
	fmt.Println(IsValid("a--laa"))
	fmt.Println(IsValid("khal*ed"))
}

func IsValid(username string) bool {
	if strings.HasPrefix(username, "-") ||
		strings.HasSuffix(username, "-") ||
		strings.Contains(username, "--") {
		return false
	}
	if match, _ := regexp.MatchString(`^[a-zA-Z0-9-]{3,39}$`, username); !match {
		return false
	}
	return true
}
