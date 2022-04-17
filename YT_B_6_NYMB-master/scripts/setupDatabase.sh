#!/bin/bash

SCRIPTS=$GOPATH/src/git.linux.iastate.edu/309Fall2017/YT_B_6_NYMB/scripts

mysql=/usr/local/mysql/bin/mysql
args='-udbu309ytb6 -prA@12XF3 db309ytb6'
if [ -f /bin/mysql ]; then
    echo "setting up remote database"
    mysql=/bin/mysql
    args='-udbu309ytb6 -hmysql.cs.iastate.edu -prA@12XF3 db309ytb6'
fi

echo -ne "initializing database\t"
$mysql $args < $SCRIPTS/init.sql

echo -ne "initializing fixtures\t"
$mysql $args < $SCRIPTS/fixtures.sql
