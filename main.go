package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"regexp"

	"github.com/foomo/simplecert"
	"github.com/foomo/tlsconfig"

	"github.com/joho/godotenv"
)

func goDotEnvVariable(key string) string {

	// load .env file
	err := godotenv.Load(".env")

	if err != nil {
		log.Fatalf("Error loading .env file")
	}

	return os.Getenv(key)
}

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

func runServer(hosts *HostMap, handler http.Handler) {

	if os.Getenv("ENVIRONMENT") == "production" {

		hostMap := hosts.GetHostsArray()

		tlsConf := tlsconfig.NewServerTLSConfig(tlsconfig.TLSModeServerStrict)

		cfg := simplecert.Default
		cfg.Domains = hostMap
		cfg.CacheDir = "letsencrypt"
		cfg.SSLEmail = goDotEnvVariable("EMAIL")
		cfg.HTTPAddress = ""

		certReloader, err := simplecert.Init(cfg, func() {
			os.Exit(0)
		})
		if err != nil {
			log.Fatal("simplecert init failed: ", err)
		}

		server := http.Server{
			Addr:      ":https",
			TLSConfig: tlsConf,
			Handler:   handler,
		}

		fmt.Println("Starting server...")

		go func() {
			http.ListenAndServe(":http", http.HandlerFunc(simplecert.Redirect))
		}()

		tlsConf.GetCertificate = certReloader.GetCertificateFunc()

		log.Fatal(server.ListenAndServeTLS("", ""))

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
