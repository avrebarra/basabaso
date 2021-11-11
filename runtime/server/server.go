package server

import (
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/avrebarra/basabaso/pkg/ctxops"
	"github.com/avrebarra/valeed"
	"github.com/gorilla/mux"
	"github.com/gravityblast/xrequestid"
	"github.com/rs/xid"
	"github.com/rs/zerolog/log"
	"github.com/unrolled/secure"
	"github.com/urfave/negroni"
)

type Config struct {
	// Dependency string `validate:"required"`
}

type Server struct {
	config Config
}

func New(cfg Config) (*Server, error) {
	if err := valeed.Validate(cfg); err != nil {
		return nil, err
	}

	e := &Server{config: cfg}

	return e, nil
}

func (e *Server) MakeHandler() (out http.Handler, err error) {
	router := mux.NewRouter()
	h := Handler{Server: e}

	// setup routes
	router.Methods("GET").Path("/").HandlerFunc(h.HealthCheck())

	router.Methods("GET").Path("/api/health").HandlerFunc(h.HealthCheck())
	// router.Methods("POST").Path("/api/users/find").HandlerFunc(h.HealthCheck())
	// router.Methods("POST").Path("/api/users/persist").HandlerFunc(h.HealthCheck())
	router.Methods("POST").Path("/api/words/find").HandlerFunc(h.HealthCheck())
	router.Methods("POST").Path("/api/words/persist").HandlerFunc(h.HealthCheck())
	router.Methods("POST").Path("/api/words/vote").HandlerFunc(h.HealthCheck())
	// router.Methods("POST").Path("/api/words/report").HandlerFunc(h.HealthCheck())

	// setup middleware
	mwglobal := negroni.New()
	mwglobal.Use(negroni.HandlerFunc(secure.New(secure.Options{
		BrowserXssFilter:      true,
		FrameDeny:             true,
		ContentSecurityPolicy: "script-src $NONCE",
	}).HandlerFuncWithNext))
	mwglobal.Use(xrequestid.New(10))
	mwglobal.UseHandler(router)

	out = mwglobal
	return
}

type HandlerFunc func(w http.ResponseWriter, r *http.Request) (err error)

func (e *Server) buildHandler(in HandlerFunc) (out http.HandlerFunc) {
	var (
		HKeyCorrelationID = "Correlation-Id"
		SystemID          = "docdeck"
	)

	type Process struct {
		Name     string                 `json:"name"`
		Data     map[string]interface{} `json:"data"`
		StartAt  time.Time              `json:"start_at"`
		FinishAt time.Time              `json:"finish_at"`
	}

	FmtProcessFunc := func(in []ctxops.Process) (out []Process) {
		out = []Process{}
		for _, v := range in {
			out = append(out, Process(v))
		}
		return
	}

	return func(w http.ResponseWriter, r *http.Request) {
		var err error

		// setup request context
		ctx := r.Context()
		ctx = ctxops.CreateWith(ctx)
		r = r.WithContext(ctx)

		// populate ops id
		opsid := r.Header.Get(HKeyCorrelationID)
		if opsid == "" {
			opsid = xid.New().String()
		}
		ctxops.SetOpsID(ctx, opsid)
		w.Header().Set(HKeyCorrelationID, opsid)

		// execute handler
		var resptime time.Duration
		func() {
			// catch panic as error
			defer func() {
				if r := recover(); r != nil {
					err = fmt.Errorf("recovered from panic: %s", fmt.Sprint(r))
					return
				}
			}()

			// perform handling
			timestamp := time.Now()
			if err = in(w, r); err != nil {
				return
			}
			resptime = time.Since(timestamp)
		}()

		// uncaught error handling
		if err != nil {
			err = fmt.Errorf("uncaught error: %w", err)
			log.Err(err).Str("opsid", opsid).Msg("req err")

			out := RespPresets[RCUnexpected]
			out.AddMsg(err.Error())

			sendJSON(ctx, w, out.Normalize())
		}

		// log request
		// ** acquire resp data
		respdata, _ := ctxops.GetData(ctx, ctxkeyrespdata).(Resp)

		// ** submit log
		log.Info().
			Str("opsid", opsid).
			Str("sysid", SystemID).
			Str("uri", fmt.Sprintf("%s:%s", strings.ToLower(r.Method), r.RequestURI)).
			Str("respcode", string(respdata.Code)).
			Str("respmsg", respdata.Message).
			Dur("rt", resptime).
			Interface("vars", ctxops.ListVars(ctx)).
			Interface("processes", FmtProcessFunc(ctxops.ListProcesses(ctx))).
			Interface("warns", ctxops.ListWarnings(ctx)).
			Msg("req done")
	}
}
