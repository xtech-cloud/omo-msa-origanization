package cache

import (
	"omo.msa.organization/config"
	"omo.msa.organization/proxy/nosql"
	"time"
)

type BaseInfo struct {
	ID         uint64 `json:"-"`
	UID        string `json:"uid"`
	Name       string `json:"name"`
	CreateTime time.Time
	UpdateTime time.Time
}

type cacheContext struct {
	scenes  []*SceneInfo
}

var cacheCtx *cacheContext

func InitData() error {
	cacheCtx = &cacheContext{}
	cacheCtx.scenes = make([]*SceneInfo, 0, 200)

	err := nosql.InitDB(config.Schema.Database.IP, config.Schema.Database.Port, config.Schema.Database.Name, config.Schema.Database.Type)
	if nil != err {
		return err
	}

	return nil
}

