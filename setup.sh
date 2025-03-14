#!/bin/bash

# Ensure the home directory and hypd binary are provided as arguments
if [ -z "$1" ] || [ -z "$2" ]; then
    printf "\033[1;31mUsage:\033[0m %s <home-directory> <hypd-path>\n" "$0"
    exit 1
fi

HOME_DIR=$1
HYPD_PATH=$2

# Function to execute a command with enhanced logging
execute_command() {
    printf "\033[1;34m[Executing]\033[0m: %s\n" "$1"
    eval "$1" > /dev/null 2>&1
    CURL_COMMAND=$(echo "$1" | grep -o 'curl.*')
    if [ ! -z "$CURL_COMMAND" ]; then
        printf "\033[1;33m[Query response]\033[0m for: %s\n" "$CURL_COMMAND"
        eval "$CURL_COMMAND | jq ."
    fi
    sleep 6
}

# Execute commands
execute_command "$HYPD_PATH tx hyperlane ism create-merkle-root-multisig 0x0c60e7eCd06429052223C78452F791AAb5C5CAc6 1 -y --chain-id hyperlane-local --keyring-backend test --from alice --fees 40000uhyp --home $HOME_DIR"
execute_command "curl localhost:1317/hyperlane/v1/isms"
execute_command "$HYPD_PATH tx hyperlane hooks igp create uhyp --chain-id hyperlane-local --keyring-backend test -y --from alice --fees 40000uhyp --home $HOME_DIR"
execute_command "curl localhost:1317/hyperlane/v1/igps"
execute_command "$HYPD_PATH tx hyperlane mailbox create 0x726f757465725f69736d00000000000000000000000000040000000000000000 75898669 -y --chain-id hyperlane-local --keyring-backend test --from alice --fees 40000uhyp --home $HOME_DIR"
execute_command "curl localhost:1317/hyperlane/v1/mailboxes"
execute_command "$HYPD_PATH tx hyperlane hooks merkle create 0x68797065726c616e650000000000000000000000000000000000000000000000 --chain-id hyperlane-local -y --keyring-backend test --from alice --fees 40000uhyp --home $HOME_DIR"
execute_command "curl localhost:1317/hyperlane/v1/merkle_tree_hooks"
execute_command "$HYPD_PATH tx hyperlane mailbox set 0x68797065726c616e650000000000000000000000000000000000000000000000 --default-hook 0x726f757465725f706f73745f6469737061746368000000040000000000000000 -y --required-hook 0x726f757465725f706f73745f6469737061746368000000030000000000000001 --chain-id hyperlane-local --keyring-backend test --from alice --fees 40000uhyp --home $HOME_DIR"
execute_command "curl localhost:1317/hyperlane/v1/mailboxes"
execute_command "$HYPD_PATH tx hyperlane-transfer create-collateral-token 0x68797065726c616e650000000000000000000000000000000000000000000000 uhyp --chain-id hyperlane-local --keyring-backend test --from alice -y --fees 40000uhyp --home $HOME_DIR"
execute_command "curl localhost:1317/hyperlane/v1/tokens"
execute_command "$HYPD_PATH tx hyperlane-transfer enroll-remote-router 0x726f757465725f61707000000000000000000000000000010000000000000000 75898669 0x726f757465725f61707000000000000000000000000000010000000000000000 50000 -y --chain-id hyperlane-local --keyring-backend test --from alice --fees 40000uhyp --home $HOME_DIR"
execute_command "curl localhost:1317/hyperlane/v1/tokens/0x726f757465725f61707000000000000000000000000000010000000000000000/remote_routers"
execute_command "$HYPD_PATH tx hyperlane hooks igp set-destination-gas-config 0x726f757465725f706f73745f6469737061746368000000040000000000000000 75898669 10000000000 1 200000 --chain-id hyperlane-local --keyring-backend test -y --from alice --fees 40000uhyp --home $HOME_DIR"
execute_command "$HYPD_PATH tx hyperlane-transfer transfer 0x726f757465725f61707000000000000000000000000000010000000000000000 75898669 0x726f757465725f61707000000000000000000000000000010000000000000000 1000000 --max-hyperlane-fee 250000uhyp -y --chain-id hyperlane-local --keyring-backend test --from alice --fees 40000uhyp --home $HOME_DIR"

printf "\033[1;32m[All Done]\033[0m All commands executed successfully!\n"
