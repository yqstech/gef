/**
 * @Author: 云起时
 * @Email: limingxiang@yqstech.com
 * @Description:
 * @File: AdminPages
 * @Version: 1.0.0
 * @Date: 2022/4/28 10:13
 */

package registry

import (
	"github.com/yqstech/gef/Handles/adminHandle"
	"github.com/yqstech/gef/builder"
)

var AdminPages = map[string]builder.NodePager{
	"index":            &adminHandle.Index{},
	"admin":            &adminHandle.Admin{},
	"account":          &adminHandle.Account{},
	"admin_rules":      &adminHandle.AdminRules{},
	"admin_group":      &adminHandle.AdminGroup{},
	"admin_log":        &adminHandle.AdminLog{},
	"configs":          &adminHandle.Configs{},         //设置项管理
	"configs_group":    &adminHandle.ConfigsGroup{},    //设置分组
	"attachment_image": &adminHandle.AttachmentImage{}, //图片附件
	"attachment_file":  &adminHandle.AttachmentFile{},  //文件附件
	"option_models":    &adminHandle.OptionModels{},    //选项集
	//!------------ 应用配置 -------------------------
	"app_configs_g1": &adminHandle.AppConfigs{GroupId: 1}, //设置分组管理页
	"app_configs_g2": &adminHandle.AppConfigs{GroupId: 2}, //设置分组管理页
	"app_configs_g3": &adminHandle.AppConfigs{GroupId: 3}, //设置分组管理页
	"app_configs_g4": &adminHandle.AppConfigs{GroupId: 4}, //设置分组管理页
	"app_configs_g5": &adminHandle.AppConfigs{GroupId: 5}, //设置分组管理页
	"app_configs_g6": &adminHandle.AppConfigs{GroupId: 6}, //设置分组管理页
	"app_configs_g7": &adminHandle.AppConfigs{GroupId: 7}, //设置分组管理页
	"app_configs_g8": &adminHandle.AppConfigs{GroupId: 8}, //设置分组管理页
	"app_configs_g9": &adminHandle.AppConfigs{GroupId: 9}, //设置分组管理页
	//!------------ EasyModelHandle 简单模型 -------------------------
	"easy_models":             &adminHandle.EasyModels{},           //模型管理
	"easy_models_fields":      &adminHandle.EasyModelsFields{},     //模型字段管理
	"easy_models_buttons":     &adminHandle.EasyModelsButtons{},    //模型按钮管理
	"easy_curd_models":        &adminHandle.EasyCurdModels{},       //接口模型 easyCurd模型管理
	"easy_curd_models_fields": &adminHandle.EasyCurdModelsFields{}, //接口模型字段管理 easyCurd模型字段管理
}
