#!/bin/sh

nohup geth --datadir data/ --syncmode 'full' --gcmode archive --nodiscover --http.api admin,eth,debug,net,web3 --http --http.addr 0.0.0.0 --http.corsdomain '*' --allow-insecure-unlock -unlock `cat data/account` --password data/password --mine > data/log 2>&1 &
