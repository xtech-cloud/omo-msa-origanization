package cache

import (
	pb "github.com/xtech-cloud/omo-msp-organization/proto/organization"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"omo.msa.organization/proxy/nosql"
	"time"
)

type MaintainInfo struct {
	Type uint8
	baseInfo
	Scene     string
	Remark    string
	Device    string
	Area      string
	Date      string
	Submitter string
}

func (mine *cacheContext) CreateMaintain(in *pb.ReqMaintainAdd, device string) (*MaintainInfo, error) {
	db := new(nosql.Maintain)
	db.UID = primitive.NewObjectID()
	db.ID = nosql.GetMaintainNextID()
	db.CreatedTime = time.Now()
	db.Creator = in.Operator
	db.Name = in.Name
	db.Remark = in.Remark
	db.Scene = in.Scene
	db.Type = uint8(in.Type)
	db.Area = in.Area
	db.Device = device
	db.Date = in.Date
	db.Submitter = in.Submitter
	err := nosql.CreateMaintain(db)
	if err == nil {
		info := new(MaintainInfo)
		info.initInfo(db)
		return info, nil
	}
	return nil, err
}

func (mine *cacheContext) GetMaintain(uid string) (*MaintainInfo, error) {
	db, err := nosql.GetMaintain(uid)
	if err == nil {
		info := new(MaintainInfo)
		info.initInfo(db)
		return info, nil
	}
	return nil, err
}

func (mine *cacheContext) GetMaintainByScene(uid string) ([]*MaintainInfo, error) {
	dbs, err := nosql.GetMaintainsByOwner(uid)
	list := make([]*MaintainInfo, 0, len(dbs))
	if err == nil {
		for _, db := range dbs {
			info := new(MaintainInfo)
			info.initInfo(db)
			list = append(list, info)
		}
	}
	return list, err
}

func (mine *cacheContext) GetMaintainByArea(scene, area string) ([]*MaintainInfo, error) {
	dbs, err := nosql.GetMaintainsByArea(scene, area)
	list := make([]*MaintainInfo, 0, len(dbs))
	if err == nil {
		for _, db := range dbs {
			info := new(MaintainInfo)
			info.initInfo(db)
			list = append(list, info)
		}
	}
	return list, err
}

func (mine *MaintainInfo) initInfo(db *nosql.Maintain) {
	mine.UID = db.UID.Hex()
	mine.ID = db.ID
	mine.Name = db.Name
	mine.Type = db.Type
	mine.CreateTime = db.CreatedTime
	mine.Creator = db.Creator
	mine.Operator = db.Operator
	mine.Scene = db.Scene
	mine.Remark = db.Remark
	mine.Date = db.Date
	mine.Area = db.Area
	mine.Submitter = db.Submitter
}
