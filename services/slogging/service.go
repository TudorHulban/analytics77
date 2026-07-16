package slogging

import (
	"context"
	"fmt"
	"io"
	"os"
	"runtime"

	"github.com/tudorhulban/arenalog"
	"github.com/tudorhulban/bytearena"
	"github.com/tudorhulban/bytearena/helpers"
	"github.com/tudorhulban/hxhelpers"
)

type ServiceLogging struct {
	Logger *arenalog.Logger
}

// NewServiceLogging provides a closing function.
func NewServiceLogging(pathLogFile string, writer ...io.Writer) (*ServiceLogging, func(), error) {
	var (
		fileHTTPServer *os.File
		errFileHTTP    error
	)

	if len(writer) != 1 {
		fileHTTPServer, errFileHTTP = os.OpenFile( //nolint:gosec
			pathLogFile,
			os.O_APPEND|os.O_CREATE|os.O_WRONLY,
			0644,
		)
		if errFileHTTP != nil {
			return nil, nil,
				errFileHTTP
		}
	}

	ingestor, errCrIngestor := bytearena.NewIngestor(
		bytearena.Size100K(),
		hxhelpers.TernaryLazy[io.Writer](
			len(writer) == 1,
			func() io.Writer {
				return writer[0]
			},
			func() io.Writer {
				return fileHTTPServer
			},
		),

		helpers.TernaryWithValueIn(
			[]int{1},
			runtime.NumCPU(),
			nil,
			bytearena.WithCounterCoreCPU(),
		),
	)
	if errCrIngestor != nil {
		return nil, nil,
			errCrIngestor
	}

	ctx, cancel := context.WithCancel(context.Background())
	chIngestionEnd := ingestor.StartIngestion(ctx)

	fnClosing := func() {
		_ = fileHTTPServer.Close()

		cancel()
		<-chIngestionEnd
	}

	l, errCrLogger := arenalog.NewLogger(
		&arenalog.ParamsNewLogger{
			Ingestor:    ingestor,
			LoggerLevel: arenalog.LevelInfo,

			WithFatalWriter: os.Stdout,
			WithJSON:        true,
		},

		arenalog.WithTimestampRFC3339UTC(ctx),
	)
	if errCrLogger != nil {
		fnClosing()

		return nil, nil,
			fmt.Errorf(
				"failed to create service logging: %w",
				errCrLogger,
			)
	}

	return &ServiceLogging{
			Logger: l,
		},
		fnClosing,
		nil
}
