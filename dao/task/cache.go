package task

import "context"

type TaskCache struct {
	client       cacheClient
	confProvider *conf.SchedulerAppConfProvider
}

type cacheClient interface {
	Transaction(ctx context.Context, commands ...*redis.Command) ([]interface{}, err)
}
