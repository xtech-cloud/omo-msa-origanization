package cache

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
	"omo.msa.organization/proxy/nosql"
	"time"
)

type DeviceInfo struct {
	baseInfo
	Remark      string
	Scene       string
	OS          string
	Quote       string
	SN          string
	RunningTime uint32
	Certificate string
}

func (mine *DeviceInfo) initInfo(db *nosql.Device) {
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
	mine.RunningTime = db.Running
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

func (mine *DeviceInfo) UpdateLength(length uint32, operator string) error {
	err := nosql.UpdateDeviceLength(mine.UID, operator, length)
	if err == nil {
		mine.RunningTime = length
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
	}
	return err
}

func (mine *DeviceInfo) Remove(operator string) error {
	return nosql.RemoveDevice(mine.UID, operator)
}

func (mine *cacheContext) CreateDevice(scene, name, sn, remark, os, quote, operator string) (*DeviceInfo, error) {
	db := new(nosql.Device)
	db.UID = primitive.NewObjectID()
	db.ID = nosql.GetRoomNextID()
	db.CreatedTime = time.Now()
	db.UpdatedTime = time.Now()
	db.Operator = operator
	db.Creator = operator
	db.Name = name
	db.Remark = remark
	db.Scene = scene
	db.OS = os
	db.SN = sn
	db.Quote = quote
	db.Running = 0
	db.Certificate = ""
	err := nosql.CreateDevice(db)
	if err == nil {
		tmp := new(DeviceInfo)
		tmp.initInfo(db)
		return tmp, nil
	}
	return nil, err
}

func (mine *cacheContext) GetDevice(sn string) (*DeviceInfo, error) {
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
