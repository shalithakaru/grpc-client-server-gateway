package main

import (
	"os"

	"google.golang.org/grpc/grpclog"
)

func init() {
	grpclog.SetLoggerV2(grpclog.NewLoggerV2WithVerbosity(os.Stdout, os.Stderr, os.Stderr, 2))
}
