package api

import (
	"context"

	"github.com/bytom/errors"
	"github.com/bytom/protocol/bc"
	"github.com/bytom/protocol/bc/types"
)

// SubmitWorkReq used to submitWork req
type GetWorkReq struct {
	UseJson bool `json:"use_json,omitempty"`
}

// func (a *API) getWork(ctx context.Context, req *GetWorkReq) Response {
func (a *API) getWork(in ...interface{}) Response {
	if len(in)>0 {
		panic("fffffffffffffffffffffffffffffffffk")
	}

	work, err := a.GetWork()
	if err != nil {
		return NewErrorResponse(err)
	}
	return NewSuccessResponse(work)
}

func (a *API) getWorkJson() Response {
	work, err := a.GetWorkJson()
	if err != nil {
		return NewErrorResponse(err)
	}
	return NewSuccessResponse(work)
}


// SubmitWorkReq used to submitWork req
type SubmitWorkReq struct {
	BlockHeader *types.BlockHeader `json:"block_header"`
}

func (a *API) submitWork(ctx context.Context, req *SubmitWorkReq) Response {
	if err := a.SubmitWork(req.BlockHeader); err != nil {
		return NewErrorResponse(err)
	}
	return NewSuccessResponse(true)
}

// GetWorkResp is resp struct for API
type GetWorkResp struct {
	BlockHeader *types.BlockHeader `json:"block_header"`
	Seed        *bc.Hash           `json:"seed"`
}

// GetWorkRespJson is resp struct for API
type GetWorkJsonResp struct {
	BlockHeaderJson *types.BlockHeaderJson 	`json:"block_header"`
	Seed        	*bc.Hash           		`json:"seed"`
}

// GetWork get work
func (a *API) GetWork() (*GetWorkResp, error) {
	bh, err := a.miningPool.GetWork()
	if err != nil {
		return nil, err
	}

	seed, err := a.chain.CalcNextSeed(&bh.PreviousBlockHash)
	if err != nil {
		return nil, err
	}

	return &GetWorkResp{
		BlockHeader: bh,
		Seed:        seed,
	}, nil
}

// GetWorkJson get work json
func (a *API) GetWorkJson() (*GetWorkJsonResp, error) {
	bhj, err := a.miningPool.GetWorkJson()
	if err != nil {
		return nil, err
	}

	seed, err := a.chain.CalcNextSeed(&bhj.PreviousBlockHash)
	if err != nil {
		return nil, err
	}

	return &GetWorkJsonResp{
		BlockHeaderJson: 	bhj,
		Seed:        		seed,
	}, nil
}

// SubmitWork submit work
func (a *API) SubmitWork(bh *types.BlockHeader) error {
	return a.miningPool.SubmitWork(bh)
}

func (a *API) setMining(in struct {
	IsMining bool `json:"is_mining"`
}) Response {
	if in.IsMining {
		return a.startMining()
	}
	return a.stopMining()
}

func (a *API) startMining() Response {
	a.cpuMiner.Start()
	if !a.IsMining() {
		return NewErrorResponse(errors.New("Failed to start mining"))
	}
	return NewSuccessResponse("")
}

func (a *API) stopMining() Response {
	a.cpuMiner.Stop()
	if a.IsMining() {
		return NewErrorResponse(errors.New("Failed to stop mining"))
	}
	return NewSuccessResponse("")
}
