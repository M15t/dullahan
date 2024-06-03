package session

import (
	"dullahan/internal/db"

	"github.com/M15t/ghoul/pkg/rbac"
	"github.com/allegro/bigcache/v3"
)

// New creates new session application service
func New(db *db.Service, rbacSvc rbac.Intf, cr Crypter, cache *bigcache.BigCache) *Session {
	return &Session{db: db, rbac: rbacSvc, cr: cr, cache: cache}
}

// Session represents latefee application service
type Session struct {
	db    *db.Service
	rbac  rbac.Intf
	cr    Crypter
	cache *bigcache.BigCache
}

// Crypter represents security interface
type Crypter interface {
	RoundFloat(f float64) float64
	Float64ToByte(f float64) []byte
}
