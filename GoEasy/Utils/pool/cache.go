package pool

import (
	"github.com/patrickmn/go-cache"
	"time"
)

var Gocache *cache.Cache

func GocacheInit() {
	Gocache = cache.New(300*time.Second, 1*time.Second)
}
