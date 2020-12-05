package limiter

import (
	"github.com/gin-gonic/gin"
	"github.com/juju/ratelimit"
	"time"
)

type LimitIfce interface {
	Key(c *gin.Context) string
	GetBucket(key string) (*ratelimit.Bucket, bool)
	AddBuckets(rules ...LimitBucketRule) LimitIfce
}

type Limiter struct {
	limiterBuckets map[string]*ratelimit.Bucket
}

type LimitBucketRule struct {
	Key string
	FillInterval time.Duration
	Capacity int64
	Quantum int64
}



