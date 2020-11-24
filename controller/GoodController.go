package controller

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"itStudioTB/common"
	"itStudioTB/model"
	"itStudioTB/response"
	"itStudioTB/utils"
	"log"
	"net/http"
	"path"
	"strings"
	"time"
)

// 创建商品
func CreateGood(ctx *gin.Context) {
	DB := common.GetDB()

	var requestGood model.Good
	err := ctx.ShouldBind(&requestGood)

	if err != nil {
		log.Printf("创建商品数据绑定失败  error : %v", err.Error())
		return
	}

	// 获取参数
	name := requestGood.Name
	price := requestGood.Price
	info := requestGood.Info
	quantity := requestGood.Quantity
	pic := requestGood.Picture

	if !utils.VerifyGoodName(name) {
		response.Response(ctx, http.StatusUnprocessableEntity, nil, "商品名称不合法,商品名称为1-20字节的中文或英文")
		return
	}

	if len(info) > 200 {
		response.Response(ctx, http.StatusUnprocessableEntity, nil, "商品简介不合法,商品简介长度为0-200字节")
		return
	}

	if quantity < 0 || quantity > 99 {
		response.Response(ctx, http.StatusUnprocessableEntity, nil, "商品数量不合法，商品数量0-99")
		return
	}

	user, _ := ctx.Get("user")
	userId := user.(model.User).UserId

	if utils.IsGoodExist(DB, userId, name) {
		response.Response(ctx, http.StatusUnprocessableEntity, nil, "该用户已经拥有该名称的商品了")
		return
	}

	newGood := model.Good{
		Name:     name,
		Price:    price,
		Info:     info,
		Quantity: quantity,
		UserId:   userId,
		Picture:  pic,
	}
	// 保存商品
	DB.Create(&newGood)

	response.Success(ctx, gin.H{"good": newGood}, "创建成功")
}

// 修改商品信息
func UpdateGood(ctx *gin.Context) {
	DB := common.GetDB()

	var requestGood model.Good
	err := ctx.ShouldBind(&requestGood)

	if err != nil {
		log.Printf("修改商品信息数据绑定失败  error : %v", err.Error())
		return
	}

	// 获取参数
	name := requestGood.Name
	info := requestGood.Info
	price := requestGood.Price
	quantity := requestGood.Quantity
	goodId := requestGood.GoodId
	pic := requestGood.Picture

	if !utils.VerifyGoodName(name) {
		response.Response(ctx, http.StatusUnprocessableEntity, nil, "店铺名不合法,店铺名长度为1-20字节英文或中文")
		return
	}

	if len(info) > 200 {
		response.Response(ctx, http.StatusUnprocessableEntity, nil, "简介不合法,简介长度为0-200字节")
		return
	}

	if quantity < 0 || quantity > 99 {
		response.Response(ctx, http.StatusUnprocessableEntity, nil, "商品数量不合法，商品数量0-99")
		return
	}

	var updateGood model.Good
	DB.Find(&updateGood, "good_id = ?", goodId)

	// 修改保存
	updateGood.Name = name
	updateGood.Info = info
	updateGood.Price = price
	updateGood.Quantity = quantity
	updateGood.Picture = pic

	DB.Save(&updateGood)

	// 返回修改后的good
	response.Success(ctx, gin.H{"good": updateGood}, "修改成功")
}

// 上传商品图片
func UploadPic(ctx *gin.Context) {
	file, err := ctx.FormFile("img")
	if err != nil {
		response.Failure(ctx, nil, "文件上传失败")
		return
	}

	// 上传文件格式不正确
	fileExt := strings.ToLower(path.Ext(file.Filename))
	if fileExt != ".png" && fileExt != ".jpg" && fileExt != ".jpeg" {
		response.Failure(ctx, nil, "上传失败,只允许png,jpg,jpeg文件")
		return
	}

	// 根据上传时间和文件名称生成对应的md5码作为图片的唯一标识符
	fileName := utils.StrToMd5(fmt.Sprintf("%s%s", file.Filename, time.Now().String()))
	// 保存文件的位置
	filepath := fmt.Sprintf("%s%s%s", filePath, fileName, fileExt)
	// 保存失败
	saveErr := ctx.SaveUploadedFile(file, filepath)
	if saveErr != nil {
		response.Response(ctx, http.StatusInternalServerError, nil, "服务端硬盘已满")
	}

	// 返回值
	response.Success(ctx, gin.H{"hash": fileName}, "上传图片成功")
}

// 商品列表页
func GetGoodsList(ctx *gin.Context) {
	DB := common.GetDB()
	tokenString := ctx.GetHeader("Authorization")

	var goods []model.Good
	if tokenString == "" {
		DB.Find(&goods)
	} else {
		// 解析token
		token, claims, err := common.ParseToken(tokenString)

		// 出现错误或者token无效
		if err != nil || !token.Valid {
			response.Response(ctx, http.StatusUnauthorized, nil, "权限不足")
			return
		}

		// 验证通过后获取claims中的userID
		userId := claims.UserId

		DB.Find(&goods, "user_id = ?", userId)
	}

	response.Success(ctx, gin.H{"goods": goods}, "查询成功")
}

// 商品详情页
func GetGoodDetail(ctx *gin.Context) {
	goodId := ctx.Param("goodId")

	var good model.Good

	log.Println(goodId)

	DB := common.GetDB()
	DB.Find(&good, "good_id = ?", goodId)

	if good.UserId == 0 {
		response.Failure(ctx, nil, "该商品不存在")
		return
	}

	// 返回
	response.Success(ctx, gin.H{"good": good}, "查询成功")
}
