# decentr
![go version](https://img.shields.io/github/go-mod/go-version/Decentr-net/decentr?color=blue) 
[![network version](https://img.shields.io/badge/network%20version-v1.1.0-blue.svg)](https://shields.io/) 
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
# Clone Decentr from the latest release found here: https://github.com/Decentr-net/decentr/releases
# Replace <latest_release> with the latest Decentr version. Looks like v1.1.0 
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

Scroll down to `seeds` in `config.toml`, and replace with

```
seeds = "77c8bddfa5715c9f1dba62e586ecf775490ec83b@ares.testnet.decentr.xyz:26656,c550b6692b9122037a2e55dd4194e328033b3836@hera.testnet.decentr.xyz:26656,6945cdeeeddf7f7cd3ae58cce9ee5dfa6811130f@hermes.testnet.decentr.xyz:26656,a224755914665558c20a6cbd6eceeefc8e0bbe79@poseidon.testnet.decentr.xyz:26656,4b63a3430e3a8f5824983ec4f8b3f136c215aed0@zeus.testnet.decentr.xyz:26656"
```

Download Genesis, Start your Node, Check your Node Status:

```bash
# Download genesis.json
wget -O $HOME/.decentrd/config/genesis.json https://raw.githubusercontent.com/Decentr-net/testnets/master/1.1.0/genesis.json
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
decentrcli tx sign unsignedTx.json --from jack --chain-id testnet > signedTx.json
```
  
And finally broadcast the signed transaction
```bash
decentrcli tx broadcast signedTx.json
```

## PDV Token
PDV tokens are assigned to the user as soon as they reveal their personal data. 
There are no transactions, only query to get PDV token balance of the specific user.

#### CLI
```bash
# Query pdv token balance
decentrcli query token balance [address]

# Query pdv token stats
decentrcli query token stats [address]
```

#### REST
```bash
# Query pdv token balance
curl -s http://localhost:1317/token/balance/{address}

# Query pdv token stats
curl -s http://localhost:1317/token/stats/{address}
```

### Profile
User profile consists of two parts: private and public. Private data is encrypted with user's private key.
Public one includes first name, last name, avatar, gender and birthday.

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
# > { "height": "0", "result": { "firstName": "foo", "lastName": "bar", "avatar": "https://avatars3.githubusercontent.com/u/1526177", "gender": "female", "birthday": "2001-02-01", "registeredAt:"1607972947"}}

# Set private profile
curl -XPOST -s http://localhost:1317/profile/private/$(decentrcli keys show jack -a) \ 
     -d '{"base_req":{"chain_id":"testnet", "from": "'$(decentrcli keys show jack -a)'"},"private": "YldWbllXaGxjbm9L"}' > unsignedTx.json

# Set public profile
curl -XPOST -s http://localhost:1317/profile/public/$(decentrcli keys show jack -a) \
     -d '{"base_req":{"chain_id":"testnet", "from": "'$(decentrcli keys show jack -a)'"},"public": { "firstName": "foo","lastName": "bar","avatar": "https://avatars3.githubusercontent.com/u/1526177","gender": "female","birthday": "2001-02-01"} }' > unsignedTx.json
```
## Bank

#### CLI
```bash
# Create and sign a send tx
decentrcli tx send [from_key_or_address] [to_address] [amount]
```

#### REST
```bash
# Get balance
curl http://localhost:1317/bank/balances/{address} 

# Transfer coins (send coins to a address)
curl -X POST http://localhost:1317/bank/accounts/{address}/transfers \ 
    -d '{"base_req":{"chain_id":"testnet", "from": "'$(decentrcli keys show jack -a)'"}, "amount": [{"denom": "udec", "amount": "15"}]}' > unsignedTx.json
```

## Community

### Categories
| Value | Description |
| --- | --- |
| 1 | World News |
| 2 | Travel & Tourism |
| 3 | Science & Technology |
| 4 | Strange World |
| 5 | Arts & Entertainment |
| 6 | Writers& Writing |
| 7 | Health & Fitness |
| 8 | Crypto & Blockchain |
| 9 | Sports |


### Likes weight
|Value | Description |
| --- | --- |
| 1 | Like |
| 0 | Delete |
| -1 | Dislike |

#### CLI
```bash
# Create post
decentrcli tx community create-post [text] --title [title] --preview-image [url to preview] --category 2 --from [account]

# Delete post
decentrcli tx community delete-post [postOwner] [postUUID] --from [account]

# Like post
decentrcli tx community like-post [postOwner] [postUUID] --weight [weight] --from [account]

# Get moderator accounts addresses
decentrcli query community moderators

# Get user's posts
decentrcli query community user-posts <account> [--from-uuid uuid] [--limit int]

# Get a single post
decentrcli query community post <owner> <uuid>

# Get the recent posts
decentrcli query community posts [--from-owner account --from-uuid uuid] [--category int] [--limit int]

# Get the most popular posts
decentrcli query community popular-posts [--from-owner account --from-uuid uuid] [--category int] [--limit int] [--interval day/week/month]

# Get user's likes
decentrcli query community user-liked-posts [owner]

# Follow
decentrcli tx community follow [whom account] --from [who account]

# Unfollow
decentrcli tx community unfollow [whom account] --from [who account]

# Get followee
decentrcli query community followee [account]   
```

#### REST
```bash
# Create post, note category is a quoted number.
curl -XPOST -s http://localhost:1317/community/posts \
     -d '{"base_req":{"chain_id":"testnet", "from": "'$(decentrcli keys show jack -a)'"},"text": "my brand new text", "title":"my first post title", "imagePreview": "https://localhost/mypicture.jpg", "category": "2"}' > unsignedTx.json

# Delete post
curl -XPOST -s http://localhost:1317/community/posts/{postOwner}/{postUUID}/delete \
     -d '{"base_req":{"chain_id":"testnet", "from": "'$(decentrcli keys show jack -a)'"}}' > unsignedTx.json

# Get moderator accouns addresses
curl -s http://localhost:1317/community/moderators

# Like post
curl -XPOST -s http://localhost:1317/community/posts/{postOwner}/{postUUID}/like\
     -d '{"base_req":{"chain_id":"testnet", "from": "'$(decentrcli keys show jack -a)'"}, "weight": 1}' > unsignedTx.json

# Get a single post
curl -s "http://localhost:1317/community/post/{owner}/{uuid}"

# Get the latest posts
curl -s "http://localhost:1317/community/posts?category={category}&limit={limit}&fromOwner={account}&fromUUID={post's uuid}"

# Get the most popular posts by day
curl -s "http://localhost:1317/community/posts/popular/byDay?category={category}&limit={limit}&fromOwner={account}&fromUUID={post's uuid}"

# Get the most popular posts by week
curl -s "http://localhost:1317/community/posts/popular/byWeek?category={category}&limit={limit}&fromOwner={account}&fromUUID={post's uuid}"

# Get the most popular posts by month
curl -s "http://localhost:1317/community/posts/popular/byMonth?category={category}&limit={limit}&fromOwner={account}&fromUUID={post's uuid}"

# Get user's posts
curl -s "http://localhost:1317/community/posts/{account}?from={postUUID}&limit={limit}"

# Get user's likes
curl -s "http://localhost:1317/community/likedPosts/{account}"


# Follow
curl -XPOST -s http://localhost:1317/community/followers/follow/{whom}\
     -d '{"base_req":{"chain_id":"testnet", "from": "'$(decentrcli keys show jack -a)'"}}' > unsignedTx.json

# Unfollow
curl -XPOST -s http://localhost:1317/community/followers/unfollow/{whom}\
     -d '{"base_req":{"chain_id":"testnet", "from": "'$(decentrcli keys show jack -a)'"}}' > unsignedTx.json

# Get followee
curl -s "http://localhost:1317/community/followers/{who}/followees"
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
