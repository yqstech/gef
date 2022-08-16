/**
 * @Author: 云起时
 * @Email: limingxiang@yqstech.com
 * @Description:
 * @File: database
 * @Version: 1.0.0
 * @Date: 2022/8/16 11:59
 */

package database

// Tables 需要维护的所有的表结构体
var Tables = []interface{}{
	&TbAdmin{},
	&TbAdminGroup{},
	&TbAdminLog{},
	&TbAdminRules{},
	&TbAdminToken{},
	&TbAppConfigs{},
	&TbAttachment{},
	&TbConfigs{},
	&TbConfigsGroup{},
	&TbOptionModels{},
	&TbEasyModels{},
	&TbEasyModelsFields{},
	&TbEasyModelsButtons{},
	&TbEasyCurdModels{},
	&TbEasyCurdModelsFields{},
}

//自动生成的表补充参数 将下边的#去掉
//正则替换，comment参数
//NOT NULL#"(.*)`(.*)// (.*)
//NOT NULL;comment:$3"$1`$2// $3

//正则替换，补充default:''
//char(.*)\);#NOT
//char$1);default:'';NOT

// TbAdmin 管理员表
type TbAdmin struct {
	ID
	GroupID  int    `gorm:"column:group_id;type:int(11);default:0;NOT NULL;comment:所属角色" json:"group_id"`      // 所属角色
	Name     string `gorm:"column:name;type:varchar(50);default:'';NOT NULL;comment:管理员昵称" json:"name"`        // 用户昵称
	Account  string `gorm:"column:account;type:varchar(15);default:'';NOT NULL;comment:登录账号" json:"account"`   // 账号
	Password string `gorm:"column:password;type:varchar(32);default:'';NOT NULL;comment:登录密码" json:"password"` // 密码
	CUSD
}

func (m *TbAdmin) TableName() string {
	return "tb_admin"
}

// TbAdminGroup 后台角色表
type TbAdminGroup struct {
	ID
	GroupName string `gorm:"column:group_name;type:varchar(50);default:'';NOT NULL;comment:角色权限组名称" json:"group_name"` // 角色权限组名称
	Rules     string `gorm:"column:rules;type:varchar(5000);default:[];NOT NULL;comment:权限组" json:"rules"`         // 权限组
	CUSD
}

func (m *TbAdminGroup) TableName() string {
	return "tb_admin_group"
}

// TbAdminLog 管理员后台日志表
type TbAdminLog struct {
	ID
	Rule        string `gorm:"column:rule;type:varchar(30);default:'';NOT NULL;comment:访问链接" json:"rule"`           // 访问链接
	RuleName    string `gorm:"column:rule_name;type:varchar(30);default:'';NOT NULL;comment:链接名称" json:"rule_name"` // 链接名称
	Url         string `gorm:"column:url;type:varchar(100);default:'';NOT NULL" json:"url"`
	Notice      string `gorm:"column:notice;type:varchar(30);default:'';NOT NULL;comment:备注" json:"notice"`               // 备注
	AccountID   int    `gorm:"column:account_id;type:int(11);default:0;NOT NULL;comment:账户ID" json:"account_id"` // 账户ID
	AccountName string `gorm:"column:account_name;type:varchar(30);default:'';NOT NULL;comment:登录账户名" json:"account_name"`   // 登录账户名
	Account     string `gorm:"column:account;type:varchar(30);default:'';NOT NULL;comment:登录账号" json:"account"`             // 登录账号
	Data        string `gorm:"column:data;type:text" json:"data"`                                   // 数据
	CD
}

func (m *TbAdminLog) TableName() string {
	return "tb_admin_log"
}

