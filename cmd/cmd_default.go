package cmd

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/avrebarra/basabaso/pkg/logutil"
	"github.com/avrebarra/basabaso/runtime/server"
	"github.com/rs/zerolog"

	zlog "github.com/rs/zerolog/log"
)

type ExecDefault struct{}

func (s ExecDefault) Run() (err error) {
	portnum := 8777
	enablePrettyLog := true

	// ---
	// setup logging
	outwriter := logutil.PrettyPrinter{Enable: enablePrettyLog, Out: os.Stdout}
	zlog.Logger = zerolog.New(outwriter).With().Timestamp().Logger()
	log.SetFlags(0)
	log.SetOutput(zerolog.New(outwriter).With().Str("level", "debug").Timestamp().Logger())
	log.Println("logging prepared")

	// ---
	// setup server
	servercore, err := server.New(server.Config{})
	if err != nil {
		err = fmt.Errorf("cannot init server core: %w", err)
		return
	}

	serverhandler, err := servercore.MakeHandler()
	if err != nil {
		err = fmt.Errorf("cannot make server handler: %w", err)
		return
	}

	serverinst := &http.Server{
		Addr:    fmt.Sprintf(":%d", portnum),
		Handler: serverhandler,
	}

	// ---
	// run runtimes
	log.Println("dispatching runtimes")

	// ** listen for sigterm signal
	done := make(chan os.Signal, 1)
	signal.Notify(done, os.Interrupt, syscall.SIGINT, syscall.SIGTERM, syscall.SIGABRT)

	// ** dispatch http server
	log.Printf("* using http://localhost:%d to start http server...\n", portnum)
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

	return
}
