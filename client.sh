#!/usr/bin/env bash

GREEN="\e[32m"
RESET="\e[0m"

CURL="curl -L -b cookie.txt"
HOST="http://localhost:5000"

green() {
    echo -e "${GREEN}$1${RESET}"
}

green "Login as user1"
curl -i -c cookie.txt $HOST/login?user=user1

green "Show information about self"
$CURL $HOST/

green "Move from sol to vega"
$CURL -X POST -d 'system=vega' $HOST/m

green "Buy 20 grains"
$CURL -X POST -d 'item=grains&quantity=20' $HOST/b