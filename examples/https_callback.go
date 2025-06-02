package main

import (
	"fmt"
	curl "github.com/BridgeSenseDev/go-curl-impersonate"
)

func main() {
	easy := curl.EasyInit()
	defer easy.Cleanup()

	easy.Setopt(curl.OPT_URL, "https://www.baidu.com/")

	// OPTIONAL - make a callback function
	fooTest := func(buf []byte, userdata any) bool {
		println("DEBUG: size=>", len(buf))
		println("DEBUG: content=>", string(buf))
		return true
	}

	easy.Setopt(curl.OPT_WRITEFUNCTION, fooTest)

	if err := easy.Perform(); err != nil {
		fmt.Printf("ERROR: %v\n", err)
	}
}
