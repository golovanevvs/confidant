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
	rt := chi.NewRouter()

	rt.Use(logger.WithLogging(hd.lg))

	rt.Route("/api", func(r chi.Router) {
		r.Post("/register", hd.accountRegisterPost)
		r.Post("/login", hd.loginPost)
		r.Get("/status", hd.getStatus)
		r.Post("/refresh_access", hd.refreshAccessTokenPost)
		r.With(hd.authByJWT).Get("/group_ids", hd.GroupIDsGet)
		r.With(hd.authByJWT).Post("/groups", hd.GroupsPost)
		r.With(hd.authByJWT).Put("/groups", hd.GroupsPut)
		r.With(hd.authByJWT).Post("/emails", hd.EmailsPost)
		r.With(hd.authByJWT).Patch("/emails", hd.EmailsPatch)
		r.With(hd.authByJWT).Get("/data_ids", hd.DataIDsGet)
		r.With(hd.authByJWT).Post("/data_dates", hd.DataDatesPost)
		r.With(hd.authByJWT).Post("/datas", hd.DatasPost)
		r.With(hd.authByJWT).Put("/datas", hd.DatasPut)
		r.With(hd.authByJWT).Post("/data_file", hd.DataFilePost)
		r.With(hd.authByJWT).Put("/data_file", hd.DataFilePut)
	})

	return rt
}
