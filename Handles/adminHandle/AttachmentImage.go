/**
 * @Author: 云起时
 * @Email: limingxiang@yqstech.com
 * @Description:
 * @File: AttachmentImage
 * @Version: 1.0.0
 * @Date: 2022/6/17 21:23
 */

package adminHandle

import (
	"github.com/gohouse/gorose/v2"
	"github.com/yqstech/gef/Utils/db"
	"github.com/yqstech/gef/builder"
	"github.com/yqstech/gef/util"
	"os"
)

type AttachmentImage struct {
	Base
}

// ImageExts 图片文件拓展名
var ImageExts = []string{".jpg", ".png", ".gif", ".jpeg", ".JPG", ".PNG", ".GIF", ".JPEG"}

// NodeBegin 开始
func (that AttachmentImage) NodeBegin(pageBuilder *builder.PageBuilder) (error, int) {
	pageBuilder.SetTitle("图片附件管理")
	pageBuilder.SetPageName("附件")
	pageBuilder.SetTbName("tb_attachment")
	return nil, 0
}

// NodeList 初始化列表
func (that AttachmentImage) NodeList(pageBuilder *builder.PageBuilder) (error, int) {
	pageBuilder.ListColumnClear()
	pageBuilder.SetListRightBtns("delete")
	pageBuilder.ListTopBtnsClear()
	pageBuilder.ListColumnAdd("src", "图片预览", "html", nil)
	pageBuilder.ListColumnAdd("file_size", "文件大小", "text", nil)
	pageBuilder.ListColumnAdd("create_time", "上传时间", "text", nil)
	return nil, 0
}

// NodeListCondition 修改查询条件
func (that AttachmentImage) NodeListCondition(pageBuilder *builder.PageBuilder, condition [][]interface{}) ([][]interface{}, error, int) {
	//追加查询条件
	condition = append(condition, []interface{}{
		"ext", "in", ImageExts,
	})
	return condition, nil, 0
}

// NodeListData 重写列表数据
func (that AttachmentImage) NodeListData(pageBuilder *builder.PageBuilder, data []gorose.Data) ([]gorose.Data, error, int) {
	for k, v := range data {
		data[k]["src"] = "<img src='" + v["src"].(string) + "' style='max-height:100px;max-width:350px'><br>" + v["src"].(string)
		data[k]["id"] = util.Int642String(v["id"].(int64))
	}
	return data, nil, 0
}

//NodeDeleteBefore 删除前操作
func (that AttachmentImage) NodeDeleteBefore(pageBuilder *builder.PageBuilder, id int64) (error, int) {
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
