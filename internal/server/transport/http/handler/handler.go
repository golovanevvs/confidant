package handler

import (
	"context"

	"github.com/go-chi/chi"
	"github.com/golovanevvs/confidant/internal/server/model"
	"github.com/golovanevvs/confidant/internal/server/transport/http/logger"
	"go.uber.org/zap"
)

type IAccountService interface {
	CreateAccount(ctx context.Context, account model.Account) (int, error)
	BuildJWTString(ctx context.Context, accountID int) (string, error)
}

type Service struct {
	IAccountService
	// IAccountService
	//IMyService
	//IYourService
}

type handler struct {
	sv *Service
	lg *zap.SugaredLogger
}

// New - the handler constructor
func New(sv *Service, lg *zap.SugaredLogger) *handler {
	return &handler{
		sv: sv,
		lg: lg,
	}
}

// InitRoutes - request routing, used as http.Handler when starting the server
func (hd *handler) InitRoutes() *chi.Mux {
	// creating a router instance
	rt := chi.NewRouter()

	// using middleware
	// logging
	rt.Use(logger.WithLogging(hd.lg))

	// routes
	rt.Post("/api/register", hd.accountRegisterPost)

	// rt.Route("/api/user", func(r chi.Router) {
	// r.Post("/register", hd.accountRegister)
	// 	r.Post("/login", hd.userLogin)
	// 	r.With(hd.authByJWT).Post("/orders", hd.userUploadOrder)
	// 	r.With(hd.authByJWT).Get("/orders", hd.getOrders)
	// 	r.With(hd.authByJWT).Get("/withdrawals", hd.withDrawals)
	// 	r.With(hd.authByJWT).Route("/balance", func(r chi.Router) {
	// 		r.Get("/", hd.getBalance)
	// 		r.Post("/withdraw", hd.withDraw)
	// 	})
	// })

	return rt
}
