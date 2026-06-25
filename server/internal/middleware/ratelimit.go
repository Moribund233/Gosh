package middleware

import (
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	"gosh/internal/database"
	"gosh/pkg/response"
)

const (
	rateLimitPrefix = "ratelimit:"
)

func RateLimit(limit int, window time.Duration) gin.HandlerFunc {
	return func(c *gin.Context) {
		if database.RedisClient == nil {
			c.Next()
			return
		}

		key := rateLimitKey(c)
		allowed, remaining, reset, err := allow(database.RedisClient, c.Request.Context(), key, limit, window)
		if err != nil {
			c.Next()
			return
		}

		c.Header("X-RateLimit-Limit", strconv.Itoa(limit))
		c.Header("X-RateLimit-Remaining", strconv.Itoa(remaining))
		c.Header("X-RateLimit-Reset", strconv.FormatInt(reset, 10))

		if !allowed {
			response.Error(c, http.StatusTooManyRequests, "rate limit exceeded")
			c.Abort()
			return
		}
		c.Next()
	}
}

func rateLimitKey(c *gin.Context) string {
	ip := c.ClientIP()
	path := c.FullPath()
	userID, _ := c.Get("user_id")
	if uid, ok := userID.(uint); ok && uid > 0 {
		return rateLimitPrefix + "user:" + strconv.FormatUint(uint64(uid), 10) + ":" + path
	}
	return rateLimitPrefix + "ip:" + ip + ":" + path
}

func allow(client *redis.Client, ctx interface{}, key string, limit int, window time.Duration) (bool, int, int64, error) {
	c, ok := ctx.(interface {
		Deadline() (time.Time, bool)
		Done() <-chan struct{}
		Err() error
		Value(interface{}) interface{}
	})
	if !ok {
		return true, 0, 0, nil
	}

	now := time.Now().Unix()
	windowSeconds := int64(window.Seconds())
	cleanupAt := now - windowSeconds

	pipe := client.Pipeline()
	pipe.ZRemRangeByScore(c, key, "0", strconv.FormatInt(cleanupAt, 10))
	count := pipe.ZCard(c, key)
	pipe.Expire(c, key, window)
	_, err := pipe.Exec(c)
	if err != nil {
		return true, 0, 0, err
	}

	current := count.Val()
	if current >= int64(limit) {
		return false, 0, now + windowSeconds, nil
	}

	pipe = client.Pipeline()
	pipe.ZAdd(c, key, redis.Z{Score: float64(now), Member: now})
	pipe.Expire(c, key, window)
	pipe.ZCard(c, key)
	_, err = pipe.Exec(c)
	if err != nil {
		return true, 0, 0, err
	}

	return true, limit - int(current) - 1, now + windowSeconds, nil
}

func RateLimitByRole(public, auth, admin int, window time.Duration) gin.HandlerFunc {
	return func(c *gin.Context) {
		if database.RedisClient == nil {
			c.Next()
			return
		}

		role, _ := c.Get("role")
		roleStr, _ := role.(string)

		limit := public
		if roleStr != "" {
			limit = auth
		}
		for _, r := range []string{"super_admin", "operator"} {
			if roleStr == r {
				limit = admin
				break
			}
		}

		key := rateLimitKey(c)
		allowed, remaining, reset, err := allow(database.RedisClient, c.Request.Context(), key, limit, window)
		if err != nil {
			c.Next()
			return
		}

		c.Header("X-RateLimit-Limit", strconv.Itoa(limit))
		c.Header("X-RateLimit-Remaining", strconv.Itoa(remaining))
		c.Header("X-RateLimit-Reset", strconv.FormatInt(reset, 10))

		if !allowed {
			response.Error(c, http.StatusTooManyRequests, "rate limit exceeded")
			c.Abort()
			return
		}
		c.Next()
	}
}


