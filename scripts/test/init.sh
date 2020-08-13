#!/bin/sh

if [[ ! -f "/root/.decentrd/config/genesis.json" ]]; then
    if [[ -d "/decentr/config/" ]]; then
        mkdir -p /root/.decentrd/ && cp -rf /decentr/config /root/.decentrd/config
    else
        decentrd init $1
        cp -f /decentr/genesis.json /root/.decentrd/config/genesis.json
    fi
fi

if [[ ! -f "/root/.decentrd/data" ]]; then
    decentrd unsafe-reset-all
fi