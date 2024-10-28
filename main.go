package main

import (
	"fmt"

	"github.com/jub0bs/namecheck/github"
)

func main() {
	fmt.Println(github.IsValid("jub0bs"))
	fmt.Println(github.IsValid("khaled"))
	fmt.Println(github.IsValid("alaa"))
	fmt.Println(github.IsValid("khaled-"))
	fmt.Println(github.IsValid("a--laa"))
	fmt.Println(github.IsValid("khal*ed"))
}
