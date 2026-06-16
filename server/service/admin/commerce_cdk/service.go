package service

import (
	"errors"

	"github.com/shopspring/decimal"

	"server/model"
	"server/pkg/utils"
	"server/repository"
	"server/service/common/commerce"
)

type service struct{}

var Service *service

type CreateReq struct {
	Type      string          `binding:"required" json:"type"`
	Value     decimal.Decimal `binding:"required" json:"value"`
	Count     int             `json:"count"`
	ExpiredAt int64           `json:"expiredAt"`
	Remark    string          `json:"remark"`
}

type PageReq struct {
	Page    int    `json:"page"`
	Size    int    `json:"size"`
	Status  int    `json:"status"`
	BatchNo string `json:"batchNo"`
}

type DisableReq struct {
	Code string `binding:"required" json:"code"`
}

type BatchCodesReq struct {
	Codes []string `binding:"required" json:"codes"`
}

type UpdateRemarkReq struct {
	Codes  []string `binding:"required" json:"codes"`
	Remark string   `json:"remark"`
}

type CreateResp struct {
	Codes   []string `json:"codes"`
	BatchNo string   `json:"batchNo"`
}

func (s *service) Create(req CreateReq) (*CreateResp, error) {
	if req.Count <= 0 {
		req.Count = 1
	}
	if req.Count > 500 {
		return nil, errors.New("单次最多生成500个")
	}
	if req.Type != model.CDK_TYPE_BALANCE && req.Type != model.CDK_TYPE_POINTS {
		return nil, errors.New("类型错误")
	}
	if req.Value.LessThanOrEqual(decimal.Zero) {
		return nil, errors.New("面值必须大于0")
	}
	db, _, _ := repository.Get("")
	batch := commerce.NewOrderNo("CDKB")
	codes := make([]string, 0, req.Count)
	for i := 0; i < req.Count; i++ {
		code := utils.RandStrPrefix(18, "CDK", utils.AllDict)
		cdk := model.CommerceCdk{Code: code, Type: req.Type, Value: req.Value, Status: model.CDK_STATUS_UNUSED, BatchNo: batch, ExpiredAt: req.ExpiredAt, Remark: req.Remark}
		if err := commerce.DB(db).Create(&cdk).Error; err != nil {
			return nil, err
		}
		codes = append(codes, code)
	}
	return &CreateResp{Codes: codes, BatchNo: batch}, nil
}

func (s *service) Page(req PageReq) (any, int64) {
	db, _, _ := repository.Get("")
	if req.Page <= 0 {
		req.Page = 1
	}
	if req.Size <= 0 {
		req.Size = 10
	}
	if req.Size > 100 {
		req.Size = 100
	}
	var list []map[string]any
	var total int64
	q := commerce.DB(db).Table("commerce_cdks")
	if req.Status > 0 {
		q = q.Where("status = ?", req.Status)
	}
	if req.BatchNo != "" {
		q = q.Where("batch_no = ?", req.BatchNo)
	}
	q.Count(&total).Order("id desc").Offset((req.Page - 1) * req.Size).Limit(req.Size).Find(&list)
	return list, total
}

func (s *service) Disable(req DisableReq) error {
	db, _, _ := repository.Get("")
	return commerce.DB(db).Model(&model.CommerceCdk{}).Where("cdk_code = ? AND status = ?", req.Code, model.CDK_STATUS_UNUSED).Update("status", model.CDK_STATUS_DISABLED).Error
}

func (s *service) BatchDisable(req BatchCodesReq) error {
	db, _, _ := repository.Get("")
	res := commerce.DB(db).Model(&model.CommerceCdk{}).Where("cdk_code IN ? AND status = ?", req.Codes, model.CDK_STATUS_UNUSED).Update("status", model.CDK_STATUS_DISABLED)
	return res.Error
}

func (s *service) BatchEnable(req BatchCodesReq) error {
	db, _, _ := repository.Get("")
	res := commerce.DB(db).Model(&model.CommerceCdk{}).Where("cdk_code IN ? AND status = ?", req.Codes, model.CDK_STATUS_DISABLED).Update("status", model.CDK_STATUS_UNUSED)
	return res.Error
}

func (s *service) BatchDelete(req BatchCodesReq) error {
	db, _, _ := repository.Get("")
	res := commerce.DB(db).Where("cdk_code IN ? AND status != ?", req.Codes, model.CDK_STATUS_USED).Delete(&model.CommerceCdk{})
	return res.Error
}

func (s *service) UpdateRemark(req UpdateRemarkReq) error {
	db, _, _ := repository.Get("")
	res := commerce.DB(db).Model(&model.CommerceCdk{}).Where("cdk_code IN ?", req.Codes).Update("remark", req.Remark)
	return res.Error
}
