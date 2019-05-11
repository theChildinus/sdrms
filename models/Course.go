package models

import (
	"github.com/astaxie/beego/orm"
)

type CourseQueryParam struct {
	BaseQueryParam
	CourseNameLike string
}

// Course 课程 实体类
type Course struct {
	Id         int
	CourseName string `orm:size(24)`
}

// 设置表名
func (a *Course) TableName() string {
	return CourseTBName()
}

func CoursePageList(params *CourseQueryParam) ([]*Course, int64) {
	query := orm.NewOrm().QueryTable(CourseTBName())
	data := make([]*Course, 0)

	// TODO: sorteder
	query = query.Filter("coursename__istartswith", params.CourseNameLike)
	total, _ := query.Count()
	query.Limit(params.Limit, params.Offset).All(&data)
	return data, total
}

func CourseDataList(params *CourseQueryParam) []*Course {
	params.Limit = -1
	data, _ := CoursePageList(params)
	return data
}

// 取单条数据
func CourseOne(id int) (*Course, error) {
	o := orm.NewOrm()
	m := Course{Id: id}
	err := o.Read(&m)
	if err != nil {
		return nil, err
	}
	return &m, err
}

// 根据课程名取数据
func CourseOneByCourseName(courseName string) (*Course, error) {
	o := orm.NewOrm()
	m := Course{}
	err := o.QueryTable(CourseTBName()).Filter("coursename", courseName).One(&m)
	if err != nil {
		return nil, err
	}
	return &m, err
}
