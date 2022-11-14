package cache

import (
	"errors"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"omo.msa.organization/proxy/nosql"
	"time"
)

type AreaInfo struct {
	baseInfo
	Remark   string
	Owner    string
	Parent   string
	Template string //产品配置模板
	Type     uint32
	Width    int32
	Height   int32
	SN       string
	Question string
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
	db.Template = ""
	db.SN = ""
	db.Question = ""
	db.Displays = make([]string, 0, 1)

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
		return nil, errors.New("the collective museum uid is empty")
	}
	db, err := nosql.GetArea(uid)
	if err != nil {
		return nil, err
	}
	info := new(AreaInfo)
	info.initInfo(db)
	return info, nil
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
	mine.Template = db.Template
	mine.Width = db.Width
	mine.Height = db.Height
	mine.Type = db.Type
	mine.SN = db.SN
	mine.Question = db.Question
	mine.Displays = db.Displays
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

func (mine *AreaInfo) UpdateDevice(sn, operator string, tp uint32) error {
	err := nosql.UpdateAreaDevice(mine.UID, sn, operator, tp)
	if err == nil {
		mine.SN = sn
		mine.Type = tp
		mine.Operator = operator
	}
	return err
}

func (mine *AreaInfo) UpdateShowings(operator string, list []string) error {
	err := nosql.UpdateAreaDisplays(mine.UID, operator, list)
	if err == nil {
		mine.Displays = list
		mine.Operator = operator
	}
	return err
}

func (mine *AreaInfo) UpdateSN(sn, operator string) error {
	err := nosql.UpdateAreaSN(mine.UID, sn, operator)
	if err == nil {
		mine.SN = sn
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
