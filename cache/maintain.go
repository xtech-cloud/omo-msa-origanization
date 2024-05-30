package cache

import (
	pb "github.com/xtech-cloud/omo-msp-organization/proto/organization"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"omo.msa.organization/proxy"
	"omo.msa.organization/proxy/nosql"
	"time"
)

type MaintainInfo struct {
	Type uint8
	baseInfo
	Scene       string //场景
	Remark      string
	Device      string //终端UID
	Area        string //区域
	Date        string
	Submitter   string
	Contacts    string                  //客户对接人
	Maintainers []string                //维护人员
	Contents    []proxy.MaintainContent //内容，原因
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
	db.Contacts = in.Contacts
	db.Maintainers = in.Maintainers
	if db.Maintainers == nil {
		db.Maintainers = make([]string, 0, 1)
	}
	db.Contents = make([]proxy.MaintainContent, 0, len(in.Contents))
	for _, item := range in.Contents {
		db.Contents = append(db.Contents, proxy.MaintainContent{Type: item.Type, Content: item.Content, Assets: item.Assets})
	}
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
	mine.Contacts = db.Contacts
	mine.Maintainers = db.Maintainers
	mine.Contents = db.Contents
}
