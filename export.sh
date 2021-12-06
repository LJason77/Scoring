#!/bin/bash

mkdir -p release

cargo build --release

cp -r target/release/scoring chaincode deploy go static .env release/
