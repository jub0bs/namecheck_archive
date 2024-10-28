package main

import (
	"fmt"
	"regexp"
	"strings"
)

var re = regexp.MustCompile(`^[a-zA-Z0-9-]{3,39}$`)

func main() {
	fmt.Println(IsValid("jub0bs"))
	fmt.Println(IsValid("khaled"))
	fmt.Println(IsValid("alaa"))
	fmt.Println(IsValid("khaled-"))
	fmt.Println(IsValid("a--laa"))
	fmt.Println(IsValid("khal*ed"))
}

func IsValid(username string) bool {
	return !strings.HasPrefix(username, "-") &&
		!strings.HasSuffix(username, "-") &&
		!strings.Contains(username, "--") &&
		re.MatchString(username)
}
