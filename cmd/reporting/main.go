package main

import (
	"context"
	"flag"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"sync"
	"sync/atomic"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"gitlabdev.vadesecure.com/engineering/app/kit/debug"
	kitgin "gitlabdev.vadesecure.com/engineering/app/kit/gin"

	"gitlabdev.vadesecure.com/engineering/app/reporting/internal/api"
	"gitlabdev.vadesecure.com/engineering/app/reporting/internal/jira"
	"gitlabdev.vadesecure.com/engineering/app/reporting/internal/storage"
)

var version = "dev"

func main() {

	var atomicReady int32

	printVersion := flag.Bool("version", false, "Show version")
	flag.Parse()
	if *printVersion {
		log.Info().Msgf("example version %s\n", version)
		os.Exit(0)
	}

	log.Info().Str("version", version).Msg("server starting")

	conf, err := ParseConfig()
	if err != nil {
		log.Fatal().Err(err).Msg("failed to load the configuration")
	}

	log.Info().Interface("config", conf).Msg("configuration loaded")

	if !conf.Log.GinDebug {
		gin.SetMode(gin.ReleaseMode)
	}

	logLevel, err := zerolog.ParseLevel(conf.Log.Level)
	if err != nil {
		log.Fatal().Err(err).Msg("failed to parse the log level")
	}
	zerolog.SetGlobalLevel(logLevel)
	zerolog.TimeFieldFormat = time.RFC3339Nano

	// Jira init
	jiraC, err := jira.NewClient(&conf.Jira)
	if err != nil {
		log.Fatal().Err(err).Msg("failed to connect the jira API")
	}
	log.Info().Msg("API jira connected")

	// mongoDB init
	var mongoC *storage.Db
	mongoC, err = storage.New(&conf.Storage)
	if err != nil {
		log.Fatal().Err(err).Msg("failed to connect mongoDB server")
	}
	log.Info().Msg("mongoDB server connected")

	// init API
	technicalServer := startAPI(conf.TechnicalAPI, technicalAPIRouter(&atomicReady))
	internalServer := startAPI(conf.InternalAPI, internalAPIRouter())
	publicServer := startAPI(conf.PublicAPI, publicAPIRouter(jiraC, mongoC))

	// The service is ready to accept requests
	atomic.StoreInt32(&atomicReady, 1)

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM)
	<-stop

	atomic.StoreInt32(&atomicReady, 0)

	log.Info().Msg("server is shutting down...")

	var wg sync.WaitGroup
	wg.Add(3)
	for _, srv := range []*http.Server{technicalServer, publicServer, internalServer} {
		go func(s *http.Server) {
			defer wg.Done()
			log.Info().Str("address", s.Addr).Msg("http server is stopping...")
			err := s.Shutdown(context.Background())
			if err != nil {
				log.Error().Err(err).Str("address", s.Addr).Msg("http server stopped with error")
			}
		}(srv)
	}
	wg.Wait()

	log.Info().Msg("reporting stopped")
}

func technicalAPIRouter(atomicReady *int32) http.Handler {
	router := kitgin.NewRouter(map[string]interface{}{"api": "technical"})

	router.GET("/info/alive", func(c *gin.Context) {})
	router.GET("/info/ready", func(c *gin.Context) {
		if atomic.LoadInt32(atomicReady) == 0 {
			c.AbortWithStatus(http.StatusServiceUnavailable)
			return
		}
	})
	router.GET("/info/operating",
		func(c *gin.Context) {
			c.String(http.StatusOK, "OK")
		})
	router.GET("/info/version",
		func(c *gin.Context) {
			c.JSON(http.StatusOK, map[string]string{"version": version})
		})
	router.GET("/info/metrics", gin.WrapH(promhttp.Handler()))
	debug.GinAddDebugPPROFHandlers(router)

	router.PUT("/log/level/:level", func(c *gin.Context) {
		logLevel, err := zerolog.ParseLevel(strings.ToLower(c.Param("level")))
		if err != nil {
			c.AbortWithStatus(http.StatusBadRequest)
			return
		}
		zerolog.SetGlobalLevel(logLevel)
		log.Log().Str("new_level", zerolog.GlobalLevel().String()).Msg("log level updated")
		c.Status(http.StatusNoContent)
	})
	router.GET("/log/level", func(c *gin.Context) {
		c.JSON(http.StatusOK, map[string]interface{}{
			"level": zerolog.GlobalLevel().String(),
			"available_levels": []string{zerolog.DebugLevel.String(), zerolog.InfoLevel.String(),
				zerolog.WarnLevel.String(), zerolog.ErrorLevel.String()},
		})
	})

	return router
}
func publicAPIRouter(jiraC *jira.Client, mongoC *storage.Db) http.Handler {
	router := kitgin.NewRouter(map[string]interface{}{"api": "public"})
	router.GET("/api/v1/ping", func(c *gin.Context) {
		c.Status(http.StatusNoContent)
	})
	router.GET("/api/v1/hello", api.Hello("public"))
	router.GET("/api/v1/jira/:key", api.IssueJira("public", jiraC))
	router.GET("/api/v1/mongodb/:key", api.IssueMongo("public", mongoC))
	router.GET("/api/v1/add/:key", api.Add("public", mongoC, jiraC))
	router.GET("/api/v1/all", api.All("public", mongoC))
	router.GET("/api/v1/delete/:key", api.Delete("public", mongoC))
	router.GET("/api/v1/history", api.History("public", mongoC))
	router.GET("/api/v1/update", api.UpdateHistory("public", mongoC))
	router.GET("/api/v1/reset", api.Reset("public", mongoC))

	return router

}

func internalAPIRouter() http.Handler {
	router := kitgin.NewRouter(map[string]interface{}{"api": "internal"})
	router.GET("/api/v1/hello", api.Hello("internal"))
	return router
}
