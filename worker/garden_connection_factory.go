package worker

import (
	"net/http"

	"code.cloudfoundry.org/clock"
	gconn "code.cloudfoundry.org/garden/client/connection"
	"code.cloudfoundry.org/garden/routes"
	"code.cloudfoundry.org/lager"
	"github.com/concourse/atc/worker/transport"
	"github.com/concourse/retryhttp"
	"github.com/tedsuo/rata"
)

//go:generate counterfeiter . GardenConnectionFactory
type GardenConnectionFactory interface {
	BuildConnection() gconn.Connection
}

type gardenConnectionFactory struct {
	db          transport.TransportDB
	logger      lager.Logger
	workerName  string
	workerHost  string
	retryPolicy transport.RetryPolicy
}

func NewGardenConnectionFactory(
	db transport.TransportDB,
	logger lager.Logger,
	workerName string,
	workerHost string,
	retryPolicy transport.RetryPolicy,
) GardenConnectionFactory {
	return &gardenConnectionFactory{
		db:          db,
		logger:      logger,
		workerName:  workerName,
		workerHost:  workerHost,
		retryPolicy: retryPolicy,
	}
}

func (gcf *gardenConnectionFactory) BuildConnection() gconn.Connection {
	httpClient := &http.Client{
		Transport: &retryhttp.RetryRoundTripper{
			Logger:       gcf.logger.Session("retryable-http-client"),
			Sleeper:      clock.NewClock(),
			RetryPolicy:  gcf.retryPolicy,
			RoundTripper: transport.NewRoundTripper(gcf.workerName, gcf.workerHost, gcf.db, &http.Transport{DisableKeepAlives: true}),
		},
	}

	hijackableClient := &retryhttp.RetryHijackableClient{
		Logger:           gcf.logger.Session("retry-hijackable-client"),
		Sleeper:          clock.NewClock(),
		RetryPolicy:      gcf.retryPolicy,
		HijackableClient: transport.NewHijackableClient(gcf.workerName, gcf.db, retryhttp.DefaultHijackableClient),
	}

	// the request generator's address doesn't matter because it's overwritten by the worker lookup clients
	hijackStreamer := &transport.WorkerHijackStreamer{
		HttpClient:       httpClient,
		HijackableClient: hijackableClient,
		Req:              rata.NewRequestGenerator("http://127.0.0.1:8080", routes.Routes),
	}

	return gconn.NewWithHijacker(hijackStreamer, gcf.logger)
}
