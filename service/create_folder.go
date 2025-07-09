package service

import (
	"filesys/dao"
	"filesys/model"
	"time"

	"github.com/gin-gonic/gin"
)

type createFolderService struct{}

var CreateFolderService = &createFolderService{}

func (s *createFolderService) Create(c *gin.Context, name string, parentID uint) (*model.File, error) {
	folder := &model.File{
		UserID:   int32(c.GetInt("user_id")),
		ParentID: int32(parentID),
		Name:     name,
		Type:     "folder",
		Ctime:    int32(time.Now().Unix()), // 创建时间
		Mtime:    int32(time.Now().Unix()), // 修改时间
	}
	if err := dao.Q.WithContext(c.Request.Context()).File.Create(folder); err != nil {
		return nil, err
	}
	return folder, nil
}
