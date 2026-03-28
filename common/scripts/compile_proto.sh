#!/bin/bash
set -euo pipefail
protoc --go_out=paths=source_relative:. --go-grpc_out=paths=source_relative:. proto/demo.proto
