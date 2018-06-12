#!/bin/bash


CURR_DIR=$(cd `dirname $0`; pwd)
cd ${CURR_DIR}

export GOPATH=$PWD:${GOPATH}

find src/DCrypto -type f -name "*.go" | while read line ;
do
    if [ "${IS_DEBUG}" = "1" ]; then
        echo "format ${line}"
        go fmt ${line}
    else
        echo ${line}
    fi
done

go install DCrypto/center
go install DCrypto/worker


[ ! -d output ] && mkdir output 
rm -rf output/*

if [ "${IS_DEBUG}" = "1" ]; then
    cp -r bin conf_test output/
else
    cp -r bin conf output/
fi

