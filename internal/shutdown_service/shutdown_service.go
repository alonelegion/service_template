package shutdown_service

import (
	"context"
	"os"
	"os/signal"
	"syscall"
)

func StopContext(c context.Context) context.Context {
	return cancelContextOnSignals(c, true, syscall.SIGINT, syscall.SIGTERM)
}

func ShutDownContext(context context.Context) context.Context {
	return cancelContextOnSignals(context, true, syscall.SIGINT, syscall.SIGTERM, syscall.SIGHUP)
}

func cancelContextOnSignals(c context.Context, osExitOnSecond bool, signals ...os.Signal) context.Context {
	ctx, cancel := context.WithCancel(c)
	done := make(chan struct{})
	go listenSignalsFunc(ctx, cancel, done, osExitOnSecond, signals...)()
	<-done
	return ctx
}

func listenSignalsFunc(
	ctx context.Context,
	cancel context.CancelFunc,
	done chan struct{},
	osExitOnSecond bool,
	signals ...os.Signal,
) func() {
	return func() {
		ch := make(chan os.Signal, 1)
		signal.Notify(ch, signals...)
		close(done)

		select {
		case <-ctx.Done():
			return
		case <-ch:
			cancel()
			if osExitOnSecond {
				<-ch
				os.Exit(1)
			}
		}
	}
}

// ReloadContext returns reload signals channel
func ReloadChannel(ctx context.Context) <-chan struct{} {
	done := make(chan struct{})
	sighups := make(chan struct{})
	go func() {
		ch := make(chan os.Signal, 1)
		signal.Notify(ch, syscall.SIGHUP)
		close(done)
		for {
			select {
			case <-ctx.Done():
				return
			case <-ch:
				sighups <- struct{}{}
			}
		}
	}()
	<-done
	return sighups
}
