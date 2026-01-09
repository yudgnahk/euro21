package sofascore

import (
	"fmt"
	"net"
	"net/http"
	"time"

	"golang.org/x/net/proxy"
)

// ProxyConfig holds proxy configuration
type ProxyConfig struct {
	Enabled        bool
	Type           string // "socks5" or "tor"
	Address        string // e.g., "localhost:9050"
	TorControlAddr string // e.g., "localhost:9051" for circuit rotation
}

// NewClientWithProxy creates a client with proxy support
func NewClientWithProxy(proxyConfig ProxyConfig) (*Client, error) {
	client := &Client{
		baseURL:   BaseURL,
		userAgent: "Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/132.0.0.0 Safari/537.36",
	}

	if !proxyConfig.Enabled {
		client.httpClient = &http.Client{
			Timeout: 15 * time.Second,
		}
		return client, nil
	}

	// Create SOCKS5 dialer
	dialer, err := proxy.SOCKS5("tcp", proxyConfig.Address, nil, proxy.Direct)
	if err != nil {
		return nil, fmt.Errorf("failed to create SOCKS5 dialer: %w", err)
	}

	// Create HTTP transport with proxy
	transport := &http.Transport{
		Dial:                dialer.Dial,
		MaxIdleConns:        10,
		IdleConnTimeout:     30 * time.Second,
		TLSHandshakeTimeout: 10 * time.Second,
	}

	client.httpClient = &http.Client{
		Transport: transport,
		Timeout:   30 * time.Second,
	}

	return client, nil
}

// RotateTorCircuit sends NEWNYM signal to Tor control port to get a new circuit
func RotateTorCircuit(controlAddr string) error {
	conn, err := net.DialTimeout("tcp", controlAddr, 5*time.Second)
	if err != nil {
		return fmt.Errorf("failed to connect to Tor control port: %w", err)
	}
	defer func() {
		if cerr := conn.Close(); cerr != nil {
			// Connection close error is not critical here
			_ = cerr
		}
	}()

	// Authenticate (no password by default)
	if _, err := conn.Write([]byte("AUTHENTICATE\r\n")); err != nil {
		return fmt.Errorf("failed to authenticate: %w", err)
	}

	// Send NEWNYM signal
	if _, err := conn.Write([]byte("SIGNAL NEWNYM\r\n")); err != nil {
		return fmt.Errorf("failed to send NEWNYM signal: %w", err)
	}

	// Wait for new circuit to establish
	time.Sleep(2 * time.Second)

	return nil
}
