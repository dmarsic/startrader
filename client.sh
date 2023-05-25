#!/usr/bin/env bash

GREEN="\e[32m"
RESET="\e[0m"

CURL="curl -L -b cookie.txt"
HOST="http://localhost:9500/api/v1"

green() {
    echo -e "\n\n${GREEN}$1${RESET}\n"
}

create_user() {
    local userfile="data/users/testuser.json"
    cat <<EOF > "$userfile"
{
	"name": "testuser",
	"credits": 1000,
	"location": "sol",
	"inventory": {
		"fuel": {
			"quantity": 250.0
		},
		"grains": {
			"quantity": 10
		}
	}
}
EOF
}

create_user

green "Login as testuser"
curl -i -c cookie.txt $HOST/login?user=testuser

green "Show information about self"
$CURL $HOST/

green "Show information about multiple users"
$CURL $HOST/u/Kofi,user1,testuser

green "Show information about the source system"
$CURL $HOST/systems/sol

green "Move from sol to vega"
$CURL -X POST -d 'system=vega' $HOST/m

green "Show information about both source and destination systems"
$CURL $HOST/systems/sol
$CURL $HOST/systems/vega

green "Buy 20 grains"
$CURL -X POST -d 'item=grains&quantity=20' $HOST/b