// TbAdminRules 后台权限菜单表
type TbAdminRules struct {
	ID
	Pid      int    `gorm:"column:pid;type:int(11);default:0;NOT NULL;comment:上级权限" json:"pid"`                // 上级权限
	Name     string `gorm:"column:name;type:varchar(50);default:'';NOT NULL;comment:权限名称" json:"name"`                    // 权限名称
	Type     int    `gorm:"column:type;type:tinyint(4);default:1;NOT NULL;comment:权限类型1菜单2按钮" json:"type"`           // 权限类型1菜单2按钮
	IsCompel int    `gorm:"column:is_compel;type:tinyint(1);default:0;NOT NULL;comment:是否必选" json:"is_compel"` // 是否必选
	Icon     string `gorm:"column:icon;type:varchar(30);default:'';NOT NULL;comment:图标字体" json:"icon"`                    // 图标字体
	Route    string `gorm:"column:route;type:varchar(100);default:'';NOT NULL;comment:路由" json:"route"`                 // 路由
	IndexNum int    `gorm:"column:index_num;type:int(11);default:200;NOT NULL;comment:排序" json:"index_num"`  // 排序
	OpenLog  int    `gorm:"column:open_log;type:tinyint(1);default:0" json:"open_log"`            // 是否开启日志
	CUSD
}

func (m *TbAdminRules) TableName() string {
	return "tb_admin_rules"
}

// TbAdminToken 登录token
type TbAdminToken struct {
	ID
	AccountID int    `gorm:"column:account_id;type:int(11);default:0;NOT NULL;comment:用户ID" json:"account_id"` // 用户ID
	Token     string `gorm:"column:token;type:varchar(40);default:'';NOT NULL;comment:token" json:"token"`                 // token
	CUSD
}

func (m *TbAdminToken) TableName() string {
	return "tb_admin_token"
}

// TbAppConfigs 应用配置表（除了应用内的配置以外的配置项）
type TbAppConfigs struct {
	ID
	GroupID int    `gorm:"column:group_id;type:int(11);default:0;NOT NULL" json:"group_id"`
	Name    string `gorm:"column:name;type:varchar(100);default:'';NOT NULL;comment:关键字" json:"name"` // 关键字
	Value   string `gorm:"column:value;type:text;NOT NULL;comment:配置内容" json:"value"`       // 配置内容
	CUSD
}

func (m *TbAppConfigs) TableName() string {
	return "tb_app_configs"
}

// TbAttachment 附件表
type TbAttachment struct {
	ID       uint64 `gorm:"column:id;type:bigint(20) unsigned;primary_key" json:"id"`                    // ID
	GroupID  uint   `gorm:"column:group_id;type:tinyint(4) unsigned;default:0;NOT NULL;comment:用户分组" json:"group_id"` // 用户分组
	UserID   uint64 `gorm:"column:user_id;type:bigint(20) unsigned;default:0;NOT NULL;comment:用户ID" json:"user_id"`   // 用户ID
	FileType int    `gorm:"column:file_type;type:int(11);default:1;NOT NULL;comment:文件类型|1图片2视频3文件" json:"file_type"`           // 文件类型|1图片2视频3文件
	FileName string `gorm:"column:file_name;type:varchar(200);default:'';NOT NULL;comment:文件名称" json:"file_name"`                // 文件名称
	Path     string `gorm:"column:path;type:varchar(255);default:'';NOT NULL;comment:文件路径" json:"path"`                          // 文件路径
	Src      string `gorm:"column:src;type:varchar(255);default:'';NOT NULL;comment:文件链接（暂时无用）" json:"src"`                            // 文件链接（暂时无用）
	Ext      string `gorm:"column:ext;type:char(8);default:'';NOT NULL;comment:文件类型" json:"ext"`                                 // 文件类型
	FileSize uint   `gorm:"column:file_size;type:int(11) unsigned;default:0;NOT NULL;comment:文件大小" json:"file_size"`  // 文件大小
	Md5      string `gorm:"column:md5;type:char(32);default:'';NOT NULL" json:"md5"`
	CUSD            // 是否删除
}

func (m *TbAttachment) TableName() string {
	return "tb_attachment"
}

