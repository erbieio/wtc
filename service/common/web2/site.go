package web2

import (
	"log"
	"math/big"
	"wtc/common/types"
)

// Site 链下网站操作对象
type Site struct {
	No int64
	// more field...
}

func NewSite() (w *Site, err error) {
	return &Site{}, nil
}

type Request struct {
	From    types.Address
	To      types.Address
	TokenId *big.Int
	Uri     types.Data
	// more field...
}

func (w *Site) LatestRequest() (r []*Request, err error) {
	// 以下是个模拟实现
	r = append(r, &Request{
		From:    "0x0000000000000000000000000000000000000000",
		To:      "0x0000000000000000001210000000121000000012",
		TokenId: nil,
		Uri:     "https://abc.com",
	})
	r = append(r, &Request{
		From:    "0x0000000000000000001210000000121000000012",
		To:      "0x0000000000000000000000000000000000000123",
		TokenId: big.NewInt(w.No),
	})
	w.No++
	return
}

type Response struct {
	From    types.Address
	To      types.Address
	TokenId *big.Int
	TxHash  types.Hash
	// more field...
}

func (w *Site) SetResponse(r *Response) (err error) {
	// 只是打印，后期需要将其发送到网站上
	log.Printf("NFT transfer from %v to %v id %v txHash %v\n", r.From, r.To, r.TokenId, r.TxHash)
	return
}
