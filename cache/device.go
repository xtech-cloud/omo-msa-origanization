package cache

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
	"omo.msa.organization/proxy/nosql"
	"time"
)

const (
	DeviceIdle    = 0  //未绑定
	DeviceBind    = 1  //已经绑定但未分配
	DeviceFill    = 2  //已经绑定也分配了场景
	DevicePend    = 3  //已经分配但未绑定
	DeviceDiscard = 99 //废弃
)

//邀请码
type DeviceInfo struct {
	Status uint8
	Type   uint8
	baseInfo
	Remark      string
	Scene       string
	OS          string
	Quote       string
	SN          string //邀请码
	Aspect      string
	ActiveTime  int64  //激活时间
	Expired     uint32 //有效时间
	Certificate string //证书
}

func (mine *DeviceInfo) initInfo(db *nosql.Invite) {
	mine.UID = db.UID.Hex()
	mine.ID = db.ID
	mine.UpdateTime = db.UpdatedTime
	mine.CreateTime = db.CreatedTime
	mine.Creator = db.Creator
	mine.Operator = db.Operator
	mine.Name = db.Name
	mine.Remark = db.Remark
	mine.Scene = db.Scene
	mine.OS = db.OS
	mine.Quote = db.Quote
	mine.SN = db.SN
	mine.Type = db.Type
	mine.Aspect = db.Aspect
	mine.ActiveTime = db.ActiveTime
	mine.Expired = db.ExpiryTime
	mine.Certificate = db.Certificate
}

func (mine *DeviceInfo) UpdateBase(name, remark, operator string) error {
	err := nosql.UpdateDeviceBase(mine.UID, name, remark, operator)
	if err == nil {
		mine.Name = name
		mine.Remark = remark
		mine.Operator = operator
	}
	return err
}

func (mine *DeviceInfo) UpdateCertificate(data, operator string) error {
	err := nosql.UpdateDeviceCertificate(mine.UID, data, operator)
	if err == nil {
		mine.Certificate = data
		mine.Operator = operator
	}
	return err
}

func (mine *DeviceInfo) UpdateScene(data, operator string) error {
	err := nosql.UpdateDeviceScene(mine.UID, data, operator)
	if err == nil {
		mine.Scene = data
		mine.Operator = operator
		mine.updateStatus(operator)
	}
	return err
}

func (mine *DeviceInfo) UpdateAspect(data, operator string) error {
	err := nosql.UpdateDeviceAspect(mine.UID, data, operator)
	if err == nil {
		mine.Aspect = data
		mine.Operator = operator
		mine.UpdateTime = time.Now()
	}
	return err
}

func (mine *DeviceInfo) updateStatus(operator string) {
	st := DeviceIdle
	if len(mine.Scene) > 2 && len(mine.Quote) > 2 {
		st = DeviceFill
	} else if len(mine.Scene) > 2 && len(mine.Quote) < 2 {
		st = DevicePend
	} else if len(mine.Scene) < 2 && len(mine.Quote) > 2 {
		st = DeviceBind
	} else {
		return
	}
	er := nosql.UpdateDeviceStatus(mine.UID, operator, uint8(st))
	if er == nil {
		mine.Status = uint8(st)
	}
}

func (mine *DeviceInfo) UpdateType(operator string, tp uint8) error {
	err := nosql.UpdateDeviceType(mine.UID, operator, tp)
	if err == nil {
		mine.Type = tp
		mine.Operator = operator
	}
	return err
}

func (mine *DeviceInfo) Bind(quote, os, operator string, act, expired uint64) error {
	err := nosql.BindDevice(mine.UID, quote, os, operator, act, expired)
	if err == nil {
		mine.Quote = quote
		mine.OS = os
		mine.ActiveTime = int64(act)
		mine.Expired = uint32(expired)
		mine.Operator = operator
		mine.updateStatus(operator)
	}
	return err
}

func (mine *DeviceInfo) Remove(operator string) error {
	//return nosql.RemoveDevice(mine.UID, operator)
	return nosql.UpdateDeviceStatus(mine.UID, operator, DeviceDiscard)
}

func (mine *cacheContext) CreateDevice(scene, name, sn, remark, operator string, tp uint8) (*DeviceInfo, error) {
	db := new(nosql.Invite)
	db.UID = primitive.NewObjectID()
	db.ID = nosql.GetRoomNextID()
	db.CreatedTime = time.Now()
	db.UpdatedTime = time.Now()
	db.Operator = operator
	db.Creator = operator
	db.Name = name
	db.Remark = remark
	db.Scene = scene
	db.OS = ""
	db.SN = sn
	db.Quote = ""
	db.ActiveTime = 0
	db.ExpiryTime = 0
	db.Certificate = ""
	db.Type = tp
	db.Status = DeviceIdle
	err := nosql.CreateDevice(db)
	if err == nil {
		tmp := new(DeviceInfo)
		tmp.initInfo(db)
		return tmp, nil
	}
	return nil, err
}

func (mine *cacheContext) GetDevice(uid string) (*DeviceInfo, error) {
	db, err := nosql.GetDevice(uid)
	if err != nil {
		return nil, err
	}
	tmp := new(DeviceInfo)
	tmp.initInfo(db)
	return tmp, nil
}

func (mine *cacheContext) GetDeviceCount() int64 {
	return nosql.GetDeviceCount()
}

func (mine *cacheContext) GetDeviceBySN(sn string) (*DeviceInfo, error) {
	db, err := nosql.GetDeviceBySN(sn)
	if err != nil {
		return nil, err
	}
	tmp := new(DeviceInfo)
	tmp.initInfo(db)
	return tmp, nil
}

func (mine *cacheContext) GetDevicesByOwner(owner string) ([]*DeviceInfo, error) {
	dbs, err := nosql.GetDevicesByScene(owner)
	if err != nil {
		return nil, err
	}
	list := make([]*DeviceInfo, 0, len(dbs))
	for _, db := range dbs {
		tmp := new(DeviceInfo)
		tmp.initInfo(db)
		list = append(list, tmp)
	}

	return list, nil
}

func (mine *cacheContext) GetDevicesByStatus(st int32) ([]*DeviceInfo, error) {
	var dbs []*nosql.Invite
	var err error
	if st < 0 {
		dbs, err = nosql.GetAllDevicesExcept(DeviceDiscard)
	} else {
		dbs, err = nosql.GetDevicesByStatus(uint8(st))
	}

	if err != nil {
		return nil, err
	}
	list := make([]*DeviceInfo, 0, len(dbs))
	for _, db := range dbs {
		tmp := new(DeviceInfo)
		tmp.initInfo(db)
		list = append(list, tmp)
	}

	return list, nil
}