// TbConfigs 应用设置项
type TbConfigs struct {
	ID
	GroupID   int    `gorm:"column:group_id;type:tinyint(4);default:1;NOT NULL;comment:设置项分组" json:"group_id"` // 设置项分组
	Name      string `gorm:"column:name;type:char(100);default:'';NOT NULL;comment:关键字" json:"name"`                    // 关键字
	Value     string `gorm:"column:value;type:varchar(1024);default:'';NOT NULL;comment:配置值" json:"value"`              // 配置值
	Title     string `gorm:"column:title;type:varchar(100);default:'';NOT NULL;comment:配置名称" json:"title"`               // 配置名称
	Notice    string `gorm:"column:notice;type:varchar(200);default:'';NOT NULL;comment:配置补充说明" json:"notice"`             // 配置补充说明
	FieldType string `gorm:"column:field_type;type:char(20);default:'';NOT NULL;comment:调用form表单项类型text html textarea" json:"field_type"`         // 调用form表单项类型text html textarea
	Options   string `gorm:"column:options;type:varchar(1024);default:[];NOT NULL" json:"options"`
	IndexNum  int    `gorm:"column:index_num;type:int(11);default:200;NOT NULL" json:"index_num"`
	If        string `gorm:"column:if;type:varchar(100);default:'';NOT NULL;comment:展示条件" json:"if"` // 展示条件
	CUSD
}

func (m *TbConfigs) TableName() string {
	return "tb_configs"
}

// TbConfigsGroup 应用设置项分组
type TbConfigsGroup struct {
	ID
	GroupName string `gorm:"column:group_name;type:char(100);default:'';NOT NULL;comment:分组名" json:"group_name"` // 分组名
	Note      string `gorm:"column:note;type:varchar(100)" json:"note"`
	CUSD
}

func (m *TbConfigsGroup) TableName() string {
	return "tb_configs_group"
}

// TbEasyCurdModels easyCurd模型
type TbEasyCurdModels struct {
	ID
	ModelKey           string `gorm:"column:model_key;type:varchar(50);default:'';NOT NULL;comment:模型关键字" json:"model_key"`                             // 模型关键字
	ModelName          string `gorm:"column:model_name;type:varchar(50);default:'';NOT NULL;comment:模型名称" json:"model_name"`                           // 模型名称
	DbTableName        string `gorm:"column:table_name;type:varchar(50);default:'';NOT NULL;comment:关联数据表名称" json:"table_name"`                           // 关联数据表名称
	AllowSelect        int    `gorm:"column:allow_select;type:int(11);default:0;NOT NULL;comment:允许查询" json:"allow_select"`                 // 允许查询
	AllowCreate        int    `gorm:"column:allow_create;type:int(11);default:0;NOT NULL;comment:允许新增" json:"allow_create"`                 // 允许新增
	AllowUpdate        int    `gorm:"column:allow_update;type:int(11);default:0;NOT NULL;comment:允许更新" json:"allow_update"`                 // 允许更新
	AllowDelete        int    `gorm:"column:allow_delete;type:int(11);default:0;NOT NULL;comment:允许删除" json:"allow_delete"`                 // 允许删除
	SoftDeleteDisable  int    `gorm:"column:soft_delete_disable;type:int(11);default:0;NOT NULL;comment:是否禁用软删除" json:"soft_delete_disable"`   // 是否禁用软删除
	CheckLogin         int    `gorm:"column:check_login;type:int(11);default:1;NOT NULL;comment:登录校验" json:"check_login"`                   // 登录校验
	UkName             string `gorm:"column:uk_name;type:varchar(30);default:user_id;NOT NULL;comment:代表用户id的字段名称" json:"uk_name"`                 // 代表用户id的字段名称
	PkName             string `gorm:"column:pk_name;type:varchar(30);default:id;NOT NULL;comment:代表主键的字段名称" json:"pk_name"`                      // 代表主键的字段名称
	SelectOrder        string `gorm:"column:select_order;type:varchar(30);default:id desc;NOT NULL;comment:查询排序方式" json:"select_order"`       // 查询排序方式
	Fields             string `gorm:"column:fields;type:varchar(1024);default:*;NOT NULL;comment:查询字段" json:"fields"`                       // 查询字段
	PrivateFields      string `gorm:"column:private_fields;type:varchar(1024);default:'';NOT NULL;comment:私密字段" json:"private_fields"`                 // 私密字段
	LockFields         string `gorm:"column:lock_fields;type:varchar(1024);default:'';NOT NULL;comment:锁定字段，不允许修改" json:"lock_fields"`                       // 锁定字段，不允许修改
	SelectWithDisabled int    `gorm:"column:select_with_disabled;type:int(11);default:0;NOT NULL;comment:禁用可查|查询数据时是否包含已禁用的数据" json:"select_with_disabled"` // 禁用可查|查询数据时是否包含已禁用的数据
	CUSD
}

