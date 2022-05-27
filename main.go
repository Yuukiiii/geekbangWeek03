package main

import (
	"context"
	"fmt"
	"golang.org/x/sync/errgroup"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

type HelloWorldHandler struct{}

func (h HelloWorldHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	_, _ = w.Write([]byte("Hello World!"))
}

func main() {
	server := http.Server{
		Addr:    ":8080",
		Handler: HelloWorldHandler{},
	}
	g, _ := errgroup.WithContext(context.Background())
	g.Go(func() error {
		return server.ListenAndServe()
	})
	g.Go(func() error {
		sigChan := make(chan os.Signal)
		defer close(sigChan)
		// 这里以 sigterm 为例，实际中可以根据不同的信号做不同的处理
		signal.Notify(sigChan, syscall.SIGTERM)
		_ = <-sigChan
		ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
		defer cancel()
		// 需要保证服务器优雅终止，但也不能超过一定时间
		return server.Shutdown(ctx)
	})
	err := g.Wait()
	fmt.Println(err)
}
