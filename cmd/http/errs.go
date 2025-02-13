package main

import "errors"

var ErrHttpRequestDecode = errors.New("cannot decode request body")
