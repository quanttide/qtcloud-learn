// Command server 是量潮学习云的服务端入口。
package main

import (
	"log"
	"net/http"

	"github.com/quanttide/qtcloud-learn-provider/internal/version"
)

func main() {
	addr := ":8080"
	log.Printf("qtcloud-learn provider %s listening on %s", version.Version, addr)
	if err := http.ListenAndServe(addr, newRouter()); err != nil {
		log.Fatal(err)
	}
}
