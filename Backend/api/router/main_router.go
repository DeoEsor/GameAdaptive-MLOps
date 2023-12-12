package router

import (
	"time"

	custommiddleware "github.com/DeoEsor/GameAdaptive-MLOp/Backend/api/middleware"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
)

const (
	defaultTimeout    = 2 * time.Minute
	swaggerUrlPattern = "http://%s:%s/docs/doc.json"
)

func NewRouter(
	errHandler custommiddleware.ErrHandler,
) chi.Router {
	root := chi.NewRouter()

	// useful log info
	root.Use(middleware.RequestID)
	root.Use(middleware.Logger)
	root.Use(middleware.Recoverer)

	// автоматический таймаут для всех запросов
	root.Use(middleware.Timeout(defaultTimeout))

	// обязательный mount на все. Нужен для создания запросов на /api
	// Например:
	// root.Mount("/api", newCusmomRouter(errHandler))
	root.Mount("/test", newTestRouter(errHandler))

	rootApi := chi.NewRouter()

	rootApi.Mount("/api", root)
	return rootApi
}
