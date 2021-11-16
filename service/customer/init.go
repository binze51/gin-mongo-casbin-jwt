package customer

import (
	"bytes"
	"context"
	"fmt"
	"io/ioutil"
	"strings"
	"time"

	"gitee.com/binze/binzekeji/model"
	"gitee.com/binze/binzekeji/pkg/engine"
	"gitee.com/binze/binzekeji/pkg/response"
	"gitee.com/binze/binzekeji/pkg/utils"
	"github.com/gin-gonic/gin"
	"github.com/h2non/filetype"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/gridfs"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.uber.org/zap"
)

//NewAccountHandler 构造系统管理
func NewManagerSvc(e *engine.Engine) *ManagerService {
	bucketOptions := options.GridFSBucket().SetName("bz_gridfs")
	bucket, _ := gridfs.NewBucket(e.DBConn.Database("binze").GetDatabase(), bucketOptions)
	SysHandler := &ManagerService{e, bucket}
	//初始化实体表字段的严格模式和相关备注
	// if conf.MongoRSConfig.CollInit {
	// 	SysHandler.UpsertCollJsonSchemaValidator()
	// }
	//初始化系统超管账号
	// if conf.HttpServConfig.InitRoot {
	// 	SysHandler.InitRoot(conf.HttpServConfig.RootPwd)
	// }
	return SysHandler
}

//SysHandler系统管理模块
type ManagerService struct {
	*engine.Engine
	FsBucket *gridfs.Bucket
}

//通用
func (a *ManagerService) Collection(collName string) *mongo.Collection {
	return a.Collection(collName)
}

//更新mongo
func (a *ManagerService) UpsertCollJsonSchemaValidator() error {
	return model.UpsertJsonSchemaValidator(a.Collection(model.TAccount), new(model.Account))
}

//初始化超级管理员
func (a *ManagerService) InitRoot(pwd string) (bool, error) {
	rootInfo := primitive.M{}
	thiscoll := a.Collection(model.TAccount)
	_ = thiscoll.FindOne(context.TODO(), bson.M{"accounName": "admin"}).Decode(&rootInfo)
	if len(rootInfo) == 0 {
		aes := utils.NewAes(model.AesKey, model.AesIv)
		encryptpwd, _ := aes.Encrypt(pwd)
		rootInfo = primitive.M{
			"name":       "admin",
			"password":   encryptpwd,
			"init_pwd":   true,
			"role":       "root",
			"email":      "better.tian@qq.com",
			"mobile":     "17722806690",
			"avatar_uri": "/api/account/v1/loadImg/6141be720749a0fb99fda397.jpeg",
			"create_at":  time.Now().Local().Unix(),
			"is_disable": false,
		}
		_, err := thiscoll.InsertOne(context.TODO(), rootInfo)
		if err != nil {
			a.Log.Error("init root account Err", zap.String("MongoERR", fmt.Sprintf("mongo table:(%s) - %s", model.TAccount, err.Error())))
			return false, err
		}
		a.Log.Info("init root account ok")
		return true, nil
	}
	return false, nil
}

// Upload 上传图片
func (a *ManagerService) Upload(c *gin.Context) {
	fileHeader, err := c.FormFile("avatar")
	if err != nil {
		response.Fail(response.ParamUploadError, c)
		return
	}
	if fileHeader.Size > 5<<20 {
		response.Fail(response.UploadFileSizeError, c)
		return
	}
	file, err := fileHeader.Open()
	if err != nil {
		response.Fail(response.FileOpenError, c)
		return
	}
	defer file.Close()
	fileBytes, err := ioutil.ReadAll(file)
	if err != nil {
		response.Fail(response.FileReadError, c)
		return
	}
	kind, _ := filetype.Match(fileBytes)
	if !strings.Contains("pngjpg", kind.Extension) {
		response.Fail(response.FileTypeError, c)
		return
	}
	fileid := primitive.NewObjectID().Hex()
	err = a.uploadWithFileID(fileid, fileHeader.Filename, fileBytes)
	if err != nil {
		response.Fail(response.FileSaveError, c)
		return
	}

	response.OkWithData(gin.H{
		"avatarUri": "/api/account/v1/loadImg/" + fileid + "." + kind.Extension,
	}, c)
}

// LoadImgByID 加载图片数据
func (a *ManagerService) LoadImgByID(c *gin.Context) {
	idArgs := c.Param("id")
	if idArgs == "" {
		response.Fail(response.ParamBindError, c)
		return
	}
	id := strings.Split(idArgs, ".")[0]
	fileBytes, err := a.loadFileByID(id)
	if err != nil {
		response.Fail(response.LoadFileError, c)
		return
	}
	response.FileLoad(fileBytes, c)
}

// uploadWithFileID 使用文件id上传图片文件
func (a *ManagerService) uploadWithFileID(fileID, fileName string, fileContent []byte) error {
	err := a.FsBucket.UploadFromStreamWithID(fileID, fileName, bytes.NewBuffer(fileContent))
	if err != nil {
		panic(err)
	}
	return nil
}

//loadFileByID - 使用文件id加载图片文件
func (a *ManagerService) loadFileByID(fileID string) (fileContent []byte, err error) {
	fileBuffer := bytes.NewBuffer(nil)
	if _, err = a.FsBucket.DownloadToStream(fileID, fileBuffer); err != nil {
		panic(err)
		// return nil, err
	}
	return fileBuffer.Bytes(), nil
}
