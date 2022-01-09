## Decentr-1.5.0 Upgrade Instructions 

1. Backup .decentrd and .decentrcli
```shell
   cp -rf $HOME/.decentrd $HOME/.decentrd.bak
   cp -rf $HOME/.decentrcli $HOME/.decentrcli.bak
```
2. Export your genesis file with the next command:
```shell
   decentrd export > $HOME/.decentrd/genesis.json
```
3. Export your key
```shell
   decentrcli keys export <name> > <name>.key
```
4. Get new version and install
5. Migrate your genesis with next command
```shell
   decentrd migrate $HOME/.decentrd/genesis.json > /tmp/genesis.json
```
6. Initialize decentr 1.5.0
```shell
   decentrd init <moniker>
```
7. Replace priv_validator_key.json with your old one
```shell
   cp -rf $HOME/.decentrd/config/priv_validator_key.json $HOME/.decentr/config/priv_validator_key.json
```
8. Replace genesis.json with migrated one
```shell
   cp -rf /tmp/genesis.json $HOME/.decentr/config/genesis.json
```
9. Patch seeds
```shell
    sed -E -i 's/seeds = \".*\"/seeds = \"f9b77dd93f28d2a45b00d4e3041b89a3c08788ef@calliope.mainnet.decentr.xyz:26656,987b5ce87b1b922793069756f594533eedf0f060@euterpe.mainnet.decentr.xyz:26656,2caebc4dad8d2ff95400918572d455392e10a63c@hera.mainnet.decentr.xyz:26656,c37f32e202e13b0725515570f794b68573a6f58c@hermes.mainnet.decentr.xyz:26656,4520b3221c91fa98a947a4c7f518ba5aab4e5b08@melpomene.mainnet.decentr.xyz:26656,c17bc88591115e52a686811630ad8c053de19f83@poseidon.mainnet.decentr.xyz:26656,c4ba719d38c871a93fb06cbfe0891ab11fedb9f7@terpsichore.mainnet.decentr.xyz:26656,9e9e0243610fadc0f65d3d927e2d682d86f71ea9@thalia.mainnet.decentr.xyz:26656,e1f3ce208776ff1fad0e8190f5475b68e841d788@zeus.mainnet.decentr.xyz:26656\"/' $HOME/.decentr/config/config.toml
```
10. Reset new version node state
```shell
    decentrd unsafe-reset-all
```
11. Start your node with
```shell
    decentrd start
```
12. Import your key back
```shell
    decentrd keys import <name> <name>.key
```
13. Verify you truly imported key
```shell
    decentrd keys show <name>
```
14. Remove exported key for security purposes
```shell
    rm <name>.key
```

