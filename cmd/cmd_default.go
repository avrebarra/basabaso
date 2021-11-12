package cmd

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/avrebarra/basabaso/pkg/fireconfig"
	"github.com/avrebarra/basabaso/pkg/logutil"
	"github.com/avrebarra/basabaso/runtime/server"
	"github.com/avrebarra/valeed"
	"github.com/caarlos0/env"
	"github.com/joho/godotenv"
	"github.com/rs/zerolog"

	zlog "github.com/rs/zerolog/log"
)

type ExecDefault struct{}

func (s ExecDefault) Run() error {
	ctx := context.Background()

	// ---
	// prepare env config
	type ConfigEnvironment struct {
		ServiceName       string `env:"SERVICE_NAME" validate:"required"`
		FirestoreCredFile string `env:"FIRESTORE_CRED_FILE" validate:"required"`
	}
	var cfgEnv = ConfigEnvironment{
		ServiceName: "basabaso",
	}
	_ = godotenv.Load()
	if err := env.Parse(&cfgEnv); err != nil {
		err = fmt.Errorf("failure parsing environment file: %w", err)
		return err
	}
	if err := valeed.Validate(cfgEnv); err != nil {
		err = fmt.Errorf("bad environment config: %w", err)
		return err
	}

	// ---
	// prepare service config
	type ConfigService struct {
		ServerPort                  int  `json:"server_port" validate:"required"`
		LoggingEnablePrettyPrinting bool `json:"logging_enable_pretty_printing" validate:"-"`
	}

	// * connect and fetch from fireconfig
	fireconf, err := fireconfig.New(fireconfig.Config{
		FirebaseConfigJSON: []byte(cfgEnv.FirestoreCredFile),
		CollecName:         "core-configs",
		ContentFieldName:   "content",
	})
	if err != nil {
		err = fmt.Errorf("fireconfig setup failed: %w", err)
		return err
	}
	outGetConfig, err := fireconf.Get(ctx, cfgEnv.ServiceName)
	if err != nil {
		err = fmt.Errorf("fireconfig fetch failed: %w", err)
		return err
	}
	fireconf.Close(ctx)

	// * parse and validate
	var cfgService = ConfigService{}
	err = json.Unmarshal(outGetConfig, &cfgService)
	if err != nil {
		return err
	}
	if err := valeed.Validate(cfgService); err != nil {
		err = fmt.Errorf("bad service config: %w", err)
		return err
	}

	// ---
	// setup logging
	outwriter := logutil.PrettyPrinter{Enable: cfgService.LoggingEnablePrettyPrinting, Out: os.Stdout}
	zlog.Logger = zerolog.New(outwriter).With().Timestamp().Logger()
	log.SetFlags(0)
	log.SetOutput(zerolog.New(outwriter).With().Str("level", "debug").Timestamp().Logger())
	log.Println("logging prepared")

	// ---
	// setup server
	servercore, err := server.New(server.Config{})
	if err != nil {
		err = fmt.Errorf("cannot init server core: %w", err)
		return err
	}

	serverhandler, err := servercore.MakeHandler()
	if err != nil {
		err = fmt.Errorf("cannot make server handler: %w", err)
		return err
	}

	serverinst := &http.Server{
		Addr:    fmt.Sprintf(":%d", cfgService.ServerPort),
		Handler: serverhandler,
	}

	// ---
	// run runtimes
	log.Println("dispatching runtimes")

	// ** listen for sigterm signal
	done := make(chan os.Signal, 1)
	signal.Notify(done, os.Interrupt, syscall.SIGINT, syscall.SIGTERM, syscall.SIGABRT)

	// ** dispatch http server
	log.Printf("* using http://localhost:%d to start http server...\n", cfgService.ServerPort)
	go func() {
		if err := serverinst.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			err = fmt.Errorf("http server listen error: %w", err)
			log.Fatal(err.Error())
		}
	}()

	log.Println("---")
	<-done // wait for sigterm and graceful shutdowns

	// ---
	// graceful kill runtimes
	// ** shutdown server
	log.Println("* shutting down server...")
	ctxShutdown, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := serverinst.Shutdown(ctxShutdown); err != nil {
		err = fmt.Errorf("http server shutdown error: %w", err)
		log.Println(err.Error())
	}

	log.Println("* exited")

	return nil
}
