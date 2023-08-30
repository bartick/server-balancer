package main

import (
	// "crypto/tls"
	"fmt"
	"net/http"
	"os"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	// "golang.org/x/crypto/acme/autocert"
)

func runServer(r *chi.Mux) {
	if os.Getenv("ENVIRONMENT") == "production" {
		fmt.Println("Starting server on port 80")
		http.ListenAndServe("0.0.0.0:80", r)
	} else {
		fmt.Println("Starting server on port 8000")
		// create a go function

		err := http.ListenAndServe("0.0.0.0:8000", r)
		if err != nil {
			fmt.Println(err)
		}
	}
}

// func redirectToTls(w http.ResponseWriter, r *http.Request) {
// 	http.Redirect(w, r, "https://"+strings.Split(r.Host, ":")[0]+r.RequestURI, http.StatusMovedPermanently)
// }

func main() {
	// certManager := autocert.Manager{
	// 	Prompt:     autocert.AcceptTOS,
	// 	HostPolicy: autocert.HostWhitelist(),   //Your domain here
	// 	Cache:      autocert.DirCache("certs"), //Folder for storing certificates
	// }

	jsonData := ReadJson()
	hosts := NewHostMap()
	hosts.SetAll(jsonData)

	r := chi.NewRouter()
	r.Use(middleware.CleanPath)
	// r.Use(middleware.Logger)
	r.Use(middleware.AllowContentEncoding("deflate", "gzip"))
	r.Use(middleware.AllowContentType("application/json", "text/xml"))
	r.Handle("/*", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		re, err := hosts.GetProxy(r.Host)
		if err != nil {
			w.Write([]byte("404 Proxy not found"))
			return
		}

		re.ServeHTTP(w, r)
	}))

	// server := http.Server{
	// 	Addr: ":https",
	// 	TLSConfig: &tls.Config{
	// 		GetCertificate: certManager.GetCertificate,
	// 		MinVersion:     tls.VersionTLS12,
	// 	},
	// 	Handler: r,
	// }

	runServer(r)

	// if err := http.ListenAndServe(":80", http.HandlerFunc(redirectToTls)); err != nil {
	// 	log.Fatalf("ListenAndServe error: %v", err)
	// }
	// server.ListenAndServeTLS("", "") //Key and cert are coming from Let's Encrypt

}
