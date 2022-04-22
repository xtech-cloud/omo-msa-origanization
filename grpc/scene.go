package grpc

import (
	"context"
	pb "github.com/xtech-cloud/omo-msp-organization/proto/organization"
	pbstatus "github.com/xtech-cloud/omo-msp-status/proto/status"
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
	tmp.Bucket = info.Bucket
	tmp.Short = info.ShortName
	tmp.Parents = info.Parents()
	tmp.Members = info.AllMembers()
	tmp.Domains = make([]*pb.ProductInfo, 0, len(info.Domains))
	for _, domain := range info.Domains {
		tmp.Domains = append(tmp.Domains, &pb.ProductInfo{Type: uint32(domain.Type), Uid: domain.UID,
			Remark: domain.Remark, Keywords: domain.Keywords, Name: domain.Name})
	}
	return tmp
}

func (mine *SceneService)AddOne(ctx context.Context, in *pb.ReqSceneAdd, out *pb.ReplySceneOne) error {
	path := "scene.add"
	inLog(path, in)
	if len(in.Name) < 1 {
		out.Status = outError(path,"the name is empty ", pbstatus.ResultStatus_Empty)
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
	info.ShortName = ""
	err := cache.Context().CreateScene(info)
	if err != nil {
		out.Status = outError(path,err.Error(), pbstatus.ResultStatus_DBException)
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
		out.Status = outError(path,"the uid is empty ", pbstatus.ResultStatus_Empty)
		return nil
	}
	info := cache.Context().GetScene(in.Uid)
	if info == nil {
		out.Status = outError(path,"the scene not found ", pbstatus.ResultStatus_NotExisted)
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
		out.Status = outError(path,"the uid is empty ", pbstatus.ResultStatus_Empty)
		return nil
	}
	info := cache.Context().GetSceneByMember(in.Uid)
	if info == nil {
		out.Status = outError(path,"the scene not found ", pbstatus.ResultStatus_NotExisted)
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
		out.Status = outError(path,"the uid is empty ", pbstatus.ResultStatus_Empty)
		return nil
	}
	err := cache.RemoveScene(in.Uid, in.Operator)
	if err != nil {
		out.Status = outError(path,err.Error(), pbstatus.ResultStatus_DBException)
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
		out.Status = outError(path,"the uid is empty ", pbstatus.ResultStatus_Empty)
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
		total,max,list = cache.Context().GetScenes(in.Page, in.Number)
	}else{
		total,max,list = cache.Context().GetScenesByParent(in.Parent, in.Page, in.Number)
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

func (mine *SceneService)GetByFilter(ctx context.Context, in *pb.RequestFilter, out *pb.ReplySceneList) error {
	path := "scene.getByFilter"
	inLog(path, in)
	var total uint32 = 0
	var max uint32 = 0
	var list []*cache.SceneInfo
	if in.Key == "shortname" {
		list = make([]*cache.SceneInfo, 0 ,1)
	}else if in.Key == "type" {
		tp,er := strconv.ParseUint(in.Scene, 10, 32)
		if er != nil {
			out.Status = outError(path, er.Error(), pbstatus.ResultStatus_DBException)
			return nil
		}
		list = cache.Context().GetScenesByType(uint8(tp))
	}else if in.Key == "parent" {
		 total, max, list = cache.Context().GetScenesByParent(in.Value, in.Page, in.Number)
	}else if in.Key == "array" {
		list = cache.Context().GetScenesByArray(in.List)
	}else{
		list = make([]*cache.SceneInfo, 0 ,1)
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

func (mine *SceneService)GetStatistic(ctx context.Context, in *pb.RequestFilter, out *pb.ReplyStatistic) error {
	path := "scene.getStatistic"
	inLog(path, in)

	out.Status = outLog(path, out)
	return nil
}

func (mine *SceneService) UpdateBase (ctx context.Context, in *pb.ReqSceneUpdate, out *pb.ReplyInfo) error {
	path := "scene.updateBase"
	inLog(path, in)
	if len(in.Uid) < 1 {
		out.Status = outError(path,"the uid is empty ", pbstatus.ResultStatus_Empty)
		return nil
	}
	info := cache.Context().GetScene(in.Uid)
	if info == nil {
		out.Status = outError(path,"the scene not found ", pbstatus.ResultStatus_NotExisted)
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
		out.Status = outError(path,err.Error(), pbstatus.ResultStatus_DBException)
		return nil
	}

	out.Status = outLog(path, out)
	return err
}

func (mine *SceneService) UpdateAddress (ctx context.Context, in *pb.RequestAddress, out *pb.ReplySceneOne) error {
	path := "scene.updateBase"
	inLog(path, in)
	if len(in.Uid) < 1 {
		out.Status = outError(path,"the uid is empty ", pbstatus.ResultStatus_Empty)
		return nil
	}
	info := cache.Context().GetScene(in.Uid)
	if info == nil {
		out.Status = outError(path,"the scene not found ", pbstatus.ResultStatus_NotExisted)
		return nil
	}
	err := info.UpdateAddress(in.Country, in.Province, in.City, in.Zone, in.Operator)
	if err != nil {
		out.Status = outError(path,err.Error(), pbstatus.ResultStatus_DBException)
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
		out.Status = outError(path,"the uid is empty ", pbstatus.ResultStatus_Empty)
		return nil
	}
	info := cache.Context().GetScene(in.Uid)
	if info == nil {
		out.Status = outError(path,"the scene not found ", pbstatus.ResultStatus_NotExisted)
		return nil
	}
	err := info.UpdateLocation(in.Flag, in.Operator)
	if err != nil {
		out.Status = outError(path,err.Error(), pbstatus.ResultStatus_DBException)
		return nil
	}

	out.Status = outLog(path, out)
	return err
}

func (mine *SceneService) UpdateStatus (ctx context.Context, in *pb.ReqSceneStatus, out *pb.ReplyInfo) error {
	path := "scene.updateStatus"
	inLog(path, in)
	if len(in.Uid) < 1 {
		out.Status = outError(path,"the uid is empty ", pbstatus.ResultStatus_Empty)
		return nil
	}
	info := cache.Context().GetScene(in.Uid)
	if info == nil {
		out.Status = outError(path,"the scene not found ", pbstatus.ResultStatus_NotExisted)
		return nil
	}

	err := info.UpdateStatus(cache.SceneStatus(in.Status),in.Operator)
	if err != nil {
		out.Status = outError(path,err.Error(), pbstatus.ResultStatus_DBException)
		return nil
	}
	out.Status = outLog(path, out)
	return nil
}

func (mine *SceneService) UpdateSupporter (ctx context.Context, in *pb.RequestFlag, out *pb.ReplyInfo) error {
	path := "scene.updateSupporter"
	inLog(path, in)
	if len(in.Flag) < 1 {
		out.Status = outError(path,"the supporter uid is empty ", pbstatus.ResultStatus_Empty)
		return nil
	}
	if len(in.Uid) < 1 {
		out.Status = outError(path,"the scene uid is empty ", pbstatus.ResultStatus_Empty)
		return nil
	}
	info := cache.Context().GetScene(in.Uid)
	if info == nil {
		out.Status = outError(path,"the scene not found ", pbstatus.ResultStatus_NotExisted)
		return nil
	}

	err := info.UpdateSupporter(in.Flag,in.Operator)
	if err != nil {
		out.Status = outError(path,err.Error(), pbstatus.ResultStatus_DBException)
		return nil
	}
	out.Status = outLog(path, out)
	return nil
}

func (mine *SceneService) UpdateDomains (ctx context.Context, in *pb.ReqSceneDomains, out *pb.ReplyInfo) error {
	path := "scene.updateDomains"
	inLog(path, in)

	if len(in.Uid) < 1 {
		out.Status = outError(path,"the scene uid is empty ", pbstatus.ResultStatus_Empty)
		return nil
	}
	info := cache.Context().GetScene(in.Uid)
	if info == nil {
		out.Status = outError(path,"the scene not found ", pbstatus.ResultStatus_NotExisted)
		return nil
	}
	arr := make([]proxy.DomainInfo, 0, len(in.List))
	for _, item := range in.List {
		arr = append(arr, proxy.DomainInfo{Type: uint8(item.Type), UID: item.Uid, Remark: item.Remark, Keywords: item.Keywords, Name: item.Name})
	}
	err := info.UpdateDomains(in.Operator, arr)
	if err != nil {
		out.Status = outError(path,err.Error(), pbstatus.ResultStatus_DBException)
		return nil
	}
	out.Status = outLog(path, out)
	return nil
}

func (mine *SceneService) UpdateBucket (ctx context.Context, in *pb.RequestFlag, out *pb.ReplyInfo) error {
	path := "scene.updateBucket"
	inLog(path, in)
	if len(in.Flag) < 1 {
		out.Status = outError(path,"the domain uid is empty ", pbstatus.ResultStatus_Empty)
		return nil
	}
	if len(in.Uid) < 1 {
		out.Status = outError(path,"the scene uid is empty ", pbstatus.ResultStatus_Empty)
		return nil
	}
	info := cache.Context().GetScene(in.Uid)
	if info == nil {
		out.Status = outError(path,"the scene not found ", pbstatus.ResultStatus_NotExisted)
		return nil
	}

	err := info.UpdateBucket(in.Flag, in.Operator)
	if err != nil {
		out.Status = outError(path,err.Error(), pbstatus.ResultStatus_DBException)
		return nil
	}
	out.Status = outLog(path, out)
	return nil
}

func (mine *SceneService) UpdateShortName (ctx context.Context, in *pb.RequestFlag, out *pb.ReplyInfo) error {
	path := "scene.updateShortName"
	inLog(path, in)
	if len(in.Uid) < 1 {
		out.Status = outError(path,"the short name is empty ", pbstatus.ResultStatus_Empty)
		return nil
	}

	info := cache.Context().GetScene(in.Uid)
	if info == nil {
		out.Status = outError(path,"the scene not found ", pbstatus.ResultStatus_NotExisted)
		return nil
	}

	err := info.UpdateShortName(in.Flag, in.Operator)
	if err != nil {
		out.Status = outError(path,err.Error(), pbstatus.ResultStatus_DBException)
		return nil
	}
	out.Status = outLog(path, out)
	return nil
}

func (mine *SceneService) AppendMember (ctx context.Context, in *pb.RequestInfo, out *pb.ReplyList) error {
	path := "scene.appendMember"
	inLog(path, in)
	if len(in.Parent) < 1 {
		out.Status = outError(path,"the parent is empty ", pbstatus.ResultStatus_Empty)
		return nil
	}
	info := cache.Context().GetScene(in.Parent)
	if info == nil {
		out.Status = outError(path,"the scene not found ", pbstatus.ResultStatus_NotExisted)
		return nil
	}

	err := info.AppendMember(in.Uid)
	if err != nil {
		out.Status = outError(path,err.Error(), pbstatus.ResultStatus_DBException)
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
		out.Status = outError(path,"the parent is empty ", pbstatus.ResultStatus_Empty)
		return nil
	}
	info := cache.Context().GetScene(in.Parent)
	if info == nil {
		out.Status = outError(path,"the scene not found ", pbstatus.ResultStatus_NotExisted)
		return nil
	}

	err := info.SubtractMember(in.Uid)
	if err != nil {
		out.Status = outError(path,err.Error(), pbstatus.ResultStatus_DBException)
		return nil
	}
	out.Uid = in.Uid
	out.List = info.AllMembers()
	out.Status = outLog(path, out)
	return nil
}

//func (mine *SceneService) UpdateDisplay (ctx context.Context, in *pb.ReqSceneDisplay, out *pb.ReplySceneDisplays) error {
//	path := "scene.updateDisplay"
//	inLog(path, in)
//	if len(in.Scene) < 1 {
//		out.Status = outError(path,"the parent is empty ", pbstatus.ResultStatus_Empty)
//		return nil
//	}
//	info := cache.Context().GetScene(in.Scene)
//	if info == nil {
//		out.Status = outError(path,"the scene not found ", pbstatus.ResultStatus_NotExisted)
//		return nil
//	}
//	if in.Slots == nil {
//		in.Slots = make([]string, 0, 1)
//	}
//	err := info.UpdateDisplay(in.Uid, in.Key, in.Skin, in.Operator, in.Slots)
//	if err != nil {
//		out.Status = outError(path,err.Error(), pbstatus.ResultStatus_DBException)
//		return nil
//	}
//	out.Uid = in.Uid
//	out.List = switchExhibitions(info.Exhibitions)
//	out.Status = outLog(path, out)
//	return nil
//}

//func (mine *SceneService) PutOnDisplay (ctx context.Context, in *pb.RequestInfo, out *pb.ReplySceneDisplays) error {
//	path := "scene.putOnDisplay"
//	inLog(path, in)
//	if len(in.Parent) < 1 {
//		out.Status = outError(path,"the parent is empty ", pbstatus.ResultStatus_Empty)
//		return nil
//	}
//	info := cache.Context().GetScene(in.Parent)
//	if info == nil {
//		out.Status = outError(path,"the scene not found ", pbstatus.ResultStatus_NotExisted)
//		return nil
//	}
//
//	err := info.PutOnDisplay(in.Uid)
//	if err != nil {
//		out.Status = outError(path,err.Error(), pbstatus.ResultStatus_DBException)
//		return nil
//	}
//	out.Uid = in.Uid
//	out.List = switchExhibitions(info.Exhibitions)
//	out.Status = outLog(path, out)
//	return nil
//}

//func (mine *SceneService) CancelDisplay (ctx context.Context, in *pb.RequestInfo, out *pb.ReplySceneDisplays) error {
//	path := "scene.cancelDisplay"
//	inLog(path, in)
//	if len(in.Parent) < 1 {
//		out.Status = outError(path,"the parent is empty ", pbstatus.ResultStatus_Empty)
//		return nil
//	}
//	info := cache.Context().GetScene(in.Parent)
//	if info == nil {
//		out.Status = outError(path,"the scene not found ", pbstatus.ResultStatus_NotExisted)
//		return nil
//	}
//
//	err := info.CancelDisplay(in.Uid)
//	if err != nil {
//		out.Status = outError(path,err.Error(), pbstatus.ResultStatus_DBException)
//		return nil
//	}
//	out.Uid = in.Uid
//	out.List = switchExhibitions(info.Exhibitions)
//	out.Status = outLog(path, out)
//	return nil
//}

func (mine *SceneService) UpdateParents (ctx context.Context, in *pb.RequestList, out *pb.ReplyList) error {
	path := "scene.updateParents"
	inLog(path, in)
	if len(in.Uid) < 1 {
		out.Status = outError(path,"the uid is empty ", pbstatus.ResultStatus_Empty)
		return nil
	}
	info := cache.Context().GetScene(in.Uid)
	if info == nil {
		out.Status = outError(path,"the scene not found ", pbstatus.ResultStatus_NotExisted)
		return nil
	}

	err := info.UpdateParents(in.Operator, in.List)
	if err != nil {
		out.Status = outError(path,err.Error(), pbstatus.ResultStatus_DBException)
		return nil
	}
	out.Uid = in.Uid
	out.List = info.Parents()
	out.Status = outLog(path, out)
	return nil
}


