rm -rf ~/.decentrd
rm -rf ~/.decentrcli

decentrd init test --chain-id=testnet

decentrcli config output json
decentrcli config indent true
decentrcli config trust-node true
decentrcli config chain-id testnet
decentrcli config keyring-backend test

decentrcli keys add jack
decentrcli keys add alice

decentrd add-genesis-account $(decentrcli keys show jack -a) 100000000stake
decentrd add-genesis-account $(decentrcli keys show alice -a) 100000000stake

decentrd gentx --name jack --keyring-backend test

echo "Collecting genesis txs..."
decentrd collect-gentxs

echo "Validating genesis file..."
decentrd validate-genesis

decentrd start
