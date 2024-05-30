package nosql

import (
	"context"
	"errors"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"omo.msa.organization/proxy"
	"time"
)

type Maintain struct {
	UID         primitive.ObjectID `bson:"_id"`
	ID          uint64             `json:"id" bson:"id"`
	Name        string             `json:"name" bson:"name"`
	CreatedTime time.Time          `json:"createdAt" bson:"createdAt"`
	UpdatedTime time.Time          `json:"updatedAt" bson:"updatedAt"`
	DeleteTime  time.Time          `json:"deleteAt" bson:"deleteAt"`
	Creator     string             `json:"creator" bson:"creator"`
	Operator    string             `json:"operator" bson:"operator"`

	Type        uint8                   `json:"type" bson:"type"`
	Remark      string                  `json:"remark" bson:"remark"`
	Scene       string                  `json:"scene" bson:"scene"`
	Area        string                  `json:"area" bson:"area"`
	Date        string                  `json:"date" bson:"date"`
	Device      string                  `json:"device" bson:"device"`
	Submitter   string                  `json:"submitter" bson:"submitter"`
	Contacts    string                  `json:"contacts" bson:"contacts"`       //客户对接人
	Maintainers []string                `json:"maintainers" json:"maintainers"` //维护人员
	Contents    []proxy.MaintainContent `json:"contents" bson:"contents"`       //内容，原因
}

func CreateMaintain(info *Maintain) error {
	_, err := insertOne(TableMaintain, &info)
	return err
}

func GetMaintainNextID() uint64 {
	num, _ := getSequenceNext(TableMaintain)
	return num
}

func RemoveMaintain(uid, operator string) error {
	if len(uid) < 2 {
		return errors.New("db Maintain uid is empty ")
	}
	_, err := removeOne(TableMaintain, uid, operator)
	return err
}

func GetMaintain(uid string) (*Maintain, error) {
	if len(uid) < 2 {
		return nil, errors.New("db Maintain uid is empty of GetMaintain")
	}

	result, err := findOne(TableMaintain, uid)
	if err != nil {
		return nil, err
	}
	model := new(Maintain)
	err1 := result.Decode(&model)
	if err1 != nil {
		return nil, err1
	}
	return model, nil
}

func GetMaintainsByOwner(owner string) ([]*Maintain, error) {
	filter := bson.M{"scene": owner}
	cursor, err1 := findMany(TableMaintain, filter, 0)
	if err1 != nil {
		return nil, err1
	}
	var items = make([]*Maintain, 0, 20)
	for cursor.Next(context.Background()) {
		var node = new(Maintain)
		if err := cursor.Decode(&node); err != nil {
			return nil, err
		} else {
			items = append(items, node)
		}
	}
	return items, nil
}

func GetMaintainsByArea(scene, area string) ([]*Maintain, error) {
	filter := bson.M{"scene": scene, "area": area}
	cursor, err1 := findMany(TableMaintain, filter, 0)
	if err1 != nil {
		return nil, err1
	}
	var items = make([]*Maintain, 0, 20)
	for cursor.Next(context.Background()) {
		var node = new(Maintain)
		if err := cursor.Decode(&node); err != nil {
			return nil, err
		} else {
			items = append(items, node)
		}
	}
	return items, nil
}

func GetMaintainCount() int64 {
	num, _ := getCount(TableMaintain)
	return num
}
