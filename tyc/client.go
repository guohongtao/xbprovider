package tyc

import (
	"net/http"
	"time"

	"github.com/hyacinthus/x/xdb"
	"github.com/hyacinthus/x/xerr"
	"github.com/jinzhu/gorm"
)

// 错误
var (
	ErrorEmptyName = xerr.New(400, "EmptyName", "企业名称不能为空")
	ErrorNotFound  = xerr.New(400, "NotFound", "在天眼查没有找到这个企业名称")
)

// Config 配置
type Config struct {
	DB          xdb.Config // 需要一个独立的数据库
	ExpiredDays int        `default:"30"` // 超过过期时间的数据我们会重新获取
	Token       string
}

// Client 维持一个持久化的 http client ，避免每次都重建
type Client struct {
	db     *gorm.DB
	httpc  *http.Client
	exp    time.Duration
	config Config
}

// NewClient create a 天眼查 client
func NewClient(db *gorm.DB, httpc *http.Client, config Config) *Client {
	go db.AutoMigrate(&Company{}, &CompanyUsedName{}, &CompanyHistory{})
	return &Client{
		db:     db,
		httpc:  httpc,
		exp:    time.Duration(config.ExpiredDays*24) * time.Hour,
		config: config,
	}
}
