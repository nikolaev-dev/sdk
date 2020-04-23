package grpc_client

import (
	"context"
	"github.com/MinterTeam/node-grpc-gateway/api_pb"
	"github.com/golang/protobuf/ptypes/empty"
	"google.golang.org/grpc"
	"strconv"
)

type Client struct {
	grpcClient api_pb.ApiServiceClient
	ctxFunc    func() context.Context
}

func New(address string) *Client {
	clientConn, err := grpc.Dial(address, grpc.WithInsecure())
	if err != nil {
		panic(err)
	}

	return &Client{grpcClient: api_pb.NewApiServiceClient(clientConn), ctxFunc: context.Background}
}

func (c *Client) WithContextFunc(contextFunc func(context.Context) func() context.Context) *Client {
	return &Client{grpcClient: c.grpcClient, ctxFunc: contextFunc(c.ctxFunc())}
}

func (c *Client) GRPCClient() api_pb.ApiServiceClient {
	return c.grpcClient
}

func (c *Client) Halts(height int) (*api_pb.HaltsResponse, error) {
	return c.grpcClient.Halts(c.ctxFunc(), &api_pb.HaltsRequest{Height: uint64(height)})
}

func (c *Client) Genesis() (*api_pb.GenesisResponse, error) {
	return c.grpcClient.Genesis(c.ctxFunc(), &empty.Empty{})
}

// Returns next transaction number (nonce) of an address.
func (c *Client) Nonce(address string) (uint64, error) {
	status, err := c.Address(address)
	if err != nil {
		return 0, err
	}

	transactionsCount, err := strconv.Atoi(status.TransactionsCount)
	if err != nil {
		return 0, err
	}

	return uint64(transactionsCount) + 1, err
}

func (c *Client) Status() (*api_pb.StatusResponse, error) {
	return c.grpcClient.Status(c.ctxFunc(), &empty.Empty{})
}

// Returns coins list, balance and transaction count of an address.
func (c *Client) Address(address string, optionalHeight ...int) (*api_pb.AddressResponse, error) {
	return c.grpcClient.Address(c.ctxFunc(), &api_pb.AddressRequest{Height: optionalInt(optionalHeight), Address: address})
}

func (c *Client) Addresses(addresses []string, optionalHeight ...int) (*api_pb.AddressesResponse, error) {
	return c.grpcClient.Addresses(c.ctxFunc(), &api_pb.AddressesRequest{Addresses: addresses, Height: optionalInt(optionalHeight)})
}

// Returns block data at given height.
func (c *Client) Block(height int) (*api_pb.BlockResponse, error) {
	return c.grpcClient.Block(c.ctxFunc(), &api_pb.BlockRequest{Height: uint64(height)})
}

// Returns candidate’s info by provided public_key. It will respond with 404 code if candidate is not found.
func (c *Client) Candidate(publicKey string, optionalHeight ...int) (*api_pb.CandidateResponse, error) {
	return c.grpcClient.Candidate(c.ctxFunc(), &api_pb.CandidateRequest{Height: optionalInt(optionalHeight), PublicKey: publicKey})
}

// Returns list of candidates.
func (c *Client) Candidates(includeStakes bool, optionalHeight ...int) (*api_pb.CandidatesResponse, error) {
	return c.grpcClient.Candidates(c.ctxFunc(), &api_pb.CandidatesRequest{Height: optionalInt(optionalHeight), IncludeStakes: includeStakes})
}

// Returns information about coin. Note: this method does not return information about base coins (MNT and BIP).
func (c *Client) CoinInfo(symbol string, optionalHeight ...int) (*api_pb.CoinInfoResponse, error) {
	return c.grpcClient.CoinInfo(c.ctxFunc(), &api_pb.CoinInfoRequest{Height: optionalInt(optionalHeight), Symbol: symbol})
}

// Return estimate of buy coin transaction.
func (c *Client) EstimateCoinBuy(coinToSell, coinToBuy, valueToBuy string, optionalHeight ...int) (*api_pb.EstimateCoinBuyResponse, error) {
	return c.grpcClient.EstimateCoinBuy(c.ctxFunc(), &api_pb.EstimateCoinBuyRequest{Height: optionalInt(optionalHeight), CoinToSell: coinToSell, CoinToBuy: coinToBuy, ValueToBuy: valueToBuy})
}

