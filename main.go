package main

import (
	"fmt"
	"net/http"
	"os"
	"regexp"

	"golang.org/x/crypto/acme/autocert"
)

// func redirectToTls(w http.ResponseWriter, r *http.Request) {
// 	http.Redirect(w, r, "https://"+strings.Split(r.Host, ":")[0]+r.RequestURI, http.StatusMovedPermanently)
// }

type route struct {
	pattern *regexp.Regexp
	handler http.Handler
}

type RegexpHandler struct {
	routes []*route
}

func (h *RegexpHandler) HandleFunc(pattern *regexp.Regexp, handler func(http.ResponseWriter, *http.Request)) {
	h.routes = append(h.routes, &route{pattern, http.HandlerFunc(handler)})
}

func ProxyHandler(hosts *HostMap) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		fmt.Print(r.Method + " " + r.URL.String() + " \n")
		fmt.Println("Request from:" + r.RemoteAddr + "\n")

		re, err := hosts.GetProxy(r.Host)
		if err != nil {
			w.Write([]byte("404 Proxy not found"))
			return
		}

		re.ServeHTTP(w, r)
	}
}

func cacheDir() string {
	home, err := os.UserHomeDir()
	if err != nil {
		panic(err)
	}
	return home + "/.cache/certs"
}

func runServer(hosts *HostMap, handler http.Handler) {

	hostMap := hosts.GetHostsArray()

	if os.Getenv("ENVIRONMENT") == "production" {

		certManager := autocert.Manager{
			Prompt:     autocert.AcceptTOS,
			HostPolicy: autocert.HostWhitelist(hostMap...), //Your domain here
			Cache:      autocert.DirCache(cacheDir()),      //Folder for storing certificates
		}

		server := http.Server{
			Addr:      ":443",
			TLSConfig: certManager.TLSConfig(),
			Handler:   handler,
		}

		fmt.Println("Starting server...")

		go func() {
			http.ListenAndServe(":80", certManager.HTTPHandler(nil))
		}()

		server.ListenAndServeTLS("", "") //Key and cert are coming from Let's Encrypt

	} else {
		fmt.Println("Starting server on port 8000")

		err := http.ListenAndServe(":8000", handler)
		if err != nil {
			fmt.Println(err)
		}
	}
}

func main() {
	const anyPattern = `^/.*`

	jsonData := ReadJson()
	hosts := NewHostMap()
	hosts.SetAll(jsonData)

	routes := RegexpHandler{}
	routes.HandleFunc(regexp.MustCompile(anyPattern), ProxyHandler(hosts))

	runServer(hosts, http.HandlerFunc(ProxyHandler(hosts)))

	// if err := http.ListenAndServe(":80", http.HandlerFunc(redirectToTls)); err != nil {
	// 	log.Fatalf("ListenAndServe error: %v", err)
	// }

}
