#!/bin/bash
rm -rf ~/.summa/cosmosrelay

if make install
then
    echo SUCCESSFULLY BUILT
else
    echo ERROR: CANNOT BUILD
    exit 1
fi

relayd init mynode --chain-id relay

echo "12345678" | relaycli keys add me

relayd add-genesis-account $(relaycli keys show me -a) 1000cbtc,100000000stake

relaycli config chain-id relay
relaycli config output json
relaycli config indent true
relaycli config trust-node true

echo "12345678" | relayd gentx --name me
relayd collect-gentxs

relayd start
