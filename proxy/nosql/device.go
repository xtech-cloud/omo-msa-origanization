package nosql

import (
	"context"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

type Device struct {
	UID         primitive.ObjectID `bson:"_id"`
	ID          uint64             `json:"id" bson:"id"`
	Name        string             `json:"name" bson:"name"`
	CreatedTime time.Time          `json:"createdAt" bson:"createdAt"`
	UpdatedTime time.Time          `json:"updatedAt" bson:"updatedAt"`
	DeleteTime  time.Time          `json:"deleteAt" bson:"deleteAt"`

	Creator  string `json:"creator" bson:"creator"`
	Operator string `json:"operator" bson:"operator"`
	Scene    string `json:"scene" bson:"scene"`
	Room     string `json:"room" bson:"room"`
	Remark   string `json:"remark" bson:"remark"`
	Type     uint32 `json:"type" bson:"type"` //产品类型
	SN 		 string `json:"sn" bson:"sn"`
	Displays  []string `json:"displays" bson:"displays"`
}

func CreateDevice(info *Device) error {
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

func GetDeviceBySN(sn string) (*Device, error) {
	msg := bson.M{"sn": sn}
	result, err := findOneBy(TableDevice, msg)
	if err != nil {
		return nil, err
	}
	model := new(Device)
	err1 := result.Decode(model)
	if err1 != nil {
		return nil, err1
	}
	return model, nil
}

func GetDeviceByID(id uint64) (*Device, error) {
	msg := bson.M{"id": id}
	result, err := findOneBy(TableDevice, msg)
	if err != nil {
		return nil, err
	}
	model := new(Device)
	err1 := result.Decode(model)
	if err1 != nil {
		return nil, err1
	}
	return model, nil
}

func RemoveDevice(uid,operator string) error {
	_, err := removeOne(TableDevice, uid, operator)
	return err
}

func GetAllDevices() ([]*Device, error) {
	cursor, err1 := findAll(TableDevice, 0)
	if err1 != nil {
		return nil, err1
	}
	var items = make([]*Device, 0, 20)
	for cursor.Next(context.Background()) {
		var node = new(Device)
		if err := cursor.Decode(node); err != nil {
			return nil, err
		} else {
			items = append(items, node)
		}
	}
	return items, nil
}

func GetDevicesByScene(scene string) ([]*Device, error) {
	cursor, err1 := findMany(TableDevice, bson.M{"scene": scene, "deleteAt": new(time.Time)}, 0)
	if err1 != nil {
		return nil, err1
	}
	var items = make([]*Device, 0, 20)
	for cursor.Next(context.Background()) {
		var node = new(Device)
		if err := cursor.Decode(node); err != nil {
			return nil, err
		} else {
			items = append(items, node)
		}
	}
	return items, nil
}

func GetDevicesByRoom(scene,room string) ([]*Device, error) {
	cursor, err1 := findMany(TableDevice, bson.M{"scene": scene, "room":room, "deleteAt": new(time.Time)}, 0)
	if err1 != nil {
		return nil, err1
	}
	var items = make([]*Device, 0, 20)
	for cursor.Next(context.Background()) {
		var node = new(Device)
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

func UpdateDeviceRoom(uid, room, operator string) error {
	msg := bson.M{"room": room, "operator": operator, "updatedAt": time.Now()}
	_, err := updateOne(TableDevice, uid, msg)
	return err
}

func UpdateDeviceDisplays(uid, operator string, displays []string) error {
	msg := bson.M{"displays": displays, "operator": operator, "updatedAt": time.Now()}
	_, err := updateOne(TableDevice, uid, msg)
	return err
}

