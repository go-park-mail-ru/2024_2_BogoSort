package entity

import "github.com/google/uuid"

type Question struct {
	ID               uuid.UUID
	Title            string
	Description      string
	Page             PageType
	TriggerValue     int
	LowerDescription string
	UpperDescription string
	ParentID         uuid.NullUUID
}

type PageType string

const (
	MainPage         PageType = "mainPage"
	AdvertPage       PageType = "advertPage"
	AdvertCreatePage PageType = "advertCreatePage"
	CartPage         PageType = "cartPage"
	CategoryPage     PageType = "categoryPage"
	AdvertEditPage   PageType = "advertEditPage"
	UserPage         PageType = "userPage"
	SellerPage       PageType = "sellerPage"
	SearchPage       PageType = "searchPage"
)

var PageTypeValues = []PageType{
	MainPage,
	AdvertPage,
	AdvertCreatePage,
	CartPage,
	CategoryPage,
	AdvertEditPage,
	UserPage,
	SellerPage,
	SearchPage,
}
