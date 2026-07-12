package ipay88

import (
	"errors"
	"net"
	"net/url"

	"github.com/openmuara/openmuara/internal/errcode"
)

// ErrUnsafeURL is returned when a URL points to a private/internal host.
var ErrUnsafeURL = errors.New("ipay88: url resolves to a private or internal host")

// IsPublicURL reports whether rawURL is an http(s) URL that does not resolve
// to a loopback, link-local, or private IP address.
func IsPublicURL(rawURL string) error {
	u, err := url.Parse(rawURL)
	if err != nil {
		return errcode.Wrap(errcode.EInvalidRequest, "ipay88: invalid url", err)
	}
	if u.Scheme != "http" && u.Scheme != "https" {
		return errcode.New(errcode.EInvalidRequest, "ipay88: url scheme must be http or https")
	}

	host := u.Hostname()
	if host == "" {
		return errcode.New(errcode.EInvalidRequest, "ipay88: url host is required")
	}

	ip := net.ParseIP(host)
	if ip != nil {
		if isPrivateIP(ip) {
			return errcode.Wrap(errcode.EInvalidRequest, "ipay88: url resolves to a private or internal host", ErrUnsafeURL)
		}
		return nil
	}

	addrs, err := net.LookupHost(host)
	if err != nil {
		return errcode.Wrap(errcode.EInvalidRequest, "ipay88: failed to lookup host", err)
	}
	for _, addr := range addrs {
		ip := net.ParseIP(addr)
		if ip != nil && isPrivateIP(ip) {
			return errcode.Wrap(errcode.EInvalidRequest, "ipay88: url resolves to a private or internal host", ErrUnsafeURL)
		}
	}
	return nil
}

func isPrivateIP(ip net.IP) bool {
	return ip.IsLoopback() || ip.IsLinkLocalUnicast() || ip.IsLinkLocalMulticast() || ip.IsPrivate()
}
