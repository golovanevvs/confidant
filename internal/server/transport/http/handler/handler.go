package handler

import (
	"net/http"

	"github.com/go-chi/chi"
	"github.com/golovanevvs/confidant/internal/server/transport"
	"github.com/golovanevvs/confidant/internal/server/transport/http/logger"
	"go.uber.org/zap"
)

type handler struct {
	sv transport.IService
	lg *zap.SugaredLogger
}

// New - the handler constructor
func New(sv transport.IService, lg *zap.SugaredLogger) *handler {
	return &handler{
		sv: sv,
		lg: lg,
	}
}

// InitRoutes - request routing, used as http.Handler when starting the server
func (hd *handler) InitRoutes() http.Handler {
	// creating a router instance
	rt := chi.NewRouter()

	// using middleware
	// logging
	rt.Use(logger.WithLogging(hd.lg))

	// routes
	rt.Post("/register", hd.accountRegisterPost)

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
