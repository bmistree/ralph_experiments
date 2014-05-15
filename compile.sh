#!/usr/bin/env bash 

root_proj=`pwd`
prev_gopath=$GOPATH
export GOPATH=$root_proj
cd src/common
go install
cd $root_proj
cd src/read_only
go install
cd $root_proj
cd src/main
go install
cd $root_proj

export GOPATH=$prev_gopath