func (m *TbEasyCurdModels) TableName() string {
	return "tb_easy_curd_models"
}

// TbEasyCurdModelsFields easyCurd模型字段管理
type TbEasyCurdModelsFields struct {
	ID
	ModelID        int    `gorm:"column:model_id;type:int(11);default:0;NOT NULL;comment:关联easyCurd模型ID" json:"model_id"`                 // 关联easyCurd模型ID
	FieldKey       string `gorm:"column:field_key;type:varchar(50);default:'';NOT NULL;comment:模型字段关键字" json:"field_key"`                     // 模型字段关键字
	FieldName      string `gorm:"column:field_name;type:varchar(50);default:'';NOT NULL;comment:模型字段名称" json:"field_name"`                   // 模型字段名称
	FieldNote      string `gorm:"column:field_note;type:varchar(50);default:'';NOT NULL;comment:模型字段备注" json:"field_note"`                   // 模型字段备注
	OptionModelsID int    `gorm:"column:option_models_id;type:int(11);default:0;NOT NULL;comment:选项集id" json:"option_models_id"` // 选项集id
	IsPrivate      int    `gorm:"column:is_private;type:int(11);default:0;NOT NULL;comment:是否私密" json:"is_private"`             // 是否私密
	IsLock         int    `gorm:"column:is_lock;type:int(11);default:0;NOT NULL;comment:是否锁定|禁止修改" json:"is_lock"`                   // 是否锁定|禁止修改
	CUSD
}

func (m *TbEasyCurdModelsFields) TableName() string {
	return "tb_easy_curd_models_fields"
}

// TbEasyModels easy模型
type TbEasyModels struct {
	ID
	ModelKey      string `gorm:"column:model_key;type:varchar(50);default:'';NOT NULL;comment:模型关键字" json:"model_key"`                   // 模型关键字
	ModelName     string `gorm:"column:model_name;type:varchar(50);default:'';NOT NULL;comment:模型名称" json:"model_name"`                 // 模型名称
	DbTableName   string `gorm:"column:table_name;type:varchar(50);default:'';NOT NULL;comment:关联数据表名称" json:"table_name"`                 // 关联数据表名称
	AllowCreate   int    `gorm:"column:allow_create;type:int(11);default:1;NOT NULL;comment:新增按钮" json:"allow_create"`       // 新增按钮
	AllowUpdate   int    `gorm:"column:allow_update;type:int(11);default:1;NOT NULL;comment:允许修改" json:"allow_update"`       // 允许修改
	AllowStatus   int    `gorm:"column:allow_status;type:int(11);default:1;NOT NULL;comment:状态按钮" json:"allow_status"`       // 状态按钮
	AllowDelete   int    `gorm:"column:allow_delete;type:int(11);default:1;NOT NULL;comment:删除按钮" json:"allow_delete"`       // 删除按钮
	DefaultColumn int    `gorm:"column:default_column;type:int(11);default:1;NOT NULL;comment:是否显示默认列(删除，程序自动同步了id列)" json:"default_column"`   // 是否显示默认列(删除，程序自动同步了id列)
	OrderType     string `gorm:"column:order_type;type:varchar(50);default:id desc;NOT NULL;comment:列表页排序方式" json:"order_type"` // 列表页排序方式
	PageSize      int    `gorm:"column:page_size;type:int(11);default:20;NOT NULL;comment:每页多少条数据" json:"page_size"`            // 每页多少条数据
	BatchAction   int    `gorm:"column:batch_action;type:tinyint(1);default:0;NOT NULL;comment:是否支持批量操作" json:"batch_action"`    // 是否支持批量操作
	PageNotice    string `gorm:"column:page_notice;type:text;NOT NULL" json:"page_notice"`
	TabsForList   string `gorm:"column:tabs_for_list;type:text;NOT NULL;comment:列表tab页" json:"tabs_for_list"`       // 列表tab页
	TopButtons    string `gorm:"column:top_buttons;type:text;NOT NULL;comment:顶部按钮" json:"top_buttons"`           // 顶部按钮
	RightButtons  string `gorm:"column:right_buttons;type:text;NOT NULL;comment:右侧按钮" json:"right_buttons"`       // 右侧按钮
	UrlParams     string `gorm:"column:url_params;type:text;NOT NULL;comment:url传参转sql查询条件" json:"url_params"`             // url传参转sql查询条件
	LevelIndent   string `gorm:"column:level_indent;type:varchar(100);default:'';NOT NULL;comment:上下级缩进|格式为:level字段:indent字段，例如pid:name" json:"level_indent"` // 上下级缩进|格式为:level字段:indent字段，例如pid:name
	Note          string `gorm:"column:note;type:varchar(30);default:'';NOT NULL;comment:备注" json:"note"`                  // 备注
	CUSD
}

