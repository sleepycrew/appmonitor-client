#!/usr/bin/env bash
# removes unwanted code that is generated, it messes json serialization, the library doesn't really support disabling that
NO_MARSHALS=$(cat pkg/data/spec.go | head -n $(cat pkg/data/spec.go | grep -m1 -n -B1 MarshalJSON | cut -d ":" -f 1 | head -n 1 |  cut -d "-" -f 1))
NO_IMPORTS=$(echo "$NO_MARSHALS" | sed -e '5,10d')
echo "$NO_IMPORTS" > pkg/data/spec.go