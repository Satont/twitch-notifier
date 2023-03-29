#!/bin/bash
cd /app

make migrate-apply
exec "$@"
