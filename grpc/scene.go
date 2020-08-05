package grpc

import (
	"context"
	"encoding/json"
	"errors"
	"github.com/micro/go-micro/v2/logger"
	pb "github.com/xtech-cloud/omo-msp-organization/proto/organization"
	"omo.msa.organization/cache"
)

type SceneService struct {}

func switchScene(info *cache.SceneInfo) *pb.SceneInfo {
	tmp := new(pb.SceneInfo)
	tmp.Uid = info.UID
	tmp.Id = info.ID
	tmp.Created = info.CreateTime.Unix()
	tmp.Updated = info.UpdateTime.Unix()
	tmp.Name = info.Name
	tmp.Type = int32(info.Type)
	tmp.Status = int32(info.Status)
	tmp.Remark = info.Remark
	tmp.Cover = info.Cover
	tmp.Location = info.Location
	tmp.Master = info.Master
	tmp.Operator = info.Operator
	tmp.Creator = info.Creator
	return tmp
}

func inLog(name, data interface{})  {
	bytes, _ := json.Marshal(data)
	msg := ByteString(bytes)
	logger.Infof("[in.%s]:data = %s", name, msg)
}

func ByteString(p []byte) string {
	for i := 0; i < len(p); i++ {
		if p[i] == 0 {
			return string(p[0:i])
		}
	}
	return string(p)
}

func (mine *SceneService)AddOne(ctx context.Context, in *pb.ReqSceneAdd, out *pb.ReplySceneOne) error {
	inLog("scene.add", in)
	if len(in.Name) < 1 {
		out.Status = pb.ResultStatus_Empty
		return errors.New("the scene name is empty")
	}
	info := new(cache.SceneInfo)
	info.Name = in.Name
	info.Remark = in.Remark
	info.Master = in.Master
	info.Location = in.Location
	info.Type = cache.SceneType(in.Type)
	info.Cover = in.Cover
	err := cache.CreateScene(info)
	if err == nil {
		out.Info = switchScene(info)
	}else{
		out.Status = pb.ResultStatus_DBException
	}

	return err
}

func (mine *SceneService)GetOne(ctx context.Context, in *pb.RequestInfo, out *pb.ReplySceneOne) error {
	inLog("scene.get", in)
	if len(in.Uid) < 1 {
		out.Status = pb.ResultStatus_Empty
		return errors.New("the scene uid is empty")
	}
	info := cache.GetScene(in.Uid)
	if info == nil {
		out.Status = pb.ResultStatus_NotExisted
		return errors.New("the scene not found")
	}
	out.Info = switchScene(info)
	return nil
}

func (mine *SceneService)GetOneByMaster(ctx context.Context, in *pb.RequestInfo, out *pb.ReplySceneOne) error {
	inLog("scene.getByMember", in)
	if len(in.Uid) < 1 {
		out.Status = pb.ResultStatus_Empty
		return errors.New("the scene uid is empty")
	}
	info := cache.GetSceneByMember(in.Uid)
	if info == nil {
		out.Status = pb.ResultStatus_NotExisted
		return errors.New("the scene not found")
	}
	out.Info = switchScene(info)
	return nil
}

func (mine *SceneService)RemoveOne(ctx context.Context, in *pb.RequestInfo, out *pb.ReplyInfo) error {
	inLog("scene.remove", in)
	if len(in.Uid) < 1 {
		out.Status = pb.ResultStatus_Empty
		return errors.New("the scene uid is empty")
	}
	err := cache.RemoveScene(in.Uid, in.Operator)
	if err != nil {
		out.Status = pb.ResultStatus_DBException
	}
	out.Uid = in.Uid
	return err
}

func (mine *SceneService)IsMasterUsed(ctx context.Context, in *pb.RequestInfo, out *pb.ReplyMasterUsed) error {
	inLog("scene.master.used", in)
	if len(in.Uid) < 1 {
		return errors.New("the scene uid is empty")
	}
	out.Master = in.Uid
	out.Used = cache.IsMasterUsed(in.Uid)
	return nil
}

func (mine *SceneService)GetList(ctx context.Context, in *pb.ReqScenePage, out *pb.ReplySceneList) error {
	total,max,list := cache.GetScenes(in.Number, in.Page)
	out.List = make([]*pb.SceneInfo, 0, len(list))
	for _, value := range list {
		out.List = append(out.List, switchScene(value))
	}
	out.PageNow = in.Page
	out.Total = total
	out.PageMax = max
	return nil
}

func (mine *SceneService) UpdateBase (ctx context.Context, in *pb.ReqSceneUpdate, out *pb.ReplySceneOne) error {
	inLog("scene.update", in)
	if len(in.Uid) < 1 {
		out.Status = pb.ResultStatus_Empty
		return errors.New("the scene uid is empty")
	}
	info := cache.GetScene(in.Uid)
	if info == nil {
		out.Status = pb.ResultStatus_NotExisted
		return errors.New("the scene not found")
	}
	var err error
	if len(in.Cover) > 0 {
		err = info.UpdateCover(in.Cover, in.Operator)
	}
	if len(in.Location) > 0 {
		err = info.UpdateLocation(in.Location, in.Operator)
	}
	if len(in.Master) > 0 {
		err = info.UpdateMaster(in.Master, in.Operator)
	}
	if len(in.Name) > 0 || len(in.Remark) > 0 {
		err = info.UpdateBase(in.Name, in.Remark, in.Operator)
	}

	if err != nil {
		out.Status = pb.ResultStatus_DBException
	}else{
		out.Info = switchScene(info)
	}
	return err
}

func (mine *SceneService) UpdateStatus (ctx context.Context, in *pb.ReqSceneStatus, out *pb.ReplySceneOne) error {
	if len(in.Uid) < 1 {
		out.Status = pb.ResultStatus_Empty
		return errors.New("the scene uid is empty")
	}
	info := cache.GetScene(in.Uid)
	if info == nil {
		out.Status = pb.ResultStatus_NotExisted
		return errors.New("the scene not found")
	}
	var err error
	err = info.UpdateStatus(cache.SceneStatus(in.Status),in.Operator)
	if err != nil {
		out.Status = pb.ResultStatus_DBException
	}else{
		out.Info = switchScene(info)
	}
	return err
}

func (mine *SceneService) AppendMember (ctx context.Context, in *pb.ReqSceneMember, out *pb.ReplySceneMembers) error {
	if len(in.Uid) < 1 {
		out.Status = pb.ResultStatus_Empty
		return errors.New("the scene uid is empty")
	}
	info := cache.GetScene(in.Uid)
	if info == nil {
		out.Status = pb.ResultStatus_NotExisted
		return errors.New("the scene not found")
	}
	var err error
	err = info.AppendMember(in.Member)
	out.Uid = in.Uid
	out.Members = info.AllMembers()
	if err != nil {
		out.Status = pb.ResultStatus_DBException
	}
	return err
}

func (mine *SceneService) SubtractMember (ctx context.Context, in *pb.ReqSceneMember, out *pb.ReplySceneMembers) error {
	if len(in.Uid) < 1 {
		out.Status = pb.ResultStatus_Empty
		return errors.New("the scene uid is empty")
	}
	info := cache.GetScene(in.Uid)
	if info == nil {
		out.Status = pb.ResultStatus_NotExisted
		return errors.New("the scene not found")
	}
	var err error
	err = info.SubtractMember(in.Member)
	out.Uid = in.Uid
	out.Members = info.AllMembers()
	if err != nil {
		out.Status = pb.ResultStatus_DBException
	}
	return err
}


