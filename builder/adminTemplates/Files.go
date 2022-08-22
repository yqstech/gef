/**
 * @Author: 云起时
 * @Email: limingxiang@yqstech.com
 * @Description:
 * @File: Templates
 * @Version: 1.0.0
 * @Date: 2022/3/16 2:52 下午
 */

package adminTemplates

import "embed"

//go:embed *
var Templates embed.FS

// AdminTemplatesAdd 自定义文件集
var AdminTemplatesAdd []embed.FS
