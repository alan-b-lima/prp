package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"

	"github.com/alan-b-lima/prp/internal/api/v1"
)

func main() {
	addr := ":8080"
	if len(os.Args) >= 2 {
		addr = os.Args[1]
	}

	ln, err := net.Listen("tcp", addr)
	if err != nil {
		fmt.Println(err)
		return
	}

	srv := http.Server{Handler: LoggingMiddleware(api.NewAPIMux())}
	EnableSignalTermination(&srv)

	addr = strings.Replace(ln.Addr().String(), "[::]", "localhost", 1)
	fmt.Printf("Server listening at \033[38;2;23;135;244m%s\033[m\n", HyperLink("http://"+addr))
	if err := srv.Serve(ln); err != nil && err != http.ErrServerClosed {
		fmt.Println(err)
		return
	}
}

type _StatusCapture struct {
	http.ResponseWriter
	status int
}

func (w *_StatusCapture) WriteHeader(statusCode int) {
	w.ResponseWriter.WriteHeader(statusCode)
	w.status = statusCode
}

func LoggingMiddleware(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		resw := &_StatusCapture{ResponseWriter: w}
		h.ServeHTTP(resw, r)
		Log(resw.status, r)
	})
}

func HyperLink(link string) string {
	return fmt.Sprintf("\033]8;;%s\033\\%s\033]8;;\033\\", link, link)
}

func EnableSignalTermination(srv *http.Server) {
	signals := make(chan os.Signal, 1)
	signal.Notify(signals, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		<-signals

		done := make(chan int, 1)
		fmt.Println("Shutting the server down...")

		go func() {
			if err := srv.Shutdown(context.Background()); err != nil {
				srv.Close()
				done <- 1
			} else {
				done <- 0
			}
		}()

		select {
		case <-signals:
			fmt.Println("Killing the server...")
			os.Exit(1)

		case <-done:
			return
		}
	}()
}

func Log(status int, r *http.Request) {
	switch {
	case 500 <= status && status <= 599:
		log.Printf("\033[38;2;235;193;193;48;2;219;9;9m %03d \033[m %s %s %s", status, get_client_ip(r), r.Method, r.URL)

	case 400 <= status && status <= 499:
		log.Printf("\033[38;2;235;235;235;48;2;219;9;9m %03d \033[m %s %s %s", status, get_client_ip(r), r.Method, r.URL)

	case 300 <= status && status <= 399:
		log.Printf("\033[38;2;235;235;235;48;2;59;143;222m %03d \033[m %s %s %s", status, get_client_ip(r), r.Method, r.URL)

	case 200 <= status && status <= 299:
		log.Printf("\033[38;2;235;235;235;48;2;59;203;91m %03d \033[m %s %s %s", status, get_client_ip(r), r.Method, r.URL)

	default:
		log.Printf("\033[38;2;0;0;0;48;2;225;225;225m %03d \033[m %s %s %s", status, get_client_ip(r), r.Method, r.URL)
	}
}

func get_client_ip(r *http.Request) string {
	if xff := r.Header.Get("X-Forwarded-For"); xff != "" {
		ips := strings.Split(xff, ",")
		return strings.TrimSpace(ips[0])
	}

	ipPort := r.RemoteAddr
	splits := strings.Split(ipPort, ":")
	if len(splits) < 1 {
		return ipPort
	}

	return strings.Join(splits[:len(splits)-1], ":")
}
