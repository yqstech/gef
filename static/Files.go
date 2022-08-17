/**
 * @Author: 云起时
 * @Email: limingxiang@yqstech.com
 * @Description:
 * @File: Files
 * @Version: 1.0.0
 * @Date: 2021/11/9 8:52 下午
 */

package static

import (
	"embed"
	"net/http"
)

//go:embed *
var Files embed.FS

var FileSystems []http.FileSystem