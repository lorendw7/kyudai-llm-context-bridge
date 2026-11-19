package main

import (
	"context"
	"io"
	"net"
	"net/http"
	"strings"
	"testing"
	"time"
)

func TestParseListenAddressAcceptsLoopback(t *testing.T) {
	tests := []struct {
		name string
		raw  string
		want string
	}{
		{"unset_uses_default", "", "127.0.0.1:8080"},
		{"ipv4_loopback", "127.0.0.1:9090", "127.0.0.1:9090"},
		{"other_ipv4_loopback", "127.0.0.2:8080", "127.0.0.2:8080"},
		{"ipv6_loopback", "[::1]:8080", "[::1]:8080"},
		{"lowest_port", "127.0.0.1:1", "127.0.0.1:1"},
		{"highest_port", "127.0.0.1:65535", "127.0.0.1:65535"},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			got, err := parseListenAddress(test.raw)
			if err != nil {
				t.Fatalf("parseListenAddress(%q) error = %v", test.raw, err)
			}

			if got != test.want {
				t.Fatalf("parseListenAddress(%q) = %q, want %q",
					test.raw, got, test.want)
			}
		})
	}
}

func TestParseListenAddressRejectsInvalidValues(t *testing.T) {
	tests := []struct {
		name string
		raw  string
	}{
		{"all_interfaces_ipv4", "0.0.0.0:8080"},
		{"all_interfaces_ipv6", "[::]:8080"},
		{"empty_host", ":8080"},
		{"private_network", "192.168.1.10:8080"},
		{"public_address", "203.0.113.5:8080"},
		{"hostname", "localhost:8080"},
		{"missing_port", "127.0.0.1"},
		{"empty_port", "127.0.0.1:"},
		{"port_zero", "127.0.0.1:0"},
		{"port_too_large", "127.0.0.1:65536"},
		{"negative_port", "127.0.0.1:-1"},
		{"named_port", "127.0.0.1:http"},
		{"signed_port", "127.0.0.1:+8080"},
		{"padded_port", "127.0.0.1:08080"},
		{"leading_space", " 127.0.0.1:8080"},
		{"trailing_space", "127.0.0.1:8080 "},
		{"url", "http://127.0.0.1:8080"},
		{"unbracketed_ipv6", "::1:8080"},
		{"space_only", " "},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			got, err := parseListenAddress(test.raw)
			if err == nil {
				t.Fatalf("parseListenAddress(%q) = %q, want error",
					test.raw, got)
			}

			if got != "" {
				t.Fatalf("parseListenAddress(%q) returned %q with an error",
					test.raw, got)
			}
		})
	}
}

func TestListenAddressFromEnvironment(t *testing.T) {
	t.Setenv(listenAddressVariable, "127.0.0.1:18080")

	got, err := listenAddressFromEnvironment()
	if err != nil || got != "127.0.0.1:18080" {
		t.Fatalf("configured address = %q, %v", got, err)
	}

	t.Setenv(listenAddressVariable, "0.0.0.0:18080")

	if got, err := listenAddressFromEnvironment(); err == nil {
		t.Fatalf("non-loopback address accepted as %q", got)
	}

	t.Setenv(listenAddressVariable, "")

	got, err = listenAddressFromEnvironment()
	if err != nil || got != defaultListenAddress {
		t.Fatalf("default address = %q, %v", got, err)
	}
}

func TestNewServerSetsEveryTimeout(t *testing.T) {
	handler := http.NewServeMux()
	server := newServer(handler)

	if server.Handler != handler {
		t.Fatal("server does not use the supplied handler")
	}

	limits := map[string]time.Duration{
		"ReadHeaderTimeout": server.ReadHeaderTimeout,
		"ReadTimeout":       server.ReadTimeout,
		"WriteTimeout":      server.WriteTimeout,
		"IdleTimeout":       server.IdleTimeout,
	}

	for name, value := range limits {
		if value <= 0 {
			t.Errorf("%s = %v, want a positive limit", name, value)
		}
	}

	if server.MaxHeaderBytes <= 0 {
		t.Errorf("MaxHeaderBytes = %d, want a positive limit",
			server.MaxHeaderBytes)
	}

	if server.Addr != "" {
		t.Errorf("Addr = %q; the listener, not the server, selects the address",
			server.Addr)
	}
}

