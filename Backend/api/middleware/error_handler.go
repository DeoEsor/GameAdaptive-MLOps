package middleware

import (
	"context"
	"encoding/json"
	"net/http"
	"time"

	"github.com/DeoEsor/GameAdaptive-MLOp/Backend/api/domain"
	customerrors "github.com/DeoEsor/GameAdaptive-MLOp/Backend/pkg/custom_errors"

	"github.com/ggicci/httpin"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
)

const contextDeadline = time.Minute * 5

// ErrHandler структура используется для того, чтобы возможно было хендлить ошибки из ручек
type ErrHandler struct {
	handlerCatcher domain.Handlers
}

// Конструктор для ErrHandler
func NewErrorHandler(handlerCatcher domain.Handlers) ErrHandler {
	return ErrHandler{
		handlerCatcher: handlerCatcher,
	}
}

// ErrMiddleware - функция-хендлер. Принимает в себя тип ручки, которая используется в хендлере
func (em ErrHandler) ErrMiddleware(handleType domain.HandlerType) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Method", string(handleType))
		logrus.Infof("method: %v", handleType)
		ctx, cancel := context.WithTimeout(r.Context(), contextDeadline)
		defer cancel()

		res, err := em.handleTypeSwitcher(ctx, r, handleType)
		if err != nil {
			switch v := err.(type) {
			case customerrors.ErrCodes:
				w.WriteHeader(v.Code)
			default:
				w.WriteHeader(http.StatusInternalServerError)
			}
			err = errors.Wrapf(err, "Method %s ->", string(handleType))
			logrus.Errorln(err.Error())

			return
		}
		em.checkExclusiveFiles(res, w, handleType)
	}
}

// checkExclusiveFiles функция делает специфическую отправку файлов (не json), если того требует логика
func (em ErrHandler) checkExclusiveFiles(res interface{}, w http.ResponseWriter, handleType domain.HandlerType) {
	if res == nil {
		w.Write(nil)
		w.WriteHeader(http.StatusOK)
		return
	}

	toSend, err := json.Marshal(res)
	if err != nil {
		err = errors.Wrapf(err, "Method %s ->", string(handleType))
		w.WriteHeader(http.StatusInternalServerError)
		logrus.Errorln("unmarshal response error ", err.Error())
		return
	}
	w.Header().Add("Content-Type", "application/json")
	w.Write(toSend)

	w.WriteHeader(http.StatusOK)
}

func (em ErrHandler) handleTypeSwitcher(ctx context.Context, r *http.Request, handleType domain.HandlerType) (interface{}, error) {
	inputQuery := ctx.Value(httpin.Input)
	switch handleType {
	case domain.SendTestTime:
		if inputQuery == nil {
			return em.handlerCatcher.SendTestTime(ctx, nil)
		}
		return em.handlerCatcher.SendTestTime(ctx, inputQuery.(*domain.SendTestTimeRequest))
	case domain.SaveRaceTime:
		if inputQuery == nil {
			return em.handlerCatcher.SaveRaceTime(ctx, nil)
		}
		return em.handlerCatcher.SaveRaceTime(ctx, inputQuery.(*domain.SaveRaceRequest))
	}
	return nil, customerrors.ErrUnknownType
}
