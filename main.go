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
		r.Error = 2
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
		r.Error = 2
		r.Text = err.Error()
		fmt.Println(r.String())
		return
	}

	// create response
	duration := time.Since(start)

	switch h.Status {
	// unknown service - warning
	case grpc_health.HealthCheckResponse_SERVICE_UNKNOWN:
		r.Error = 1
		r.Text = fmt.Sprintf("Service %v not known", serv)
		fmt.Println(r.String())
		return
	// not service - error
	case grpc_health.HealthCheckResponse_NOT_SERVING:
		r.Error = 3
		r.Text = fmt.Sprintf("Service health failing")
		fmt.Println(r.String())
		return
	}

	r.AddChannel(prtg.SensorChannel{
		Name:  "Response time",
		Value: float64(duration.Milliseconds()),
		Unit:  prtg.UnitTimeResponse,
		Float: 1,
	})

	fmt.Println(r.String())
}
