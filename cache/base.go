package cache

import (
	"github.com/micro/go-micro/v2/logger"
	pb "github.com/xtech-cloud/omo-msp-organization/proto/organization"
	"omo.msa.organization/config"
	"omo.msa.organization/proxy/nosql"
	"reflect"
	"time"
)

type baseInfo struct {
	ID         uint64 `json:"-"`
	UID        string `json:"uid"`
	Name       string `json:"name"`
	Creator    string
	Operator   string
	CreateTime time.Time
	UpdateTime time.Time
}

type cacheContext struct {
	scenes []*SceneInfo
}

var cacheCtx *cacheContext

func InitData() error {
	cacheCtx = &cacheContext{}
	cacheCtx.scenes = make([]*SceneInfo, 0, 200)

	err := nosql.InitDB(config.Schema.Database.IP, config.Schema.Database.Port, config.Schema.Database.Name, config.Schema.Database.Type)
	if nil != err {
		return err
	}

	scenes, err := nosql.GetAllScenes()
	if err == nil {
		for _, scene := range scenes {
			tmp := new(SceneInfo)
			tmp.initInfo(scene)
			cacheCtx.scenes = append(cacheCtx.scenes, tmp)
		}
	}
	logger.Infof("init scenes that number = %d", len(cacheCtx.scenes))

	return nil
}

func Context() *cacheContext {
	return cacheCtx
}

func checkPage(page, number uint32, all interface{}) (uint32, uint32, interface{}) {
	if number < 1 {
		number = 10
	}
	array := reflect.ValueOf(all)
	total := uint32(array.Len())
	maxPage := total / number
	if total%number != 0 {
		maxPage = total/number + 1
	}
	if page < 1 {
		return total, maxPage, all
	}
	if page > maxPage {
		page = maxPage
	}

	var start = (page - 1) * number
	var end = start + number
	if end > total {
		end = total
	}

	list := array.Slice(int(start), int(end))
	return total, maxPage, list.Interface()
}

func SwitchAreaToProduct(info *AreaInfo) *pb.ProductInfo {
	tmp := new(pb.ProductInfo)
	tmp.Sn = info.Device
	tmp.Room = info.Parent
	tmp.Template = info.Template
	tmp.Name = info.Name
	tmp.Area = info.UID
	tmp.Type = info.Type
	tmp.Remark = info.Remark
	tmp.Catalog = info.Catalog
	tmp.Question = info.Question
	tmp.Displays = Context().switchDisplays(info.Type, info.Displays)
	return tmp
}
