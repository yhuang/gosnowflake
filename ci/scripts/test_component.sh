#!/bin/bash
#
# Build and Test Golang driver
#
set -e
set -o pipefail
CI_SCRIPTS_DIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )"
TOPDIR=$( cd $CI_SCRIPTS_DIR/../.. && pwd )
# eval $( jq -r '.testconnection | to_entries | map("export \(.key)=\(.value|tostring)")|.[]' $TOPDIR/parameters.json )
SNOWFLAKE_CONFIG=$( jq -r '.testconnection | to_entries | map("\(.key)=\(.value|tostring)")|.[]' $TOPDIR/parameters.json )

for i in ${SNOWFLAKE_CONFIG[0]}; do
	export $i
done

if [[ -n "$GITHUB_WORKFLOW" ]]; then
	export SNOWFLAKE_TEST_PRIVATE_KEY=$TOPDIR/rsa-2048-private-key.p8
fi

env | grep SNOWFLAKE | grep -v PASS | sort

cd $TOPDIR
go test -timeout 30m -race $COVFLAGS -v .
