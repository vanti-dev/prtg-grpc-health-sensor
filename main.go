package main

import (
	"context"
	"flag"
	"fmt"
	"time"

	"github.com/PRTG/go-prtg-sensor-api"
	"google.golang.org/grpc"
	grpc_health "google.golang.org/grpc/health/grpc_health_v1"
)

var (
	addr = flag.String("address", ":9001", "host:port address of grpc server")
	serv = flag.String("service", "", "Service name to check (defaults to \"\")")
)

func main() {
	flag.Parse()

	// create a response and log start time
	r := &prtg.SensorResponse{}
	start := time.Now()

	// setup connection
	conn, err := grpc.Dial(
		*addr,
		grpc.WithInsecure(),
	)

	if err != nil {
		r.Error = 1
		r.Text = err.Error()
		fmt.Println(r.String())
		return
	}
	defer conn.Close()

	// call the health check rpc
	h, err := grpc_health.NewHealthClient(conn).Check(context.Background(),
		&grpc_health.HealthCheckRequest{
			Service: *serv,
		},
	)

	if err != nil {
		r.Error = 1
		r.Text = err.Error()
		fmt.Println(r.String())
		return
	}

	if h.Status == grpc_health.HealthCheckResponse_SERVICE_UNKNOWN {
		r.Error = 1
		r.Text = fmt.Sprintf("Service %v not known", serv)
		fmt.Println(r.String())
		return
	}

	// create response
	duration := time.Since(start)

	r.AddChannel(prtg.SensorChannel{
		Name:  "Response time",
		Value: float64(duration.Milliseconds()),
		Unit:  prtg.UnitTimeResponse,
		Float: 1,
	})

	statusChan := prtg.SensorChannel{
		Name:  "Service status",
		Float: 0,
	}

	switch h.Status {
	case grpc_health.HealthCheckResponse_SERVING:
		statusChan.Value = 0
	case grpc_health.HealthCheckResponse_NOT_SERVING:
		statusChan.Value = 1
	}

	r.AddChannel(statusChan)

	fmt.Println(r.String())
}