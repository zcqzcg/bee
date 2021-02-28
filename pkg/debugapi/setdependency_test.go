package debugapi

import (
	"context"
	"crypto/ecdsa"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethersphere/bee/pkg/jsonhttp"
	"github.com/ethersphere/bee/pkg/jsonhttp/jsonhttptest"
	"github.com/ethersphere/bee/pkg/logging"
	pingpongmock "github.com/ethersphere/bee/pkg/pingpong/mock"
	"github.com/ethersphere/bee/pkg/swarm"
	"resenje.org/web"
)

func TestSetDependency(t *testing.T) {
	//topologyDriver := topologymock.NewTopologyDriver(o.TopologyOpts...)
	//acc := accountingmock.NewAccounting(o.AccountingOpts...)
	//settlement := swapmock.New(o.SettlementOpts...)
	//chequebook := chequebookmock.NewChequebook(o.ChequebookOpts...)
	//swapserv := swapmock.NewApiInterface(o.SwapOpts...)
	s := New(swarm.ZeroAddress, ecdsa.PublicKey{}, ecdsa.PublicKey{}, common.Address{}, nil, nil, nil, nil, logging.New(ioutil.Discard, 0), nil, nil, nil, nil, true, nil, nil, Options{})
	ts := httptest.NewServer(s)
	t.Cleanup(ts.Close)

	client := &http.Client{
		Transport: web.RoundTripperFunc(func(r *http.Request) (*http.Response, error) {
			u, err := url.Parse(ts.URL + r.URL.String())
			if err != nil {
				return nil, err
			}
			r.URL = u
			return ts.Client().Transport.RoundTrip(r)
		}),
	}
	rtt := 1 * time.Second

	// set the dependency
	pingpongService := pingpongmock.New(func(ctx context.Context, address swarm.Address, msgs ...string) (time.Duration, error) {
		fmt.Println(1)
		return rtt, nil
	})

	s.SetDependency(pingpongService)

	t.Run("ok", func(t *testing.T) {
		jsonhttptest.Request(t, client, http.MethodPost, "/pingpong/abcd", http.StatusNotFound,
			jsonhttptest.WithExpectedJSONResponse(jsonhttp.StatusResponse{
				Code:    http.StatusNotFound,
				Message: "peer not found",
			}),
		)
	})

}
