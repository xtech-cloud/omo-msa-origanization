package cache

import (
	"errors"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"omo.msa.organization/proxy"
	"omo.msa.organization/proxy/nosql"
	"time"
)

type AreaInfo struct {
	baseInfo
	Remark   string
	Owner    string
	Parent   string
	Template string //产品配置模板
	Type     uint32 //产品类型

	Width    int32
	Height   int32  //
	LimitNum uint32 //展览限制的数量
	//UrgentPage string            //紧急播放页面
	//PlaySheet  string            //
	Device   string //设备UID
	Question string //终端使用的答题类型
	Catalog  string //终端定制目录base64加密

	deviceInfo *DeviceInfo

	Modules  []*proxy.PairInfo //模块配置
	Sources  []*proxy.PairInfo //定制资源配置
	Displays []string
}

func (mine *cacheContext) CreateArea(name, remark, owner, parent, operator string) (*AreaInfo, error) {
	db := new(nosql.Area)
	db.UID = primitive.NewObjectID()
	db.ID = nosql.GetAreaNextID()
	db.CreatedTime = time.Now()
	db.Creator = operator
	db.Name = name
	db.Remark = remark
	db.Scene = owner
	db.Parent = parent
	db.Type = 0
	db.Width = 0
	db.Height = 0
	db.Limit = 0
	db.Template = ""
	db.Device = ""
	db.Question = ""
	db.Catalog = ""
	db.Displays = make([]string, 0, 1)
	db.Modules = make([]*proxy.PairInfo, 0, 1)
	db.Sources = make([]*proxy.PairInfo, 0, 1)

	err := nosql.CreateArea(db)
	if err != nil {
		return nil, err
	}
	info := new(AreaInfo)
	info.initInfo(db)
	return info, nil
}

func (mine *cacheContext) GetArea(uid string) (*AreaInfo, error) {
	if len(uid) < 2 {
		return nil, errors.New("the area uid is empty")
	}
	db, err := nosql.GetArea(uid)
	if err != nil {
		return nil, err
	}
	info := new(AreaInfo)
	info.initInfo(db)
	return info, nil
}

func (mine *cacheContext) GetAreaByDevice(uid string) (*AreaInfo, error) {
	if len(uid) < 2 {
		return nil, errors.New("the area device uid is empty")
	}
	db, err := nosql.GetAreaByDevice(uid)
	if err != nil {
		return nil, err
	}
	info := new(AreaInfo)
	info.initInfo(db)
	return info, nil
}

func (mine *cacheContext) GetAreaBySN(sn string) (*AreaInfo, error) {
	if len(sn) < 2 {
		return nil, errors.New("the area sn is empty")
	}
	db, err := nosql.GetDeviceBySN(sn)
	if err != nil {
		return nil, err
	}
	return mine.GetAreaByDevice(db.UID.Hex())
}

func (mine *cacheContext) GetAreas(parent string) []*AreaInfo {
	list := make([]*AreaInfo, 0, 20)
	array, err := nosql.GetAreasByParent(parent)
	if err != nil {
		return list
	}
	for _, item := range array {
		info := new(AreaInfo)
		info.initInfo(item)
		list = append(list, info)
	}
	return list
}

func (mine *cacheContext) GetAreasByOwner(uid string) []*AreaInfo {
	array, err := nosql.GetAreasByOwner(uid)
	if err != nil {
		return make([]*AreaInfo, 0, 0)
	}
	list := make([]*AreaInfo, 0, len(array))
	for _, item := range array {
		info := new(AreaInfo)
		info.initInfo(item)
		list = append(list, info)
	}
	return list
}

func (mine *cacheContext) GetAreasByTemplate(owner, template string) []*AreaInfo {
	array, err := nosql.GetAreasByTemplate(owner, template)
	if err != nil {
		return make([]*AreaInfo, 0, 0)
	}
	list := make([]*AreaInfo, 0, len(array))
	for _, item := range array {
		info := new(AreaInfo)
		info.initInfo(item)
		list = append(list, info)
	}
	return list
}

func (mine *cacheContext) GetAreaList(array []string) []*AreaInfo {
	if array == nil || len(array) < 1 {
		return make([]*AreaInfo, 0, 0)
	}
	list := make([]*AreaInfo, 0, len(array))
	for i := 0; i < len(array); i += 1 {
		db, err := nosql.GetArea(array[i])
		if err == nil {
			info := new(AreaInfo)
			info.initInfo(db)
			list = append(list, info)
		}
	}
	return list
}

