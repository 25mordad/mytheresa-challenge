package controller

import (
	cache "github.com/patrickmn/go-cache"
)

// BaseHandler will hold everything that controller needs
type BaseHandler struct {
	cache *cache.Cache
}

// NewBaseHandler returns a new BaseHandler
func NewBaseHandler(c *cache.Cache) *BaseHandler {
	return &BaseHandler{
		cache: c,
	}
}
