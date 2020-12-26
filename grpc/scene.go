package grpc

import (
	"context"
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
	tmp.Entity = info.Entity
	tmp.Operator = info.Operator
	tmp.Creator = info.Creator
	return tmp
}

func (mine *SceneService)AddOne(ctx context.Context, in *pb.ReqSceneAdd, out *pb.ReplySceneOne) error {
	path := "scene.add"
	inLog(path, in)
	if len(in.Name) < 1 {
		out.Status = outError(path,"the name is empty ", pb.ResultStatus_Empty)
		return nil
	}
	info := new(cache.SceneInfo)
	info.Name = in.Name
	info.Remark = in.Remark
	info.Master = in.Master
	info.Location = in.Location
	info.Type = cache.SceneType(in.Type)
	info.Cover = in.Cover
	info.Entity = in.Entity
	err := cache.Context().CreateScene(info)
	if err != nil {
		out.Status = outError(path,err.Error(), pb.ResultStatus_DBException)
		return nil
	}
	out.Info = switchScene(info)
	out.Status = outLog(path, out)
	return nil
}

func (mine *SceneService)GetOne(ctx context.Context, in *pb.RequestInfo, out *pb.ReplySceneOne) error {
	path := "scene.getOne"
	inLog(path, in)
	if len(in.Uid) < 1 {
		out.Status = outError(path,"the uid is empty ", pb.ResultStatus_Empty)
		return nil
	}
	info := cache.Context().GetScene(in.Uid)
	if info == nil {
		out.Status = outError(path,"the scene not found ", pb.ResultStatus_NotExisted)
		return nil
	}
	out.Info = switchScene(info)
	out.Status = outLog(path, out)
	return nil
}

func (mine *SceneService)GetOneByMaster(ctx context.Context, in *pb.RequestInfo, out *pb.ReplySceneOne) error {
	path := "scene.getByMaster"
	inLog(path, in)
	if len(in.Uid) < 1 {
		out.Status = outError(path,"the uid is empty ", pb.ResultStatus_Empty)
		return nil
	}
	info := cache.Context().GetSceneByMember(in.Uid)
	if info == nil {
		out.Status = outError(path,"the scene not found ", pb.ResultStatus_NotExisted)
		return nil
	}
	out.Info = switchScene(info)
	out.Status = outLog(path, out)
	return nil
}

func (mine *SceneService)RemoveOne(ctx context.Context, in *pb.RequestInfo, out *pb.ReplyInfo) error {
	path := "scene.remove"
	inLog(path, in)
	if len(in.Uid) < 1 {
		out.Status = outError(path,"the uid is empty ", pb.ResultStatus_Empty)
		return nil
	}
	err := cache.RemoveScene(in.Uid, in.Operator)
	if err != nil {
		out.Status = outError(path,err.Error(), pb.ResultStatus_DBException)
		return nil
	}
	out.Uid = in.Uid
	out.Status = outLog(path, out)
	return nil
}

func (mine *SceneService)IsMasterUsed(ctx context.Context, in *pb.RequestInfo, out *pb.ReplyMasterUsed) error {
	path := "scene.isUsed"
	inLog(path, in)
	if len(in.Uid) < 1 {
		out.Status = outError(path,"the uid is empty ", pb.ResultStatus_Empty)
		return nil
	}
	out.Master = in.Uid
	out.Used = cache.IsMasterUsed(in.Uid)
	out.Status = outLog(path, out)
	return nil
}

func (mine *SceneService)GetList(ctx context.Context, in *pb.RequestPage, out *pb.ReplySceneList) error {
	path := "scene.getList"
	inLog(path, in)
	total,max,list := cache.Context().GetScenes(in.Number, in.Page)
	out.List = make([]*pb.SceneInfo, 0, len(list))
	for _, value := range list {
		out.List = append(out.List, switchScene(value))
	}
	out.PageNow = in.Page
	out.Total = total
	out.PageMax = max
	out.Status = &pb.ReplyStatus{
		Code: 0,
		Error: "",
	}
	return nil
}

