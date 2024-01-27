package server

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"go.uber.org/zap"
	"golang.org/x/sync/errgroup"
)

// A Server defines parameters for running an HTTP server.
type Server struct {
	Addr               string
	SecureAddr         string
	Cert               string
	Key                string
	Host               string
	ReadTimeout        time.Duration
	WriteTimeout       time.Duration
	TerminationTimeout time.Duration
	Handler            http.Handler
}

// NewServer returns an http server
func NewServer(
	handler http.Handler,
	ipaddr, port, sport, scert, skey, shost string,
	readTimeout, writeTimeout, terminationTimeout time.Duration,
) *Server {
	return &Server{
		Addr:               fmt.Sprintf("%s:%s", ipaddr, port),
		SecureAddr:         fmt.Sprintf("%s:%s", ipaddr, sport),
		Cert:               scert,
		Key:                skey,
		Host:               shost,
		Handler:            handler,
		ReadTimeout:        readTimeout,
		WriteTimeout:       writeTimeout,
		TerminationTimeout: terminationTimeout,
	}
}

// ListenAndServe initializes a server to respond to HTTP network requests.
func (s Server) ListenAndServe(ctx context.Context) error {
	if s.Key != "" {
		// TODO: 对证书进行检查
		return s.listenAndServeTLS(ctx)
	}
	return s.listenAndServe(ctx)
}

func (s Server) listenAndServe(ctx context.Context) error {
	g, ctx := errgroup.WithContext(ctx)
	s1 := &http.Server{
		Addr:         s.Addr,
		Handler:      s.Handler,
		ReadTimeout:  s.ReadTimeout,
		WriteTimeout: s.WriteTimeout,
	}
	g.Go(func() error {
		// 表示接收到信号
		<-ctx.Done()
		zap.L().Named("default").WithOptions(zap.AddCallerSkip(-1)).Error("shutting down")
		// 一定要使用一个新的ctx
		ctxShutdown, cancelFunc := context.WithTimeout(context.Background(), s.TerminationTimeout)
		defer cancelFunc()
		return s1.Shutdown(ctxShutdown)
	})
	g.Go(func() error {
		return s1.ListenAndServe()
	})
	return g.Wait()
}

func (s Server) listenAndServeTLS(ctx context.Context) error {
	g, ctx := errgroup.WithContext(ctx)
	s1 := &http.Server{
		Addr:    s.Addr,
		Handler: http.HandlerFunc(redirect),
	}
	s2 := &http.Server{
		Addr:         s.SecureAddr,
		Handler:      s.Handler,
		ReadTimeout:  s.ReadTimeout,
		WriteTimeout: s.WriteTimeout,
	}
	g.Go(func() error {
		return s1.ListenAndServe()
	})
	g.Go(func() error {
		return s2.ListenAndServeTLS(
			s.Cert,
			s.Key,
		)
	})
	g.Go(func() error {
		// 表示接收到信号
		<-ctx.Done()
		var gShutdown errgroup.Group
		ctxShutdown, cancelFunc := context.WithTimeout(context.Background(), s.TerminationTimeout)
		defer cancelFunc()

		gShutdown.Go(func() error {
			zap.L().Named("default").WithOptions(zap.AddCallerSkip(-1)).Error("http shutting down")
			return s1.Shutdown(ctxShutdown)
		})
		gShutdown.Go(func() error {
			zap.L().Named("default").WithOptions(zap.AddCallerSkip(-1)).Error("https shutting down")
			return s2.Shutdown(ctxShutdown)
		})
		return gShutdown.Wait()
	})
	return g.Wait()
}

// redirect http to https with a 307 Temporary Redirect
func redirect(w http.ResponseWriter, req *http.Request) {
	target := "https://" + req.Host + req.URL.Path
	http.Redirect(w, req, target, http.StatusTemporaryRedirect)
}
