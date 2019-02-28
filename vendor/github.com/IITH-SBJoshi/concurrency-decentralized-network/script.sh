#!/bin/bash


for f in $PWD/src/*
do
	cd $f
	godoc -html -goroot=$HOME/concurrency-decentralized-network cmd/${PWD##*/} > $HOME/concurrency-decentralized-network/docs/${PWD##*/}/index.html

done