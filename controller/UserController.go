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

// 文件的默认文件
var filePath = "/Users/MaLiangji/Desktop/itStudioTB/picFile/"

// 登录
func Login(ctx *gin.Context) {
	DB := common.GetDB()
	var requestUser = model.User{}

	// 绑定对象，参数为指针
	err := ctx.ShouldBind(&requestUser)
	if err != nil {
		log.Printf("登录绑定对象失败 error : %v", err.Error())
		return
	}

	// 获取参数
	account := requestUser.Account
	password := requestUser.Password

	// 数据验证
	if !utils.VerifyAccount(account) {
		response.Response(ctx, http.StatusUnprocessableEntity, nil, "账号不合法，账号长度为6-20字符，可以使用英文和数字")
		return
	}

	if !utils.VerifyPassword(password) {
		response.Response(ctx, http.StatusUnprocessableEntity, nil, "密码长度为8-16字符，可以使用英文、数字和特殊字符")
		return
	}

	// 判断账号是否已经存在
	if !utils.IsAccountExist(DB, account) {
		response.Response(ctx, http.StatusUnprocessableEntity, nil, "该账号不存在")
		return
	}

	// 判断密码是否正确
	if !utils.IsRightPassword(DB, account, password) {
		response.Failure(ctx, nil, "密码不正确")
		return
	}

	// 账号密码都正确返回的用户
	var user model.User
	DB.Where("account = ?", account).First(&user)

	// 生成token
	token, err := common.IssueToken(user)
	if err != nil {
		response.Response(ctx, http.StatusInternalServerError, nil, "系统异常,生成token失败")
		log.Printf("token 获取异常 error : %v", err.Error())
		return
	}

	// 登录成功
	response.Success(ctx, gin.H{"token": token}, "登录成功")
}

// 注册
func Register(ctx *gin.Context) {
	DB := common.GetDB()
	var requestUser = model.User{}

	// 绑定数据
	err := ctx.ShouldBind(&requestUser)

	if err != nil {
		log.Printf("注册数据绑定失败 error : %v", err.Error())
		return
	}

	// 获取数据
	account := requestUser.Account
	password := requestUser.Password

	// 数据验证
	if !utils.VerifyAccount(account) {
		response.Response(ctx, http.StatusUnprocessableEntity, nil, "账号不合法，账号长度为6-20字符，可以使用英文和数字")
		return
	}

	if !utils.VerifyPassword(password) {
		response.Response(ctx, http.StatusUnprocessableEntity, nil, "密码不合法，账号长度为8-16字符，可以使用英文、数字和特殊字符")
		return
	}

	// 验证账号是否已经存在
	if utils.IsAccountExist(DB, account) {
		response.Response(ctx, http.StatusBadRequest, nil, "该账号已经存在")
		return
	}

	// 保存新用户到数据库
	newUser := model.User{
		Account:  account,
		Password: password,
		Name:     utils.RandomString(5), //随机生成昵称
	}
	DB.Create(&newUser)

	// 发放token
	token, err := common.IssueToken(newUser)

	if err != nil {
		response.Response(ctx, http.StatusInternalServerError, nil, "服务端出错")
		log.Printf("注册时生成token失败 error : %v", err.Error())
		return
	}

	// 注册成功
	response.Success(ctx, gin.H{"token": token}, "注册成功")
}

// 修改用户信息
func UpdateUser(ctx *gin.Context) {
	DB := common.GetDB()
	var requestUser model.User

	// 数据绑定
	err := ctx.ShouldBind(&requestUser)

	if err != nil {
		log.Printf("更新用户绑定出错 error : %v", err.Error())
		return
	}

	// 获取数据
	name := requestUser.Name
	sex := requestUser.Sex
	info := requestUser.Info

	// 昵称长度 1-20
	if len(name) < 1 || len(name) > 20 {
		response.Response(ctx, http.StatusUnprocessableEntity, nil, "昵称格式不正确，昵称长度为1-20字符")
		return
	}

	// 简介长度0-200
	if len(info) > 200 {
		response.Response(ctx, http.StatusUnprocessableEntity, nil, "个人简介长度太长，应为0-200字符")
		return
	}
	// 通过token获得用户信息
	user, _ := ctx.Get("user")

	// 将该用户取出
	var updateUser model.User
	DB.Find(&updateUser, "user_id = ?", user.(model.User).UserId)

	// 修改该用户的个人信息
	updateUser.Name = name
	updateUser.Sex = sex
	updateUser.Info = info

	// 保存记录
	DB.Save(&updateUser)

	// 更新用户信息成功
	response.Success(ctx, gin.H{"user": updateUser}, "更新用户信息成功")
}

// 获取用户信息
func GetUserData(ctx *gin.Context) {
	user, isExist := ctx.Get("user")

	if isExist {
		response.Success(ctx, gin.H{"user": user}, "获取成功")
	} else {
		return
	}
}

// 上传头像
func UploadHead(ctx *gin.Context) {
	f, err := ctx.FormFile("img")

	// 上传失败
	if err != nil {
		response.Failure(ctx, nil, "上传失败")
		return
	}

	// 上传文件格式不正确
	fileExt := strings.ToLower(path.Ext(f.Filename))
	if fileExt != ".png" && fileExt != ".jpg" && fileExt != ".jpeg" {
		response.Failure(ctx, nil, "上传失败,只允许png,jpg,jpeg文件")
		return
	}

	// 根据上传时间和文件名称生成对应的md5码作为图片的唯一标识符
	fileName := utils.StrToMd5(fmt.Sprintf("%s%s", f.Filename, time.Now().String()))
	// 保存文件的位置
	filepath := fmt.Sprintf("%s%s%s", filePath, fileName, fileExt)
	// 保存失败
	saveErr := ctx.SaveUploadedFile(f, filepath)
	if saveErr != nil {
		response.Response(ctx, http.StatusInternalServerError, nil, "服务端硬盘已满")
	}

	// 保存成功，修改用户头像字段
	user, _ := ctx.Get("user")
	userId := user.(model.User).UserId

	// 更新到数据库
	db := common.GetDB()
	var updateUser model.User
	db.Find(&updateUser, "user_id = ?", userId)
	updateUser.Head = fileName
	db.Save(&updateUser)

	// 返回值
	response.Success(ctx, nil, "上传图片成功")
}

// 显示图片
func ShowPic(ctx *gin.Context) {
	// 文件MD5名
	fileName := ctx.Param("name")
	filePath := fmt.Sprintf("%s%s", filePath, fileName)

	// 获取文件类型
	file, err := utils.GetFilePath(filePath)
	if err != nil {
		response.Failure(ctx, nil, err.Error())
	}

	// 显示图片
	ctx.File(string(file))
}
