package conf

import (
	"log"
	"os"
	"strconv"
	"time"

	"github.com/fsnotify/fsnotify"
	"github.com/joho/godotenv"
)

// default allocation
var (
	ChainUrl = "http://127.0.0.1:8545"
	HexKey   = "7bbfec284ee43e328438d46ec803863c8e1367ab46072f7864c07e0a03ba61fd"
	HexAddr  = "0xfffffffffffffffffffffffffffffffffffffff0"
	Interval = 4 * time.Second
)

const confFile = ".env"

func init() {
	Overload()
}

func NewWatcher() (chan time.Time, error) {
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		return nil, err
	}
	ch := make(chan time.Time)
	watcher.Add(confFile)
	go func() {
		defer watcher.Close()
		lastTime := time.Now()
		for {
			select {
			case e := <-watcher.Events:
				// anti-shake
				if now := time.Now(); e.Op == fsnotify.Write && now.Sub(lastTime) > time.Second {
					Overload()
					lastTime = now
					ch <- now
				}
			case err := <-watcher.Errors:
				log.Println("watch conf file err: ", err)
			}
		}
	}()
	return ch, nil
}

// Overload read configuration to override default value
func Overload() {
	err := godotenv.Overload(confFile)
	if err != nil {
		log.Println("Failed to load environment variables from env file,", err)
	}

	// Parse the basic configuration of the server
	if chainUrl := os.Getenv("CHAIN_URL"); chainUrl != "" {
		ChainUrl = chainUrl
	}
	if hexKey := os.Getenv("HEX_KEY"); hexKey != "" {
		HexKey = hexKey
	}
	if hexAddr := os.Getenv("HEX_ADDR"); hexAddr != "" {
		HexAddr = hexAddr
	}
	if interval := os.Getenv("INTERVAL"); interval != "" {
		num, err := strconv.ParseUint(interval, 10, 64)
		if err == nil && num > 0 {
			Interval = time.Duration(num) * time.Second
		} else {
			log.Println("Failed to set interval from env file,", err)
		}
	}
}
