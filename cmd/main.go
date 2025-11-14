package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"

	"github.com/alan-b-lima/prp/internal/api/v1"

	"github.com/alan-b-lima/ansi-escape-sequences"
)

func main() {
	mux:= http.NewServeMux()

	mux.Handle("/", http.FileServer(http.Dir("../ui/web/")))
	mux.Handle("/api/", api.New())

	ln, err := net.Listen("tcp", ":4545")
	if err != nil {
		log.Println(err)
		return
	}

	url := "http://" + strings.Replace(ln.Addr().String(), "[::]", "localhost", 1)
	log.Printf("Server listening at %s\n", HyperLink(url))

	srv := http.Server{Handler: LogMiddleware(mux)}

	done := EnableGracefulShutdown(func() {
		log.Println("Closing...")
		srv.Shutdown(context.Background())
	})

	if err := srv.Serve(ln); err != nil && err != http.ErrServerClosed {
		log.Println(err)
	}

	<-done
}

func LogMiddleware(handler http.Handler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		rw := &responseWriter{w, 200}

		handler.ServeHTTP(rw, r)

		var pen ansi.Pen
		switch rw.StatusCode / 100 {
		case 5:
			pen = ServerError
		case 4:
			pen = ClientError
		case 3:
			pen = Redirect
		case 2:
			pen = Success

		default:
			pen.SetStyle(false)
		}

		var b strings.Builder
		pen.Writer = &b

		io.WriteString(&pen, fmt.Sprintf(" %03d ", rw.StatusCode))
		log.Printf("%s %s %s %s\n", b.String(), r.RemoteAddr, r.Method, r.URL)
	}
}

type responseWriter struct {
	http.ResponseWriter
	StatusCode int
}

func (rw *responseWriter) WriteHeader(status int) {
	rw.ResponseWriter.WriteHeader(status)
	rw.StatusCode = status
}

var Ansi bool

var (
	Success     ansi.Pen
	Redirect    ansi.Pen
	ClientError ansi.Pen
	ServerError ansi.Pen
)

func init() {
	Ansi = EnableVirtualTerminal()

	if Ansi {
		Success.BGColor(ansi.RGBFromHex(0x0ed145))
		Success.FGColor(ansi.RGBFromHex(0xffffff))

		Redirect.BGColor(ansi.RGBFromHex(0x4b53cc))
		Redirect.FGColor(ansi.RGBFromHex(0xffffff))

		ClientError.BGColor(ansi.RGBFromHex(0xea1d1d))
		ClientError.FGColor(ansi.RGBFromHex(0xffffff))

		ServerError.BGColor(ansi.RGBFromHex(0x88001b))
		ServerError.FGColor(ansi.RGBFromHex(0xffffff))
	}

	Success.SetStyle(Ansi)
	Redirect.SetStyle(Ansi)
	ClientError.SetStyle(Ansi)
	ServerError.SetStyle(Ansi)
}

func EnableVirtualTerminal() bool {
	return ansi.EnableVirtualTerminal(os.Stdout.Fd()) == nil
}

func EnableGracefulShutdown(fn func()) <-chan struct{} {
	signals := make(chan os.Signal, 1)
	signal.Notify(signals, syscall.SIGTERM, syscall.SIGINT)

	done := make(chan struct{}, 1)

	go func() {
		<-signals
		fn()
		done <- struct{}{}
	}()

	return done
}

func Pen() (pen ansi.Pen) {
	pen.SetStyle(Ansi)
	return pen
}

func HyperLink(link string) string {
	if !Ansi {
		return link
	}

	pen := Pen()
	pen.FGColor(ansi.RGBFromHex(0x4e8597))

	return pen.Sprint(ansi.HyperLinkP(link))
}
