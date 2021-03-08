package core

import (
	"encoding/json"
	"fmt"
	"time"
)

/**
@author zhangchengji
*/
type SysLog struct {
	Id         uint       `json:"id" form:"id" gorm:"column:id;comment:id;"`
	Title      string     `json:"title" form:"title" gorm:"column:title;comment:日志标题;type:varchar(255);size:255;"`
	CreateBy   string     `json:"createBy" form:"createBy" gorm:"column:create_by;comment:日志标题;type:varchar(255);size:255;"`
	CreateTime *time.Time `json:"createTime" form:"createTime" gorm:"column:create_time;comment:创建时间;default:null;"`
	RemoteAddr string     `json:"remoteAddr" form:"remoteAddr" gorm:"column:remote_addr;comment:操作ip地址;type:varchar(30);size:30;"`
	UserAgent  string     `json:"userAgent" form:"userAgent" gorm:"column:user_agent;comment:用户代理;type:varchar(255);size:255;"`
	RequestUri string     `json:"requestUri" form:"requestUri" gorm:"column:request_uri;comment:请求uri;type:varchar(60);size:60;"`
	Method     string     `json:"method" form:"method" gorm:"column:method;comment:操作方式;type:varchar(10);size:10;"`
	MethodName string     `json:"methodName" form:"methodName" gorm:"column:method_name;comment:操作方法;type:varchar(255);size:255;"`
	ClassName  string     `json:"className" form:"className" gorm:"column:class_name;comment:操作类;type:varchar(255);size:255;"`
	Params     string     `json:"params" form:"params" gorm:"column:params;comment:数据;type:varchar(255);size:255;"`
	Time       uint64     `json:"time" form:"time" gorm:"column:params;comment:方法执行时间;type:bigint"`
	ServiceId  string     `json:"serviceId" form:"serviceId" gorm:"column:service_id;comment:应用标识;type:varchar(50);size:50;"`
}

func (SysLog) TableName() string {
	return "sys_log"
}
func (sysLog *SysLog) Save() (err error) {
	err = DB.Save(sysLog).Error
	return
}

var (
	logData chan *SysLog
)

func Init() {
	logData = make(chan *SysLog, 10000)
	go saveLog()

}
func SendChanLog(data string) {
	var log SysLog
	if err := json.Unmarshal([]byte(data), &log); err != nil {
		fmt.Println("json  Unmarshal syslog failed,err", err)
		return
	}
	logData <- &log
}
func saveLog() {
	for {
		select {
		case log := <-logData:
			log.Save()
		default:
			time.Sleep(time.Millisecond * 100)
		}
	}
}