// Return estimate of sell coin transaction.
func (c *Client) EstimateCoinSell(coinToBuy, coinToSell, valueToBuy string, optionalHeight ...int) (*api_pb.EstimateCoinSellResponse, error) {
	return c.grpcClient.EstimateCoinSell(c.ctxFunc(), &api_pb.EstimateCoinSellRequest{Height: optionalInt(optionalHeight), CoinToBuy: coinToBuy, CoinToSell: coinToSell, ValueToSell: valueToBuy})
}

// Return estimate of sell all coin transaction.
func (c *Client) EstimateCoinSellAll(coinToBuy, coinToSell, valueToBuy string, gasPrice int, optionalHeight ...int) (*api_pb.EstimateCoinSellAllResponse, error) {
	return c.grpcClient.EstimateCoinSellAll(c.ctxFunc(), &api_pb.EstimateCoinSellAllRequest{Height: optionalInt(optionalHeight), CoinToBuy: coinToBuy, CoinToSell: coinToSell, ValueToSell: valueToBuy, GasPrice: uint64(gasPrice)})
}

// Return estimate of transaction.
func (c *Client) EstimateTxCommission(tx string, optionalHeight ...int) (*api_pb.EstimateTxCommissionResponse, error) {
	return c.grpcClient.EstimateTxCommission(c.ctxFunc(), &api_pb.EstimateTxCommissionRequest{Height: optionalInt(optionalHeight), Tx: tx})
}

// Returns events at given height.
func (c *Client) Events(optionalHeight ...int) (*api_pb.EventsResponse, error) {
	return c.grpcClient.Events(c.ctxFunc(), &api_pb.EventsRequest{Height: optionalInt(optionalHeight)})
}

// Returns current max gas.
func (c *Client) MaxGas(optionalHeight ...int) (*api_pb.MaxGasResponse, error) {
	return c.grpcClient.MaxGas(c.ctxFunc(), &api_pb.MaxGasRequest{Height: optionalInt(optionalHeight)})
}

// Returns current min gas price.
func (c *Client) MinGasPrice() (*api_pb.MinGasPriceResponse, error) {
	return c.grpcClient.MinGasPrice(c.ctxFunc(), &empty.Empty{})
}

// Returns missed blocks by validator public key.
func (c *Client) MissedBlocks(publicKey string, optionalHeight ...int) (*api_pb.MissedBlocksResponse, error) {
	return c.grpcClient.MissedBlocks(c.ctxFunc(), &api_pb.MissedBlocksRequest{Height: optionalInt(optionalHeight), PublicKey: publicKey})
}

// Returns network info
func (c *Client) NetInfo() (*api_pb.NetInfoResponse, error) {
	return c.grpcClient.NetInfo(c.ctxFunc(), &empty.Empty{})
}

// Returns the result of sending signed tx.
func (c *Client) SendTransaction(tx string) (*api_pb.SendTransactionResponse, error) {
	return c.grpcClient.SendGetTransaction(c.ctxFunc(), &api_pb.SendTransactionRequest{Tx: tx})
}

// Returns transaction info.
func (c *Client) Transaction(hash string) (*api_pb.TransactionResponse, error) {
	return c.grpcClient.Transaction(c.ctxFunc(), &api_pb.TransactionRequest{Hash: hash})
}

// Return transactions by query.
func (c *Client) Transactions(query string, page, perPage int) (*api_pb.TransactionsResponse, error) {
	return c.grpcClient.Transactions(c.ctxFunc(), &api_pb.TransactionsRequest{Query: query, Page: int32(page), PerPage: int32(perPage)})
}

// Returns unconfirmed transactions.
func (c *Client) UnconfirmedTxs(limit ...int) (*api_pb.UnconfirmedTxsResponse, error) {
	return c.grpcClient.UnconfirmedTxs(c.ctxFunc(), &api_pb.UnconfirmedTxsRequest{Limit: int32(optionalInt(limit))})
}

// Returns list of active validators.
func (c *Client) Validators(page, perPage int, limit ...int) (*api_pb.ValidatorsResponse, error) {
	return c.grpcClient.Validators(c.ctxFunc(), &api_pb.ValidatorsRequest{Height: optionalInt(limit), Page: int32(page), PerPage: int32(perPage)})
}

// Returns a subscription for events by query
func (c *Client) Subscribe(query string) (api_pb.ApiService_SubscribeClient, error) {
	return c.grpcClient.Subscribe(c.ctxFunc(), &api_pb.SubscribeRequest{Query: query})
}

func optionalInt(height []int) uint64 {
	if len(height) == 1 {
		return uint64(height[0])
	}
	return 0
}
