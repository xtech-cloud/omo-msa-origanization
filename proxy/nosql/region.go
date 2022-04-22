package nosql

import (
	"context"
	"errors"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

// 行政区域
type Region struct {
	UID         primitive.ObjectID `bson:"_id"`
	ID          uint64             `json:"id" bson:"id"`
	Name        string             `json:"name" bson:"name"`
	CreatedTime time.Time          `json:"createdAt" bson:"createdAt"`
	UpdatedTime time.Time          `json:"updatedAt" bson:"updatedAt"`
	DeleteTime  time.Time          `json:"deleteAt" bson:"deleteAt"`

	Creator  string `json:"creator" bson:"creator"`
	Operator string `json:"operator" bson:"operator"`
	Scene    string `json:"scene" bson:"scene"`
	Entity   string `json:"entity" bson:"entity"`
	Remark   string `json:"remark" bson:"remark"`
	Code     string `json:"code" bson:"code"`
	Parent   string `json:"parent" bson:"parent"`
	Master   string `json:"master" bson:"master"`
	Location  string      `json:"location" bson:"location"`
	Address   AddressInfo `json:"address" bson:"address"`
	Members  []string `json:"members" bson:"members"`
}

func CreateRegion(info *Region) error {
	_, err := insertOne(TableRegion, info)
	if err != nil {
		return err
	}
	return nil
}

func GetRegionNextID() uint64 {
	num, _ := getSequenceNext(TableRegion)
	return num
}

func GetRegionCount() int64 {
	num, _ := getCount(TableRegion)
	return num
}

func GetRegion(uid string) (*Region, error) {
	result, err := findOne(TableRegion, uid)
	if err != nil {
		return nil, err
	}
	model := new(Region)
	err1 := result.Decode(model)
	if err1 != nil {
		return nil, err1
	}
	return model, nil
}

func GetRegionByID(id uint64) (*Region, error) {
	msg := bson.M{"id": id}
	result, err := findOneBy(TableRegion, msg)
	if err != nil {
		return nil, err
	}
	model := new(Region)
	err1 := result.Decode(model)
	if err1 != nil {
		return nil, err1
	}
	return model, nil
}

func RemoveRegion(uid,operator string) error {
	_, err := removeOne(TableRegion, uid, operator)
	return err
}

func GetAllRegions() ([]*Region, error) {
	cursor, err1 := findAll(TableRegion, 0)
	if err1 != nil {
		return nil, err1
	}
	var items = make([]*Region, 0, 20)
	for cursor.Next(context.Background()) {
		var node = new(Region)
		if err := cursor.Decode(node); err != nil {
			return nil, err
		} else {
			items = append(items, node)
		}
	}
	return items, nil
}

func GetRegionsByScene(scene string) ([]*Region, error) {
	cursor, err1 := findMany(TableRegion, bson.M{"scene": scene, "deleteAt": new(time.Time)}, 0)
	if err1 != nil {
		return nil, err1
	}
	var items = make([]*Region, 0, 20)
	for cursor.Next(context.Background()) {
		var node = new(Region)
		if err := cursor.Decode(node); err != nil {
			return nil, err
		} else {
			items = append(items, node)
		}
	}
	return items, nil
}

func GetRegionsByParent(parent string) ([]*Region, error) {
	cursor, err1 := findMany(TableRegion, bson.M{"parent": parent, "deleteAt": new(time.Time)}, 0)
	if err1 != nil {
		return nil, err1
	}
	var items = make([]*Region, 0, 20)
	for cursor.Next(context.Background()) {
		var node = new(Region)
		if err := cursor.Decode(node); err != nil {
			return nil, err
		} else {
			items = append(items, node)
		}
	}
	return items, nil
}

func UpdateRegionBase(uid, name, remark, operator string) error {
	msg := bson.M{"name": name, "remark": remark, "operator": operator, "updatedAt": time.Now()}
	_, err := updateOne(TableRegion, uid, msg)
	return err
}

func UpdateRegionMaster(uid, master, operator string) error {
	msg := bson.M{"master": master, "operator": operator, "updatedAt": time.Now()}
	_, err := updateOne(TableRegion, uid, msg)
	return err
}

func UpdateRegionEntity(uid, entity, operator string) error {
	msg := bson.M{"entity": entity, "operator": operator, "updatedAt": time.Now()}
	_, err := updateOne(TableRegion, uid, msg)
	return err
}

func UpdateRegionParent(uid, parent, operator string) error {
	msg := bson.M{"parent": parent, "operator": operator, "updatedAt": time.Now()}
	_, err := updateOne(TableRegion, uid, msg)
	return err
}

func UpdateRegionAddress(uid, operator string, address AddressInfo) error {
	msg := bson.M{"operator": operator, "address": address, "updatedAt": time.Now()}
	_, err := updateOne(TableRegion, uid, msg)
	return err
}

func UpdateRegionLocation(uid, location, operator string) error {
	msg := bson.M{"location": location,"operator": operator, "updatedAt": time.Now()}
	_, err := updateOne(TableRegion, uid, msg)
	return err
}

func AppendRegionMember(uid string, member string) error {
	if len(uid) < 1 {
		return errors.New("the uid is empty")
	}
	msg := bson.M{"members": member}
	_, err := appendElement(TableRegion, uid, msg)
	return err
}

func SubtractRegionMember(uid, member string) error {
	if len(uid) < 1 {
		return errors.New("the uid is empty")
	}
	msg := bson.M{"members": member}
	_, err := removeElement(TableRegion, uid, msg)
	return err
}


