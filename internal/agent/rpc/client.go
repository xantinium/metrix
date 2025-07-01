// Package rpc содержит RPC-клиент.
package rpc

import (
	"context"
	"net"
	"time"

	grpcpool "github.com/processout/grpc-go-pool"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	"github.com/xantinium/metrix/internal/agent/rpc/interceptors"
	"github.com/xantinium/metrix/internal/presentation/rpc/gen"
)

var (
	poolSize            = 5
	rpcPoolConnWaitTime = 5 * time.Second
	connMaxIdleTime     = 5 * time.Minute
	connMaxLifeTime     = 15 * time.Minute
)

// ClientOptions параметры RPC-клиента.
type ClientOptions struct {
	ServerAddr      string
	PrivateKey      string
	CryptoPublicKey string
	Addr            net.IP
}

// NewClient создаёт новый RPC-клиент.
func NewClient(opts ClientOptions) *Client {
	client := &Client{
		serverAddr:      opts.ServerAddr,
		privateKey:      opts.PrivateKey,
		cryptoPublicKey: opts.CryptoPublicKey,
		addr:            opts.Addr,
	}

	// При init=0, err всегда будет nil.
	client.pool, _ = grpcpool.New(client.connFactory, 0, poolSize, connMaxIdleTime, connMaxLifeTime)

	return client
}

// Client структура, описывающая RPC-клиент.
type Client struct {
	serverAddr      string
	privateKey      string
	cryptoPublicKey string

	pool *grpcpool.Pool
	addr net.IP
}

// Close закрывает клиент.
func (client *Client) Close() {
	client.pool.Close()
}

func (client *Client) connFactory() (*grpc.ClientConn, error) {
	return grpc.NewClient(
		client.serverAddr,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithChainUnaryInterceptor(client.prepareInterceptors()...),
	)
}

func (client *Client) prepareInterceptors() []grpc.UnaryClientInterceptor {
	result := []grpc.UnaryClientInterceptor{}

	if client.addr != nil {
		result = append(result, interceptors.XRealIPInterceptor(client.addr))
	}

	// TODO: реализовать перехватчики:
	// 1) HashInterceptor
	// 2) EncryptInterceptor

	return result
}

func (client *Client) getMetricsClient() gen.MetricsClient {
	return gen.NewMetricsClient(newConn(client))
}

func newConn(client *Client) *conn {
	return &conn{client: client}
}

// conn является обёрткой над gRPC-подключением.
// Нужен для удобной работы с пулом подключений.
// Реализует интерфейс grpc.ClientConnInterface.
type conn struct {
	client *Client
}

func (c *conn) Invoke(
	ctx context.Context,
	method string,
	args, reply any,
	opts ...grpc.CallOption,
) error {
	poolCtx, cancel := context.WithTimeout(ctx, rpcPoolConnWaitTime)
	defer cancel()

	conn, err := c.client.pool.Get(poolCtx)
	if err != nil {
		return err
	}
	defer conn.Close()

	return conn.Invoke(ctx, method, args, reply, opts...)
}

func (c *conn) NewStream(
	_ context.Context,
	_ *grpc.StreamDesc,
	_ string,
	_ ...grpc.CallOption,
) (grpc.ClientStream, error) {
	return nil, nil
}
