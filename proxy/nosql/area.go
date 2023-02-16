package nosql

import (
	"context"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

type Area struct {
	UID         primitive.ObjectID `bson:"_id"`
	ID          uint64             `json:"id" bson:"id"`
	CreatedTime time.Time          `json:"createdAt" bson:"createdAt"`
	UpdatedTime time.Time          `json:"updatedAt" bson:"updatedAt"`
	DeleteTime  time.Time          `json:"deletedAt" bson:"deletedAt"`
	Creator     string             `json:"creator" bson:"creator"`
	Operator    string             `json:"operator" bson:"operator"`

	Type     uint32   `json:"type" bson:"type"`
	Name     string   `json:"name" bson:"name"`
	Remark   string   `json:"remark" bson:"remark"`
	Scene    string   `json:"scene" bson:"scene"`
	Parent   string   `json:"parent" bson:"parent"`
	Template string   `json:"template" bson:"template"`
	Width    int32    `json:"width" bson:"width"`
	Height   int32    `json:"height" bson:"height"`
	Device   string   `json:"device" bson:"device"`
	Question string   `json:"question" bson:"question"`
	Catalog  string   `json:"catalog" bson:"catalog"`
	Displays []string `json:"displays" bson:"displays"` // 正在设备展览的列表
}

func CreateArea(info *Area) error {
	_, err := insertOne(TableArea, info)
	if err != nil {
		return err
	}
	return nil
}

func GetAreaNextID() uint64 {
	num, _ := getSequenceNext(TableArea)
	return num
}

func GetArea(uid string) (*Area, error) {
	result, err := findOne(TableArea, uid)
	if err != nil {
		return nil, err
	}
	model := new(Area)
	err1 := result.Decode(model)
	if err1 != nil {
		return nil, err1
	}
	return model, nil
}

func GetAreasByOwner(owner string) ([]*Area, error) {
	msg := bson.M{"scene": owner, "deletedAt": new(time.Time)}
	cursor, err1 := findMany(TableArea, msg, 0)
	if err1 != nil {
		return nil, err1
	}
	defer cursor.Close(context.Background())
	var items = make([]*Area, 0, 50)
	for cursor.Next(context.Background()) {
		var node = new(Area)
		if err := cursor.Decode(node); err != nil {
			return nil, err
		} else {
			items = append(items, node)
		}
	}
	return items, nil
}

func GetAreasByParent(parent string) ([]*Area, error) {
	msg := bson.M{"parent": parent, "deletedAt": new(time.Time)}
	cursor, err1 := findMany(TableArea, msg, 0)
	if err1 != nil {
		return nil, err1
	}
	defer cursor.Close(context.Background())
	var items = make([]*Area, 0, 50)
	for cursor.Next(context.Background()) {
		var node = new(Area)
		if err := cursor.Decode(node); err != nil {
			return nil, err
		} else {
			items = append(items, node)
		}
	}
	return items, nil
}

func GetAreasBy(owner, parent string) ([]*Area, error) {
	msg := bson.M{"scene": owner, "parent": parent, "deletedAt": new(time.Time)}
	cursor, err1 := findMany(TableArea, msg, 0)
	if err1 != nil {
		return nil, err1
	}
	defer cursor.Close(context.Background())
	var items = make([]*Area, 0, 50)
	for cursor.Next(context.Background()) {
		var node = new(Area)
		if err := cursor.Decode(node); err != nil {
			return nil, err
		} else {
			items = append(items, node)
		}
	}
	return items, nil
}

func GetAllAreas() ([]*Area, error) {
	cursor, err1 := findAll(TableArea, 0)
	if err1 != nil {
		return nil, err1
	}
	defer cursor.Close(context.Background())
	var items = make([]*Area, 0, 100)
	for cursor.Next(context.Background()) {
		var node = new(Area)
		if err := cursor.Decode(node); err != nil {
			return nil, err
		} else {
			items = append(items, node)
		}
	}
	return items, nil
}

func UpdateAreaBase(uid, name, remark, operator string) error {
	msg := bson.M{"name": name, "remark": remark, "operator": operator, "updatedAt": time.Now()}
	_, err := updateOne(TableArea, uid, msg)
	return err
}

func UpdateAreaTemplate(uid, template, operator string) error {
	msg := bson.M{"template": template, "operator": operator, "updatedAt": time.Now()}
	_, err := updateOne(TableArea, uid, msg)
	return err
}

func UpdateAreaDevice(uid, device, operator string, tp uint32) error {
	msg := bson.M{"device": device, "type": tp, "operator": operator, "updatedAt": time.Now()}
	_, err := updateOne(TableArea, uid, msg)
	return err
}

func UpdateAreaCatalog(uid, catalog, operator string) error {
	msg := bson.M{"catalog": catalog, "operator": operator, "updatedAt": time.Now()}
	_, err := updateOne(TableArea, uid, msg)
	return err
}

func UpdateAreaDevice2(uid, device, operator string) error {
	msg := bson.M{"device": device, "operator": operator, "updatedAt": time.Now()}
	_, err := updateOne(TableArea, uid, msg)
	return err
}

func UpdateAreaType(uid, operator string, tp uint32) error {
	msg := bson.M{"type": tp, "operator": operator, "updatedAt": time.Now()}
	_, err := updateOne(TableArea, uid, msg)
	return err
}

func UpdateAreaQuestion(uid, question, operator string) error {
	msg := bson.M{"question": question, "operator": operator, "updatedAt": time.Now()}
	_, err := updateOne(TableArea, uid, msg)
	return err
}

func UpdateAreaDisplays(uid, operator string, displays []string) error {
	msg := bson.M{"displays": displays, "operator": operator, "updatedAt": time.Now()}
	_, err := updateOne(TableArea, uid, msg)
	return err
}

func RemoveArea(uid, operator string) error {
	_, err := removeOne(TableArea, uid, operator)
	return err
}
