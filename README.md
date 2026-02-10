go-curl-impersonate
=======

> **Note:** This is a temporary fork of [github.com/BridgeSenseDev/go-curl-impersonate](https://github.com/BridgeSenseDev/go-curl-impersonate). Please refer to the original repository for the latest updates. This fork exists only until the PR fixing windows build & header issues is accepted upstream.

go-curl-impersonate is a GoLang interface to [libcurl-impersonate](https://github.com/lwthiker/curl-impersonate#libcurl-impersonate),
a special build of the multiprotocol file transfer library that can impersonate major web browsers.
Similar to the HTTP support in [net/http](https://pkg.go.dev/net/http), go-curl-impersonate can be used to
fetch objects from a Go program. While go-curl-impersonate can provide simple fetches,
it also exposes most of the functionality of libcurl, including:

 * Browser Impersonation: Perform TLS and HTTP handshakes identical to real browsers.
 * Speed - libcurl is very fast.
 * Multiple protocols (not just HTTP).
 * SSL, authentication and proxy support.
 * Support for libcurl's callbacks.

This said, libcurl API can be less easy to learn than net/http.

LICENSE
-------

go-curl-impersonate is licensed under the Apache License, Version 2.0 (http://www.apache.org/licenses/LICENSE-2.0.html).

Current Development Status
--------------------------

 * currently *unstable*
 * websocket support soon
 * READ, WRITE, HEADER, PROGRESS function callback
 * a Multipart Form supports file uploading
 * Most curl_easy_setopt option
 * partly implement share & multi interface
 * new callback function prototype

Requirements
------------
 * Any version of Go
 * libcurl-impersonate binaries are included in the libs directory
 * Python 3 (used only by configure scripts)


How to Install
--------------

    $ go get -u github.com/BridgeSenseDev/go-curl-impersonate

Current Status
--------------

 * Linux x64
   * passed go1 (ArchLinux)
 * Windows x86
   * passed go1 (win7, mingw-gcc 4.5.2, curl 7.22.0)
 * Mac OS
   * passed go1 (Mac OS X 10.7.3, curl 7.21.4)

NOTE: Above information is outdated ("help wanted")

Sample Program
--------------

Following comes from [examples/https_impersonate.go](./examples/https_impersonate.go) and demonstrates impersonation of tls and http handshakes on [browserleaks](https://browserleaks.com/tls).
Simply type `go run ./examples/https_impersonate.go` to execute.
```go
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
```

See also the [examples](./examples/) directory!

Acknowledgements
----------------

 * [curl-cffi](https://github.com/lexiforest/curl_cffi): Python binding for curl-impersonate fork via cffi. Inspiration and blueprint for this package.
  * [go-curl](https://github.com/andelf/go-curl): The original go bindings for libcurl, from which this package is forked.
 * [curl-impersonate](https://github.com/lwthiker/curl-impersonate): A special build of curl that can impersonate Chrome & Firefox
 * [curl](https://curl.se/): A command line tool and library for transferring data with URL syntax,