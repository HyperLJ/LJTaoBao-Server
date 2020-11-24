package main

import (
	"github.com/gin-gonic/gin"
	"itStudioTB/controller"
	"itStudioTB/middleware"
)

func CollectRoute(r *gin.Engine) *gin.Engine  {
	// 用户
	r.POST("/api/login",controller.Login)
	r.POST("/api/register",controller.Register)
	r.PUT("/api/user",middleware.AuthMiddleware(),controller.UpdateUser)
	r.POST("/api/upload/head",middleware.AuthMiddleware(),controller.UploadHead)
	r.GET("/api/user",middleware.AuthMiddleware(),controller.GetUserData)

	// 商品
	r.POST("/api/good",middleware.AuthMiddleware(),controller.CreateGood)
	r.GET("/api/goods",controller.GetGoodsList)
	r.GET("/api/goods/:goodId",controller.GetGoodDetail)
	r.PUT("/api/good",middleware.AuthMiddleware(),controller.UpdateGood)
	r.POST("/api/upload/good",middleware.AuthMiddleware(),controller.UploadPic)

	// 订单
	r.POST("/api/order",middleware.AuthMiddleware(),controller.BuyGood)
	r.GET("/api/order",middleware.AuthMiddleware(),controller.GetOrdersList)

	// 查看图片
	r.GET("/api/img/:name",controller.ShowPic)

	return r
}