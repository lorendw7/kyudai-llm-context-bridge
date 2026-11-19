package main

import (
	"context"
	"errors"
	"fmt"
	"net"
	"net/http"
	"os"
	"strconv"
	"time"
)

const (
	listenAddressVariable = "BRIDGE_LISTEN_ADDRESS"
	defaultListenAddress  = "127.0.0.1:8080"

	// Request bodies are small JSON documents and turn submission answers
	// with 202 before generation starts, so no handler needs a long exchange.
	readHeaderTimeout = 5 * time.Second
	readTimeout       = 15 * time.Second
	writeTimeout      = 30 * time.Second
	idleTimeout       = 60 * time.Second
	maxHeaderBytes    = 1 << 16

	// How long in-flight requests may finish after a stop signal.
	shutdownTimeout = 10 * time.Second
)

// listenAddressFromEnvironment returns the address to listen on. An unset
// variable selects the default. The bridge has no transport encryption and a
// single shared token, so only loopback IP addresses are accepted.
func listenAddressFromEnvironment() (string, error) {
	return parseListenAddress(os.Getenv(listenAddressVariable))
}

func parseListenAddress(raw string) (string, error) {
	if raw == "" {
		return defaultListenAddress, nil
	}

	host, port, err := net.SplitHostPort(raw)
	if err != nil {
		return "", errors.New("listen address must have the form host:port")
	}

	ip := net.ParseIP(host)
	if ip == nil || !ip.IsLoopback() {
		return "", errors.New("listen address must use a loopback IP address")
	}

	number, err := strconv.Atoi(port)
	if err != nil || number < 1 || number > 65535 || strconv.Itoa(number) != port {
		return "", errors.New("listen port must be a number from 1 to 65535")
	}

	return net.JoinHostPort(host, port), nil
}

func newServer(handler http.Handler) *http.Server {
	return &http.Server{
		Handler:           handler,
		ReadHeaderTimeout: readHeaderTimeout,
		ReadTimeout:       readTimeout,
		WriteTimeout:      writeTimeout,
		IdleTimeout:       idleTimeout,
		MaxHeaderBytes:    maxHeaderBytes,
	}
}

// serve accepts connections on listener until ctx is cancelled, then stops
// accepting and waits up to timeout for in-flight requests to finish. It
// returns nil after a clean stop. serve owns the goroutine it starts and does
// not return before that goroutine has ended.
func serve(
	ctx context.Context,
	server *http.Server,
	listener net.Listener,
	timeout time.Duration,
) error {
	serveResult := make(chan error, 1)

	go func() {
		serveResult <- server.Serve(listener)
	}()

	select {
	case err := <-serveResult:
		// Serve returned without a stop request, for example because the
		// listener failed.
		return err
	case <-ctx.Done():
	}

	shutdownContext, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	if err := server.Shutdown(shutdownContext); err != nil {
		// Requests are still running after the grace period. Close their
		// connections so the process can exit.
		_ = server.Close()
		<-serveResult

		return fmt.Errorf("shutdown did not complete: %w", err)
	}

	if err := <-serveResult; !errors.Is(err, http.ErrServerClosed) {
		return err
	}

	return nil
}
