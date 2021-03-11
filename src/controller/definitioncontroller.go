/**
 * @Author: lzw5399
 * @Date: 2021/01/15
 * @Desc: process控制器
 */
package controller

import (
	"log"
	"net/http"

	"workflow/src/global/response"
	"workflow/src/model/request"
	"workflow/src/service"
	"workflow/src/util"

	"github.com/labstack/echo/v4"
)

var (
	definitionService service.DefinitionService = service.NewDefinitionService()
)

// 创建流程process
func CreateProcessDefinition(c echo.Context) error {
	var (
		r   request.ProcessDefinitionRequest
		err error
	)

	if err = c.Bind(&r); err != nil {
		return response.BadRequest(c)
	}

	// 验证
	err = definitionService.Validate(&r, 0)
	if err != nil {
		return response.BadRequestWithMessage(c, err)
	}

	// 创建
	processDefinition, err := definitionService.CreateDefinition(&r)
	if err != nil {
		log.Printf("CreateProcess错误，原因: %s", err.Error())
		return response.Failed(c, http.StatusInternalServerError)
	}

	return response.OkWithData(c, processDefinition)
}

// 更新模板
func UpdateProcessDefinition(c echo.Context) error {
	var (
		r   request.ProcessDefinitionRequest
		err error
	)

	if err = c.Bind(&r); err != nil {
		return response.BadRequest(c)
	}

	// 验证
	err = definitionService.Validate(&r, r.Id)
	if err != nil {
		return response.BadRequestWithMessage(c, err)
	}

	err = definitionService.UpdateDefinition(&r)
	if err != nil {
		log.Printf("UpdateProcessDefinition错误，原因: %s", err.Error())
		return response.Failed(c, http.StatusInternalServerError)
	}

	return response.Ok(c)
}

// 删除流程
func DeleteProcessDefinition(c echo.Context) error {
	definitionId := c.Param("id")
	if definitionId == "" {
		return response.BadRequestWithMessage(c, "参数不正确，请确定参数processDefinitionId是否传递")
	}

	err := definitionService.DeleteDefinition(util.StringToUint(definitionId))
	if err != nil {
		return response.Failed(c, http.StatusInternalServerError)
	}

	return response.Ok(c)
}

// 流程详情
func GetProcessDefinition(c echo.Context) error {
	definitionId := c.Param("id")
	if definitionId == "" {
		return response.BadRequestWithMessage(c, "参数不正确，请确定参数processDefinitionId是否传递")
	}

	definition, err := definitionService.GetDefinition(util.StringToUint(definitionId))
	if err != nil {
		return response.Failed(c, http.StatusNotFound)
	}

	return response.OkWithData(c, definition)
}

//// 分类流程列表
//func ClassifyProcessList(c echo.Context) error {
//	var (
//		err            error
//		classifyIdList []int
//		classifyList   []*struct {
//			process2.Classify
//			ProcessList []*process2.Info `json:"process_list"`
//		}
//	)
//
//	processName := c.DefaultQuery("name", "")
//	if processName == "" {
//		err = orm.Eloquent.Model(&process2.Classify{}).Find(&classifyList).Error
//		if err != nil {
//			app.Error(c, -1, err, fmt.Sprintf("获取分类列表失败，%v", err.Error()))
//			return
//		}
//	} else {
//		err = orm.Eloquent.Model(&process2.Info{}).
//			Where("name LIKE ?", fmt.Sprintf("%%%v%%", processName)).
//			Pluck("distinct classify", &classifyIdList).Error
//		if err != nil {
//			app.Error(c, -1, err, fmt.Sprintf("获取分类失败，%v", err.Error()))
//			return
//		}
//
//		err = orm.Eloquent.Model(&process2.Classify{}).
//			Where("id in (?)", classifyIdList).
//			Find(&classifyList).Error
//		if err != nil {
//			app.Error(c, -1, err, fmt.Sprintf("获取分类失败，%v", err.Error()))
//			return
//		}
//	}
//
//	for _, item := range classifyList {
//		err = orm.Eloquent.Model(&process2.Info{}).
//			Where("classify = ? and name LIKE ?", item.Id, fmt.Sprintf("%%%v%%", processName)).
//			Select("id, create_time, update_time, name, icon, remarks").
//			Find(&item.ProcessList).Error
//		if err != nil {
//			app.Error(c, -1, err, fmt.Sprintf("获取流程失败，%v", err.Error()))
//			return
//		}
//	}
//
//	app.OK(c, classifyList, "成功获取数据")
//}