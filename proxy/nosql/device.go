package nosql

import (
	"context"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"omo.msa.organization/proxy"
	"time"
)

type Invite struct {
	UID         primitive.ObjectID `bson:"_id"`
	ID          uint64             `json:"id" bson:"id"`
	Name        string             `json:"name" bson:"name"`
	CreatedTime time.Time          `json:"createdAt" bson:"createdAt"`
	UpdatedTime time.Time          `json:"updatedAt" bson:"updatedAt"`
	DeleteTime  time.Time          `json:"deleteAt" bson:"deleteAt"`

	Creator     string         `json:"creator" bson:"creator"`
	Operator    string         `json:"operator" bson:"operator"`
	Scene       string         `json:"scene" bson:"scene"` // 所属场景
	Type        uint8          `json:"type" bson:"type"`   //类型
	Status      uint8          `json:"status" bson:"status"`
	Remark      string         `json:"remark" bson:"remark"` //备注
	SN          string         `json:"sn" bson:"sn"`         //设备SN或者邀请码
	OS          string         `json:"os" bson:"os"`         //操作系统
	Aspect      string         `json:"aspect" bson:"aspect"`
	ExpiryTime  uint32         `json:"expiry" bson:"expiry"`           //有效时长
	ActiveTime  int64          `json:"activated" bson:"activated"`     //激活时间
	Quote       string         `json:"quote" bson:"quote"`             //
	Certificate string         `json:"certificate" bson:"certificate"` //激活证书
	Meta        string         `json:"meta" bson:"meta"`
	Auto        proxy.AutoInfo `json:"auto" bson:"auto"`
}

func CreateDevice(info *Invite) error {
	_, err := insertOne(TableDevice, info)
	if err != nil {
		return err
	}
	return nil
}

func GetDeviceNextID() uint64 {
	num, _ := getSequenceNext(TableDevice)
	return num
}

func GetDeviceCount() int64 {
	num, _ := getCount(TableDevice)
	return num
}

func GetDeviceBySN(sn string) (*Invite, error) {
	msg := bson.M{"sn": sn}
	result, err := findOneBy(TableDevice, msg)
	if err != nil {
		return nil, err
	}
	model := new(Invite)
	err1 := result.Decode(model)
	if err1 != nil {
		return nil, err1
	}
	return model, nil
}

func GetDevice(uid string) (*Invite, error) {
	result, err := findOne(TableDevice, uid)
	if err != nil {
		return nil, err
	}
	model := new(Invite)
	err1 := result.Decode(model)
	if err1 != nil {
		return nil, err1
	}
	return model, nil
}

func GetDeviceByID(id uint64) (*Invite, error) {
	msg := bson.M{"id": id}
	result, err := findOneBy(TableDevice, msg)
	if err != nil {
		return nil, err
	}
	model := new(Invite)
	err1 := result.Decode(model)
	if err1 != nil {
		return nil, err1
	}
	return model, nil
}

func RemoveDevice(uid, operator string) error {
	_, err := removeOne(TableDevice, uid, operator)
	return err
}

func GetAllDevices() ([]*Invite, error) {
	cursor, err1 := findAll(TableDevice, 0)
	if err1 != nil {
		return nil, err1
	}
	var items = make([]*Invite, 0, 20)
	for cursor.Next(context.Background()) {
		var node = new(Invite)
		if err := cursor.Decode(node); err != nil {
			return nil, err
		} else {
			items = append(items, node)
		}
	}
	return items, nil
}

func GetAllDevicesExcept(st uint8) ([]*Invite, error) {
	cursor, err1 := findMany(TableDevice, bson.M{"status": bson.M{"$ne": st}, "deleteAt": new(time.Time)}, 0)
	if err1 != nil {
		return nil, err1
	}
	var items = make([]*Invite, 0, 20)
	for cursor.Next(context.Background()) {
		var node = new(Invite)
		if err := cursor.Decode(node); err != nil {
			return nil, err
		} else {
			items = append(items, node)
		}
	}
	return items, nil
}

func GetDevicesByScene(scene string) ([]*Invite, error) {
	cursor, err1 := findMany(TableDevice, bson.M{"scene": scene, "deleteAt": new(time.Time)}, 0)
	if err1 != nil {
		return nil, err1
	}
	var items = make([]*Invite, 0, 20)
	for cursor.Next(context.Background()) {
		var node = new(Invite)
		if err := cursor.Decode(node); err != nil {
			return nil, err
		} else {
			items = append(items, node)
		}
	}
	return items, nil
}

func GetDevicesByStatus(st uint8) ([]*Invite, error) {
	cursor, err1 := findMany(TableDevice, bson.M{"status": st, "deleteAt": new(time.Time)}, 0)
	if err1 != nil {
		return nil, err1
	}
	var items = make([]*Invite, 0, 20)
	for cursor.Next(context.Background()) {
		var node = new(Invite)
		if err := cursor.Decode(node); err != nil {
			return nil, err
		} else {
			items = append(items, node)
		}
	}
	return items, nil
}

func UpdateDeviceBase(uid, name, remark, operator string) error {
	msg := bson.M{"name": name, "remark": remark, "operator": operator, "updatedAt": time.Now()}
	_, err := updateOne(TableDevice, uid, msg)
	return err
}

func UpdateDeviceTime(uid, operator string, act, expiry uint64) error {
	msg := bson.M{"activated": act, "expiry": expiry, "operator": operator, "updatedAt": time.Now()}
	_, err := updateOne(TableDevice, uid, msg)
	return err
}

func UpdateDeviceCertificate(uid, data, operator string) error {
	msg := bson.M{"certificate": data, "operator": operator, "updatedAt": time.Now()}
	_, err := updateOne(TableDevice, uid, msg)
	return err
}

func UpdateDeviceScene(uid, data, operator string) error {
	msg := bson.M{"scene": data, "operator": operator, "updatedAt": time.Now()}
	_, err := updateOne(TableDevice, uid, msg)
	return err
}

func UpdateDeviceAspect(uid, data, operator string) error {
	msg := bson.M{"aspect": data, "operator": operator, "updatedAt": time.Now()}
	_, err := updateOne(TableDevice, uid, msg)
	return err
}

func BindDevice(uid, quote, os, operator string, act, expiry uint64) error {
	msg := bson.M{"quote": quote, "os": os, "activated": act, "expiry": expiry, "operator": operator, "updatedAt": time.Now()}
	_, err := updateOne(TableDevice, uid, msg)
	return err
}

func UpdateDeviceStatus(uid, operator string, st uint8) error {
	msg := bson.M{"status": st, "operator": operator, "updatedAt": time.Now()}
	_, err := updateOne(TableDevice, uid, msg)
	return err
}

func UpdateDeviceMeta(uid, meta, operator string) error {
	msg := bson.M{"meta": meta, "operator": operator, "updatedAt": time.Now()}
	_, err := updateOne(TableDevice, uid, msg)
	return err
}

func UpdateDeviceAuto(uid, operator string, auto proxy.AutoInfo) error {
	msg := bson.M{"auto": auto, "operator": operator, "updatedAt": time.Now()}
	_, err := updateOne(TableDevice, uid, msg)
	return err
}

func UpdateDeviceType(uid, operator string, tp uint8) error {
	msg := bson.M{"type": tp, "operator": operator, "updatedAt": time.Now()}
	_, err := updateOne(TableDevice, uid, msg)
	return err
}