func (m *TbEasyModels) TableName() string {
	return "tb_easy_models"
}

// TbEasyModelsButtons easy模型按钮
type TbEasyModelsButtons struct {
	ID
	ButtonKey   string `gorm:"column:button_key;type:varchar(50);default:'';NOT NULL;comment:按钮唯一识别" json:"button_key"`              // 按钮唯一识别
	ButtonName  string `gorm:"column:button_name;type:varchar(20);default:'';NOT NULL;comment:按钮名称" json:"button_name"`            // 按钮名称
	ButtonNote  string `gorm:"column:button_note;type:varchar(50);default:'';NOT NULL;comment:按钮备注" json:"button_note"`            // 按钮备注
	ButtonIcon  string `gorm:"column:button_icon;type:varchar(50);default:'';NOT NULL;comment:按钮图标" json:"button_icon"`            // 按钮图标
	ClassName   string `gorm:"column:class_name;type:varchar(50);default:'';NOT NULL;comment:按钮class样式类" json:"class_name"`              // 按钮class样式类
	Display     string `gorm:"column:display;type:varchar(50);default:'';NOT NULL;comment:按钮显示条件" json:"display"`                    // 按钮显示条件
	ActionType  int    `gorm:"column:action_type;type:int(11);default:200;NOT NULL;comment:按钮类型:1、ajax操作 2、弹出页面 3、执行javascript" json:"action_type"`    // 按钮类型:1、ajax操作 2、弹出页面 3、执行javascript
	Action      string `gorm:"column:action;type:varchar(50);default:'';NOT NULL;comment:动作值，校验权限使用，格式例如: add 或 /user/add" json:"action"`                      // 动作值，校验权限使用，格式例如: add 或 /user/add
	ActionUrl   string `gorm:"column:action_url;type:varchar(50);default:'';NOT NULL;comment:动作url地址，相对路径或绝对路径,使用/开头自动添加后台url地址" json:"action_url"`              // 动作url地址，相对路径或绝对路径,使用/开头自动添加后台url地址
	ConfirmMsg  string `gorm:"column:confirm_msg;type:varchar(50);default:'';NOT NULL;comment:确认对话框信息，ActionType=1有效" json:"confirm_msg"`            // 确认对话框信息，ActionType=1有效
	LayerTitle  string `gorm:"column:layer_title;type:varchar(50);default:'';NOT NULL;comment:弹窗标题" json:"layer_title"`            // 弹窗标题
	LayerWidth  string `gorm:"column:layer_width;type:varchar(50);default:'';NOT NULL;comment:弹窗宽度，支持px或%" json:"layer_width"`            // 弹窗宽度，支持px或%
	LayerHeight string `gorm:"column:layer_height;type:varchar(50);default:'';NOT NULL;comment:弹窗高度，支持px或%" json:"layer_height"`          // 弹窗高度，支持px或%
	BatchAction int    `gorm:"column:batch_action;type:tinyint(1);default:0;NOT NULL;comment:是否支持批量操作" json:"batch_action"` // 是否支持批量操作
	CUSD
}

func (m *TbEasyModelsButtons) TableName() string {
	return "tb_easy_models_buttons"
}

