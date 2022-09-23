package cache

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
	"omo.msa.organization/proxy/nosql"
	"time"
)

type DeviceInfo struct {
	Type uint32
	baseInfo
	Remark   string
	Scene    string
	Room     string
	Area     string
	SN       string
	Question string
	Displays []string
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
	mine.Room = db.Room
	mine.Type = db.Type
	mine.SN = db.SN
	mine.Question = db.Question
	mine.Area = db.Area
	mine.Displays = db.Displays
	if mine.Displays == nil {
		mine.Displays = make([]string, 0, 1)
	}
}

func (mine *DeviceInfo) UpdateRoom(room, area, operator string) error {
	err := nosql.UpdateDeviceRoom(mine.UID, room, area, operator)
	if err == nil {
		mine.Room = room
		mine.Area = area
		mine.Operator = operator
	}
	return err
}

func (mine *DeviceInfo) UpdateShowings(operator string, list []string) error {
	err := nosql.UpdateDeviceDisplays(mine.UID, operator, list)
	if err == nil {
		mine.Displays = list
		mine.Operator = operator
	}
	return err
}

func (mine *cacheContext) createDevice(scene, room, area, sn, remark, operator string, tp uint32) (*DeviceInfo, error) {
	db := new(nosql.Device)
	db.UID = primitive.NewObjectID()
	db.ID = nosql.GetRoomNextID()
	db.CreatedTime = time.Now()
	db.UpdatedTime = time.Now()
	db.Operator = operator
	db.Creator = operator
	db.Name = ""
	db.Remark = remark
	db.Scene = scene
	db.Room = room
	db.Type = tp
	db.SN = sn
	db.Question = ""
	db.Area = area
	db.Displays = make([]string, 0, 1)
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

func (mine *cacheContext) checkDevice(scene, room, area, sn, remark, operator string, tp uint32) (*DeviceInfo, error) {
	info, err := mine.GetDevice(sn)
	if err == nil {
		return info, nil
	}
	return mine.createDevice(scene, room, area, sn, remark, operator, tp)
}

func (mine *DeviceInfo) UpdateQuestion(question, operator string) error {
	err := nosql.UpdateDeviceQuestion(mine.UID, question, operator)
	if err == nil {
		mine.Question = question
		mine.Operator = operator
	}
	return err
}
