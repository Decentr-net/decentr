decentr init test --chain-id=local

decentr config output json
decentr config trust-node true
decentr config chain-id local
decentr config keyring-backend test

decentr keys add jack
decentr keys add alice

jack=$(decentr keys show jack -a)
alice=$(decentr keys show alice -a)

decentr add-genesis-account "$jack" 1000000000udec
decentr add-genesis-account "$alice" 1000000000udec
decentr add-genesis-moderator "$jack"
decentr add-genesis-supervisor "$jack"

decentr --keyring-backend=test --chain-id=local gentx jack 1000000udec

echo "Collecting genesis txs..."
decentr collect-gentxs

# Replace all "stake" denom with "udec"
sed -i -e 's/"stake"/"udec"/g' ~/.decentr/config/genesis.json

echo "Validating genesis file..."
decentr validate-genesis

decentr start
