package main

import (
    "fmt"
     curlCmd "github.com/golang-basic/go-curl"
)

func main() {

 easy := curlCmd.EasyInit()
 defer easy.Cleanup()

 easy.Setopt(curlCmd.OPT_URL, "http://golang-basic.blogspot.com/2014/05/why-go-dont-use-space-indentation-as.html?view=sidebar")

 // make a callback function

  fooTest := func (buf []byte, userdata interface{}) bool {
      println("DEBUG: size=>", len(buf))
      println("DEBUG: content=>", string(buf))
  
      return true
  }
 
  easy.Setopt(curlCmd.OPT_WRITEFUNCTION, fooTest)
  if err := easy.Perform(); err != nil {
      fmt.Printf("ERROR: %v\n", err)
  }

}
