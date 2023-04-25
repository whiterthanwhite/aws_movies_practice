package main

import (
	"bytes"
	"encoding/base64"
	"fmt"
	"io"
	"log"
	"os"
)

func main() {
	body, err := io.ReadAll(os.Stdin)
	if err != nil {
		log.Fatal(err.Error())
	}
	var dBody []byte
	b := bytes.NewBuffer(dBody)
	encoder := base64.NewEncoder(base64.StdEncoding, b)
	_, err = encoder.Write(body)
	if err != nil {
		log.Fatal(err.Error())
	}
	fmt.Println(b)
}
