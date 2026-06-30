#!/bin/bash

find internal/config -type f \
 -not -name '*_test.go' \
 -exec echo "=== {} ===" \; \
 -exec cat {} \; > /tmp/dump.tmp && mv /tmp/dump.tmp ./dump.txt