# Decentr
![go version](https://img.shields.io/github/go-mod/go-version/Decentr-net/decentr?color=blue) 
[![testnet version](https://img.shields.io/badge/testnet%20version-v1.3.8-blue.svg)](https://shields.io/) 
[![mainnet version](https://img.shields.io/badge/mainnet%20version-v1.4.6-brightgreen.svg)](https://shields.io/) 
![latest version](https://img.shields.io/github/v/tag/Decentr-net/decentr?label=latest%20version&color=yellow)

Decentr blockchain

## Run Node Quick Start
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

### Mainnet

Build, Install, and Name your Node:

```bash
# Clone Decentr from the latest release
git clone -b v1.4.6 https://github.com/Decentr-net/decentr
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
seeds = "f9b77dd93f28d2a45b00d4e3041b89a3c08788ef@calliope.mainnet.decentr.xyz:26656,987b5ce87b1b922793069756f594533eedf0f060@euterpe.mainnet.decentr.xyz:26656,2caebc4dad8d2ff95400918572d455392e10a63c@hera.mainnet.decentr.xyz:26656,c37f32e202e13b0725515570f794b68573a6f58c@hermes.mainnet.decentr.xyz:26656,4520b3221c91fa98a947a4c7f518ba5aab4e5b08@melpomene.mainnet.decentr.xyz:26656,c17bc88591115e52a686811630ad8c053de19f83@poseidon.mainnet.decentr.xyz:26656,c4ba719d38c871a93fb06cbfe0891ab11fedb9f7@terpsichore.mainnet.decentr.xyz:26656,9e9e0243610fadc0f65d3d927e2d682d86f71ea9@thalia.mainnet.decentr.xyz:26656,e1f3ce208776ff1fad0e8190f5475b68e841d788@zeus.mainnet.decentr.xyz:26656"
```

Download Genesis, Start your Node, Check your Node Status:

```bash
# Download genesis.json
wget -O $HOME/.decentrd/config/genesis.json https://raw.githubusercontent.com/Decentr-net/mainnets/master/1.0/genesis.json
# Start Decentrd
decentrd start
# Check your node's status with decentrcli
decentrcli status
```

Welcome to the Decentr Mainnet!

### Testnet

Build, Install, and Name your Node:

```bash
# Clone Decentr from the latest release
git clone -b v1.3.8 https://github.com/Decentr-net/decentr
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
wget -O $HOME/.decentrd/config/genesis.json https://raw.githubusercontent.com/Decentr-net/testnets/master/1.3.8/genesis.json
# Start Decentrd
decentrd start
# Check your node's status with decentrcli
decentrcli status
```

Welcome to the Decentr Testnet!

## Follow us!
Your data is value. Decentr makes your data payable and tradeable online.
* [Medium](https://medium.com/@DecentrNet)
* [Twitter](https://twitter.com/DecentrNet)
* [Telegram](https://t.me/DecentrNet)
