# Decentr
![go version](https://img.shields.io/github/go-mod/go-version/Decentr-net/decentr?color=blue)
[![testnet version](https://img.shields.io/badge/testnet%20version-v1.5.6-blue.svg)](https://shields.io/)
[![mainnet version](https://img.shields.io/badge/mainnet%20version-v1.5.8-brightgreen.svg)](https://shields.io/)
![latest version](https://img.shields.io/github/v/tag/Decentr-net/decentr?label=latest%20version&color=yellow)

Decentr blockchain

## Run Local Node Quick Start
This assumes that you're running Linux or MacOS and have installed [Go 1.16+](https://golang.org/dl/).  This guide helps you:

* build and install Decentr
* allow you to name your node
* add seeds to your config file
* download genesis state
* start your node
* use decentrdcli to check the status of your node.


If you already have a previous version of Decentr installed:
```
rm -rf ~/.decentr
```

### Mainnet

Build, Install, and Name your Node:

```bash
# Clone Decentr from the latest release
git clone https://github.com/Decentr-net/decentr
# Enter the folder Decentr was cloned into
cd decentr && git checkout v1.5.8
# Compile and install Decentr
make install
# Initialize decentrd in ~/.decentrd and name your node
decentrd init <yournodenamehere>
```

Patch Seeds:

```bash
sed -E -i 's/seeds = \".*\"/seeds = \"7708addcfb9d4ff394b18fbc6c016b4aaa90a10a@ares.mainnet.decentr.xyz:26656,8a3485f940c3b2b9f0dd979a16ea28de154f14dd@calliope.mainnet.decentr.xyz:26656,87490fd832f3226ac5d090f6a438d402670881d0@euterpe.mainnet.decentr.xyz:26656,3261bff0b7c16dcf6b5b8e62dd54faafbfd75415@hera.mainnet.decentr.xyz:26656,5f3cfa2e3d5ed2c2ef699c8593a3d93c902406a9@hermes.mainnet.decentr.xyz:26656,a529801b5390f56d5c280eaff4ae95b7163e385f@melpomene.mainnet.decentr.xyz:26656,385129dbe71bceff982204afa11ed7fa0ee39430@poseidon.mainnet.decentr.xyz:26656,35a934228c32ad8329ac917613a25474cc79bc08@terpsichore.mainnet.decentr.xyz:26656,0fd62bcd1de6f2e3cfc15852cdde9f3f8a7987e4@thalia.mainnet.decentr.xyz:26656,bd99693d0dbc855b0367f781fb48bf1ca6a6a58b@zeus.mainnet.decentr.xyz:26656\"/' $HOME/.decentr/config/config.toml
```

Download Genesis, Start your Node, Check your Node Status:

```bash
# Download genesis.json
wget -O $HOME/.decentr/config/genesis.json https://raw.githubusercontent.com/Decentr-net/mainnets/master/3.0/genesis.json
# Start Decentrd
decentrd start
# Check your node's status
decentrd status
```

Welcome to the Decentr Mainnet!

### Testnet

Build, Install, and Name your Node:

```bash
# Clone Decentr from the latest release
git clone -b v1.5.6 https://github.com/Decentr-net/decentr
# Enter the folder Decentr was cloned into
cd decentr
# Compile and install Decentr
make install
# Initialize decentrd in ~/.decentrd and name your node
decentrd init <yournodenamehere>
```

Patch Seeds:

```bash
sed -E -i 's/seeds = \".*\"/seeds = \"ff015e51856e77f5bad502dc8c38bcefeeffe55e@ares.testnet.decentr.xyz:26656,6a7e6b3db1dccffa9cc89b71216461aecac9c5f0@hera.testnet.decentr.xyz:26656,facd1cec2d372faa0381ae206aa2c13c16f4ea48@hermes.testnet.decentr.xyz:26656,26e21815870cc2ea62136b0cfd70184dfec0634f@poseidon.testnet.decentr.xyz:26656,7f338a46cc175248ea8933875df092cdaea8a0e9@zeus.testnet.decentr.xyz:26656\"/' $HOME/.decentr/config/config.toml
```

Download Genesis, Start your Node, Check your Node Status:

```bash
# Download genesis.json
wget -O $HOME/.decentr/config/genesis.json https://raw.githubusercontent.com/Decentr-net/testnets/master/1.5.0/genesis.json
# Start Decentrd
decentrd start
# Check your node's status
decentrd status
```

Welcome to the Decentr Testnet!

## Dev tools

### Requirements
To build project you should have:
- go >= 1.16
- docker

### Guide

To fetch last proto 3rd party
```
make proto-update-deps
```

To generate go models from proto
```
make proto-gen
```

To generate swagger from proto
```
make proto-swagger-gen
```

### Scripts
- [scripts/protocgen.sh](scripts/protocgen.sh)
generates goproto
- [scripts/protoc-swagger-gen.sh](scripts/protoc-swagger-gen.sh)
  generates swagger  
- [Dockerfile](scripts/Dockerfile)
  decentr docker image
- [buildtools.Dockerfile](scripts/buildtools.Dockerfile)
  docker image used in makefile (contains proto compilers and swagger-combine)
  
## Follow us!
Your data is value. Decentr makes your data payable and tradeable online.
* [Medium](https://medium.com/@DecentrNet)
* [Twitter](https://twitter.com/DecentrNet)
* [Telegram](https://t.me/DecentrNet)
* [Discord](https://discord.gg/9cSxwKyEjR)
