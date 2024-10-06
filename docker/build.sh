#!/bin/bash

cd ..
rm ./oe3xorstatus || true
go build .

cd frontend || exit
npm install
npm run build

cd ../docker || exit
rm -rf tmp/ || true
mkdir tmp

cp ../oe3xorstatus ./tmp
cp -r ../pb_public ./tmp/

docker build . -t oe3xorstatus
