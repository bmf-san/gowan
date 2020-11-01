package main

import (
	"fmt"
	"net/http"
	"time"

	gowan "github.com/bmf-san/gowan"
)

func main() {
	cache := gowan.New()

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fk := "first-key"
		sk := "second-key"

		cache.Put(fk, "first-value", time.Now().Add(2*time.Second).UnixNano())
		fmt.Printf("Put first-key: %v\n", fk)
		fmt.Printf("Get first-key: %v\n", cache.Get(fk))

		time.Sleep(10 * time.Second)

		fmt.Print("10 seconds have passed")
		fmt.Printf("Get first-key: %v\n", cache.Get(fk))

		if cache.Get(fk) == nil {
			cache.Put(sk, "second-value", time.Now().Add(100*time.Second).UnixNano())
			fmt.Printf("Put second-key: %v\n", sk)
		}
		fmt.Printf("Get second-key: %v\n", cache.Get(sk))
	})
	http.ListenAndServe(":8080", nil)
}