func (mine *AreaInfo) initInfo(db *nosql.Area) {
	mine.Name = db.Name
	mine.UID = db.UID.Hex()
	mine.ID = db.ID
	mine.Remark = db.Remark
	mine.CreateTime = db.CreatedTime
	mine.UpdateTime = db.UpdatedTime
	mine.Creator = db.Creator
	mine.Operator = db.Operator
	mine.Parent = db.Parent
	mine.Owner = db.Scene
	mine.LimitNum = db.Limit
	mine.Template = db.Template
	mine.Width = db.Width
	mine.Height = db.Height
	mine.Type = db.Type
	mine.Device = db.Device
	mine.Catalog = db.Catalog
	mine.Question = db.Question
	mine.Displays = db.Displays
	mine.Modules = db.Modules
	mine.Sources = db.Sources
}

func (mine *AreaInfo) DeviceInfo() (*DeviceInfo, error) {
	return cacheCtx.GetDevice(mine.Device)
}

func (mine *AreaInfo) DeviceSN() string {
	if mine.deviceInfo == nil {
		mine.deviceInfo, _ = cacheCtx.GetDevice(mine.Device)
	}

	if mine.deviceInfo == nil {
		return ""
	}
	return mine.deviceInfo.SN
}

func (mine *AreaInfo) GetAspect() string {
	if mine.deviceInfo == nil {
		mine.deviceInfo, _ = cacheCtx.GetDevice(mine.Device)
	}

	if mine.deviceInfo == nil {
		return ""
	}
	return mine.deviceInfo.Aspect
}

func (mine *AreaInfo) UpdateBase(name, remark, operator string) error {
	err := nosql.UpdateAreaBase(mine.UID, name, remark, operator)
	if err == nil {
		mine.Name = name
		mine.Remark = remark
		mine.Operator = operator
	}
	return err
}

func (mine *AreaInfo) UpdateTemplate(template, operator string) error {
	err := nosql.UpdateAreaTemplate(mine.UID, template, operator)
	if err == nil {
		mine.Template = template
		mine.Operator = operator
	}
	return err
}

func (mine *AreaInfo) UpdateLimitCount(operator string, num uint32) error {
	if mine.LimitNum == num {
		return nil
	}
	err := nosql.UpdateAreaLimit(mine.UID, operator, num)
	if err == nil {
		mine.LimitNum = num
		mine.Operator = operator
	}
	return err
}

func (mine *AreaInfo) UpdateDevice(device, operator string, tp uint32) error {
	err := nosql.UpdateAreaDevice(mine.UID, device, operator, tp)
	if err == nil {
		mine.Device = device
		mine.Type = tp
		mine.Operator = operator
	}
	return err
}

func (mine *AreaInfo) UpdateDisplays(operator string, list []string) error {
	err := nosql.UpdateAreaDisplays(mine.UID, operator, list)
	if err == nil {
		mine.Displays = list
		mine.Operator = operator
	}
	return err
}

func (mine *AreaInfo) UpdateDevice2(sn, operator string) error {
	err := nosql.UpdateAreaDevice2(mine.UID, sn, operator)
	if err == nil {
		mine.Device = sn
		mine.Operator = operator
	}
	return err
}

func (mine *AreaInfo) UpdateType(tp uint32, operator string) error {
	err := nosql.UpdateAreaType(mine.UID, operator, tp)
	if err == nil {
		mine.Type = tp
		mine.Operator = operator
	}
	return err
}

func (mine *AreaInfo) UpdateCatalog(catalog, operator string) error {
	err := nosql.UpdateAreaCatalog(mine.UID, catalog, operator)
	if err == nil {
		mine.Catalog = catalog
		mine.Operator = operator
	}
	return err
}

func (mine *AreaInfo) UpdateQuestion(question, operator string) error {
	err := nosql.UpdateAreaQuestion(mine.UID, question, operator)
	if err == nil {
		mine.Question = question
		mine.Operator = operator
	}
	return err
}

func (mine *AreaInfo) Remove(operator string) error {
	return nosql.RemoveArea(mine.UID, operator)
}

func (mine *AreaInfo) UpdateModule(key, value, operator string) error {
	arr := make([]*proxy.PairInfo, 0, len(mine.Modules))
	arr = append(arr, mine.Modules...)
	ok := false
	for _, info := range arr {
		if info.Key == key {
			info.Value = value
			ok = true
			break
		}
	}
	if !ok {
		arr = append(arr, &proxy.PairInfo{Key: key, Value: value})
	}
	return nosql.UpdateAreaModules(mine.UID, operator, arr)
}

func (mine *AreaInfo) UpdateCustomSource(key, value, operator string) error {
	arr := make([]*proxy.PairInfo, 0, len(mine.Sources))
	arr = append(arr, mine.Sources...)
	ok := false
	for _, info := range arr {
		if info.Key == key {
			info.Value = value
			ok = true
			break
		}
	}
	if !ok {
		arr = append(arr, &proxy.PairInfo{Key: key, Value: value})
	}
	return nosql.UpdateAreaSources(mine.UID, operator, arr)
}
