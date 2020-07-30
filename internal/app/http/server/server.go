package server

import (
	"bytes"
	"context"
	"github.com/alonelegion/service_template/internal/app/config"
	"go.uber.org/zap"
	"io/ioutil"
	"net/http"
	"strconv"

	"github.com/chapsuk/wait"
	"github.com/julienschmidt/httprouter"

	"go.elastic.co/apm/module/apmhttp"
)

type (
	Server struct {
		ctx context.Context

		logger *zap.Logger
		config *config.AppConfig
	}

	panicHandler struct {
		logger *zap.Logger

		handler http.Handler
	}
)

func (h panicHandler) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	// Query проверяет правильность формирования пар значений в запросе
	urlQuery := req.URL.Query()
	// Считываем все данные
	data, err := ioutil.ReadAll(req.Body)
	if err != nil {
		h.logger.Error("error while read body for recover", zap.Error(err))

		w.WriteHeader(http.StatusNoContent)
		return
	}

	defer func() {
		if err := recover(); err != nil {
			h.logger.Error(
				"http request recovered",
				zap.Any("error", err),
				zap.Any("query", urlQuery),
				zap.ByteString("body", data),
			)
		}
	}()

	req.Body = ioutil.NopCloser(bytes.NewReader(data))
	h.handler.ServeHTTP(w, req)
}

func NewServer(
	ctx context.Context,
	logger *zap.Logger,
	config *config.AppConfig,
) <-chan error {
	return Register(func() error {
		return (&Server{
			ctx: ctx,

			logger: logger,
			config: config,
		}).Start(ctx)
	})
}

func (as *Server) Start(ctx context.Context) error {
	server := http.Server{
		Addr:    ":" + strconv.Itoa(as.config.HTTP.Port),
		Handler: as.handlers(),
	}

	hf := func() error {
		return server.ListenAndServe()
	}

	as.logger.Info("http server is running", zap.String("host", as.config.HTTP.Host),
		zap.Int("port", as.config.HTTP.Port))
	select {
	case err := <-Register(hf):
		as.logger.Info("Shutdown http server", zap.String("by", "error"), zap.Error(err))
		return server.Shutdown(ctx)
	case <-ctx.Done():
		as.logger.Info("Shutdown http server", zap.String("by", "context.Done"))
		return server.Shutdown(ctx)
	}
}

func (as *Server) ph(handler http.HandlerFunc) http.Handler {
	return apmhttp.Wrap(panicHandler{
		logger:  as.logger,
		handler: handler,
	})
}

func (as *Server) handlers() *httprouter.Router {
	return func() *httprouter.Router {
		router := httprouter.New()

		router.Handler(http.MethodGet, "/", as.ph(hw))

		return router
	}()
}

func Register(fn func() error) <-chan error {
	ch := make(chan error)
	wg := wait.Group{}
	wg.Add(func() {
		ch <- fn()
	})

	go func() {
		wg.Wait()
		close(ch)
	}()

	return ch
}

func hw(w http.ResponseWriter, _ *http.Request) {
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write([]byte("Hello World!"))
}

func check(w http.ResponseWriter, _ *http.Request) {
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write([]byte("OK"))
}
