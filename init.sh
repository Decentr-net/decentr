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

decentrd add-genesis-account "$(decentrcli keys show jack -a)" 1000000000udec
decentrd add-genesis-account "$(decentrcli keys show alice -a)" 1000000000udec
decentrd add-genesis-community-moderators "$(decentrcli keys show jack -a)"
decentrd add-genesis-pdv-supervisors "$(decentrcli keys show jack -a)"

decentrd gentx --name jack --keyring-backend test --amount 1000000udec

echo "Collecting genesis txs..."
decentrd collect-gentxs

# Replace all "stake" denom with "udec"
sed -i -e 's/"stake"/"udec"/g' ~/.decentrd/config/genesis.json

echo "Validating genesis file..."
decentrd validate-genesis

decentrd start
