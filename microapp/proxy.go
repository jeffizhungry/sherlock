package microapp

import (
	"fmt"
	"log"
	"net/http"
	"regexp"

	"github.com/elazarl/goproxy"
	"github.com/jeffizhungry/sherlock/pkg/rawhttp"
)

// Proxy acts as an Transparent HTTP Proxy and forwards raw HTTP request response
// paris to it's subscription channel.
func Proxy(port string, subscription chan<- rawhttp.Pair) {

	// Init proxy
	proxy := goproxy.NewProxyHttpServer()

	// Setup proxy handler
	r := proxy.OnRequest(goproxy.ReqHostMatches(regexp.MustCompile("^.*$")))
	r.HandleConnect(goproxy.AlwaysMitm)
	r.DoFunc(func(req *http.Request, ctx *goproxy.ProxyCtx) (*http.Request, *http.Response) {
		resp, err := ctx.RoundTrip(req)
		if err != nil {
			panic(err)
		}
		subscription <- rawhttp.Pair{
			Request:  req,
			Response: resp,
		}
		return req, resp
	})

	// Start proxy server
	fmt.Println("HTTP Transparent Proxy. Listening on localhost:" + port)
	log.Fatal(http.ListenAndServe("localhost:"+port, proxy))
}
