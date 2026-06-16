package commerce_cdk

import (
	"github.com/gin-gonic/gin"
	"server/pkg/bean"
	"server/service/admin/commerce_cdk"
)

var svr = service.Service

func Create(c *gin.Context) {
	var req service.CreateReq
	if err := c.ShouldBindJSON(&req); err != nil {
		bean.Response.Param(c, err)
		return
	}
	codes, err := svr.Create(req)
	if err != nil {
		bean.Response.Fail(c, err.Error())
		return
	}
	bean.Response.OkData(c, codes)
}

func Page(c *gin.Context) {
	var req service.PageReq
	if err := c.ShouldBindJSON(&req); err != nil {
		bean.Response.Param(c, err)
		return
	}
	list, total := svr.Page(req)
	bean.Response.OkData(c, bean.NewPage(list, total))
}

func Disable(c *gin.Context) {
	var req service.DisableReq
	if err := c.ShouldBindJSON(&req); err != nil {
		bean.Response.Param(c, err)
		return
	}
	if err := svr.Disable(req); err != nil {
		bean.Response.Fail(c, err.Error())
		return
	}
	bean.Response.Ok(c)
}

func BatchDisable(c *gin.Context) {
	var req service.BatchCodesReq
	if err := c.ShouldBindJSON(&req); err != nil {
		bean.Response.Param(c, err)
		return
	}
	if err := svr.BatchDisable(req); err != nil {
		bean.Response.Fail(c, err.Error())
		return
	}
	bean.Response.Ok(c)
}

func BatchEnable(c *gin.Context) {
	var req service.BatchCodesReq
	if err := c.ShouldBindJSON(&req); err != nil {
		bean.Response.Param(c, err)
		return
	}
	if err := svr.BatchEnable(req); err != nil {
		bean.Response.Fail(c, err.Error())
		return
	}
	bean.Response.Ok(c)
}

func BatchDelete(c *gin.Context) {
	var req service.BatchCodesReq
	if err := c.ShouldBindJSON(&req); err != nil {
		bean.Response.Param(c, err)
		return
	}
	if err := svr.BatchDelete(req); err != nil {
		bean.Response.Fail(c, err.Error())
		return
	}
	bean.Response.Ok(c)
}

func UpdateRemark(c *gin.Context) {
	var req service.UpdateRemarkReq
	if err := c.ShouldBindJSON(&req); err != nil {
		bean.Response.Param(c, err)
		return
	}
	if err := svr.UpdateRemark(req); err != nil {
		bean.Response.Fail(c, err.Error())
		return
	}
	bean.Response.Ok(c)
}
