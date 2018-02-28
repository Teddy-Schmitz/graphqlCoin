// Package rpcclient contains the types and functions used to connect to a bitcoin compliant jsonrpc server.
package rpcclient

import (
	"bytes"
	"encoding/json"
	"net/http"

	"fmt"

	"context"

	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
)

// JSONRPCRequest is used to create a jsonrpc request object to send to the server.
type JSONRPCRequest struct {
	// By default always set to 1
	JSONRPC int `json:"jsonrpc"`

	// The id of the client making the request
	ID string `json:"id"`

	// The jsonrpc method being called
	Method string `json:"method"`

	// A set of named variables passed to the jsonrpc server
	Params map[string]interface{} `json:"params"`
}

// JSONRPCResponse is used to hold the response sent back from the server and is useful for checking for errors before parsing the result.
type JSONRPCResponse struct {
	Result json.RawMessage `json:"result"`
	Error  *JSONRPCError   `json:"error,omitempty"`
	ID     string          `json:"id"`
}

// JSONRPCError is populated if the server returns an error
type JSONRPCError struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

// Makes a new request object with the default variables
func makeNewRequest() *JSONRPCRequest {
	j := &JSONRPCRequest{
		JSONRPC: 1,
		ID:      "graphqlCoin-v0.1",
	}

	return j
}

// HTTPClient is an interface for any valid http client.
type HttpClient interface {
	Do(req *http.Request) (*http.Response, error)
}

type Client struct {
	// The hostname or IP of the jsonrpc server.
	Host string

	// The user to set for basic auth to the server
	User string

	// The password to set for basic auth to the server
	Password string

	// The HttpClient to use for making requests.  Most users can set `http.DefaultClient`
	// but this allows for custom clients to allow for customized retry logic, etc.
	Client HttpClient
}

// Create and send a request to the jsonrpc server
func (r *Client) newHTTPRequest(method string, params map[string]interface{}) (*JSONRPCResponse, error) {

	j := makeNewRequest()
	j.Method = method
	j.Params = params
	body := new(bytes.Buffer)
	json.NewEncoder(body).Encode(j)
	logrus.Debugln("making new http request")
	req, err := http.NewRequest("POST", fmt.Sprintf("http://%s", r.Host), body)
	if err != nil {
		return nil, err
	}
	req.SetBasicAuth(r.User, r.Password)
	req.Header.Set("Content-Type", "text/plain")

	resp, err := r.Client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	jResponse := &JSONRPCResponse{}

	if err := json.NewDecoder(resp.Body).Decode(jResponse); err != nil {
		return nil, err
	}
	logrus.WithFields(logrus.Fields{"response": string(jResponse.Result), "error": jResponse.Error}).Debugln("jsonrpc response")

	if jResponse.Error != nil {
		return nil, errors.Errorf("jsonrpc error code: %v msg: %v", jResponse.Error.Code, jResponse.Error.Message)
	}

	return jResponse, nil
}

// GetBlock is used to retrieve a blocks information by its hash
func (r *Client) GetBlock(blockhash string) (*Block, error) {
	jResponse, err := r.newHTTPRequest("getblock",
		map[string]interface{}{"blockhash": blockhash, "verbosity": 1})
	if err != nil {
		return nil, err
	}

	block := &Block{}

	if err := json.Unmarshal(jResponse.Result, block); err != nil {
		return nil, err
	}

	return block, nil
}

// GetBlockHash is used to retrieve the hash of a block at a certain block height
// useful when the hash isn't known.
func (r *Client) GetBlockHash(height int) (string, error) {
	resp, err := r.newHTTPRequest("getblockhash", map[string]interface{}{"height": height})
	if err != nil {
		return "", err
	}

	var hash string
	if err := json.Unmarshal(resp.Result, &hash); err != nil {
		return "", err
	}

	return hash, nil
}

// GetDifficulty returns the current difficulty seen by the server
func (r *Client) GetDifficulty(ctx context.Context) (float64, error) {
	resp, err := r.newHTTPRequest("getdifficulty", nil)
	if err != nil {
		return 0, err
	}

	var diff float64
	if err := json.Unmarshal(resp.Result, &diff); err != nil {
		return 0, err
	}

	return diff, nil
}

// GetTransaction retrieves a transactions information by its txid
func (r *Client) GetTransaction(txid string) (*Transaction, error) {
	resp, err := r.newHTTPRequest("getrawtransaction", map[string]interface{}{"txid": txid, "verbose": 1})
	if err != nil {
		return nil, err
	}

	trx := &Transaction{}
	if err := json.Unmarshal(resp.Result, trx); err != nil {
		return nil, err
	}
	return trx, nil
}

// GetEstimateFee gets the current fee estimate to allow a transaction to be added within n blocks.
func (r *Client) GetEstimateFee(ctx context.Context, blocks int) (*FeeEstimate, error) {
	resp, err := r.newHTTPRequest("estimatesmartfee", map[string]interface{}{"conf_target": blocks})
	if err != nil {
		return nil, err
	}

	fee := FeeEstimate{}
	if err := json.Unmarshal(resp.Result, &fee); err != nil {
		return nil, err
	}

	return &fee, nil
}

// GetMempool returns information on transactions in the servers memory pool
func (r *Client) GetMempool(ctx context.Context) ([]MemPoolTrx, error) {
	jResponse, err := r.newHTTPRequest("getrawmempool", map[string]interface{}{"verbose": true})
	if err != nil {
		return nil, err
	}

	trxs := map[string]MemPoolTrx{}
	if err := json.Unmarshal(jResponse.Result, &trxs); err != nil {
		return nil, err
	}

	res := []MemPoolTrx{}
	for id, t := range trxs {
		t.ID = id
		res = append(res, t)
	}

	return res, nil
}
