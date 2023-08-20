package tcp

import (
	"context"
	"io"
	"net"
	"os"
	"os/signal"
	"sync"
	"syscall"

	"github.com/sunflower10086/go-redis/interface/tcp"
	"github.com/sunflower10086/go-redis/lib/logger"
)

type Config struct {
	Address string
}

func ListenAndServeWithSignal(
	conf *Config,
	handler tcp.Handler) error {

	closeChan := make(chan struct{})
	sigChan := make(chan os.Signal)
	signal.Notify(sigChan, syscall.SIGHUP, syscall.SIGQUIT, syscall.SIGTERM, syscall.SIGINT)
	go func() {
		<-sigChan
		closeChan <- struct{}{}
	}()

	listener, err := net.Listen("tcp", conf.Address)
	if err != nil {
		return err
	}
	logger.Info("start listen")

	ListenAndServe(listener, handler, closeChan)
	return nil
}

func ListenAndServe(
	listener net.Listener,
	handler tcp.Handler,
	closeChan <-chan struct{}) {

	go func() {
		<-closeChan
		logger.Info("shutting down")
		_ = listener.Close()
		_ = handler.Close()
	}()

	defer func() {
		_ = listener.Close()
		_ = handler.Close()
	}()

	var wg sync.WaitGroup
	for {
		conn, err := listener.Accept()
		if err != nil {
			if err == io.EOF {
				logger.Info("Server close")
			} else {
				logger.Error(err)
			}
			break
		}

		logger.Info("accepted link")
		ctx := context.Background()
		wg.Add(1)
		go func() {
			defer wg.Done()
			handler.Handle(ctx, conn)
		}()
	}

	wg.Wait()
}
