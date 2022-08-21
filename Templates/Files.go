/**
 * @Author: 云起时
 * @Email: limingxiang@yqstech.com
 * @Description:
 * @File: Files
 * @Version: 1.0.0
 * @Date: 2022/3/16 2:52 下午
 */

package Templates

import "embed"

//go:embed *
var Files embed.FS

// FilesAdds 自定义文件集
var FilesAdds []embed.FS

// FilesSelf 自定义文件
var FilesSelf embed.FS
