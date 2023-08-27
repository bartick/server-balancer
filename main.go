package main

import (
	// "crypto/tls"
	"fmt"
	"net/http"
	"net/http/httputil"
	"net/url"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	// "golang.org/x/crypto/acme/autocert"
)

// func redirectToTls(w http.ResponseWriter, r *http.Request) {
// 	http.Redirect(w, r, "https://"+r.Host+":443"+r.RequestURI, http.StatusMovedPermanently)
// }

func main() {
	// certManager := autocert.Manager{
	// 	Prompt:     autocert.AcceptTOS,
	// 	HostPolicy: autocert.HostWhitelist(),   //Your domain here
	// 	Cache:      autocert.DirCache("certs"), //Folder for storing certificates
	// }

	re := httputil.NewSingleHostReverseProxy(&url.URL{Scheme: "http", Host: "localhost:8080"})

	r := chi.NewRouter()
	r.Use(middleware.CleanPath)
	r.Use(middleware.Logger)
	r.Use(middleware.AllowContentEncoding("deflate", "gzip"))
	r.Use(middleware.AllowContentType("application/json", "text/xml"))
	r.Get("/*", func(w http.ResponseWriter, r *http.Request) {
		fmt.Println(r.Host)
		re.ServeHTTP(w, r)
		// w.Write([]byte(r.RequestURI))
	})

	// server := http.Server{
	// 	Addr: ":https",
	// 	TLSConfig: &tls.Config{
	// 		GetCertificate: certManager.GetCertificate,
	// 		MinVersion:     tls.VersionTLS12,
	// 	},
	// 	Handler: r,
	// }

	http.ListenAndServe("0.0.0.0:8000", r)
	// if err := http.ListenAndServe(":80", http.HandlerFunc(redirectToTls)); err != nil {
	// 	log.Fatalf("ListenAndServe error: %v", err)
	// }
	// server.ListenAndServeTLS("", "") //Key and cert are coming from Let's Encrypt
}
