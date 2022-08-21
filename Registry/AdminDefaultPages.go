/**
 * @Author: 云起时
 * @Email: limingxiang@yqstech.com
 * @Description:
 * @File: AdminPages
 * @Version: 1.0.0
 * @Date: 2022/4/28 10:13
 */

package Registry

import (
	"github.com/yqstech/gef/EasyApp"
	adminHandle2 "github.com/yqstech/gef/Handles/adminHandle"
)

var AdminPages = map[string]EasyApp.AppPage{
	"index":            adminHandle2.Index{},
	"admin":            adminHandle2.Admin{},
	"account":          adminHandle2.Account{},
	"admin_rules":      adminHandle2.AdminRules{},
	"admin_group":      adminHandle2.AdminGroup{},
	"admin_log":        adminHandle2.AdminLog{},
	"configs":          adminHandle2.Configs{},         //设置项管理
	"configs_group":    adminHandle2.ConfigsGroup{},    //设置分组
	"attachment_image": adminHandle2.AttachmentImage{}, //图片附件
	"attachment_file":  adminHandle2.AttachmentFile{},  //文件附件
	"option_models":    adminHandle2.OptionModels{},    //选项集
	//!------------ 应用配置 -------------------------
	"app_configs_g1": adminHandle2.AppConfigs{GroupId: 1}, //设置分组管理页
	"app_configs_g2": adminHandle2.AppConfigs{GroupId: 2}, //设置分组管理页
	"app_configs_g3": adminHandle2.AppConfigs{GroupId: 3}, //设置分组管理页
	"app_configs_g4": adminHandle2.AppConfigs{GroupId: 4}, //设置分组管理页
	"app_configs_g5": adminHandle2.AppConfigs{GroupId: 5}, //设置分组管理页
	"app_configs_g6": adminHandle2.AppConfigs{GroupId: 6}, //设置分组管理页
	"app_configs_g7": adminHandle2.AppConfigs{GroupId: 7}, //设置分组管理页
	"app_configs_g8": adminHandle2.AppConfigs{GroupId: 8}, //设置分组管理页
	"app_configs_g9": adminHandle2.AppConfigs{GroupId: 9}, //设置分组管理页
	//!------------ EasyModel 简单模型 -------------------------
	"easy_models":             adminHandle2.EasyModels{},           //模型管理
	"easy_models_fields":      adminHandle2.EasyModelsFields{},     //模型字段管理
	"easy_models_buttons":     adminHandle2.EasyModelsButtons{},    //模型按钮管理
	"easy_curd_models":        adminHandle2.EasyCurdModels{},       //接口模型 easyCurd模型管理
	"easy_curd_models_fields": adminHandle2.EasyCurdModelsFields{}, //接口模型字段管理 easyCurd模型字段管理
}
