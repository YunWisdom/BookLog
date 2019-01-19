

package util

import (
	"crypto/tls"
	"net"
	"net/http"
	"net/http/httputil"
	"net/url"

	"github.com/gin-gonic/gin"
)

// HacPaiURL is the URL of HacPai community.
const HacPaiURL = "https://hacpai.com"

// HacPaiAPI is a reverse proxy for https://hacpai.com.
func HacPaiAPI() gin.HandlerFunc {
	return func(c *gin.Context) {
		proxy := httputil.NewSingleHostReverseProxy(&url.URL{
			Scheme: "https",
			Host:   "hacpai.com",
		})

		proxy.Transport = &http.Transport{DialTLS: dialTLS}
		director := proxy.Director
		proxy.Director = func(req *http.Request) {
			director(req)
			req.Host = req.URL.Host
			req.URL.Path = req.URL.Path[len("api/hp/"):]
		}

		proxy.ServeHTTP(c.Writer, c.Request)
	}
}

func dialTLS(network, addr string) (net.Conn, error) {
	conn, err := net.Dial(network, addr)
	if err != nil {
		return nil, err
	}

	host, _, err := net.SplitHostPort(addr)
	if err != nil {
		return nil, err
	}
	cfg := &tls.Config{ServerName: host}

	tlsConn := tls.Client(conn, cfg)
	if err := tlsConn.Handshake(); err != nil {
		conn.Close()
		return nil, err
	}

	cs := tlsConn.ConnectionState()
	cert := cs.PeerCertificates[0]

	cert.VerifyHostname(host)

	return tlsConn, nil
}
