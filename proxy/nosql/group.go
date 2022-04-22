package nosql

import (
	"context"
	"errors"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

// 部门或者虚拟组
type Group struct {
	UID         primitive.ObjectID `bson:"_id"`
	ID          uint64             `json:"id" bson:"id"`
	Name        string             `json:"name" bson:"name"`
	CreatedTime time.Time          `json:"createdAt" bson:"createdAt"`
	UpdatedTime time.Time          `json:"updatedAt" bson:"updatedAt"`
	DeleteTime  time.Time          `json:"deleteAt" bson:"deleteAt"`

	Creator  string `json:"creator" bson:"creator"`
	Operator string `json:"operator" bson:"operator"`
	Scene    string `json:"scene" bson:"scene"`
	Remark   string `json:"remark" bson:"remark"`
	Contact  string `json:"contact" bson:"contact"`
	Cover    string `json:"cover" bson:"cover"`

	Master    string      `json:"master" bson:"master"`
	Assistant string      `json:"assistant" bson:"assistant"`
	Address   AddressInfo `json:"address" bson:"address"`
	Location  string      `json:"location" bson:"location"`
	Members   []string    `json:"members" bson:"members"`
}

func CreateGroup(info *Group) error {
	_, err := insertOne(TableGroup, info)
	if err != nil {
		return err
	}
	return nil
}

func GetGroupNextID() uint64 {
	num, _ := getSequenceNext(TableGroup)
	return num
}

func GetGroupCount() int64 {
	num, _ := getCount(TableGroup)
	return num
}

func GetGroup(uid string) (*Group, error) {
	result, err := findOne(TableGroup, uid)
	if err != nil {
		return nil, err
	}
	model := new(Group)
	err1 := result.Decode(model)
	if err1 != nil {
		return nil, err1
	}
	return model, nil
}

func GetGroupByID(id uint64) (*Group, error) {
	msg := bson.M{"id": id}
	result, err := findOneBy(TableGroup, msg)
	if err != nil {
		return nil, err
	}
	model := new(Group)
	err1 := result.Decode(model)
	if err1 != nil {
		return nil, err1
	}
	return model, nil
}

func RemoveGroup(uid,operator string) error {
	_, err := removeOne(TableGroup, uid, operator)
	return err
}

func GetAllGroups() ([]*Group, error) {
	cursor, err1 := findAll(TableGroup, 0)
	if err1 != nil {
		return nil, err1
	}
	var items = make([]*Group, 0, 20)
	for cursor.Next(context.Background()) {
		var node = new(Group)
		if err := cursor.Decode(node); err != nil {
			return nil, err
		} else {
			items = append(items, node)
		}
	}
	return items, nil
}

func GetGroupsByScene(scene string) ([]*Group, error) {
	cursor, err1 := findMany(TableGroup, bson.M{"scene": scene, "deleteAt": new(time.Time)}, 0)
	if err1 != nil {
		return nil, err1
	}
	var items = make([]*Group, 0, 20)
	for cursor.Next(context.Background()) {
		var node = new(Group)
		if err := cursor.Decode(node); err != nil {
			return nil, err
		} else {
			items = append(items, node)
		}
	}
	return items, nil
}

func UpdateGroupBase(uid, name, remark, operator string) error {
	msg := bson.M{"name": name, "remark": remark, "operator": operator, "updatedAt": time.Now()}
	_, err := updateOne(TableGroup, uid, msg)
	return err
}

func UpdateGroupCover(uid, cover, operator string) error {
	msg := bson.M{"cover": cover,"operator": operator, "updatedAt": time.Now()}
	_, err := updateOne(TableGroup, uid, msg)
	return err
}

func UpdateGroupMembers(uid, operator string, members []string) error {
	msg := bson.M{"operator": operator, "members": members, "updatedAt": time.Now()}
	_, err := updateOne(TableGroup, uid, msg)
	return err
}

func UpdateGroupAddress(uid, operator string, address AddressInfo) error {
	msg := bson.M{"operator": operator, "address": address, "updatedAt": time.Now()}
	_, err := updateOne(TableGroup, uid, msg)
	return err
}

func UpdateGroupLocation(uid, location, operator string) error {
	msg := bson.M{"location": location,"operator": operator, "updatedAt": time.Now()}
	_, err := updateOne(TableGroup, uid, msg)
	return err
}

func UpdateGroupContact(uid, phone, operator string) error {
	msg := bson.M{"contact": phone,"operator": operator, "updatedAt": time.Now()}
	_, err := updateOne(TableGroup, uid, msg)
	return err
}

func UpdateGroupMaster(uid, member, operator string) error {
	msg := bson.M{"master": member,"operator": operator, "updatedAt": time.Now()}
	_, err := updateOne(TableGroup, uid, msg)
	return err
}

func UpdateGroupAssistant(uid, member, operator string) error {
	msg := bson.M{"assistant": member,"operator": operator, "updatedAt": time.Now()}
	_, err := updateOne(TableGroup, uid, msg)
	return err
}

func AppendGroupMember(uid, member string) error {
	if len(member) < 1 {
		return errors.New("the member uid is empty")
	}
	msg := bson.M{"members": member}
	_, err := appendElement(TableGroup, uid, msg)
	return err
}

func SubtractGroupMember(uid string, member string) error {
	if len(member) < 1 {
		return errors.New("the member uid is empty")
	}
	msg := bson.M{"members": member}
	_, err := removeElement(TableGroup, uid, msg)
	return err
}
