## Decentr-1.5.7 Upgrade Instructions

decentr-2 Upgrade Instructions

### Software Version and Key Dates*

We will be upgrading from chain-id "mainnet-1" to chain-id "mainnet-2".
The version of decentr for mainnet-2 is v1.5.7
The mainnet-1 chain will be shutdown with a SoftwareUpgradeProposal that activates at block height 1688950, which is approximately 13:00 UTC on January, 24 2022.
mainnet-2 genesis time is set to January 24th, 2022 at 15:00 UTC
The version of cosmos-sdk for mainnet-2 is v0.44.3
The version of tendermint for mainnet-2 is v0.34.14
The recommended version of golang for mainnet-2 is 1.16.

### Risks

As a validator, performing the upgrade procedure on your consensus nodes carries a heightened risk of double-signing and being slashed. The most important piece of this procedure is verifying your software version and genesis file hash before starting your validator and signing.

The riskiest thing a validator can do is discover that they made a mistake and repeat the upgrade procedure again during the network startup. If you discover a mistake in the process, the best thing to do is wait for the network to start before correcting it. If the network is halted and you have started with a different genesis file than the expected one, seek advice from a Decentr developer before resetting your validator.

### Recovery

Prior to exporting mainnet-1 state, validators are encouraged to take a full data snapshot at the export height before proceeding. Snap-shotting depends heavily on infrastructure, but generally this can be done by backing up the .decentrd and .decentrcli directories.

It is critically important to back-up the .decentrd/data/priv_validator_state.json file after stopping your decentrd process. This file is updated every block as your validator participates in consensus rounds. It is a critical file needed to prevent double-signing, in case the upgrade fails and the previous chain needs to be restarted.

In the event that the upgrade does not succeed, validators and operators must restore their nodes from backup and upgrade to v1.4.8 of the decentr software.

### Upgrade Procedure

#### Before the upgrade

Decentr has submitted a SoftwareUpgradeProposal that specifies block height 1688950 as the final block height for mainnet-1. This height corresponds to approximately 13:00 UTC on January 24th. Once the proposal passes, the chain will shutdown automatically at the specified height and does not require manual intervention by validators.

#### On the day of the upgrade

The decentr chain is expected to halt at block height 1688950, at approximately 13:00 UTC, and restart with new software at 15:00 UTC January 24th. Do not stop your node and begin the upgrade before 13:00UTC on January 24th, or you may go offline and be unable to recover until after the upgrade!

Make sure the decentrd process is stopped before proceeding and that you have backed up your validator. Failure to backup your validator could make it impossible to restart your node if the upgrade fails.

#### Guide

1. Stop the service that's running the node
```shell
sudo systemctl stop decentr_node
```

2. Backup .decentrd and .decentrcli
```shell
cp -rf $HOME/.decentrd $HOME/.decentrd.bak
cp -rf $HOME/.decentrcli $HOME/.decentrcli.bak
```

3. Export your genesis file with the next command
```shell
decentrd export > $HOME/.decentrd/genesis.json
```

4. Export your key
```shell
decentrcli keys export <name>
```
save the output to <name>.key file. The easiest way is to execute 
```shell
nano <name>.key
```
paste and save.
  
The file shold start with `-----BEGIN TENDERMINT PRIVATE KEY-----` and ends with `-----END TENDERMINT PRIVATE KEY-----`
  
5. Clone Decentr from the latest release
```shell
git clone -b v1.5.7 https://github.com/Decentr-net/decentr
cd decentr
```
  
6. Compile and install new version of Decentr
```shell
make install
```
and check version
  
```shell
decentrd version
```
It has to be `1.5.7`

7. Migrate your genesis with next command
```shell
decentrd migrate $HOME/.decentrd/genesis.json > /tmp/genesis.json
```

8. Compare hashsum of migrated genesis
```shell
md5sum /tmp/genesis.json
```
It has to be `35bdd2a3fff849e3e0eba7d849764126`

9. Initialize mainnet-2
```shell
decentrd init <moniker>
```

10. Replace priv_validator_key.json with your old one
```shell
cp -rf $HOME/.decentrd/config/priv_validator_key.json $HOME/.decentr/config/priv_validator_key.json
```

11. Replace genesis.json with migrated one
```shell
cp -rf /tmp/genesis.json $HOME/.decentr/config/genesis.json
```

12. Patch seeds
```shell
sed -E -i 's/seeds = \".*\"/seeds = \"7708addcfb9d4ff394b18fbc6c016b4aaa90a10a@ares.mainnet.decentr.xyz:26656,8a3485f940c3b2b9f0dd979a16ea28de154f14dd@calliope.mainnet.decentr.xyz:26656,87490fd832f3226ac5d090f6a438d402670881d0@euterpe.mainnet.decentr.xyz:26656,3261bff0b7c16dcf6b5b8e62dd54faafbfd75415@hera.mainnet.decentr.xyz:26656,5f3cfa2e3d5ed2c2ef699c8593a3d93c902406a9@hermes.mainnet.decentr.xyz:26656,a529801b5390f56d5c280eaff4ae95b7163e385f@melpomene.mainnet.decentr.xyz:26656,385129dbe71bceff982204afa11ed7fa0ee39430@poseidon.mainnet.decentr.xyz:26656,35a934228c32ad8329ac917613a25474cc79bc08@terpsichore.mainnet.decentr.xyz:26656,0fd62bcd1de6f2e3cfc15852cdde9f3f8a7987e4@thalia.mainnet.decentr.xyz:26656,bd99693d0dbc855b0367f781fb48bf1ca6a6a58b@zeus.mainnet.decentr.xyz:26656\"/' $HOME/.decentr/config/config.toml
```

13. Reset the state
```shell
decentrd unsafe-reset-all
```

14. Import your key back
```shell
decentrd keys import <name> <name>.key
```

15. Verify you truly imported key
```shell
decentrd keys show <name>
```

16. Remove exported key for security purposes
```shell
rm <name>.key
```

17. Start your node back
```shell
sudo systemctl start decentr_node
```

18. Validate your node is up
```shell
sudo journalctl -u decentr_node.service -f
```

## Coordination

If the mainnet-2 chain does not launch by January 24, 2022 at 17:00 UTC, the launch should be considered a failure and validators should refer to the rollback instructions to restart the previous mainnet-1 chain. In the event of launch failure, coordination will occur in the Decentr discord.
