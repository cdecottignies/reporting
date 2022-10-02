package main

import (
	"net"
	"net/http"

	"github.com/rs/cors"
	"github.com/rs/zerolog/log"
)

// startAPI could move in a lib
func startAPI(apiConfig APIConfig, handler http.Handler) *http.Server {
	options := &cors.Options{}
	options.AllowedOrigins = []string{"*"}
	options.AllowedMethods = []string{http.MethodGet, http.MethodPost, http.MethodHead, http.MethodOptions, http.MethodDelete, http.MethodPut}
	options.AllowedHeaders = []string{"*"}
	srv := http.Server{
		Addr:         apiConfig.BindAddr,
		Handler:      cors.New(*options).Handler(handler),
		ReadTimeout:  apiConfig.ReadTimeout,
		WriteTimeout: apiConfig.WriteTimeout,
		IdleTimeout:  apiConfig.IdleTimeout,
	}
	// We don't use srv.ListenAndServe to ensure this service is able to accept new connexions
	// when this function returns. Running srv.ListenAndServe in a goroutine doesn't ensure
	// the port is correctly opened at the end of this function.
	ln, err := net.Listen("tcp", srv.Addr)
	if err != nil {
		log.Fatal().Err(err).Str("address", srv.Addr).Msg("failed to listen on tcp port")
	}
	go func() {
		log.Info().Str("address", srv.Addr).Msg("serve API")
		if err := srv.Serve(ln); err != http.ErrServerClosed {
			log.Fatal().Err(err).Str("address", srv.Addr).Msg("http server stopped with error")
		}
	}()

	return &srv
}
