#!/bin/bash

NYMB=$GOPATH/src/git.linux.iastate.edu/309Fall2017/YT_B_6_NYMB

SCRIPTS=$NYMB/scripts
TESTS=$SCRIPTS/tests

$SCRIPTS/setupDatabase.sh

$TESTS/account.sh
$TESTS/balance.sh
$TESTS/currency.sh
$TESTS/transaction.sh
$TESTS/user.sh
$TESTS/vault.sh

$SCRIPTS/setupDatabase.sh
