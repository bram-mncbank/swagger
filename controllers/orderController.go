package controller

import (
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

type Request struct {
	CustomerName string  `json:"customerName" example:"Tom Jerry"`
	Items        []Items `json:"items"`
}

type Response struct {
	CustomerName string    `json:"customerName" example:"Tom Jerry"`
	Items        []Items   `json:"items"`
	OrderedAt    time.Time `json:"orderedAt" example:"2019-11-09T21:21:46+07:00"`
	OrderId      int       `json:"orderID" example:"1"`
}

type Orders struct {
	OrderId      int
	CustomerName string
	OrderedAt    time.Time
}

type Items struct {
	ItemId      int    `json:"itemId":omitempty`
	ItemCode    string `json:"itemCode" example:"DJS 4123"`
	Description string `json:"description":omitempty`
	Quantity    int    `json:"quantity" example:"2"`
}

type OrderItems struct {
	ItemId   int
	OrderId  int
	Quantity int
}

type OrderRepository struct {
	orders     []*Orders
	items      []*Items
	orderItems []*OrderItems
}

// CreateItem implements OrderIface
func (o *OrderRepository) CreateItem(c *gin.Context) {
	var requestItem []Items
	var h bool
	if err := c.ShouldBindJSON(&requestItem); err != nil {
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}
	fmt.Printf("Create item \n")
	for _, rit := range requestItem {
		h = false
		fmt.Printf("Sent item %v\n", rit.ItemCode)
		for _, it := range o.items {
			fmt.Printf("Cmp: %v - %v\n", it.ItemCode, rit.ItemCode)
			if it.ItemCode == rit.ItemCode {
				h = true
				fmt.Printf("Update Item %v\n", rit.ItemCode)
				it.Description = rit.Description
				it.Quantity = it.Quantity + rit.Quantity
				break
			}
		}
		if h == false {
			fmt.Println("Add _ Item")
			rit = Items{
				ItemId:      getLastItemId(o.items) + 1,
				ItemCode:    rit.ItemCode,
				Description: rit.Description,
				Quantity:    rit.Quantity,
			}
			o.items = append(o.items, &rit)
		}
	}
	c.JSON(http.StatusOK, o.items)
}

// Create One Order
// @Summary Create New Order
// @Description Create New Order
// @Tags Orders
// @Accept  */*
// @Produce  json
// @Param data body Request true "Order"
// @Success 200 {object} Response
// @Failure 500 {object} string "error"
// @Router /orders [post]
func (o *OrderRepository) CreateOrder(c *gin.Context) {
	// var item *Items
	var request Request
	var response Response
	var newOrder Orders
	var newItems []Items

	fmt.Printf("Creating order...\n")
	if err := c.ShouldBindJSON(&request); err != nil {
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}
	fmt.Printf("ID: %d\n", getLastOrderId(o.orders))
	fmt.Printf("CustomerName: %s\n", request.CustomerName)
	currentId := getLastOrderId(o.orders) + 1

	newOrder.CustomerName = request.CustomerName
	newOrder.OrderId = currentId
	newOrder.OrderedAt = time.Now()

	for _, ritems := range request.Items {
		for _, it := range o.items {
			if it.ItemId == ritems.ItemId {
				it.Quantity = it.Quantity - ritems.Quantity
				if it.Quantity < 0 {
					it.Quantity = 0
				}
				newItems = append(newItems, Items{
					ItemId:      ritems.ItemId,
					Quantity:    ritems.Quantity,
					ItemCode:    it.ItemCode,
					Description: it.Description,
				})
				o.orderItems = append(o.orderItems, &OrderItems{
					ItemId:   ritems.ItemId,
					OrderId:  currentId,
					Quantity: ritems.Quantity,
				})
			}
		}
	}
	o.orders = append(o.orders, &newOrder)
	response.CustomerName = newOrder.CustomerName
	response.OrderId = currentId
	response.OrderedAt = newOrder.OrderedAt
	response.Items = newItems
	c.JSON(http.StatusOK, response)
}

// Delete Order
// @Summary Delete Order
// @Description Delete Order
// @Tags Orders
// @Accept  */*
// @Produce  json
// @Success 200 {object} Response
// @Failure 500 {object} string "error"
// @Router /order/:id [delete]
func (o *OrderRepository) DeleteOrder(c *gin.Context) {
	id := c.Params.ByName("orderId")
	orderId, _ := strconv.Atoi(id)

	for i, oi := range o.orderItems {
		if oi.OrderId == orderId {
			o.orderItems = append(o.orderItems[:i], o.orderItems[i+1:]...)
			break
		}
	}
	c.JSON(http.StatusOK, nil)
}

// Get Order
// @Summary Get All Order
// @Description Get All Order
// @Tags Orders
// @Accept  */*
// @Produce  json
// @Success 200 {object} Response
// @Failure 500 {object} string "error"
// @Router /orders [get]
func (o *OrderRepository) GetOrders(c *gin.Context) {
	var r Response
	var response []Response
	var listItems []Items
	c.Header("Content-Type", "application/json")
	for _, od := range o.orders {
		r.CustomerName = od.CustomerName
		r.OrderedAt = od.OrderedAt
		r.OrderId = od.OrderId
		listItems = nil
		for _, oi := range o.orderItems {
			if oi.OrderId != od.OrderId {
				continue
			}
			for _, it := range o.items {
				if it.ItemId != oi.ItemId {
					continue
				}
				listItems = append(listItems, *it)
			}
		}
		r.Items = listItems
		response = append(response, r)
	}
	c.JSON(http.StatusOK, response)
}

// Update Order
// @Summary Update Order
// @Description Update Order
// @Tags Orders
// @Accept  */*
// @Produce  json
// @Success 200 {object} Response
// @Failure 500 {object} string "error"
// @Router /order/:id [put]
func (o *OrderRepository) UpdateOrder(c *gin.Context) {
	panic("unimplemented")
}

type OrderIface interface {
	CreateOrder(c *gin.Context)
	GetOrders(c *gin.Context)
	UpdateOrder(c *gin.Context)
	DeleteOrder(c *gin.Context)
	CreateItem(c *gin.Context)
}

func OrderService(orders []*Orders, items []*Items, orderItems []*OrderItems) OrderIface {
	return &OrderRepository{
		orders:     orders,
		items:      items,
		orderItems: orderItems,
	}
}

func getLastOrderId(o []*Orders) int {
	var idx int
	for _, v := range o {
		idx = v.OrderId
	}
	return idx
}

func getLastItemId(i []*Items) int {
	var idx int
	for _, v := range i {
		idx = v.ItemId
	}
	return idx
}
