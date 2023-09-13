package server

import (
	"log/slog"
	"net/http"
	"strconv"

	"github.com/F0rzend/simple-go-webserver/app/common"

	"github.com/F0rzend/simple-go-webserver/app/aggregate/bitcoin/handlers"
	"github.com/F0rzend/simple-go-webserver/app/aggregate/bitcoin/service"
	"github.com/F0rzend/simple-go-webserver/app/aggregate/user/handlers"
	"github.com/F0rzend/simple-go-webserver/app/aggregate/user/service"
	"github.com/F0rzend/simple-go-webserver/pkg/hlog"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

type Server struct {
	userRoutes    *userhandlers.UserHTTPHandlers
	bitcoinRoutes *bitcoinhandlers.BitcoinHTTPHandlers
}

func getUserIDFromURL(r *http.Request) (uint64, error) {
	const userIDURLKey = "id"

	return strconv.ParseUint(chi.URLParam(r, userIDURLKey), 10, 64) //nolint:gomnd
}

func NewServer(
	userRepository userservice.UserRepository,
	btcPriceGetter userservice.BTCPriceGetter,
	bitcoinRepository bitcoinservice.BTCRepository,
) *Server {
	bitcoinRoutes := bitcoinhandlers.NewBitcoinHTTPHandlers(bitcoinservice.NewBitcoinService(bitcoinRepository))
	userRoutes := userhandlers.NewUserHTTPHandlers(
		userservice.NewUserService(userRepository, btcPriceGetter),
		getUserIDFromURL,
	)

	return &Server{
		userRoutes:    userRoutes,
		bitcoinRoutes: bitcoinRoutes,
	}
}

func (s *Server) GetHTTPHandler(
	logger *slog.Logger,
) http.Handler {
	r := chi.NewRouter()

	r.Use(
		middleware.Recoverer,
		middleware.AllowContentType("application/json"),

		hlog.LoggerInjectionMiddleware(logger),
		hlog.RequestID,
		hlog.RequestMiddleware,
	)

	r.Route("/users", func(r chi.Router) {
		r.Post("/", common.ErrorHandler(s.userRoutes.CreateUser))

		r.Route("/{id}", func(r chi.Router) {
			r.Get("/", common.ErrorHandler(s.userRoutes.GetUser))
			r.Put("/", common.ErrorHandler(s.userRoutes.UpdateUser))
			r.Get("/balance", common.ErrorHandler(s.userRoutes.GetUserBalance))

			r.Post("/usd", common.ErrorHandler(s.userRoutes.ChangeUSDBalance))
			r.Post("/btc", common.ErrorHandler(s.userRoutes.ChangeBTCBalance))
		})
	})

	r.Route("/bitcoin", func(r chi.Router) {
		r.Get("/", common.ErrorHandler(s.bitcoinRoutes.GetBTCPrice))
		r.Put("/", common.ErrorHandler(s.bitcoinRoutes.SetBTCPrice))
	})

	return r
}
