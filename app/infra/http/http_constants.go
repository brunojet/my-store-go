package http

import httptypes "github.com/brunojet/my-store-go/app/infra/http/types"

// Deprecated: driver enum lives in infra/http/types.
// This file exists to keep existing imports stable.

type HTTPDriver = httptypes.HTTPDriver

const (
	HTTPDriverGin HTTPDriver = httptypes.HTTPDriverGin
	HTTPDriverChi HTTPDriver = httptypes.HTTPDriverChi
)

func NormalizeHTTPDriver(raw string) HTTPDriver {
	return httptypes.NormalizeHTTPDriver(raw)
}
