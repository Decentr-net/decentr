# decentr
Decentr blockchain


### Node
```bash
rm -rf ~/.decentrd
rm -rf ~/.decentrcli

decentrd help
decentrd init test --chain-id testnet

decentrd add-genesis-account $(decentrcli keys show jack -a) 100000000stake
decentrd gentx --name jack --keyring-backend test

decentrd collect-gentxs
decentrd validate-genesis

## start the node
decentrd start
```

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

decentrcli rest-server
# > I[2020-07-31|13:50:22.088] Starting application REST service (chain-id: "testnet")... module=rest-server 
# > I[2020-07-31|13:50:22.088] Starting RPC HTTP server on 127.0.0.1:1317   module=rest-server 
```

### Profile

#### CLI
```bash
# Query private profile. Returns base64 encode string.
decentrcli query profile private [address]

# Query public profile
decentrcli query profile public [address] 

# Set private profile data that you own. The data should be encrypted with your private key beforehead.
decentrcli tx profile set-private [data] --from [account]

# Set public profile data that you own. Public profile are attributes: gender, birth date.
# Birthday date format is yyyy-mm-dd. Gender: male, female, custom
decentrcli tx profile set-public '{"gender": "female", "birthday": "2019-12-12"}' --from [account]
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
# > {"height": "0", "result": { "gender": "female", "birthday": "2019-12-12"}}

# Set private profile
curl -XPOST -s http://localhost:1317/profile/private/$(decentrcli keys show jack -a) \ 
     -d '{"base_req":{"chain_id":"testnet", "from": "'$(decentrcli keys show jack -a)'"},"private": "YldWbllXaGxjbm9L"}' > unsignedTx.json
# > {"type":"cosmos-sdk/StdTx","value":{"msg":[{"type":"profile/SetPrivate","value":{"owner":"decentr1z4z94y4lf33tdk4qvwh237ly8ngyjv5my6xqrw","private":"YldWbllXaGxjbm9L"}}],"fee":{"amount":[],"gas":"200000"},"signatures":null,"memo":""}}

# Then sign this transaction
decentrcli tx sign unsignedTx.json --from jack --offline --chain-id testnet --sequence 1 --account-number 3 > signedTx.json

# And finally broadcast the signed transaction
decentrcli tx broadcast signedTx.json


# Set public profile
curl -XPOST -s http://localhost:1317/profile/public/$(decentrcli keys show jack -a) \
     -d '{"base_req":{"chain_id":"testnet", "from": "'$(decentrcli keys show jack -a)'"},"public": {"gender":"female", "birthday": "2001-02-01"} }' > unsignedTx.json

# > {"type":"cosmos-sdk/StdTx","value":{"msg":[{"type":"profile/SetPublic","value":{"owner":"decentr1z4z94y4lf33tdk4qvwh237ly8ngyjv5my6xqrw","public":{"gender":"female","birthday":"2001-02-01"}}}],"fee":{"amount":[],"gas":"200000"},"signatures":null,"memo":""}}

# Then sign this transaction
decentrcli tx sign unsignedTx.json --from jack --offline --chain-id testnet --sequence 1 --account-number 3 > signedTx.json

# And finally broadcast the signed transaction
decentrcli tx broadcast signedTx.json
```


### Build
```bash
make install
```
creates two binaries: decentrd (node) and decentrcli (cli)
