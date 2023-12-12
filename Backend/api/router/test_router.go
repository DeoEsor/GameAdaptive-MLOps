package router

import (
	"github.com/DeoEsor/GameAdaptive-MLOp/Backend/api/domain"
	"github.com/DeoEsor/GameAdaptive-MLOp/Backend/api/middleware"
	"github.com/ggicci/httpin"
	"github.com/go-chi/chi"
)

func newTestRouter(errHandler middleware.ErrHandler) chi.Router {
	testRouter := chi.NewRouter()

	testRouter.With(httpin.NewInput(domain.SendTestTimeRequest{})).Post("/SendTestTime",
		errHandler.ErrMiddleware(domain.SendTestTime),
	)
	return testRouter
}
