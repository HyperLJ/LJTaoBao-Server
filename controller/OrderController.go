package controller

import (
	"github.com/gin-gonic/gin"
	"itStudioTB/common"
	"itStudioTB/model"
	"itStudioTB/response"
	"log"
	"net/http"
)

// 购买商品
func BuyGood(ctx *gin.Context) {
	// 数据绑定
	var requestOrder model.Order
	err := ctx.ShouldBind(&requestOrder)

	if err != nil {
		log.Printf("购买商品数据绑定出错 error : %v", err.Error())
		return
	}

	// 获取参数
	goodId := requestOrder.GoodId
	goodsCount := requestOrder.GoodsCount

	var good model.Good
	DB := common.GetDB()

	user, _ := ctx.Get("user")
	userId := user.(model.User).UserId

	DB.Find(&good, "good_id = ?", goodId)

	// 判断该商品是否属于该用户，属于则不能购买
	if good.UserId == userId {
		response.Response(ctx, http.StatusUnprocessableEntity, nil, "不能购买自己的商品")
		return
	}

	// 查询是否存在该商品
	if good.GoodId == 0 {
		response.Response(ctx, http.StatusUnprocessableEntity, nil, "不存在该商品")
		return
	}

	// 判断商品数量是否正确
	if good.Quantity < goodsCount {
		response.Response(ctx, http.StatusBadRequest, nil, "数据错误，该商品数量不够")
		return
	}

	// 商品数量减少 保存
	good.Quantity -= goodsCount
	DB.Save(&good)

	// 保存订单
	order := model.Order{
		UserId:     userId,
		GoodId:     goodId,
		GoodsCount: goodsCount,
		GoodsPrice: float32(goodsCount) * good.Price, //计算总价
	}

	DB.Create(&order)

	response.Success(ctx, gin.H{"order": order}, "购买成功")
}

// 历史订单
func GetOrdersList(ctx *gin.Context) {
	user, _ := ctx.Get("user")
	userId := user.(model.User).UserId

	var orders []model.Order

	DB := common.GetDB()
	DB.Find(&orders, "user_id = ?", userId)

	response.Success(ctx, gin.H{"orders": orders}, "查询成功")
}
