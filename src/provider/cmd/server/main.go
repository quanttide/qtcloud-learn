// Command server 是量潮学习云的服务端入口。
package main

import (
	"log"
	"net/http"
)

func main() {
	addr := ":8080"
	log.Printf("qtcloud-learn provider listening on %s", addr)
	if err := http.ListenAndServe(addr, handler()); err != nil {
		log.Fatal(err)
	}
}
