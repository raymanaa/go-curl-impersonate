package main

import (
	"fmt"
	curl "github.com/BridgeSenseDev/go-curl-impersonate"
	"io"
	"os"
)

func writeData(ptr []byte, userdata interface{}) bool {
	writer, ok := userdata.(io.Writer)
	if !ok {
		fmt.Println("WriteData: userdata is not an io.Writer")
		return false
	}
	_, err := writer.Write(ptr)
	return err == nil
}

func main() {
	easy := curl.EasyInit()
	if easy == nil {
		fmt.Println("EasyInit failed")
		return
	}
	defer easy.Cleanup()

	err := easy.Setopt(curl.OPT_URL, "https://tls.browserleaks.com/json")
	if err != nil {
		fmt.Printf("Setopt URL failed: %v\n", err)
		return
	}

	easy.Setopt(curl.OPT_VERBOSE, true)

	err = easy.Setopt(curl.OPT_ACCEPT_ENCODING, "")
	if err != nil {
		fmt.Printf("Setopt OPT_ACCEPT_ENCODING failed: %v\n", err)
		return
	}

	err = easy.Impersonate("chrome136", true)
	if err != nil {
		fmt.Printf("Impersonate failed: %v\n", err)
	}

	easy.Setopt(curl.OPT_WRITEFUNCTION, writeData)
	easy.Setopt(curl.OPT_WRITEDATA, os.Stdout)

	fmt.Println("Performing request...")
	err = easy.Perform()
	if err != nil {
		fmt.Printf("Perform failed: %v\n", err)
	} else {
		fmt.Println("\nRequest performed successfully.")
	}
}