func (mine *SceneService) UpdateBase (ctx context.Context, in *pb.ReqSceneUpdate, out *pb.ReplyInfo) error {
	path := "scene.updateBase"
	inLog(path, in)
	if len(in.Uid) < 1 {
		out.Status = outError(path,"the uid is empty ", pb.ResultStatus_Empty)
		return nil
	}
	info := cache.Context().GetScene(in.Uid)
	if info == nil {
		out.Status = outError(path,"the scene not found ", pb.ResultStatus_NotExisted)
		return nil
	}
	var err error
	if len(in.Cover) > 0 {
		err = info.UpdateCover(in.Cover, in.Operator)
	}
	if len(in.Master) > 0 {
		err = info.UpdateMaster(in.Master, in.Operator)
	}
	if len(in.Name) > 0 || len(in.Remark) > 0 {
		err = info.UpdateBase(in.Name, in.Remark, in.Operator)
	}
	if err != nil {
		out.Status = outError(path,err.Error(), pb.ResultStatus_DBException)
		return nil
	}

	out.Status = outLog(path, out)
	return err
}

func (mine *SceneService) UpdateAddress (ctx context.Context, in *pb.RequestAddress, out *pb.ReplySceneOne) error {
	path := "scene.updateBase"
	inLog(path, in)
	if len(in.Uid) < 1 {
		out.Status = outError(path,"the uid is empty ", pb.ResultStatus_Empty)
		return nil
	}
	info := cache.Context().GetScene(in.Uid)
	if info == nil {
		out.Status = outError(path,"the scene not found ", pb.ResultStatus_NotExisted)
		return nil
	}
	err := info.UpdateAddress(in.Country, in.Province, in.City, in.Zone, in.Operator)
	if err != nil {
		out.Status = outError(path,err.Error(), pb.ResultStatus_DBException)
		return nil
	}

	out.Info = switchScene(info)
	out.Status = outLog(path, out)
	return err
}

func (mine *SceneService) UpdateLocation (ctx context.Context, in *pb.RequestFlag, out *pb.ReplyInfo) error {
	path := "scene.updateBase"
	inLog(path, in)
	if len(in.Uid) < 1 {
		out.Status = outError(path,"the uid is empty ", pb.ResultStatus_Empty)
		return nil
	}
	info := cache.Context().GetScene(in.Uid)
	if info == nil {
		out.Status = outError(path,"the scene not found ", pb.ResultStatus_NotExisted)
		return nil
	}
	err := info.UpdateLocation(in.Flag, in.Operator)
	if err != nil {
		out.Status = outError(path,err.Error(), pb.ResultStatus_DBException)
		return nil
	}

	out.Status = outLog(path, out)
	return err
}

func (mine *SceneService) UpdateStatus (ctx context.Context, in *pb.ReqSceneStatus, out *pb.ReplyInfo) error {
	path := "scene.updateStatus"
	inLog(path, in)
	if len(in.Uid) < 1 {
		out.Status = outError(path,"the uid is empty ", pb.ResultStatus_Empty)
		return nil
	}
	info := cache.Context().GetScene(in.Uid)
	if info == nil {
		out.Status = outError(path,"the scene not found ", pb.ResultStatus_NotExisted)
		return nil
	}

	err := info.UpdateStatus(cache.SceneStatus(in.Status),in.Operator)
	if err != nil {
		out.Status = outError(path,err.Error(), pb.ResultStatus_DBException)
		return nil
	}
	out.Status = outLog(path, out)
	return nil
}

func (mine *SceneService) AppendMember (ctx context.Context, in *pb.RequestMember, out *pb.ReplyMembers) error {
	path := "scene.appendMember"
	inLog(path, in)
	if len(in.Uid) < 1 {
		out.Status = outError(path,"the uid is empty ", pb.ResultStatus_Empty)
		return nil
	}
	info := cache.Context().GetScene(in.Uid)
	if info == nil {
		out.Status = outError(path,"the scene not found ", pb.ResultStatus_NotExisted)
		return nil
	}

	err := info.AppendMember(in.Member)
	if err != nil {
		out.Status = outError(path,err.Error(), pb.ResultStatus_DBException)
		return nil
	}
	out.Uid = in.Uid
	out.Members = info.AllMembers()
	out.Status = outLog(path, out)
	return nil
}

func (mine *SceneService) SubtractMember (ctx context.Context, in *pb.RequestMember, out *pb.ReplyMembers) error {
	path := "scene.subtractMember"
	inLog(path, in)
	if len(in.Uid) < 1 {
		out.Status = outError(path,"the uid is empty ", pb.ResultStatus_Empty)
		return nil
	}
	info := cache.Context().GetScene(in.Uid)
	if info == nil {
		out.Status = outError(path,"the scene not found ", pb.ResultStatus_NotExisted)
		return nil
	}

	err := info.SubtractMember(in.Member)
	if err != nil {
		out.Status = outError(path,err.Error(), pb.ResultStatus_DBException)
		return nil
	}
	out.Uid = in.Uid
	out.Members = info.AllMembers()
	out.Status = outLog(path, out)
	return nil
}


