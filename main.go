package main

import (
	"fmt"
	"github.com/HenryHK/http-cache-middleware/adapter"
	api "github.com/HenryHK/http-cache-middleware/api"
	"github.com/HenryHK/http-cache-middleware/cache"
	"net/http"
	"os"
	"time"
)

func example(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Ok"))
}

func main() {
	fmt.Println("Hello, World")

	opt := &adapter.RingOptions{
		Addrs: map[string]string{
			"server": "localhost:6379",
		},
	}

	client, err := cache.New(
		cache.ClientWithAdapter(adapter.New(opt)),
		cache.ClientWithTTL(10*time.Minute),
	)

	if err != nil {
		fmt.Printf("error when create cache client %v\n", err)
		os.Exit(1)
	}

	getHandler := http.HandlerFunc(api.GetHandler)
	postHandler := http.HandlerFunc(api.PostHandler)

	// GET Handler
	http.Handle("/contact/", client.Middleware(getHandler))
	// POST/DELETE/PUT handler
	http.Handle("/contact", client.Middleware(postHandler))
	http.ListenAndServe(":8080", nil)
}
