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

	"github.com/alan-b-lima/ansi-escape-sequences"
	vterm "github.com/alan-b-lima/ansi-escape-sequences/vterminal"

	"github.com/alan-b-lima/prp/internal/api/v1"
)

func main() {
	if err := vterm.EnableVirtualTerminal(os.Stdout.Fd()); err != nil {
		fmt.Println(err)
		return
	}
	defer vterm.DisableVirtualTerminal(os.Stdout.Fd())

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
	fmt.Printf("Server listening at %s%s%s\n",
		ansi.FGColor(ansi.RGB{23, 135, 244}),
		ansi.HyperLinkP("http://"+addr),
		ansi.Reset(),
	)

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
	var pen ansi.Pen

	switch {
	case 500 <= status && status <= 599:
		pen.BGColor(ansi.RGB{219, 9, 9})
		pen.FGColor(ansi.RGB{235, 193, 193})

	case 400 <= status && status <= 499:
		pen.BGColor(ansi.RGB{219, 9, 9})
		pen.FGColor(ansi.RGB{235, 235, 235})

	case 300 <= status && status <= 399:
		pen.BGColor(ansi.RGB{59, 143, 222})
		pen.FGColor(ansi.RGB{235, 235, 235})

	case 200 <= status && status <= 299:
		pen.BGColor(ansi.RGB{59, 203, 91})
		pen.FGColor(ansi.RGB{235, 235, 235})

	default:
		pen.BGColor(ansi.RGB{255, 255, 255})
		pen.FGColor(ansi.RGB{0, 0, 0})
	}

	log.Printf("%s %03d %s %s %s %s",
		pen.Style(), status, ansi.Reset(),
		get_client_ip(r), r.Method, r.URL,
	)
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
