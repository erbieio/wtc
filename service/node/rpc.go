package node

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"reflect"

	"github.com/INFURA/go-ethlibs/jsonrpc"
	"github.com/INFURA/go-ethlibs/node"
)

type RPC struct {
	node.Client
}

// NewRPC connects RPC client to the given URL.
func NewRPC(rawurl string) (*RPC, error) {
	client, err := node.NewClient(context.Background(), rawurl)
	if err != nil {
		return nil, err
	}
	return &RPC{client}, nil
}

func (c *RPC) Call(result interface{}, method string, args ...interface{}) error {
	ctx := context.Background()
	return c.CallContext(ctx, result, method, args...)
}

func (c *RPC) CallContext(ctx context.Context, result interface{}, method string, args ...interface{}) error {
	if result != nil && reflect.TypeOf(result).Kind() != reflect.Ptr {
		return fmt.Errorf("call result parameter must be pointer or nil interface: %v", result)
	}

	request := jsonrpc.Request{
		ID:     jsonrpc.ID{Num: 1},
		Method: method,
		Params: jsonrpc.MustParams(args...),
	}

	response, err := c.Request(ctx, &request)
	if err != nil {
		return err
	}

	if response.Error != nil {
		return errors.New(string(*response.Error))
	}
	return json.Unmarshal(response.Result, &result)
}
