package routers

import (
	controller "swaggo/controllers"
	docs "swaggo/docs"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

var orders []*controller.Orders
var items []*controller.Items
var orderItems []*controller.OrderItems

func StartServer() *gin.Engine {
	orderService := controller.OrderService(orders, items, orderItems)
	router := gin.Default()
	docs.SwaggerInfo.BasePath = "/api/v1"
	router.POST("/orders", orderService.CreateOrder)
	router.GET("/orders", orderService.GetOrders)
	router.PUT("/orders/:orderId", orderService.UpdateOrder)
	router.DELETE("/orders/:orderId", orderService.DeleteOrder)
	router.POST("/items", orderService.CreateItem)
	url := ginSwagger.URL("http://localhost:8181/swagger/doc.json")
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler, url))

	return router
}
