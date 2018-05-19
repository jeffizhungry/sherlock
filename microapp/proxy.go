package microapp

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"regexp"

	"github.com/elazarl/goproxy"
	"github.com/jeffizhungry/sherlock/pkg/rawhttp"
)

func drainBody(b io.ReadCloser) (r1, r2 io.ReadCloser, err error) {
	if b == http.NoBody {
		// No copying needed. Preserve the magic sentinel meaning of NoBody.
		return http.NoBody, http.NoBody, nil
	}
	var buf bytes.Buffer
	if _, err = buf.ReadFrom(b); err != nil {
		return nil, b, err
	}
	if err = b.Close(); err != nil {
		return nil, b, err
	}
	return ioutil.NopCloser(&buf), ioutil.NopCloser(bytes.NewReader(buf.Bytes())), nil
}

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

		// Copy request
		copyReq := *req
		copyReq.Body, req.Body, err = drainBody(req.Body)
		if err != nil {
			panic(err)
		}

		// Copy response
		copyResp := *resp
		copyResp.Body, resp.Body, err = drainBody(resp.Body)
		if err != nil {
			panic(err)
		}

		subscription <- rawhttp.Pair{
			Request:  &copyReq,
			Response: &copyResp,
		}

		// Pass through
		return req, resp
	})

	// Start proxy server
	fmt.Println("HTTP Transparent Proxy. Listening on localhost:" + port)
	log.Fatal(http.ListenAndServe("localhost:"+port, proxy))
}
