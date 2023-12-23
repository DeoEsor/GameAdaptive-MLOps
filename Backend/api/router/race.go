package router

import (
	"github.com/DeoEsor/GameAdaptive-MLOp/Backend/api/domain"
	"github.com/DeoEsor/GameAdaptive-MLOp/Backend/api/middleware"
	"github.com/ggicci/httpin"
	"github.com/go-chi/chi"
)

func newRaceRouter(parentRoute chi.Router, errHandler middleware.ErrHandler) {
	parentRoute.With(httpin.NewInput(domain.SaveRaceRequest{})).Post("/SaveRaceTime",
		errHandler.ErrMiddleware(domain.SaveRaceTime),
	)
}
