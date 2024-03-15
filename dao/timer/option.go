package timer

import (
	"context"
	"gorm.io/gorm"
	"gotimer/common/model/po"
	"gotimer/pkg/log"
	"gotimer/pkg/mysql"
)

type TimerDAO struct {
	client *mysql.Client
}

func NewTimerDAO(client *mysql.Client) *TimerDAO {
	return &TimerDAO{
		client: client,
	}
}

func (t *TimerDAO) CreateTimer(ctx context.Context, timer *po.Timer) (uint, error) {
	return timer.ID, t.client.DB.WithContext(ctx).Create(timer).Error
}

func (t *TimerDAO) UpdateTimer(ctx context.Context, timer *po.Timer) error {
	return t.client.DB.WithContext(ctx).Updates(timer).Error
}

func (t *TimerDAO) DeleteTimer(ctx context.Context, id uint) error {
	return t.client.DB.WithContext(ctx).Delete(&po.Timer{Model: gorm.Model{ID: id}}).Error
}

func (t *TimerDAO) GetTimer(ctx context.Context, opts ...Option) (*po.Timer, error) {
	db := t.client.DB.WithContext(ctx)
	for _, opt := range opts {
		db = opt(db)
	}
	var timer po.Timer
	return &timer, db.First(&timer).Error
}

func (t *TimerDAO) GetTimers(ctx context.Context, opts ...Option) ([]*po.Timer, error) {
	db := t.client.DB.WithContext(ctx).Model(&po.Timer{})
	for _, opt := range opts {
		db = opt(db)
	}
	var timers []*po.Timer
	return timers, db.Scan(&timers).Error
}

func (t *TimerDAO) Count(ctx context.Context, opts ...Option) (int64, error) {
	db := t.client.DB.WithContext(ctx).Model(&po.Timer{})
	for _, opt := range opts {
		db = opt(db)
	}
	var cnt int64
	return cnt, db.Debug().Count(&cnt).Error
}
func (t *TimerDAO) Transaction(ctx context.Context, do func(ctx context.Context, dao *TimerDAO) error) error {
	return t.client.Transaction(func(tx *gorm.DB) error {
		//defer 在这里使用是为了确保即使事务中发生 panic，
		//也会执行 tx.Rollback() 和 log.ErrorContextf() 语句。
		//这有助于保持事务的完整性并记录任何发生的错误。
		defer func() {
			if err := recover(); err != nil {
				tx.Rollback()
				log.ErrorContextf(ctx, "transaction err: %v", err)
			}
		}()
		return do(ctx, NewTimerDAO(mysql.NewClient(tx)))
	})
}

// 将task切片从po -> mysql
func (t *TimerDAO) BatchCreateRecords(ctx context.Context, tasks []*po.Task) error {
	return t.client.DB.Model(&po.Task{}).WithContext(ctx).CreateInBatches(tasks, len(tasks)).Error
}

// 查询时用FOR UPDATE上锁
func (t *TimerDAO) DoWithLock(ctx context.Context, id uint, do func(ctx context.Context, dao *TimerDAO, timer *po.Timer) error) error {
	return t.client.Transaction(func(tx *gorm.DB) error {
		defer func() {
			if err := recover(); err != nil {
				tx.Rollback()
				log.ErrorContextf(ctx, "transaction with lock err: %v,timer id : %d", err, id)
			}
		}()
		var timer po.Timer
		if err := tx.Set("gorm:query_option", "FOR UPDATE").WithContext(ctx).First(&timer, id).Error; err != nil {
			return err
		}
		return do(ctx, NewTimerDAO(mysql.NewClient(tx)), &timer)
	})
}
