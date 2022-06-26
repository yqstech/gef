/**
 * @Author: 云起时
 * @Email: limingxiang@yqstech.com
 * @Description:
 * @File: Files
 * @Version: 1.0.0
 * @Date: 2021/11/9 8:52 下午
 */

package static

import "embed"

//go:embed *
var Files embed.FS

// FilesAdd 自定义文件
var FilesAdd embed.FS
