package utils

import (
	"crypto/md5"
	"encoding/hex"
	"errors"
	"fmt"
	"github.com/jinzhu/gorm"
	"itStudioTB/model"
	"math/rand"
	"os"
	"time"
	"unicode"
)

// 随机生成 user + n位随机数字
func RandomString(n int) string {
	var letters = []byte("1234567890")
	result := make([]byte, n)

	rand.Seed(time.Now().Unix())
	for i := range result {
		result[i] = letters[rand.Intn(len(letters))]
	}

	return "user" + string(result)
}

// 验证账号是否合法 长度6-20字符 只能用英文和数字
func VerifyAccount(account string) bool {
	// 若长度不符合直接返回false
	if len(account) < 6 || len(account) > 20 {
		return false
	} else {
		flag := true
		// 判断每一个字符 不是英文并且也不是数字 则不合法
		for _, r := range account {
			if !unicode.IsLetter(r) && !unicode.IsNumber(r) {
				flag = false
			}
		}
		return flag
	}
}

// 验证密码是否合法 长度 8-16字符 只能英文 数字 特殊字符
func VerifyPassword(password string) bool {
	if len(password) < 8 || len(password) > 16 {
		return false
	} else {
		flag := true
		// 判断密码是否由 数字 英文 特殊字符
		for _, r := range password {
			if !unicode.IsLetter(r) && !unicode.IsNumber(r) && (r < 33 || r > 47) {
				flag = false
			}
		}
		return flag
	}
}

// 验证该账号是否已经存在
func IsAccountExist(db *gorm.DB, account string) bool {
	var user model.User
	db.Find(&user, "account = ?", account)

	if user.UserId != 0 {
		return true
	}
	return false
}

// 验证密码是否正确
func IsRightPassword(db *gorm.DB, account string, password string) bool {
	var user model.User
	db.Find(&user, "account = ?", account)
	if user.Password == password {
		return true
	}
	return false
}

// 验证商品名是否合法，商品名只能由1-20位中文或英文构成
func VerifyGoodName(name string) bool {
	if len(name) < 1 || len(name) > 20 {
		return false
	} else {
		flag := true
		// 店铺名只能由中文和英文构成
		for _, r := range name {
			if !unicode.IsLetter(r) {
				flag = false
			}
		}
		return flag
	}
}

// 验证该用户是否已经存在该商品名的商品了
func IsGoodExist(db *gorm.DB, userId uint, name string) bool {
	var good model.Good
	db.Find(&good, "user_id = ? AND name = ?", userId, name)

	// 不存在则返回false
	if good.GoodId != 0 {
		return true
	}
	return false
}

// 生成MD5码
func StrToMd5(str string) string {
	h := md5.New()
	h.Write([]byte(str))
	return hex.EncodeToString(h.Sum(nil))
}

// 判断是否存在该文件
func isFileExist(filePath string) bool {
	_, err := os.Stat(filePath)
	if err == nil {
		return true
	}
	if os.IsNotExist(err) {
		return false
	}
	return false
}

// 返回文件名（添加文件类型）
func GetFilePath(path string) (string,error) {
	filePathPng := fmt.Sprintf("%s%s",path,".png")
	filePathJpg := fmt.Sprintf("%s%s",path,".jpg")
	filePathJpeg := fmt.Sprintf("%s%s",path,".jpeg")

	if isFileExist(filePathPng) {
		return filePathPng,nil
	}else if isFileExist(filePathJpg) {
		return filePathJpg,nil
	}else if isFileExist(filePathJpeg) {
		return filePathJpeg,nil
	}else{
		return "",errors.New("不存在该文件")
	}
}