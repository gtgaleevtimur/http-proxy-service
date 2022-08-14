package main

import (
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
)

//Локальный адрес прокси
const proxyAddr string = "localhost:8082"

var (
	counter        int    = 0
	firstHostAddr  string = "http://localhost:8080"
	secondHostAddr string = "http://localhost:8081"
)

//Запуск прокси-сервера
func main() {
	http.HandleFunc("/", handleProxy)
	log.Fatalln(http.ListenAndServe(proxyAddr, nil))
}

//Хэндлер для прокси-сервера
func handleProxy(w http.ResponseWriter, r *http.Request) {
	if counter == 0 {
		urlAdr, err := url.Parse(firstHostAddr)
		if err != nil {
			log.Fatalln(err)
		}

		proxy := httputil.NewSingleHostReverseProxy(urlAdr)
		proxy.ServeHTTP(w, r)

		counter++

		return
	}

	secondUrlAdr, err := url.Parse(secondHostAddr)
	if err != nil {
		log.Fatalln(err)
	}

	proxy := httputil.NewSingleHostReverseProxy(secondUrlAdr)
	proxy.ServeHTTP(w, r)

	counter--
}
