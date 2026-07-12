package httputil

import (
	"errors"
	"fmt"
	"net"
	"net/url"
	"strings"
)

// ErrUnsafeURL indicates that a URL resolves to a private or internal host.
var ErrUnsafeURL = errors.New("url resolves to a private or internal host")

// ValidateWebhookURL checks that rawURL is a valid webhook callback URL.
// Non-HTTP(S) schemes are always rejected. When hardened is true, loopback,
// link-local, and private IP addresses (including hostnames that resolve to
// them) are also rejected so that webhooks cannot be abused for SSRF against
// internal services.
func ValidateWebhookURL(rawURL string, hardened bool) error {
	if strings.TrimSpace(rawURL) == "" {
		return errors.New("webhook url is required")
	}

	u, err := url.Parse(rawURL)
	if err != nil {
		return fmt.Errorf("invalid webhook url: %w", err)
	}
	if u.Scheme != "http" && u.Scheme != "https" {
		return fmt.Errorf("webhook url scheme %q is not allowed; use http or https", u.Scheme)
	}
	if u.Host == "" {
		return errors.New("webhook url host is required")
	}

	if !hardened {
		return nil
	}

	host := u.Hostname()
	if host == "" {
		return errors.New("webhook url host is required")
	}

	ip := net.ParseIP(host)
	if ip != nil {
		if isPrivateOrInternalIP(ip) {
			return ErrUnsafeURL
		}
		return nil
	}

	addrs, err := net.LookupHost(host)
	if err != nil {
		return fmt.Errorf("webhook url host lookup failed: %w", err)
	}
	for _, addr := range addrs {
		ip := net.ParseIP(addr)
		if ip != nil && isPrivateOrInternalIP(ip) {
			return ErrUnsafeURL
		}
	}
	return nil
}

func isPrivateOrInternalIP(ip net.IP) bool {
	return ip.IsLoopback() || ip.IsLinkLocalUnicast() || ip.IsLinkLocalMulticast() || ip.IsPrivate()
}
