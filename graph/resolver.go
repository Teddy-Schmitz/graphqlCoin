//The graph package contains the resolver as well as generated graph server and custom scalars.
package graph

import (
	"context"

	"github.com/Teddy-Schmitz/graphqlCoin/rpcclient"
	"github.com/pkg/errors"
)

type Resolver struct {
	Client *rpcclient.Client
}

func (r *Resolver) Query_block(ctx context.Context, hash *string, height *uint64) (rpcclient.Block, error) {
	var b *rpcclient.Block
	var err error

	if hash != nil {
		b, err = r.Client.GetBlock(*hash)
		if err != nil {
			return rpcclient.Block{}, err
		}
	}

	if height != nil {
		h, err := r.Client.GetBlockHash(int(*height))
		if err != nil {
			return rpcclient.Block{}, err
		}

		b, err = r.Client.GetBlock(h)
		if err != nil {
			return rpcclient.Block{}, err
		}
	}

	if b == nil {
		return rpcclient.Block{}, errors.New("block not found")
	}

	return *b, nil
}

func (r *Resolver) Query_difficulty(ctx context.Context) (float64, error) {
	return r.Client.GetDifficulty(ctx)
}

func (r *Resolver) Query_transaction(ctx context.Context, id string) (rpcclient.Transaction, error) {
	t, err := r.Client.GetTransaction(id)
	if err != nil {
		return rpcclient.Transaction{}, err
	}

	if t == nil {
		return rpcclient.Transaction{}, errors.New("trx not found")
	}

	return *t, nil
}

func (r *Resolver) Query_estimatefee(ctx context.Context, blocks int) (rpcclient.FeeEstimate, error) {

	if blocks < 1 || blocks > 1008 {
		return rpcclient.FeeEstimate{}, errors.New("blocks must be between 1 and 1008")
	}

	f, err := r.Client.GetEstimateFee(ctx, blocks)
	if err != nil {
		return rpcclient.FeeEstimate{}, err
	}

	if f == nil {
		return rpcclient.FeeEstimate{}, errors.New("estimate fee failed")
	}

	return *f, nil
}

func (r *Resolver) Query_mempool(ctx context.Context) ([]rpcclient.MemPoolTrx, error) {
	return r.Client.GetMempool(ctx)
}

func (r *Resolver) Block_transactions(ctx context.Context, obj *rpcclient.Block) ([]rpcclient.Transaction, error) {

	// genesis block doesnt technically have a transaction ( though u can retrieve it with a more verbose call to getblock)
	// https://github.com/bitcoin/bitcoin/issues/3303
	if obj.Height == 0 {
		return nil, nil
	}

	res := []rpcclient.Transaction{}
	for _, id := range obj.TrxIDs {
		trx, err := r.Client.GetTransaction(id)
		if err != nil {
			return nil, err
		}
		if trx == nil {
			return nil, errors.Errorf("transaction %s not found", id)
		}
		res = append(res, *trx)
	}

	return res, nil
}

func (r *Resolver) Transaction_block(ctx context.Context, obj *rpcclient.Transaction) (rpcclient.Block, error) {
	b, err := r.Client.GetBlock(obj.BlockHash)
	if err != nil {
		return rpcclient.Block{}, err
	}

	if b == nil {
		return rpcclient.Block{}, errors.Errorf("block %s not found", obj.BlockHash)
	}

	return *b, nil
}
