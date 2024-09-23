package rate

import (
	"context"
	"fmt"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	descRate "github.com/VadimGossip/drs_data_loader/pkg/rate_v1"
	descGrpc "github.com/VadimGossip/drs_storage_tester/internal/client/grpc"
	"github.com/VadimGossip/drs_storage_tester/internal/config"
)

type client struct {
	cl descRate.RateV1Client
}

func NewClient(authGRPCClientConfig config.RateGrpcConfig) (descGrpc.RateClient, error) {
	conn, err := grpc.NewClient(authGRPCClientConfig.Address(),
		grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, fmt.Errorf("failed to connect to grpc server: %v", err)
	}
	return &client{cl: descRate.NewRateV1Client(conn)}, nil
}

func (c *client) FindRate(ctx context.Context, gwgrId, dateAt int64, dir uint8, aNumber, bNumber string) (int64, float64, time.Duration, error) {
	ts := time.Now()
	res, err := c.cl.FindRate(ctx, &descRate.FindRateRequest{
		GwgrId:  gwgrId,
		DateAt:  dateAt,
		Dir:     uint32(dir),
		ANumber: aNumber,
		BNumber: bNumber,
	})
	if err != nil {
		return 0, 0, time.Since(ts), err
	}
	return res.RmsrId, res.PriceBase, time.Since(ts), nil
}

func (c *client) FindSupRates(ctx context.Context, gwgrIds []int64, dateAt int64, aNumber, bNumber string) (int64, time.Duration, error) {
	return 0, 0, fmt.Errorf("unimplemented")
}