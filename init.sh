decentrd init test --chain-id=local -o

decentrd config output json
decentrd config chain-id local
decentrd config keyring-backend test

decentrd keys add jack
decentrd keys add alice

jack=$(decentrd keys show jack -a)
alice=$(decentrd keys show alice -a)

decentrd add-genesis-account "$jack" 1000000000udec
decentrd add-genesis-account "$alice" 1000000000udec
decentrd add-genesis-moderator "$jack"
decentrd add-genesis-supervisor "$jack"

decentrd --keyring-backend=test --chain-id=local gentx jack 1000000udec

echo "Collecting genesis txs..."
decentrd collect-gentxs

# Replace all "stake" denom with "udec"
sed -i -e 's/"stake"/"udec"/g' ~/.decentr/config/genesis.json

echo "Validating genesis file..."
decentrd validate-genesis

decentrd start
