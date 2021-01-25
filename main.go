
package main

import (
	"flag"
	"log"
	"net/http"
	"strings"

	"github.com/elazarl/goproxy"
)

func main() {
	verbose := flag.Bool("v", false, "should every proxy request be logged to stdout")
	addr := flag.String("a", ":8080", "proxy listen address")
	flag.Parse()
	setCA(caCert, caKey)

	proxy := goproxy.NewProxyHttpServer()
	proxy.OnRequest().HandleConnect(goproxy.AlwaysMitm)

	proxy.Verbose = *verbose

	proxy.OnRequest().DoFunc(
		func(r *http.Request, ctx *goproxy.ProxyCtx) (*http.Request, *http.Response) {

			//log.Printf("%s %s\n", r.Method, r.URL)
			if !strings.HasPrefix(r.Header.Get("Content-Type"), "text/x-gwt-rpc") {
				return r, nil
			}

			return r, nil
		})
	proxy.OnResponse().DoFunc(
		func(r *http.Response, ctx *goproxy.ProxyCtx) *http.Response {
			log.Printf("%d : %-6s %s (Server: %s)\n", r.StatusCode, r.Request.Method, r.Request.URL, r.Header.Get("Server"))
			return r
		})

	log.Printf("Listening on %s\n", *addr)
	log.Fatal(http.ListenAndServe(*addr, proxy))
}
