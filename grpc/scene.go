package grpc

import (
	"context"
	pb "github.com/xtech-cloud/omo-msp-organization/proto/organization"
	"omo.msa.organization/cache"
	"omo.msa.organization/proxy"
	"strconv"
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
	tmp.Supporter = info.Supporter
	tmp.Domain = info.Domain
	tmp.Parents = info.Parents()
	tmp.Members = info.AllMembers()
	tmp.Devices = info.Devices()
	tmp.Exhibitions = switchExhibitions(info.Exhibitions)
	return tmp
}

func switchExhibitions(list []proxy.ShowingInfo) []*pb.ExhibitInfo {
	array := make([]*pb.ExhibitInfo, 0, len(list))
	for _, info := range list {
		array = append(array, &pb.ExhibitInfo{Uid: info.UID, Effect: info.Effect, Skin: info.Skin, Slots:info.Slots})
	}
	return array
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
	info.Creator = in.Operator
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
	path := "scene.isMasterUsed"
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
	var total uint32 = 0
	var max uint32 = 0
	var list []*cache.SceneInfo
	if in.Parent == "" {
		total,max,list = cache.Context().GetScenes(in.Number, in.Page)
	}else{
		tp,er := strconv.ParseUint(in.Parent, 10, 32)
		if er != nil {
			out.Status = outError(path, er.Error(), pb.ResultStatus_DBException)
			return nil
		}
		list = cache.Context().GetScenesByType(uint8(tp))
	}

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
	if in.Type > 0 {
		err = info.UpdateType(in.Operator, uint8(in.Type))
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

func (mine *SceneService) UpdateSupporter (ctx context.Context, in *pb.RequestInfo, out *pb.ReplyInfo) error {
	path := "scene.updateSupporter"
	inLog(path, in)
	if len(in.Uid) < 1 {
		out.Status = outError(path,"the supporter uid is empty ", pb.ResultStatus_Empty)
		return nil
	}
	if len(in.Parent) < 1 {
		out.Status = outError(path,"the scene uid is empty ", pb.ResultStatus_Empty)
		return nil
	}
	info := cache.Context().GetScene(in.Parent)
	if info == nil {
		out.Status = outError(path,"the scene not found ", pb.ResultStatus_NotExisted)
		return nil
	}

	err := info.UpdateSupporter(in.Uid,in.Operator)
	if err != nil {
		out.Status = outError(path,err.Error(), pb.ResultStatus_DBException)
		return nil
	}
	out.Status = outLog(path, out)
	return nil
}

func (mine *SceneService) UpdateDomain (ctx context.Context, in *pb.RequestInfo, out *pb.ReplyInfo) error {
	path := "scene.updateSupporter"
	inLog(path, in)
	if len(in.Uid) < 1 {
		out.Status = outError(path,"the domain uid is empty ", pb.ResultStatus_Empty)
		return nil
	}
	if len(in.Parent) < 1 {
		out.Status = outError(path,"the scene uid is empty ", pb.ResultStatus_Empty)
		return nil
	}
	info := cache.Context().GetScene(in.Parent)
	if info == nil {
		out.Status = outError(path,"the scene not found ", pb.ResultStatus_NotExisted)
		return nil
	}

	err := info.UpdateDomain(in.Uid, in.Operator)
	if err != nil {
		out.Status = outError(path,err.Error(), pb.ResultStatus_DBException)
		return nil
	}
	out.Status = outLog(path, out)
	return nil
}

func (mine *SceneService) AppendMember (ctx context.Context, in *pb.RequestInfo, out *pb.ReplyList) error {
	path := "scene.appendMember"
	inLog(path, in)
	if len(in.Parent) < 1 {
		out.Status = outError(path,"the parent is empty ", pb.ResultStatus_Empty)
		return nil
	}
	info := cache.Context().GetScene(in.Parent)
	if info == nil {
		out.Status = outError(path,"the scene not found ", pb.ResultStatus_NotExisted)
		return nil
	}

	err := info.AppendMember(in.Uid)
	if err != nil {
		out.Status = outError(path,err.Error(), pb.ResultStatus_DBException)
		return nil
	}
	out.Uid = in.Uid
	out.List = info.AllMembers()
	out.Status = outLog(path, out)
	return nil
}

func (mine *SceneService) SubtractMember (ctx context.Context, in *pb.RequestInfo, out *pb.ReplyList) error {
	path := "scene.subtractMember"
	inLog(path, in)
	if len(in.Parent) < 1 {
		out.Status = outError(path,"the parent is empty ", pb.ResultStatus_Empty)
		return nil
	}
	info := cache.Context().GetScene(in.Parent)
	if info == nil {
		out.Status = outError(path,"the scene not found ", pb.ResultStatus_NotExisted)
		return nil
	}

	err := info.SubtractMember(in.Uid)
	if err != nil {
		out.Status = outError(path,err.Error(), pb.ResultStatus_DBException)
		return nil
	}
	out.Uid = in.Uid
	out.List = info.AllMembers()
	out.Status = outLog(path, out)
	return nil
}

func (mine *SceneService) UpdateDisplay (ctx context.Context, in *pb.ReqSceneDisplay, out *pb.ReplySceneDisplays) error {
	path := "scene.updateDisplay"
	inLog(path, in)
	if len(in.Scene) < 1 {
		out.Status = outError(path,"the parent is empty ", pb.ResultStatus_Empty)
		return nil
	}
	info := cache.Context().GetScene(in.Scene)
	if info == nil {
		out.Status = outError(path,"the scene not found ", pb.ResultStatus_NotExisted)
		return nil
	}
	if in.Slots == nil {
		in.Slots = make([]string, 0, 1)
	}
	err := info.UpdateDisplay(in.Uid, in.Key, in.Skin, in.Operator, in.Slots)
	if err != nil {
		out.Status = outError(path,err.Error(), pb.ResultStatus_DBException)
		return nil
	}
	out.Uid = in.Uid
	out.List = switchExhibitions(info.Exhibitions)
	out.Status = outLog(path, out)
	return nil
}

func (mine *SceneService) PutOnDisplay (ctx context.Context, in *pb.RequestInfo, out *pb.ReplySceneDisplays) error {
	path := "scene.putOnDisplay"
	inLog(path, in)
	if len(in.Parent) < 1 {
		out.Status = outError(path,"the parent is empty ", pb.ResultStatus_Empty)
		return nil
	}
	info := cache.Context().GetScene(in.Parent)
	if info == nil {
		out.Status = outError(path,"the scene not found ", pb.ResultStatus_NotExisted)
		return nil
	}

	err := info.PutOnDisplay(in.Uid)
	if err != nil {
		out.Status = outError(path,err.Error(), pb.ResultStatus_DBException)
		return nil
	}
	out.Uid = in.Uid
	out.List = switchExhibitions(info.Exhibitions)
	out.Status = outLog(path, out)
	return nil
}

func (mine *SceneService) CancelDisplay (ctx context.Context, in *pb.RequestInfo, out *pb.ReplySceneDisplays) error {
	path := "scene.cancelDisplay"
	inLog(path, in)
	if len(in.Parent) < 1 {
		out.Status = outError(path,"the parent is empty ", pb.ResultStatus_Empty)
		return nil
	}
	info := cache.Context().GetScene(in.Parent)
	if info == nil {
		out.Status = outError(path,"the scene not found ", pb.ResultStatus_NotExisted)
		return nil
	}

	err := info.CancelDisplay(in.Uid)
	if err != nil {
		out.Status = outError(path,err.Error(), pb.ResultStatus_DBException)
		return nil
	}
	out.Uid = in.Uid
	out.List = switchExhibitions(info.Exhibitions)
	out.Status = outLog(path, out)
	return nil
}

func (mine *SceneService) UpdateParents (ctx context.Context, in *pb.RequestList, out *pb.ReplyList) error {
	path := "scene.updateChildren"
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

	err := info.UpdateParents(in.Operator, in.List)
	if err != nil {
		out.Status = outError(path,err.Error(), pb.ResultStatus_DBException)
		return nil
	}
	out.Uid = in.Uid
	out.List = info.Parents()
	out.Status = outLog(path, out)
	return nil
}

func (mine *SceneService) AppendDevice (ctx context.Context, in *pb.ReqSceneDevice, out *pb.ReplySceneDevices) error {
	path := "scene.appendDevice"
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

	err := info.AppendDevice(in.Device, in.Remark, in.Type)
	if err != nil {
		out.Status = outError(path,err.Error(), pb.ResultStatus_DBException)
		return nil
	}
	out.List = info.Devices()
	out.Status = outLog(path, out)
	return nil
}

func (mine *SceneService) SubtractDevice (ctx context.Context, in *pb.ReqSceneDevice, out *pb.ReplySceneDevices) error {
	path := "scene.subtractDevice"
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

	err := info.SubtractDevice(in.Device)
	if err != nil {
		out.Status = outError(path,err.Error(), pb.ResultStatus_DBException)
		return nil
	}
	out.List = info.Devices()
	out.Status = outLog(path, out)
	return nil
}



