#!/bin/bash

set -e

if ! command -v buf &> /dev/null; then
  echo "Buf is not installed. Installing..."
  curl -sSL https://github.com/bufbuild/buf/releases/download/v1.16.0/buf-Linux-x86_64 -o /usr/local/bin/buf
  chmod +x /usr/local/bin/buf
fi

buf generate