// TbEasyModelsFields 模型字段
type TbEasyModelsFields struct {
	ID
	ModelID               int    `gorm:"column:model_id;type:int(11);default:0;NOT NULL;comment:模型ID" json:"model_id"`                                 // 模型ID
	FieldKey              string `gorm:"column:field_key;type:varchar(50);default:'';NOT NULL;comment:模型字段关键字" json:"field_key"`                                     // 模型字段关键字
	FieldName             string `gorm:"column:field_name;type:varchar(50);default:'';NOT NULL;comment:模型字段名称" json:"field_name"`                                   // 模型字段名称
	FieldNameReset        string `gorm:"column:field_name_reset;type:varchar(50);default:'';NOT NULL;comment:重置字段名称（列表顶部）" json:"field_name_reset"`                       // 重置字段名称（列表顶部）
	FieldNotice           string `gorm:"column:field_notice;type:varchar(50);default:'';NOT NULL;comment:字段提示" json:"field_notice"`                               // 字段提示
	IndexNum              int    `gorm:"column:index_num;type:int(11);default:200;NOT NULL;comment:排序" json:"index_num"`                             // 排序
	OptionModelsID        int    `gorm:"column:option_models_id;type:int(11);default:0;NOT NULL;comment:选项集id" json:"option_models_id"`                 // 选项集id
	OptionBeautify        int    `gorm:"column:option_beautify;type:tinyint(1);default:1;NOT NULL;comment:选项美化" json:"option_beautify"`                // 选项美化
	OptionIndent          int    `gorm:"column:option_indent;type:tinyint(1);default:0;NOT NULL;comment:选项按照上下级缩进" json:"option_indent"`                    // 选项按照上下级缩进
	DynamicOptionModelsID int    `gorm:"column:dynamic_option_models_id;type:int(11);default:0;NOT NULL;comment:动态选项集" json:"dynamic_option_models_id"` // 动态选项集
	WatchFields           string `gorm:"column:watch_fields;type:varchar(200);default:'';NOT NULL;comment:监听字段|多个字段用英文逗号分割" json:"watch_fields"`                              // 监听字段|多个字段用英文逗号分割
	SetAsTabs             int    `gorm:"column:set_as_tabs;type:tinyint(1);default:0;NOT NULL;comment:是否将字段设置成列表Tabs" json:"set_as_tabs"`                        // 是否将字段设置成列表Tabs
	IsShowOnList          int    `gorm:"column:is_show_on_list;type:int(11);default:1;NOT NULL;comment:列表页显示" json:"is_show_on_list"`                   // 列表页显示
	DataTypeOnList        string `gorm:"column:data_type_on_list;type:varchar(50);default:'';NOT NULL;comment:列表页数据类型（组件）" json:"data_type_on_list"`                     // 列表页数据类型（组件）
	DataTypeCommandOnList string `gorm:"column:data_type_command_on_list;type:varchar(50);default:'';NOT NULL;comment:数据指令（switch）" json:"data_type_command_on_list"`     // 数据指令（switch）
	AllowCreate           int    `gorm:"column:allow_create;type:int(11);default:1;NOT NULL;comment:可新增" json:"allow_create"`                         // 可新增
	AllowUpdate           int    `gorm:"column:allow_update;type:int(11);default:1;NOT NULL;comment:可修改" json:"allow_update"`                         // 可修改
	DataTypeOnCreate      string `gorm:"column:data_type_on_create;type:varchar(50);default:'';NOT NULL;comment:新增页数据类型" json:"data_type_on_create"`                 // 新增页数据类型
	DataTypeOnUpdate      string `gorm:"column:data_type_on_update;type:varchar(50);default:'';NOT NULL;comment:编辑页数据类型（组件）" json:"data_type_on_update"`                 // 编辑页数据类型（组件）
	IsMust                int    `gorm:"column:is_must;type:int(11);default:1;NOT NULL;comment:必填项" json:"is_must"`                                   // 必填项
	DefaultValue          string `gorm:"column:default_value;type:varchar(50);default:'';NOT NULL;comment:默认值" json:"default_value"`                             // 默认值
	ExpandIf              string `gorm:"column:expand_if;type:varchar(255);default:'';NOT NULL;comment:拓展数据，查询条件" json:"expand_if"`                                    // 拓展数据，查询条件
	GroupTitle            string `gorm:"column:group_title;type:varchar(100);default:'';NOT NULL;comment:表单分组名称" json:"group_title"`                                // 表单分组名称
	FieldAugment          string `gorm:"column:field_augment;type:varchar(1024);default:'';NOT NULL;comment:字段值扩充" json:"field_augment"`                           // 字段值扩充
	AttachToField         string `gorm:"column:attach_to_field;type:varchar(50);default:'';NOT NULL;comment:数据附加到此字段下" json:"attach_to_field"`                         // 数据附加到此字段下
	SaveTransRule         string `gorm:"column:save_trans_rule;type:varchar(100);default:'';NOT NULL;comment:保存时的数据转换规则" json:"save_trans_rule"`                        // 保存时的数据转换规则
	FieldStyleReset       string `gorm:"column:field_style_reset;type:varchar(50);default:'';NOT NULL;comment:重置字段样式列表" json:"field_style_reset"`                     // 重置字段样式列表
	CUSD
}

