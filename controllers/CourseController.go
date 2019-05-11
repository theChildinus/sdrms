package controllers

import (
	"encoding/json"
	"fmt"
	"strconv"
	"strings"

	"github.com/astaxie/beego/orm"
	"github.com/lhtzbj12/sdrms/enums"

	"github.com/lhtzbj12/sdrms/models"
)

type CourseController struct {
	BaseController
}

// 课程资源
func (c *CourseController) Prepare() {
	c.BaseController.Prepare()
	c.checkAuthor("DataGrid")
	// c.checkLogin()
}

func (c *CourseController) Index() {
	c.Data["showMoreQuery"] = false
	c.Data["activeSidebarUrl"] = c.URLFor(c.controllerName + "." + c.actionName)
	c.setTpl()
	c.LayoutSections = make(map[string]string)
	c.LayoutSections["headcssjs"] = "course/index_headcssjs.html"
	c.LayoutSections["footerjs"] = "course/index_footerjs.html"

	c.Data["canEdit"] = c.checkActionAuthor("CourseController", "Edit")
	c.Data["canDelete"] = c.checkActionAuthor("CourseController", "Delete")
}

// 获取 课程信息
func (c *CourseController) DataGrid() {
	var params models.CourseQueryParam
	json.Unmarshal(c.Ctx.Input.RequestBody, &params)

	data, total := models.CoursePageList(&params)
	result := make(map[string]interface{})
	result["total"] = total
	result["rows"] = data
	c.Data["json"] = result
	c.ServeJSON()
}

// 编辑课程信息
func (c *CourseController) Edit() {
	if c.Ctx.Request.Method == "POST" {
		c.Save()
	}
	Id, _ := c.GetInt(":id", 0)
	m := &models.Course{}
	var err error
	//  对已有的课程进行修改
	if Id > 0 {
		m, err = models.CourseOne(Id)
		if err != nil {
			c.pageError("数据无效，请刷新后重试")
		}
		//  TODO: o := orm.NewOrm()
		// TODO
	}
	c.Data["m"] = m
	fmt.Println("cData: ", c.Data["m"])
	c.setTpl("course/edit.html", "shared/layout_pullbox.html")
	c.LayoutSections = make(map[string]string)
	c.LayoutSections["footerjs"] = "course/edit_footerjs.html"
}

// 保存修改
func (c *CourseController) Save() {
	m := models.Course{}
	o := orm.NewOrm()
	var err error
	// 获取表单数据
	fmt.Println(m)
	if err = c.ParseForm(&m); err != nil {
		c.jsonResult(enums.JRCodeFailed, "获取数据失败", m.Id)
	}
	if m.Id == 0 {
		// 添加新数据
		if _, err := o.Insert(&m); err != nil {
			c.jsonResult(enums.JRCodeFailed, "获取数据失败", m.Id)
		}
	} else {
		// 更新旧数据
		if _, err := o.Update(&m); err != nil {
			c.jsonResult(enums.JRCodeFailed, "编辑失败", m.Id)
		}
	}
	c.jsonResult(enums.JRCodeSucc, "保存成功", m.Id)
}

// 删除课程信息
func (c *CourseController) Delete() {
	strs := c.GetString("ids")
	ids := make([]int, 0, len(strs))
	for _, str := range strings.Split(strs, ",") {
		if id, err := strconv.Atoi(str); err == nil {
			ids = append(ids, id)
		}
	}
	query := orm.NewOrm().QueryTable(models.CourseTBName())
	if num, err := query.Filter("id__in", ids).Delete(); err == nil {
		c.jsonResult(enums.JRCodeSucc, fmt.Sprintf("成功删除 %d 项", num), 0)
	} else {
		c.jsonResult(enums.JRCodeFailed, "删除失败", 0)
	}
}