// startTestServer runs serve on an ephemeral loopback port and returns the
// base URL, the function that requests a stop, and the channel carrying
// serve's result.
func startTestServer(
	t *testing.T,
	handler http.Handler,
	timeout time.Duration,
) (string, context.CancelFunc, <-chan error) {
	t.Helper()

	listener, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatalf("listen: %v", err)
	}

	ctx, stop := context.WithCancel(context.Background())
	result := make(chan error, 1)

	go func() {
		result <- serve(ctx, newServer(handler), listener, timeout)
	}()

	t.Cleanup(stop)

	return "http://" + listener.Addr().String(), stop, result
}

func newTestClient(t *testing.T) *http.Client {
	t.Helper()

	transport := &http.Transport{Proxy: nil, DisableKeepAlives: true}
	t.Cleanup(transport.CloseIdleConnections)

	return &http.Client{Transport: transport, Timeout: 5 * time.Second}
}

func waitForServeResult(t *testing.T, result <-chan error) error {
	t.Helper()

	select {
	case err := <-result:
		return err
	case <-time.After(5 * time.Second):
		t.Fatal("serve did not return")
		return nil
	}
}

func TestServeAnswersThenStopsCleanly(t *testing.T) {
	base, stop, result := startTestServer(t, newHandler(), time.Second)
	client := newTestClient(t)

	response, err := client.Get(base + "/health")
	if err != nil {
		t.Fatalf("health request: %v", err)
	}

	body, _ := io.ReadAll(response.Body)
	response.Body.Close()

	if response.StatusCode != http.StatusOK ||
		strings.TrimSpace(string(body)) != `{"status":"ok","service":"context-bridge"}` {
		t.Fatalf("health = %d %q", response.StatusCode, body)
	}

	stop()

	if err := waitForServeResult(t, result); err != nil {
		t.Fatalf("serve after a clean stop = %v, want nil", err)
	}

	if _, err := client.Get(base + "/health"); err == nil {
		t.Fatal("server still accepts connections after serve returned")
	}
}

func TestServeLetsInFlightRequestFinish(t *testing.T) {
	entered := make(chan struct{})
	release := make(chan struct{})

	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		close(entered)
		<-release
		_, _ = io.WriteString(w, "finished")
	})

	base, stop, result := startTestServer(t, handler, 5*time.Second)
	client := newTestClient(t)

	type outcome struct {
		body string
		err  error
	}

	responses := make(chan outcome, 1)

	go func() {
		response, err := client.Get(base + "/slow")
		if err != nil {
			responses <- outcome{err: err}
			return
		}

		defer response.Body.Close()
		body, err := io.ReadAll(response.Body)
		responses <- outcome{body: string(body), err: err}
	}()

	<-entered
	stop()

	// The stop request alone must not end serve while a request is running.
	select {
	case err := <-result:
		t.Fatalf("serve returned %v while a request was in flight", err)
	case <-time.After(100 * time.Millisecond):
	}

	close(release)

	got := <-responses
	if got.err != nil || got.body != "finished" {
		t.Fatalf("in-flight response = %q, %v", got.body, got.err)
	}

	if err := waitForServeResult(t, result); err != nil {
		t.Fatalf("serve after draining = %v, want nil", err)
	}
}

func TestServeReportsRequestsOutlastingGracePeriod(t *testing.T) {
	entered := make(chan struct{})
	release := make(chan struct{})
	defer close(release)

	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		close(entered)

		select {
		case <-release:
		case <-r.Context().Done():
		}
	})

	base, stop, result := startTestServer(t, handler, 50*time.Millisecond)
	client := newTestClient(t)

	go func() {
		response, err := client.Get(base + "/stuck")
		if err == nil {
			response.Body.Close()
		}
	}()

	<-entered
	stop()

	if err := waitForServeResult(t, result); err == nil {
		t.Fatal("serve reported a clean stop although a request was cut off")
	}

	if _, err := client.Get(base + "/health"); err == nil {
		t.Fatal("server still accepts connections after a forced stop")
	}
}

func TestServeReturnsWhenListenerFails(t *testing.T) {
	listener, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatalf("listen: %v", err)
	}

	listener.Close()

	result := make(chan error, 1)

	go func() {
		result <- serve(context.Background(), newServer(newHandler()),
			listener, time.Second)
	}()

	if err := waitForServeResult(t, result); err == nil {
		t.Fatal("serve returned nil for a listener that cannot accept")
	}
}