func (m *TbEasyModelsFields) TableName() string {
	return "tb_easy_models_fields"
}

// TbOptionModels 模型字段选项数据源
type TbOptionModels struct {
	ID
	Name                  string `gorm:"column:name;type:varchar(30);default:未命名;NOT NULL;comment:选项名称" json:"name"`                                   // 选项名称
	DataType              int    `gorm:"column:data_type;type:int(11);default:0;NOT NULL;comment:数据类型" json:"data_type"`                               // 数据类型
	StaticData            string `gorm:"column:static_data;type:varchar(1024);default:'';NOT NULL;comment:静态数据(json)" json:"static_data"`                               // 静态数据(json)
	DbTableName           string `gorm:"column:table_name;type:varchar(50);default:'';NOT NULL;comment:数据表" json:"table_name"`                                   // 数据表
	ValueField            string `gorm:"column:value_field;type:varchar(50);default:'';NOT NULL;comment:value字段" json:"value_field"`                                 // value字段
	NameField             string `gorm:"column:name_field;type:varchar(50);default:'';NOT NULL;comment:name字段" json:"name_field"`                                   // name字段
	ParentField           string `gorm:"column:parent_field;type:varchar(50);default:'';NOT NULL;comment:上级字段" json:"parent_field"`                               // 上级字段
	ToTreeArray           int    `gorm:"column:to_tree_array;type:tinyint(1);default:0;NOT NULL;comment:选项集根据pid值转多维数组" json:"to_tree_array"`                    // 选项集根据pid值转多维数组
	ChildrenOptionModelID int    `gorm:"column:children_option_model_id;type:int(11);default:0;NOT NULL;comment:下级选项集" json:"children_option_model_id"` // 下级选项集
	OptionsDisable        int    `gorm:"column:options_disable;type:tinyint(1);default:0;NOT NULL;comment:当前选项集禁选" json:"options_disable"`                // 当前选项集禁选
	ColorField            string `gorm:"column:color_field;type:varchar(50);default:'';NOT NULL;comment:颜色字段" json:"color_field"`                                 // 颜色字段
	ColorArray            string `gorm:"column:color_array;type:varchar(1024);default:'';NOT NULL;comment:颜色集" json:"color_array"`                               // 颜色集
	IconField             string `gorm:"column:icon_field;type:varchar(50);default:'';NOT NULL;comment:Icon字段" json:"icon_field"`                                   // Icon字段
	IconArray             string `gorm:"column:icon_array;type:varchar(2048);default:'';NOT NULL;comment:图标集" json:"icon_array"`                                 // 图标集
	SelectWhere           string `gorm:"column:select_where;type:varchar(100);default:'';NOT NULL;comment:查询条件" json:"select_where"`                              // 查询条件
	DynamicParams         string `gorm:"column:dynamic_params;type:varchar(1024);default:'';NOT NULL;comment:动态参数" json:"dynamic_params"`                         // 动态参数
	SelectOrder           string `gorm:"column:select_order;type:varchar(100);default:'';NOT NULL;comment:查询排序条件" json:"select_order"`                              // 查询排序条件
	IndexNum              int    `gorm:"column:index_num;type:int(11);default:200;NOT NULL;comment:显示排序" json:"index_num"`                             // 显示排序
	MatchFields           string `gorm:"column:match_fields;type:varchar(1024);default:'';NOT NULL;comment:匹配字段" json:"match_fields"`                             // 匹配字段
	DefaultData           string `gorm:"column:default_data;type:varchar(1024);default:'';NOT NULL;comment:默认数据" json:"default_data"`                             // 默认数据
	CUSD
}

func (m *TbOptionModels) TableName() string {
	return "tb_option_models"
}
