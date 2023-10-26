package service

import (
	"log"
	"time"
	"wtc/common/types"
	"wtc/common/web2"
	"wtc/common/web3"

	. "wtc/conf"
)

func init() {
	wtc, err := NewWTC()
	if err != nil {
		panic(err)
	}
	if err := wtc.Start(); err != nil {
		panic(err)
	}
}

// WTC web2到区块链的NFT服务
type WTC struct {
	contract *web3.Contract
	site     *web2.Site
	chanExit chan int
	chanConf chan time.Time
}

func NewWTC() (wtc *WTC, err error) {
	contract, err := web3.NewContract(ChainUrl, types.Address(HexAddr), HexKey)
	if err != nil {
		return
	}
	site, err := web2.NewSite()
	if err != nil {
		return
	}
	return &WTC{
		contract: contract,
		site:     site,
		chanExit: make(chan int, 1),
		chanConf: make(chan time.Time, 1),
	}, nil
}

func (w *WTC) loop() {
	ticker := time.NewTicker(Interval)
	for {
		select {
		case code := <-w.chanExit:
			ticker.Stop()
			log.Println("wtc stop code: ", code)
		case now := <-w.chanConf:

			log.Println("Update conf at time: ", now)
		case _ = <-ticker.C:
			// 获取最新的web2网站请求，铸造或转移NFT
			request, err := w.site.LatestRequest()
			if err != nil {
				log.Println("site service err:", err)
				continue
			}
			// 发送交易
			for _, r := range request {
				if r.TokenId == nil || r.TokenId.Uint64() == 0 {
					hash, err := w.contract.MintNFT(r.To, r.Uri)
					if err != nil {
						log.Println("contract service err:", err)
						continue
					}
					log.Println("contract service mint tx hash:", hash)
				} else {
					hash, err := w.contract.TransferNFT(r.From, r.To, r.TokenId)
					if err != nil {
						log.Println("contract service err:", err)
						continue
					}
					log.Println("contract service transfer tx hash:", hash)
				}

			}

			// 获取最新的链上NFT铸造或转移事件
			events, err := w.contract.LatestEvents()
			if err != nil {
				log.Println("contract service err:", err)
				continue
			}
			for _, e := range events {
				err := w.site.SetResponse(&web2.Response{
					From:    e.From,
					To:      e.To,
					TokenId: e.Id,
					TxHash:  e.TxHash,
				})
				if err != nil {
					log.Println("site service err:", err)
					continue
				}
			}
			log.Println("service running once")
		}
	}
}

func (w *WTC) Start() (err error) {
	w.chanConf, err = NewWatcher()
	if err != nil {
		return err
	}
	w.chanExit = make(chan int)
	go w.loop()
	return nil
}

func (w *WTC) Stop() {
	w.chanExit <- 0
	close(w.chanExit)
	close(w.chanConf)
}
