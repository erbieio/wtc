#!/bin/sh

rm data/log
rm -r data/geth
geth --datadir data/ init ./gensis.json
