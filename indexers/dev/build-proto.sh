#!/bin/bash

set -e

protoc --go_out=. exporter_protocol/*.proto
