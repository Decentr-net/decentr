# decentr
![go version](https://img.shields.io/github/go-mod/go-version/Decentr-net/decentr?color=blue) 
[![network version](https://img.shields.io/badge/network%20version-v1.2.5-blue.svg)](https://shields.io/) 
![candidate version](https://img.shields.io/github/v/tag/Decentr-net/decentr?label=candidate%20version&color=green)

Decentr blockchain


## Testnet Full Node Quick Start
This assumes that you're running Linux or MacOS and have installed [Go 1.15+](https://golang.org/dl/).  This guide helps you:

* build and install Decentr
* allow you to name your node
* add seeds to your config file
* download genesis state
* start your node
* use decentrdcli to check the status of your node.

If you already have a previous version of Decentr installed:
```
rm -rf ~/.decentrd
rm -rf ~/.decentrcli
```

Build, Install, and Name your Node:

```bash
# Clone Decentr from the latest release
git clone -b v1.2.5 https://github.com/Decentr-net/decentr
# Enter the folder Decentr was cloned into
cd decentr
# Compile and install Decentr
make install
# Initialize decentrd in ~/.decentrd and name your node
decentrd init yournodenamehere
```

Add Seeds:

```bash
# Edit config.toml
nano ~/.decentrd/config/config.toml
```

Scroll down to `seeds` in `config.toml`, and replace with

```
seeds = "95a70f0119af52e54697fa7feb8b09b4e7c7ec21@ares.testnet.decentr.xyz:26656,b6d499b2b0146627b9bf6f33a9a7e4013312c6d1@hera.testnet.decentr.xyz:26656,576d044b24cc449366850a95f7616f03ab8d14b3@hermes.testnet.decentr.xyz:26656,c98511455134b4450ebb20fce57308a9fb300b89@poseidon.testnet.decentr.xyz:26656,acc5524b4ff34591357a28d5fccf4efb5ad883c5@zeus.testnet.decentr.xyz:26656"
```

Download Genesis, Start your Node, Check your Node Status:

```bash
# Download genesis.json
wget -O $HOME/.decentrd/config/genesis.json https://raw.githubusercontent.com/Decentr-net/testnets/master/1.2.5/genesis.json
# Start Decentrd
decentrd start
# Check your node's status with decentrcli
decentrcli status
```

Welcome to the Decentr!

To start LCD (light-client daemon), a local REST server
```bash
decentrcli rest-server
# > I[2020-07-31|13:50:22.088] Starting application REST service (chain-id: "testnet")... module=rest-server 
# > I[2020-07-31|13:50:22.088] Starting RPC HTTP server on 127.0.0.1:1317   module=rest-server 
``` 
The server is available at `127.0.0.1:1317`

### CLI
```bash
decentrcli help
decentrcli config chain-id testnet
decentrcli config keyring-backend test 

decentrcli keys add megaherz
# > 
# {
#   "name": "megaherz",
#   "type": "local",
#   "address": "decentr1m8k9dy3962v8km0d5jwsqanwvf0h5fmj6f5zyp",
#   "pubkey": "decentrpub1addwnpepq2yrdqzcnleu2gr69c5zkw7laa4el7mtj8ala97s648wzlvegk7vcpsh6kg",
#   "mnemonic": "crouch goddess pass cigar conduct odor beach coil hole enroll fringe crane witness squeeze mention pioneer inmate wink concert laugh segment abuse tomorrow amused"
#  }
```


## Build
```bash
make install
```
creates two binaries: decentrd (node) and decentrcli (cli)

#### Build local image image
```
docker build -t decentr-local -f scripts/Dockerfile .
```
#### Start local testnet
```
cd scripts/test && docker-compose up
```

## Follow us!
Your data is value. Decentr makes your data payable and tradeable online.
* [Medium](https://medium.com/@DecentrNet)
* [Twitter](https://twitter.com/DecentrNet)
* [Telegram](https://t.me/DecentrNet)
