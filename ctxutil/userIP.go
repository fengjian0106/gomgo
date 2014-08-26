package ctxutil

import (
	"fmt"
	"net"
	"net/http"
	"strings"

	"code.google.com/p/go.net/context"
)

// request and associating it with a Context.

// FromRequest extracts the user IP address from req, if present.
func IpFromRequest(req *http.Request) (net.IP, error) {
	s := strings.SplitN(req.RemoteAddr, ":", 2)
	userIP := net.ParseIP(s[0])
	if userIP == nil {
		return nil, fmt.Errorf("userip: %q is not IP:port", req.RemoteAddr)
	}
	return userIP, nil
}

// userIPkey is the context key for the user IP address.
var userIPKey context.Key = context.NewKey("context.userIPKey")

// NewContext returns a new Context carrying userIP.
func NewContextWithIp(ctx context.Context, userIP net.IP) context.Context {
	return context.WithValue(ctx, userIPKey, userIP)
}

// FromContext extracts the user IP address from ctx, if present.
func IpFromContext(ctx context.Context) (net.IP, bool) {
	// ctx.Value returns nil if ctx has no value for the key;
	// the net.IP type assertion returns ok=false for nil.
	userIP, ok := ctx.Value(userIPKey).(net.IP)
	return userIP, ok
}
