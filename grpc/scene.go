package grpc

import (
	"context"
	"errors"
	"omo.msa.organization/cache"
	pb "omo.msa.organization/grpc/proto"
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
	return tmp
}

func (mine *SceneService)AddOne(ctx context.Context, in *pb.ReqSceneAdd, out *pb.ReplySceneOne) error {
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

func (mine *SceneService)RemoveOne(ctx context.Context, in *pb.RequestInfo, out *pb.ReplyInfo) error {
	if len(in.Uid) < 1 {
		out.Status = pb.ResultStatus_Empty
		return errors.New("the scene uid is empty")
	}
	err := cache.RemoveScene(in.Uid)
	if err != nil {
		out.Status = pb.ResultStatus_DBException
	}
	out.Uid = in.Uid
	return err
}

func (mine *SceneService)GetList(ctx context.Context, in *pb.ReqSceneList, out *pb.ReplySceneList) error {
	out.List = make([]*pb.SceneInfo, 0, len(in.List))
	for _, value := range in.List {
		info := cache.GetScene(value)
		if info != nil {
			out.List = append(out.List, switchScene(info))
		}
	}
	return nil
}

func (mine *SceneService) UpdateBase (ctx context.Context, in *pb.ReqSceneUpdate, out *pb.ReplySceneOne) error {
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
		err = info.UpdateCover(in.Cover)
	}else if len(in.Location) > 0 {
		err = info.UpdateLocation(in.Location)
	}else {
		err = info.UpdateBase(in.Name, in.Remark)
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
	err = info.UpdateStatus(cache.SceneStatus(in.Status))
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


