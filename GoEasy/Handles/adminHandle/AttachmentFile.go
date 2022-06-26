/**
 * @Author: 云起时
 * @Email: limingxiang@yqstech.com
 * @Description:
 * @File: AttachmentFile
 * @Version: 1.0.0
 * @Date: 2022/6/17 21:23
 */

package adminHandle

import (
	"github.com/gef/GoEasy/EasyApp"
	"github.com/gef/GoEasy/Utils/db"
	"github.com/gef/GoEasy/Utils/util"
	"github.com/gohouse/gorose/v2"
	"os"
	"strings"
)

type AttachmentFile struct {
	Base
}

// FileExtImage 文件拓展名图片
var FileExtImage = map[string]string{
	".doc":  "/static/images/exts/doc.png",
	".docx": "/static/images/exts/doc.png",
	".dotx": "/static/images/exts/doc.png",
	".dot":  "/static/images/exts/doc.png",
	".docm": "/static/images/exts/doc.png",
	".ppt":  "/static/images/exts/ppt.png",
	".pptx": "/static/images/exts/ppt.png",
	".pps":  "/static/images/exts/ppt.png",
	".ppsx": "/static/images/exts/ppt.png",
	".xls":  "/static/images/exts/xls.png",
	".xlsx": "/static/images/exts/xls.png",
	".pdf":  "/static/images/exts/pdf.png",
	".txt":  "/static/images/exts/txt.png",
	".zip":  "/static/images/exts/zip.png",
	".rar":  "/static/images/exts/zip.png",
	".7z":   "/static/images/exts/zip.png",
	".mp4":  "/static/images/exts/video.png",
	".avi":  "/static/images/exts/video.png",
	".mov":  "/static/images/exts/video.png",
	".rmvb": "/static/images/exts/video.png",
	".flv":  "/static/images/exts/video.png",
	".3gp":  "/static/images/exts/video.png",
	"other": "/static/images/exts/other.png",
}

// NodeBegin 开始
func (that AttachmentFile) NodeBegin(pageData *EasyApp.PageData) (error, int) {
	pageData.SetTitle("文件附件管理")
	pageData.SetPageName("附件")
	pageData.SetTbName("tb_attachment")
	return nil, 0
}

// NodeList 初始化列表
func (that AttachmentFile) NodeList(pageData *EasyApp.PageData) (error, int) {
	pageData.ListColumnClear()
	pageData.SetListRightBtns("delete")
	pageData.ListTopBtnsClear()
	pageData.ListColumnAdd("ext_image", "文件类型", "image60", nil)
	pageData.ListColumnAdd("file_name", "文件名称", "text", nil)
	pageData.ListColumnAdd("ext", "文件后缀", "text", nil)
	pageData.ListColumnAdd("file_size", "文件大小", "text", nil)
	pageData.ListColumnAdd("create_time", "上传时间", "text", nil)
	return nil, 0
}

// NodeListCondition 修改查询条件
func (that AttachmentFile) NodeListCondition(pageData *EasyApp.PageData, condition [][]interface{}) ([][]interface{}, error, int) {
	//追加查询条件
	condition = append(condition, []interface{}{
		"ext", "not in", ImageExts,
	})
	return condition, nil, 0
}

// NodeListData 重写列表数据
func (that AttachmentFile) NodeListData(pageData *EasyApp.PageData, data []gorose.Data) ([]gorose.Data, error, int) {
	for k, v := range data {
		data[k]["id"] = util.Int642String(v["id"].(int64))
		if src, ok := FileExtImage[strings.ToLower(v["ext"].(string))]; ok {
			data[k]["ext_image"] = src
		} else {
			data[k]["ext_image"] = FileExtImage["other"]
		}
	}
	return data, nil, 0
}

//NodeDeleteBefore 删除前操作
func (that AttachmentFile) NodeDeleteBefore(pageData *EasyApp.PageData, id int64) (error, int) {
	first, err := db.New().Table("tb_attachment").Where("id", id).First()
	if err != nil {
		return err, 0
	}
	if first != nil {
		path := first["path"].(string)
		if _, err := os.Stat(path); os.IsNotExist(err) {
			return nil, 0
		}
		err := os.Remove(path)
		if err != nil {
			return err, 0
		}
	}
	return nil, 0
}
