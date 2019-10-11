package main

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"os"
	"time"

	"github.com/PRTG/go-prtg-sensor-api"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	grpc_health "google.golang.org/grpc/health/grpc_health_v1"
	"google.golang.org/grpc/status"
)

var (
	addr    = flag.String("address", ":9001", "host:port address of grpc server")
	serv    = flag.String("service", "", "Service name to check (defaults to \"\")")
	timeout = flag.Duration("timeout", 20*time.Second, "Configure the timeout for the health API request")
)

func main() {
	flag.Parse()

	// validate arguments
	if *timeout < 0 {
		_, _ = fmt.Fprintf(os.Stderr, "-timeout cannot be negative: %v\n", *timeout)
		flag.Usage()
		os.Exit(1)
	}

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
	ctx := context.Background()
	if *timeout > 0 {
		ctx, _ = context.WithTimeout(ctx, *timeout)
	}
	h, err := grpc_health.NewHealthClient(conn).Check(ctx,
		&grpc_health.HealthCheckRequest{
			Service: *serv,
		},
	)

	if err != nil {
		r.Error = 2
		r.Text = textFromError(err)
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

func textFromError(err error) string {
	// check gRPC error codes
	if s, ok := status.FromError(err); ok {
		switch s.Code() {
		case codes.DeadlineExceeded:
			return fmt.Sprintf("Service %v did not respond within %v (deadline exceeded)", *addr, *timeout)
		}
	}

	// check go error types
	switch {
	case errors.Is(err, context.DeadlineExceeded):
		return fmt.Sprintf("Service %v did not respond within %v (deadline exceeded)", *addr, *timeout)
	default:
		return err.Error()
	}
}
