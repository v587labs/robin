package gorm

import (
	"context"
	"errors"
	"fmt"
	"gorm.io/gorm"
	"strings"
	"sync"
)

var clients = &sync.Map{}

func Start(cfgs []Config) error {
	for i, cfg := range cfgs {

		group := cfg.Group

		if len(group) == 0 {
			group = GroupDefault
		}
		_, ok := clients.Load(group)
		if ok {
			continue
		}

		db, err := gorm.Open(cfg.GetDialector(), cfg.BuildGormConfig())
		if err != nil {
			return err
		}
		// db setting
		sqlDB, err := db.DB()

		if cfg.MaxIdleCount > 0 {
			// SetMaxIdleConns sets the maximum number of connections in the idle connection pool.
			sqlDB.SetMaxIdleConns(cfg.MaxIdleCount)
		}

		if cfg.MaxOpenConns > 0 {
			// SetMaxOpenConns sets the maximum number of open connections to the database.
			sqlDB.SetMaxOpenConns(cfg.MaxOpenConns)
		}

		if cfg.MaxLifetime > 0 {
			// SetConnMaxLifetime sets the maximum amount of time a connection may be reused.
			sqlDB.SetConnMaxLifetime(cfg.GetMaxLifetime())
		}

		c := &Client{
			db:  db,
			cfg: cfg,
		}
		if i == 0 || group == GroupDefault {
			clients.Store(group, c)
		}
		clients.Store(cfg.Group, c)
	}
	return nil
}

func DB(ctx context.Context, group ...string) (*gorm.DB, error) {
	g := GroupDefault
	if len(group) > 0 && len(group[0]) > 0 {
		g = group[0]
	}
	c, ok := clients.Load(g)
	if !ok {
		if li := strings.LastIndexByte(g, ':'); li >= 0 {
			return DB(ctx, g[0:li])
		}
		if g != GroupDefault {
			return DB(ctx)
		}
		return nil, errors.New(fmt.Sprintf("not found group %s", g))
	}
	cc, ok := c.(*Client)
	if !ok {
		return nil, errors.New(fmt.Sprintf("group connect error %s", g))
	}
	return cc.db, nil
}
