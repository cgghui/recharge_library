package fulu

import (
	"context"
	"errors"
	"fmt"
	"testing"
)

func TestNewApp(t *testing.T) {

	app := NewApp(
		OpenApiURL,
		"3+UltK9rfqxYQhzDWxgxsb+EHPK9Va/vO18FJ8UGsxIQ0SC/GFtSRexWciFlnGlO",
		"08f73b22cf9a4a9995dd651f6c578c06",
	)

	//{
	//	ret, err := app.GoodsListGet(&GoodsListGetParam{
	//		FirstCategoryID:  0,
	//		SecondCategoryID: 0,
	//		ThirdCategoryID:  0,
	//		ProductID:        0,
	//		ProductName:      "",
	//		ProductType:      "",
	//		FaceValue:        0,
	//	})
	//	fmt.Println(ret)
	//	fmt.Println(err)
	//}

	//{
	//	ret, err := app.GoodsStockCheck(&GoodsStockCheckParam{
	//		BuyNum:    100,
	//		ProductID: "15988986",
	//	})
	//	fmt.Println(ret)
	//	fmt.Println(err)
	//}

	//{
	//	ret, err := app.GoodsInfoGet(&GoodsInfoGetParam{
	//		ProductID: "16832264",
	//	})
	//	fmt.Println(ret)
	//	fmt.Println(err)
	//}

	{
		ret, err := app.OrderInfoGet(&OrderInfoGetParam{
			CustomerOrderNo: "22072616043839960131",
		})
		fmt.Println(ret)
		fmt.Println(err)
	}

	{
		// 多商品轮询购买
		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()
		channel := app.OrderDirectAddChannel(ctx, []int{1, 2, 3, 4, 5, 6, 7}, &OrderDirectAddParam{
			ChargeAccount: "207408420951111",
			BuyNum:        1,
			ShopType:      "其他",
		})
		for range channel.GetWait() {
			ret := channel.GetData()
			err := channel.GetError()
			if ret == nil && errors.Is(err, ChannelCancel) {
				return
			}
			fmt.Println(*ret.ParentInfo)
			fmt.Println(err)
		}
		fmt.Println("1")
	}

}
