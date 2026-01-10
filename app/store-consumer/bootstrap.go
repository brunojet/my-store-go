package storeconsumer

import (
	"github.com/brunojet/my-store-go/app/infra/http/contracts"
	"gorm.io/gorm"
)

// Register wires store-consumer routes and services into the given router.
//
// This module currently has no public HTTP endpoints wired.
func Register(_ contracts.Router, _ *gorm.DB) {}
