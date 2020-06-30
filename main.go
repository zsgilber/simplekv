package main

import (
	"fmt"
	"os"

	"github.com/zsgilber/simplekv/pkg/kv"
	"github.com/zsgilber/simplekv/pkg/resp"
)

func main() {
	m := make(map[string]string)
	store := &kv.MapStore{
		Map: m,
	}
	server := resp.NewServer(store)
	fmt.Println("starting server...")
	if err := server.ListenAndServe("localhost:3003"); err != nil {
		fmt.Println("error listening")
		os.Exit(1)
	}
}
