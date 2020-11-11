# decentr
![img](https://img.shields.io/docker/cloud/build/decentr/decentr.svg)

Decentr blockchain


## Testnet Full Node Quick Start
This assumes that you're running Linux or MacOS and have installed [Go 1.14+](https://golang.org/dl/).  This guide helps you:

* build and install Decentr
* allow you to name your node
* add seeds to your config file
* download genesis state
* start your node
* use decentrdcli to check the status of your node.

Build, Install, and Name your Node:

```bash
# Clone Decentr from the latest release found here: https://github.com/Decentr-net/decentr/releases
git clone -b <latest_release> https://github.com/Decentr-net/decentr
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

Scroll down to `seeds` in `config.toml`, and add some of these seeds as a comma-separated list:

```c
253f5dd5309e47cda8c1ad4db9c5cd056ec54f88@ares.testnet.decentr.xyz:26656
25635874579827697da420227dcebe7424238eaa@hera.testnet.decentr.xyz:26656
b98336c214cd7319fa79cef40f52f681d71fbfbb@hermes.testnet.decentr.xyz:26656
3bddcf4f6856b3546a10c7574c614f47c499ca43@poseidon.testnet.decentr.xyz:26656
37c2b10599ed9926aacb96639f2c8f54dc674db2@zeus.testnet.decentr.xyz:26656
```

Download Genesis, Start your Node, Check your Node Status:

```bash
# Download genesis.json
wget -O $HOME/.decentrd/config/genesis.json https://raw.githubusercontent.com/Decentr-net/testnets/master/1.0/genesis.json
# Start Decentrd
decentrd start --cerberus-addr https://cerberus.testnet.decentr.xyz
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

### REST transactions
If you want to use REST to create tx, you should get it from rest service, then sign and broadcast it.

#### Example
Get tx body
```bash
curl -XPOST -s http://localhost:1317/profile/public/$(decentrcli keys show jack -a) \
     -d '{"base_req":{"chain_id":"testnet", "from": "'$(decentrcli keys show jack -a)'"},"public": { "firstName": "foo","lastName": "bar","avatar": "https://avatars3.githubusercontent.com/u/1526177","gender": "female","birthday": "2001-02-01"} }' > unsignedTx.json
```
  
unsignedTx.json will contain
```json
{"type":"cosmos-sdk/StdTx","value":{"msg":[{"type":"profile/SetPublic","value":{"owner":"decentr1z4z94y4lf33tdk4qvwh237ly8ngyjv5my6xqrw","public":{ "firstName": "foo","lastName": "bar","avatar": "https://avatars3.githubusercontent.com/u/1526177","gender": "female","birthday": "2001-02-01"} }}],"fee":{"amount":[],"gas":"200000"},"signatures":null,"memo":""}}
```
  
Then sign this transaction
```bash
decentrcli tx sign unsignedTx.json --from jack --offline --chain-id testnet --sequence 1 --account-number 3 > signedTx.json
```
  
And finally broadcast the signed transaction
```bash
decentrcli tx broadcast signedTx.json
```

## PDV (Personal Data Value) Data

#### CLI
```bash
# Query pdv owner by its address
decentrcli query pdv owner <address>

# Query pdv full
decentrcli query pdv show <address>

# List account's pdv
decentrcli query pdv list <owner> [page] [limit]

# Get cerberus address
decentrcli query pdv cerberus

# Create pdv
decentrcli tx pdv create [pdv] --from [account]
```

#### REST
```bash
# Query pdv owner by its address
curl -s http://localhost:1317/pdv/{address}/owner

# Query pdv full
curl -s http://localhost:1317/pdv/{address}/show

# List account's pdv
curl -s http://localhost:1317/pdv/{owner}/list

# List account's daily stats
curl -s http://localhost:1317/pdv/{owner}/stats

# Get cerberus address
curl -s http://localhost:1317/pdv/cerberus-addr

# Create PDV
curl -XPOST -s http://localhost:1317/pdv \ 
     -d '{"base_req":{"chain_id":"testnet", "from": "'$(decentrcli keys show jack -a)'"},"address": "address from cerberus"}' > unsignedTx.json
```

## PDV Token
PDV tokens are assigned to the user as soon as they reveal their personal data. 
There are no transactions, only query to get PDV token balance of the specific user.

#### CLI
```bash
# Query pdv token balance
decentrcli query token balance [address]
```

#### REST
```bash
curl -s http://localhost:1317/token/balance/{address}
```

### Profile
User profile consists of two parts: private and public. Private data is encrypted with user's private key.
Public one includes gender and birthday.

#### CLI
```bash
# Query private profile. Returns base64 encode string.
decentrcli query profile private [address]

# Query public profile
decentrcli query profile public [address] 

# Set private profile data that you own. The data should be encrypted with your private key beforehead.
decentrcli tx profile set-private [data] --from [account]

# Set public profile data that you own. Public profile are attributes: gender, birth date.
# Birthday date format is yyyy-mm-dd. Gender: male, female
decentrcli tx profile set-public '{ "firstName": "foo", "lastName": "bar", "avatar": "https://avatars3.githubusercontent.com/u/1526177", "gender": "female", "birthday": "2001-02-01"}' --from [account]
```

#### REST
To execute REST command decentrcli has to be run as a REST server `decentrcli rest-server` 

```bash
### Get account info
curl -s http://localhost:1317/auth/accounts/$(decentrcli keys show jack -a)
# > {"value": { "address": "decentr1d7narytgsy5lj2agt0t8sntaq3p8ucjhermqjj","coins": [], "public_key": "decentrpub1addwnpepq2jqxxu853rh0pa0agnkaxwaz6qdz6kpd4esqpw33sz3mp3a6mwh5eejl8q", "account_number": 3,"sequence": 6 }}

# Query private profile. Returns base64 encode string.
curl -s http://localhost:1317/profile/private/$(decentrcli keys show jack -a)
# > {"height": "0", "result": "YldWbllXaGxjbm9L"}

# Query public profile.
curl -s http://localhost:1317/profile/public/$(decentrcli keys show jack -a)
# > { "height": "0", "result": { "firstName": "foo", "lastName": "bar", "avatar": "https://avatars3.githubusercontent.com/u/1526177", "gender": "female", "birthday": "2001-02-01"}}

# Set private profile
curl -XPOST -s http://localhost:1317/profile/private/$(decentrcli keys show jack -a) \ 
     -d '{"base_req":{"chain_id":"testnet", "from": "'$(decentrcli keys show jack -a)'"},"private": "YldWbllXaGxjbm9L"}' > unsignedTx.json

# Set public profile
curl -XPOST -s http://localhost:1317/profile/public/$(decentrcli keys show jack -a) \
     -d '{"base_req":{"chain_id":"testnet", "from": "'$(decentrcli keys show jack -a)'"},"public": { "firstName": "foo","lastName": "bar","avatar": "https://avatars3.githubusercontent.com/u/1526177","gender": "female","birthday": "2001-02-01"} }' > unsignedTx.json
```

## Community module

#### CLI
```bash
# Create post
decentrcli tx community create-post [text] --title [title] --image-preview [url to preview] --from [account]

# Delete post
decentrcli tx community delete-post [uuid] --from [account]
```

#### REST
```bash
# Create post
curl -XPOST -s http://localhost:1317/post \ 
     -d '{"base_req":{"chain_id":"testnet", "from": "'$(decentrcli keys show jack -a)'"},"text": "my new post's text", "title":"my first post title", "imagePreview": "https://localhost/mypicture.jpg"}' > unsignedTx.json

# Delete post
curl -XPOST -s http://localhost:1317/post/{address}/{uuid}/delete \ 
     -d '{"base_req":{"chain_id":"testnet", "from": "'$(decentrcli keys show jack -a)'"}}' > unsignedTx.json
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
