package service

import (
	"log"
	"time"

	. "wtc/conf"
)

var wtc *WTC

func init() {
	wtc = NewWTC()
	if err := wtc.Start(); err != nil {
		panic(err)
	}
}

// WTC web2到区块链的NFT服务
type WTC struct {
	contract *Contract
	chanExit chan int
	chanConf chan time.Time
}

func NewWTC() *WTC {
	var wtc = WTC{}

	return &wtc
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
		case now := <-ticker.C:
			log.Println(now)
			// todo: parse and send tx
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
