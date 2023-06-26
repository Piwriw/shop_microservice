package model

// Category 分类信息
type Category struct {
	BaseModel
	Name             string `gorm:"type:varchar(20);not null"`
	ParentCategoryID int32
	ParentCategory   *Category   `json:"-"`
	SubCategory      []*Category `gorm:"foreignKey:ParentCategoryID;references:ID" json:"subCategory"`
	Level            int32       `gorm:"type:int;not null;default 1;comment:'1 一级类目 2 二级类目 3 三级类目'" json:"level"`
	IsTab            bool        `gorm:"default:false;not null" json:"isTab"`
}

func (c *Category) TableName() string {
	return "category"
}

// Brand 品牌表
type Brand struct {
	BaseModel
	Name string `gorm:"type:varchar(20);not null"`
	Logo string `gorm:"type:varchar(200);not null;default:''"`
}

// GoodsCategoryBrand 手动模式 映射category 和brand
type GoodsCategoryBrand struct {
	BaseModel
	CategoryID int32 `gorm:"type:int;index:idx_category_brand,unique"`
	Category   Category
	BrandID    int32 `gorm:"type:int;index:idx_category_brand,unique"`
	Brand      Brand
}

// Banner 轮播图
type Banner struct {
	BaseModel
	Image string `gorm:"type:varchar(200);not null"`
	Url   string `gorm:"type:varchar(200);not null"`
	Index int32  `gorm:"type:int;default 1;not null"`
}
type Good struct {
	BaseModel
	CategoryID int32 `gorm:"column:category_id;type:int"`
	Category   Category
	BrandID    int32 `gorm:"column:brand_id;type:int"`
	Brand      Brand
	OnSale     bool  `gorm:"column:on_sale;default:false;not null"`
	FavNum     int32 `gorm:"column:fav_num;type:int;default:0;not null"`
	ShipFree   bool  `gorm:"column:ship_free;default:false;not null"`
	IsNew      bool  `gorm:"column:is_new;default:false;not null"`
	// 是否热卖
	IsHot bool `gorm:"column:is_hot;default:false;not null;comment '是否热卖'"`
	// 商品名字
	Name string `gorm:"column:name;type:varchar(50);not null;comment '商品名字'"`
	// 商品序号
	GoodSn         string   `gorm:"column:good_sn;type:varchar(50);not null;comment:'商品序号'"`
	ClickNum       int32    `gorm:"column:click_num;type:int;default 0;not null;comment:'点击数'"`
	SoldNum        int32    `gorm:"column:sold_num;type:int;default 0;not null;comment:'销售数'"`
	MarketPrice    float32  `gorm:"column:market_price;not null;comment:'活动价'"`
	ShopPrice      float32  `gorm:"column:shop_price;not null;comment:'原价'"`
	GoodBrief      string   `gorm:"column:good_brief;type:varchar(100);not null;comment:'商品简介'"`
	Images         GormList `gorm:"column:images;type:varchar(1000);not null;"`
	DescImages     GormList `gorm:"column:desc_images;type:varchar(1000);not null;"`
	GoodFrontImage string   `gorm:"column:good_front_image;type:varchar(200);not null"`
}
