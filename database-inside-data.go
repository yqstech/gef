/**
 * @Author: 云起时
 * @Email: limingxiang@yqstech.com
 * @Description:
 * @File: database-inside-data
 * @Version: 1.0.0
 * @Date: 2022/8/16 22:20
 */

package gef

var insideData = []InsideData{
	{TableName: "tb_admin", Condition: [][]interface{}{{"id", "1"}}, Data: map[string]interface{}{
		"group_id": 0,
		"name":     "系统管理员",
		"account":  "root",
		"password": "99b74a36ec072636d319ff7b49917d95",
	}},
	{TableName: "tb_configs_group", Condition: [][]interface{}{{"id", "1"}}, Data: map[string]interface{}{
		"id":         1,
		"group_name": "系统配置",
		"note":       "系统全局配置项",
	}},
	{TableName: "tb_configs", Condition: [][]interface{}{{"name", "upload_extension"}}, Data: map[string]interface{}{
		"group_id":   1,
		"name":       "upload_extension",
		"value":      "jpg,png,gif,jpeg,JPG,PNG,GIF,JPEG,xml,XML",
		"title":      "可用附件类型",
		"notice":     "允许上传的附件拓展名，多个拓展名用英文逗号分割，不用加点",
		"field_type": "text",
		"index_num":  1,
	}},
	//设置项
	{TableName: "tb_configs_group", Condition: [][]interface{}{{"id", "2"}}, Data: map[string]interface{}{
		"id":         2,
		"group_name": "应用配置",
		"note":       "应用配置项",
	}},
	{TableName: "tb_configs", Condition: [][]interface{}{{"name", "app_name"}}, Data: map[string]interface{}{
		"group_id":   2,
		"name":       "app_name",
		"value":      "Gef开发框架",
		"title":      "应用名称",
		"notice":     "多个地方会显示，勿删！",
		"field_type": "text-sm",
		"index_num":  1,
	}},
	{TableName: "tb_app_configs", Condition: [][]interface{}{{"name", "app_name"}}, Data: map[string]interface{}{
		"group_id": 2,
		"name":     "app_name",
		"value":    "Gef开发框架",
	}},
	{TableName: "tb_app_configs", Condition: [][]interface{}{{"name", "upload_extension"}}, Data: map[string]interface{}{
		"group_id": 1,
		"name":     "upload_extension",
		"value":    "jpg,png,gif,jpeg,JPG,PNG,GIF,JPEG,xml,XML",
	}},
	//选项集
	{TableName: "tb_option_models", Condition: [][]interface{}{{"unique_key", "is"}},
		Data: map[string]interface{}{"children_option_model_key":"","color_array":"","color_field":"","data_type":0,"default_data":"","dynamic_params":"","icon_array":"","icon_field":"","index_num":1,"match_fields":"is_*\nallow_*","name":"是|否","name_field":"","options_disable":0,"parent_field":"","select_order":"","select_where":"","static_data":"[{\"name\":\"是\",\"value\":\"1\",\"color\":\"#66CC00\",\"icon\":\"ri-checkbox-circle-fill\"},{\"name\":\"否\",\"value\":\"0\",\"color\":\"#FF6666\",\"icon\":\"ri-forbid-fill\"}]","status":1,"table_name":"","to_tree_array":0,"unique_key":"is","value_field":""}},
	{TableName: "tb_option_models", Condition: [][]interface{}{{"unique_key", "status"}},
		Data: map[string]interface{}{"children_option_model_key":"","color_array":"","color_field":"","data_type":0,"default_data":"","dynamic_params":"","icon_array":"","icon_field":"","index_num":2,"match_fields":"status","name":"启用|禁用","name_field":"","options_disable":0,"parent_field":"","select_order":"","select_where":"","static_data":"[{\"name\":\"启用\",\"value\":\"1\",\"color\":\"#66CC00\",\"icon\":\"ri-checkbox-circle-fill\"},{\"name\":\"禁用\",\"value\":\"0\",\"color\":\"#FF6666\",\"icon\":\"ri-forbid-fill\"}]","status":1,"table_name":"","to_tree_array":0,"unique_key":"status","value_field":""}},
	{TableName: "tb_option_models", Condition: [][]interface{}{{"unique_key", "is_show"}},
		Data: map[string]interface{}{"children_option_model_key":"","color_array":"","color_field":"","data_type":0,"default_data":"","dynamic_params":"","icon_array":"","icon_field":"","index_num":3,"match_fields":"is_show","name":"显示|隐藏","name_field":"","options_disable":0,"parent_field":"","select_order":"","select_where":"","static_data":"[{\"name\":\"显示\",\"value\":\"1\"},{\"name\":\"隐藏\",\"value\":\"0\"}]","status":1,"table_name":"","to_tree_array":0,"unique_key":"is_show","value_field":""}},
	{TableName: "tb_option_models", Condition: [][]interface{}{{"unique_key", "is_open"}},
		Data: map[string]interface{}{"children_option_model_key":"","color_array":"","color_field":"","data_type":0,"default_data":"","dynamic_params":"","icon_array":"","icon_field":"","index_num":4,"match_fields":"is_open","name":"开|关","name_field":"","options_disable":0,"parent_field":"","select_order":"","select_where":"","static_data":"[{\"name\":\"开\",\"value\":\"1\",\"color\":\"#66CC00\",\"icon\":\"ri-checkbox-circle-fill\"},{\"name\":\"关\",\"value\":\"0\",\"color\":\"#FF6666\",\"icon\":\"ri-forbid-fill\"}]","status":1,"table_name":"","to_tree_array":0,"unique_key":"is_open","value_field":""}},
	{TableName: "tb_option_models", Condition: [][]interface{}{{"unique_key", "rule_type"}},
		Data: map[string]interface{}{"children_option_model_key":"","color_array":"","color_field":"","data_type":0,"default_data":"","dynamic_params":"","icon_array":"","icon_field":"","index_num":5,"match_fields":"","name":"后台权限类型","name_field":"","options_disable":0,"parent_field":"","select_order":"","select_where":"","static_data":"[{\"name\":\"菜单\",\"value\":\"1\",\"icon\":\"ri-menu-fill\",\"color\":\"#0099CC\"},{\"name\":\"操作\",\"value\":\"2\",\"icon\":\"ri-cursor-line\",\"color\":\"#CC6699\"}]","status":1,"table_name":"","to_tree_array":0,"unique_key":"rule_type","value_field":""}},
}
