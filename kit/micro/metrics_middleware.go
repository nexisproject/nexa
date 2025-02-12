// Copyright (C) micros. 2025-present.
//
// Created at 2025-02-10, by liasica

package micro

import (
	"context"
	"time"

	"github.com/go-kit/kit/endpoint"
	"github.com/go-kit/kit/metrics"
	"github.com/go-kit/kit/metrics/prometheus"
	stdprometheus "github.com/prometheus/client_golang/prometheus"
)

func InstrumentingMiddleware(next endpoint.Endpoint) endpoint.Endpoint {
	return func(ctx context.Context, request any) (response any, err error) {
		var dur metrics.Histogram = prometheus.NewSummaryFrom(stdprometheus.SummaryOpts{
			Namespace: ctx.Value("namespace").(string),
			Subsystem: "api",
			Name:      "request_duration_seconds",
			Help:      "Total time spent serving requests.",
		}, []string{})

		begin := time.Now()
		response, err = next(ctx, request)
		dur.Observe(time.Since(begin).Seconds())

		return
	}
}
