package webserver

import (
	"context"
	"fmt"
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

type Server struct {
	blocker chan struct{}
	closer  chan struct{}
	addr    net.Addr
}

func (s *Server) URL() string {
	return "http://" + s.addr.String()
}

func (s *Server) Block() {
	<-s.blocker
}

func (s *Server) AsyncClose() {
	close(s.closer)
}

func (s *Server) Close() {
	close(s.closer)
	<-s.blocker
}

func Port(port int) optionSetter {
	return func(o *option) {
		o.Port = port
	}
}

func ShutdownTimeout(t time.Duration) optionSetter {
	return func(o *option) {
		o.ShutdownTimeout = t
	}
}

type option struct {
	Port            int
	ShutdownTimeout time.Duration
}

type optionSetter func(*option)

func StartHTTPServer(ctx context.Context, router http.Handler, options ...optionSetter) (*Server, error) {

	opt := &option{
		Port:            0,
		ShutdownTimeout: 10 * time.Second,
	}
	for _, o := range options {
		o(opt)
	}

	closeNotifier := make(chan struct{})
	blocker := make(chan struct{})

	srv := &http.Server{
		Handler: router,
	}

	listener, err := net.Listen("tcp", fmt.Sprintf(":%d", opt.Port))
	if err != nil {
		return nil, err
	}

	go func() {
		if err := srv.Serve(listener); err != nil && err != http.ErrServerClosed {
			log.Fatalf("listen: %s\n", err)
		}
	}()

	go func() {
		// Wait for interrupt signal to gracefully shutdown the server
		quit := make(chan os.Signal, 1)
		signal.Notify(quit, os.Interrupt, syscall.SIGTERM)
		signal.Notify(quit, syscall.SIGQUIT, syscall.SIGTERM, syscall.SIGINT)

		select {
		case <-quit:
		case <-closeNotifier:
		case <-ctx.Done():
		}

		ctx, cancel := context.WithTimeout(context.Background(), opt.ShutdownTimeout)
		defer cancel()
		if err := srv.Shutdown(ctx); err != nil {
			log.Fatal("Server Shutdown:", err)
		}
		log.Println("Server exited")

		close(blocker)
	}()

	return &Server{
		closer:  closeNotifier,
		blocker: blocker,
		addr:    listener.Addr(),
	}, nil
}
