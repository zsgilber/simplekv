package main

import (
	"fmt"
	"os"

	"github.com/zsgilber/simplekv/pkg/resp"
)

func main() {
	server := resp.NewServer()
	if err := server.ListenAndServe("localhost:3003"); err != nil {
		fmt.Println("error listening")
		os.Exit(1)
	}
}
