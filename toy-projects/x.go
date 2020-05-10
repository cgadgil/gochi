package main

import (
    "bytes"
    "fmt"
    "strings"
    "unicode/utf8"
)

func main() {
    b := "म्हणाला"
    len1 := len([]rune(b))
    len2 := bytes.Count([]byte(b), nil) -1
    len3 := strings.Count(b, "") - 1
    len4 := utf8.RuneCountInString(b)
    fmt.Println(len1)
    fmt.Println(len2)
    fmt.Println(len3)
    fmt.Println(len4)

}

