# Decentr
![go version](https://img.shields.io/github/go-mod/go-version/Decentr-net/decentr?color=blue)
[![testnet version](https://img.shields.io/badge/testnet%20version-v1.4.6-blue.svg)](https://shields.io/)
[![mainnet version](https://img.shields.io/badge/mainnet%20version-v1.4.6-brightgreen.svg)](https://shields.io/)
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
rm -rf ~/.decentrd
```

### Mainnet

Build, Install, and Name your Node:

```bash
# Clone Decentr from the latest release
git clone https://github.com/Decentr-net/decentr
# Enter the folder Decentr was cloned into
cd decentr && git checkout v1.4.5
# Compile and install Decentr
make install
# Initialize decentrd in ~/.decentrd and name your node
decentrd init <yournodenamehere>
```

Add Seeds:

```bash
sed -E -i 's/seeds = \".*\"/seeds = \"f9b77dd93f28d2a45b00d4e3041b89a3c08788ef@calliope.mainnet.decentr.xyz:26656,987b5ce87b1b922793069756f594533eedf0f060@euterpe.mainnet.decentr.xyz:26656,2caebc4dad8d2ff95400918572d455392e10a63c@hera.mainnet.decentr.xyz:26656,c37f32e202e13b0725515570f794b68573a6f58c@hermes.mainnet.decentr.xyz:26656,4520b3221c91fa98a947a4c7f518ba5aab4e5b08@melpomene.mainnet.decentr.xyz:26656,c17bc88591115e52a686811630ad8c053de19f83@poseidon.mainnet.decentr.xyz:26656,c4ba719d38c871a93fb06cbfe0891ab11fedb9f7@terpsichore.mainnet.decentr.xyz:26656,9e9e0243610fadc0f65d3d927e2d682d86f71ea9@thalia.mainnet.decentr.xyz:26656,e1f3ce208776ff1fad0e8190f5475b68e841d788@zeus.mainnet.decentr.xyz:26656\"/' $HOME/.decentrd/config/config.toml

```

Download Genesis, Start your Node, Check your Node Status:

```bash
# Download genesis.json
wget -O $HOME/.decentrd/config/genesis.json https://raw.githubusercontent.com/Decentr-net/mainnets/master/1.0/genesis.json
# Start Decentrd
decentrd start
# Check your node's status with decentrcli
decentrd status
```

At block 145000, you will need to update your binary to a new version:

```bash
#Enter the folder Decentr and change to the new version
cd $HOME/decentr && git checkout v1.4.6

# Compile and install Decentr
make install
``` 

Start your Node, Check your Node Status:

```bash
# Start Decentrd
decentrd start
# Check your node's status with decentrcli
decentrd status
```

Welcome to the Decentr Mainnet!

### Testnet

Build, Install, and Name your Node:

```bash
# Clone Decentr from the latest release
git clone -b v1.5.5 https://github.com/Decentr-net/decentr
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
nano ~/.decentr/config/config.toml
```

Scroll down to `seeds` in `config.toml`, and replace with

```
seeds = "6ae322ed10db6af2c9178ae62ef8d667cb42d23b@ares.testnet.decentr.xyz:26656,6d75c934b5eec42d3b6cabe648604a8354c87a76@hera.testnet.decentr.xyz:26656,8836831f518c23b127e4ebc032f457a43461778a@hermes.testnet.decentr.xyz:26656,a09f312c09ca158f70b025ee38396744e861543f@poseidon.testnet.decentr.xyz:26656,9a2509f640aefd0100e9e58ca1a0c4362815722d@zeus.testnet.decentr.xyz:26656"
```

Download Genesis, Start your Node, Check your Node Status:

```bash
# Download genesis.json
wget -O $HOME/.decentr/config/genesis.json https://raw.githubusercontent.com/Decentr-net/testnets/master/1.5.0/genesis.json
# Start Decentrd
decentrd start
# Check your node's status with decentrcli
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
