package nosql

import (
	"context"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"omo.msa.organization/proxy"
	"time"
)

type Room struct {
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
	Quotes   []string `json:"quotes" bson::"quotes"`
	//Displays  []proxy.DisplayInfo `json:"displays" bson:"displays"`
}

func CreateRoom(info *Room) error {
	_, err := insertOne(TableRoom, info)
	if err != nil {
		return err
	}
	return nil
}

func GetRoomNextID() uint64 {
	num, _ := getSequenceNext(TableRoom)
	return num
}

func GetRoomCount() int64 {
	num, _ := getCount(TableRoom)
	return num
}

func GetRoom(uid string) (*Room, error) {
	result, err := findOne(TableRoom, uid)
	if err != nil {
		return nil, err
	}
	model := new(Room)
	err1 := result.Decode(model)
	if err1 != nil {
		return nil, err1
	}
	return model, nil
}

func GetRoomByID(id uint64) (*Room, error) {
	msg := bson.M{"id": id}
	result, err := findOneBy(TableRoom, msg)
	if err != nil {
		return nil, err
	}
	model := new(Room)
	err1 := result.Decode(model)
	if err1 != nil {
		return nil, err1
	}
	return model, nil
}

func RemoveRoom(uid,operator string) error {
	_, err := removeOne(TableRoom, uid, operator)
	return err
}

func GetAllRooms() ([]*Room, error) {
	cursor, err1 := findAll(TableRoom, 0)
	if err1 != nil {
		return nil, err1
	}
	var items = make([]*Room, 0, 20)
	for cursor.Next(context.Background()) {
		var node = new(Room)
		if err := cursor.Decode(node); err != nil {
			return nil, err
		} else {
			items = append(items, node)
		}
	}
	return items, nil
}

func GetRoomsByScene(scene string) ([]*Room, error) {
	cursor, err1 := findMany(TableRoom, bson.M{"scene": scene, "deleteAt": new(time.Time)}, 0)
	if err1 != nil {
		return nil, err1
	}
	var items = make([]*Room, 0, 20)
	for cursor.Next(context.Background()) {
		var node = new(Room)
		if err := cursor.Decode(node); err != nil {
			return nil, err
		} else {
			items = append(items, node)
		}
	}
	return items, nil
}

func UpdateRoomBase(uid, name, remark, operator string) error {
	msg := bson.M{"name": name, "remark": remark, "operator": operator, "updatedAt": time.Now()}
	_, err := updateOne(TableRoom, uid, msg)
	return err
}

func UpdateRoomDisplays(uid, operator string, list []*proxy.DisplayInfo) error {
	msg := bson.M{"displays": list, "operator": operator, "updatedAt": time.Now()}
	_, err := updateOne(TableRoom, uid, msg)
	return err
}

func UpdateRoomQuotes(uid, operator string, arr []string) error {
	msg := bson.M{"operator": operator, "quotes": arr, "updatedAt": time.Now()}
	_, err := updateOne(TableRoom, uid, msg)
	return err
}


