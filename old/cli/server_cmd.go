package cli

import (
	"context"
	"encoding/json"
	"errors"
	"net"
	"net/http"
	"strings"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/jpillora/ipfilter"
	"github.com/ogen-go/ogen/ogenerrors"
	"github.com/ory/graceful"
	"github.com/spf13/cobra"
	"github.com/ugent-library/people-service/api/v1"
	"github.com/ugent-library/people-service/public"
	"github.com/ugent-library/zaphttp"
	"github.com/ugent-library/zaphttp/zapchi"
)

func init() {
	rootCmd.AddCommand(serverCmd)
}

type ErrorMessage struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

type apiSecurityHandler struct {
	APIKey string
}

func (s *apiSecurityHandler) HandleApiKey(ctx context.Context, operationName string, t api.ApiKey) (context.Context, error) {
	if t.APIKey == s.APIKey {
		return ctx, nil
	}
	return ctx, errors.New("unauthorized")
}

var serverCmd = &cobra.Command{
	Use:   "server",
	Short: "start the openapi server",
	RunE: func(cmd *cobra.Command, args []string) error {
		repo, err := newRepository()
		if err != nil {
			return err
		}

		mux := chi.NewMux()

		mux.Use(middleware.RequestID)
		mux.Use(middleware.RealIP)
		mux.Use(zaphttp.SetLogger(logger.Desugar(), zapchi.RequestID))
		mux.Use(middleware.RequestLogger(zapchi.LogFormatter()))
		mux.Use(middleware.Recoverer)

		if config.IPRanges != "" {
			ipFilter := ipfilter.New(ipfilter.Options{
				AllowedIPs:     strings.Split(config.IPRanges, ","),
				BlockByDefault: true,
			})
			mux.Use(func(next http.Handler) http.Handler {
				return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
					remoteIP := r.Header.Get("X-Real-IP")
					// middleware realip set RemoteAddr to ip address only (no :port), so below split leaves an empty string
					if remoteIP == "" {
						rip, _, _ := net.SplitHostPort(r.RemoteAddr)
						remoteIP = rip
					}
					if ipFilter.Allowed(remoteIP) {
						next.ServeHTTP(w, r)
						return
					}
					data, _ := json.Marshal(ErrorMessage{
						Code:    http.StatusForbidden,
						Message: "forbidden",
					})
					w.Header().Set("Content-Type", "application/json")
					w.WriteHeader(http.StatusForbidden)
					w.Write(data)
				})
			})
		}

		apiServer, err := api.NewServer(
			api.NewService(repo),
			&apiSecurityHandler{APIKey: config.Api.Key},
			api.WithErrorHandler(func(ctx context.Context, w http.ResponseWriter, r *http.Request, err error) {
				status := ogenerrors.ErrorCode(err)
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(status)
				data, _ := json.Marshal(ErrorMessage{
					Code:    status,
					Message: err.Error(),
				})
				w.Write(data)
			}),
		)
		if err != nil {
			return err
		}

		mux.Mount("/api/v1/openapi.yaml", http.StripPrefix("/api/v1/", api.OpenapiFileServer()))
		mux.Mount("/swagger/", http.StripPrefix("/swagger/", public.SwaggerFileServer()))
		mux.Mount("/api/v1", http.StripPrefix("/api/v1", apiServer))
		mux.Get("/info", func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(200)
			bytes, _ := json.Marshal(&version)
			w.Write(bytes)
		})

		srv := graceful.WithDefaults(&http.Server{
			Addr:         config.Api.Addr(),
			Handler:      mux,
			ReadTimeout:  10 * time.Second,
			WriteTimeout: 10 * time.Second,
		})

		logger.Infof("starting server at %s", config.Api.Addr())
		if err := graceful.Graceful(srv.ListenAndServe, srv.Shutdown); err != nil {
			return err
		}
		logger.Info("gracefully stopped server")
		return nil
	},
}
