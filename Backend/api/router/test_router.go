package router

import (
	"github.com/DeoEsor/GameAdaptive-MLOp/Backend/api/domain"
	"github.com/DeoEsor/GameAdaptive-MLOp/Backend/api/middleware"
	"github.com/ggicci/httpin"
	"github.com/go-chi/chi"
)

func newTestRouter(parentRoute chi.Router, errHandler middleware.ErrHandler) {
	parentRoute.With(httpin.NewInput(domain.SendTestTimeRequest{})).Post("/SendTestTime",
		errHandler.ErrMiddleware(domain.SendTestTime),
	)
}
