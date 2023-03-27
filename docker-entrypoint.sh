#!/bin/bash
cd /app

make migrate-apply
./build-out
