#!/usr/bin/env bash

if [ -z "$MKCERT_REPO" ]; then
  echo "Specify mkcert repo with 'export MKCERT_REPO=/path/to/repo/FiloSottile/mkcert'"
  exit 1
fi

mkdir -p ./mkcert

cp $MKCERT_REPO/*.go ./mkcert/

for gofile in $(find ./mkcert/ -name '*.go' -print); do
  sed -i'.bak' "s/^package main$/package mkcert/" "$gofile"
  rm "${gofile}.bak"
done