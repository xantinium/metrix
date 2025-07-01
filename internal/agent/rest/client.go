// Package rpc содержит REST-клиент.
package rest

import (
	"bytes"
	"context"
	"net"
	"net/http"

	"github.com/mailru/easyjson"

	"github.com/xantinium/metrix/internal/logger"
	"github.com/xantinium/metrix/internal/tools"
)

// ClientOptions параметры REST-клиента.
type ClientOptions struct {
	ServerAddr      string
	PrivateKey      string
	CryptoPublicKey string
	Addr            net.IP
}

// NewClient создаёт новый REST-клиент.
func NewClient(opts ClientOptions) *Client {
	client := &Client{
		retrier:         tools.DefaulRetrier,
		serverAddr:      opts.ServerAddr,
		privateKey:      opts.PrivateKey,
		cryptoPublicKey: opts.CryptoPublicKey,
		addr:            opts.Addr,
	}

	return client
}

// Client структура, описывающая REST-клиент.
type Client struct {
	retrier *tools.Retrier

	serverAddr      string
	privateKey      string
	cryptoPublicKey string

	addr net.IP
}

func (client *Client) sendRequest(ctx context.Context, url string, req easyjson.Marshaler) error {
	var (
		err      error
		httpReq  *http.Request
		reqBytes []byte
	)

	reqBytes, err = easyjson.Marshal(req)
	if err != nil {
		return err
	}

	reqBytes, err = tools.Compress(reqBytes)
	if err != nil {
		return err
	}

	if client.cryptoPublicKey != "" {
		reqBytes, err = tools.Encrypt(client.cryptoPublicKey, reqBytes)
		if err != nil {
			return err
		}
	}

	reqBody := bytes.NewBuffer(reqBytes)
	httpReq, err = http.NewRequestWithContext(ctx, http.MethodPost, url, reqBody)
	if err != nil {
		return err
	}

	httpReq.Header.Set(tools.HeaderAcceptEncoding, "gzip")
	httpReq.Header.Set(tools.HeaderContentEncoding, "gzip")
	httpReq.Header.Set(tools.HeaderContentType, "application/json")
	if client.addr != nil {
		httpReq.Header.Set(tools.HeaderXRealIP, client.addr.String())
	}

	if client.privateKey != "" {
		var hashedReq string
		hashedReq, err = tools.CalcSHA256(reqBytes, client.privateKey)
		if err != nil {
			return err
		}

		httpReq.Header.Set(tools.HeaderHashSHA256, hashedReq)
	}

	client.retrier.Exec(func() bool {
		var resp *http.Response
		resp, err = http.DefaultClient.Do(httpReq)
		if resp != nil {
			resp.Body.Close()
		}
		return err != nil
	})

	return err
}

func (client *Client) logError(err error) {
	field := logger.Field{
		Name:  "entity",
		Value: "agent-worker",
	}

	logger.Error(err.Error(), field)
}
