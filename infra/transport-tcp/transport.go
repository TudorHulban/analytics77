package transporttcp

import (
	"encoding/gob"
	"errors"
	"fmt"
	"net"
	"sync"

	"github.com/tudorhulban/analytics77/services/sanalytics"
	"github.com/tudorhulban/analytics77/services/slogging"
	"github.com/tudorhulban/analytics77/shared"
	"github.com/tudorhulban/arenalog"
	"github.com/tudorhulban/hxhelpers"
	"github.com/tudorhulban/hxhelpers/piers"
)

type TransportTCP struct {
	listener net.Listener

	serviceLogging   *slogging.ServiceLogging
	serviceAnalytics *sanalytics.ServiceAnalytics

	logContext *arenalog.LogContext

	chQuit             chan struct{}
	wgHandleConnection sync.WaitGroup
}

type PiersNewTransportTCP struct {
	ServiceLogging   *slogging.ServiceLogging
	ServiceAnalytics *sanalytics.ServiceAnalytics
}

func NewTransportTCP(l net.Listener, dependencies *PiersNewTransportTCP) (*TransportTCP, error) {
	if errValidate := piers.ValidateDependencies(dependencies); errValidate != nil {
		return nil,
			errValidate
	}

	logContext := arenalog.
		NewLogContext(dependencies.ServiceLogging.Logger).
		WithRoot("transport", "TCP")

	return &TransportTCP{
			listener:         l,
			serviceAnalytics: dependencies.ServiceAnalytics,
			serviceLogging:   dependencies.ServiceLogging,
			logContext:       logContext,

			chQuit: make(chan struct{}),
		},
		nil
}

func (s *TransportTCP) GetListeningAddress() string {
	return s.listener.Addr().String()
}

func (s *TransportTCP) handleConnection(conn net.Conn) {
	defer s.wgHandleConnection.Done()
	defer conn.Close()

	// Create decoder for this specific payload.
	decoder := gob.NewDecoder(conn)

	var batch shared.Requests

	// Expect exactly ONE decode operation.
	if errDecode := decoder.Decode(&batch); errDecode != nil {
		s.logContext.Print(
			fmt.Sprintf(
				"failed to decode payload from %s: %s\n",

				conn.RemoteAddr(),
				errDecode.Error(),
			),
		)

		return // Exit immediately, defer handles the connection close
	}

	s.logContext.Print(
		fmt.Sprintf(
			"received %d request(s) from %s",
			len(batch),
			conn.RemoteAddr(),
		),
	)

	// Process the data.
	errsValidationEvents, errsProcessEvents := s.
		serviceAnalytics.
		RecordEvents(batch)
	if len(errsValidationEvents) > 0 {
		s.logContext.Print(
			fmt.Sprintf(
				"handleConnection - validation error(s)(%d) from %s: %v",

				len(errsValidationEvents),
				conn.RemoteAddr(),
				errsValidationEvents,
			),
		)

		return
	}

	if len(errsProcessEvents) > 0 {
		s.logContext.Print(
			fmt.Sprintf(
				"handleConnection - processing error(s)(%d) from %s: %v",

				len(errsProcessEvents),
				conn.RemoteAddr(),
				errsProcessEvents,
			),
		)
	}

	s.logContext.Print(
		fmt.Sprintf(
			"processed with no errors %d request(s) from %s",

			len(batch),
			conn.RemoteAddr(),
		),
	)
}

func (s *TransportTCP) Start() error {
	s.logContext.Print(
		hxhelpers.Sprintf(
			"TCP transport started, listening on %s",
			s.listener.Addr().String()),
	)

	go func() {
		<-s.chQuit

		_ = s.listener.Close()
	}()

	for {
		conn, errAccept := s.listener.Accept()
		if errAccept != nil {
			// When s.listener.Close() is called,
			// Accept() instantly wakes up and returns net.ErrClosed
			if errors.Is(errAccept, net.ErrClosed) {
				return nil // Clean, graceful shutdown
			}

			return errAccept
		}

		s.wgHandleConnection.Add(1)

		go s.handleConnection(conn)
	}
}

func (s *TransportTCP) Stop() {
	close(s.chQuit)

	_ = s.listener.Close()

	s.wgHandleConnection.Wait() // Waits for active connections to finish processing
}
