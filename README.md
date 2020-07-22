# decentr
Decentr blockchain

### Build
```bash
make install
```
creates two binaries: decentrd (node) and decentrcli (cli)

```bash
rm -rf ~/.decentrd
rm -rf ~/.decentrcli
```

### Node
```bash
decentrd help
decentrd init test --chain-id testnet

decentrd add-genesis-account $(decentrcli keys show megaherz -a) 1000pdvtoken,100000000stake
decentrd gentx --name megaherz --keyring-backend test

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
```
