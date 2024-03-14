package task

import (
	"context"
	"gotimer/common/conf"
	"gotimer/common/model/po"
	"gotimer/pkg/redis"
	"time"
)

type TaskCache struct {
	client       cacheClient
	confProvider *conf.SchedulerAppConfProvider
}

type cacheClient interface {
	Transaction(ctx context.Context, commands ...*redis.Command) ([]interface{}, error)
	ZrangeByScore(ctx context.Context, table string, score1, score2 int64) ([]string, error)
	Expire(ctx context.Context, key string, expireSeconds int64) error
	MGet(ctx context.Context, keys ...interface{}) ([]string, error)
}

func NewTaskCache(client cacheClient, confProvider *conf.SchedulerAppConfProvider) *TaskCache {
	return &TaskCache{
		client:       client,
		confProvider: confProvider,
	}
}

func (t *TaskCache) BatchCreateBucket(ctx context.Context, cntByMins []*po.MinuteTaskCnt, end time.Time) error {
	conf := t.confProvider.Get()

	expireSeconds := int64(time.Until(end) / time.Second)
	commands := make([]*redis.Command, 0, 2*len(cntByMins))
	for _, detail := range cntByMins {
		commands = append(commands, redis.NewSetCommand(utils.GetBu))
	}
}
