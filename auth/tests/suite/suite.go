package suite

import (
	"context"
	"fmt"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"net"
	"os"
	"strconv"
	"testing"

	"github.com/viacheslavek/grpcauth/auth/internal/config"

	authv1 "github.com/viacheslavek/grpcauth/api/gen/go/auth"
)

type Suite struct {
	*testing.T
	Ctx         context.Context
	Cfg         *config.Config
	OwnerClient authv1.OwnerControllerClient
}

const (
	grpcHost = "localhost"
	testPath = "config/local.yaml"
)

func New(t *testing.T) *Suite {
	t.Helper()
	t.Parallel()

	dir, errG := os.Getwd()
	if errG != nil {
		t.Fatalf("work dir error: %v", errG)
	}
	fmt.Println("Current working directory:", dir)

	cfg := config.MustLoadPath(testPath)

	ctx, cancel := context.WithTimeout(context.Background(), cfg.GRPC.Timeout)

	t.Cleanup(func() {
		t.Helper()
		cancel()
	})

	target := net.JoinHostPort(grpcHost, strconv.Itoa(cfg.GRPC.Port))
	opts := []grpc.DialOption{
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	}

	clientConn, err := grpc.NewClient(target, opts...)
	if err != nil {
		t.Fatalf("grpc server connection failed: %v", err)
	}

	return &Suite{
		T:           t,
		Ctx:         ctx,
		Cfg:         cfg,
		OwnerClient: authv1.NewOwnerControllerClient(clientConn),
	}
}
