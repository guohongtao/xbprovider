package tyc

import (
	"time"

	"github.com/hyacinthus/x/model"
	"github.com/mitchellh/hashstructure"
)

// CompanyRest 天眼查返回结构
type CompanyResp struct {
	Result    *Company `json:"result"`
	Reason    string   `json:"reason"`
	ErrorCode int      `json:"error_code"`
}

// Company 天眼查返回企业数据
type Company struct {
	TycID                     int64     `json:"id" gorm:"primary_key"`                   // 企业天眼查ID
	Name                      string    `json:"name" gorm:"unique_index"`                // 企业名
	TycUpdatedAt              time.Time `json:"-" hash:"-"`                              // 我们追加的拉取时间
	Hash                      uint64    `json:"-" hash:"-"`                              // 用来判断数据是否变化
	StaffNumRange             string    `json:"staffNumRange" gorm:"type:varchar(200)"`  // 人员规模
	FromTimeUnix              int64     `json:"fromTime"`                                // 经营开始时间（毫秒数）
	LegalPersonType           int       `json:"type"`                                    // 法人类型，1 人 2 公司
	BondName                  string    `json:"bondName"`                                // 股票名
	IsMicroEnt                int       `json:"isMicroEnt"`                              // 是否是小微企业 0不是 1是
	UsedBondName              string    `json:"usedBondName"`                            // 股票曾用名
	RegNumber                 string    `json:"regNumber"`                               // 注册号
	PercentileScore           int       `json:"percentileScore"`                         // 企业评分：万分制
	RegCapitalText            string    `json:"regCapital"`                              // 注册资本
	RegInstitute              string    `json:"regInstitute"`                            // 登记机关
	RegLocation               string    `json:"regLocation"`                             // 注册地址
	Industry                  string    `json:"industry"`                                // 行业
	ApprovedTimeUnix          int64     `json:"approvedTime"`                            // 核准时间（毫秒数）
	SocialStaffNum            int       `json:"socialStaffNum"`                          // 参保人数
	Tags                      string    `json:"tags"`                                    // 企业标签
	TaxNumber                 string    `json:"taxNumber"`                               // 纳税人识别号
	BusinessScope             string    `json:"businessScope" gorm:"type:varchar(4096)"` // 经营范围
	NameEn                    string    `json:"property3"`                               // 英文名
	Alias                     string    `json:"alias"`                                   // 简称
	OrgNumber                 string    `json:"orgNumber"`                               // 组织机构代码
	RegStatus                 string    `json:"regStatus"`                               // 企业状态
	EstablishTimeUnix         int64     `json:"estiblishTime"`                           // 成立日期（毫秒数） 天眼查写错了 我们纠正
	BondTypeText              string    `json:"bondType"`                                // 股票类型
	LegalPersonName           string    `json:"legalPersonName"`                         // 法人
	ToTimeUnix                int64     `json:"toTime"`                                  // 经营结束时间（毫秒数）
	ActualCapitalText         string    `json:"actualCapital"`                           // 实收注册资金
	CompanyOrgType            string    `json:"companyOrgType"`                          // 企业类型
	Base                      string    `json:"base"`                                    // 省份简称
	CreditCode                string    `json:"creditCode"`                              // 统一社会信用代码
	HistoryNames              string    `json:"historyNames"`                            // 曾用名
	BondNum                   string    `json:"bondNum"`                                 // 股票号
	RegCapitalCurrencyText    string    `json:"regCapitalCurrency"`                      // 注册资本币种 人民币 美元 欧元 等(暂未使用)
	ActualCapitalCurrencyText string    `json:"actualCapitalCurrency"`                   // 实收注册资本币种 人民币 美元 欧元 等(暂未使用)
	RevokeDateUnix            int64     `json:"revokeDate"`                              // 吊销日期（毫秒数）
	RevokeReason              string    `json:"revokeReason" gorm:"type:varchar(500)"`   // 吊销原因
	CancelDateUnix            int64     `json:"cancelDate"`                              // 注销日期（毫秒数）
	CancelReason              string    `json:"cancelReason" gorm:"type:varchar(500)"`   // 注销原因
	// 嵌入的国民行业分类与在数据库中的展开
	IndustryAll    IndustryAll `json:"industryAll" gorm:"-"`
	IndustryL1Text string      `gorm:"type:varchar(30);default:''"`
	IndustryL2Text string      `gorm:"type:varchar(30);default:''"`
	IndustryL3Text string      `gorm:"type:varchar(30);default:''"`
	IndustryL4Text string      `gorm:"type:varchar(30);default:''"`
}

// IndustryAll 天眼查模型中嵌入的国民行业分类
type IndustryAll struct {
	Category       string `json:"category"`
	CategoryBig    string `json:"categoryBig"`
	CategoryMiddle string `json:"categoryMiddle"`
	CategorySmall  string `json:"categorySmall"`
}

// 收到数据后平面化以存储在数据库
func (c Company) Plane() {
	c.IndustryL1Text = c.IndustryAll.Category
	c.IndustryL2Text = c.IndustryAll.CategoryBig
	c.IndustryL3Text = c.IndustryAll.CategoryMiddle
	c.IndustryL4Text = c.IndustryAll.CategorySmall
}

// GenHash 生成当前的哈希用以对比历史变化
func (c Company) GenHash() (uint64, error) {
	return hashstructure.Hash(c, nil)
}

// IsMicro 是否小微企业 复制时从数字转换为bool
func (c *Company) IsMicro() bool {
	return c.IsMicroEnt == 1
}

// EstablishTime copier
func (c *Company) EstablishTime() *time.Time {
	if c.EstablishTimeUnix == 0 {
		return nil
	}
	establishTime := time.Unix(0, c.EstablishTimeUnix*1e6)
	return &establishTime
}

// FromTime copier
func (c *Company) FromTime() *time.Time {
	if c.FromTimeUnix == 0 {
		return nil
	}
	fromTime := time.Unix(0, c.FromTimeUnix*1e6)
	return &fromTime
}

// ApprovedTime copier
func (c *Company) ApprovedTime() *time.Time {
	if c.ApprovedTimeUnix == 0 {
		return nil
	}
	approvedTime := time.Unix(0, c.ApprovedTimeUnix*1e6)
	return &approvedTime
}

// RevokeDate copier
func (c *Company) RevokeDate() *time.Time {
	if c.RevokeDateUnix == 0 {
		return nil
	}
	revokeDate := time.Unix(0, c.RevokeDateUnix*1e6)
	return &revokeDate
}

// CancelDate copier
func (c *Company) CancelDate() *time.Time {
	if c.CancelDateUnix == 0 {
		return nil
	}
	cancelDate := time.Unix(0, c.CancelDateUnix*1e6)
	return &cancelDate
}

// ToTime copier
func (c *Company) ToTime() *time.Time {
	if c.ToTimeUnix == 0 {
		return nil
	}
	toTime := time.Unix(0, c.ToTimeUnix*1e6)
	return &toTime
}

// Base copier
func (c *Company) ProvinceBase() string {
	return c.Base
}

// RegLocationText copier
func (c *Company) RegLocationText() string {
	return c.RegLocation
}

// CompanyHistory 天眼查企业数据历史
type CompanyHistory struct {
	model.Log
	TycID                     int64     `json:"id" gorm:"index"`                         // 企业天眼查ID
	Name                      string    `json:"name" gorm:"index"`                       // 企业名
	TycUpdatedAt              time.Time `json:"-"`                                       // 我们追加的拉取时间
	StaffNumRange             string    `json:"staffNumRange" gorm:"type:varchar(200)"`  // 人员规模
	FromTimeUnix              int64     `json:"fromTime"`                                // 经营开始时间（毫秒数）
	LegalPersonType           int       `json:"type"`                                    // 法人类型，1 人 2 公司
	BondName                  string    `json:"bondName"`                                // 股票名
	IsMicroEnt                int       `json:"isMicroEnt"`                              // 是否是小微企业 0不是 1是
	UsedBondName              string    `json:"usedBondName"`                            // 股票曾用名
	RegNumber                 string    `json:"regNumber"`                               // 注册号
	PercentileScore           int       `json:"percentileScore"`                         // 企业评分：万分制
	RegCapitalText            string    `json:"regCapital"`                              // 注册资本
	RegInstitute              string    `json:"regInstitute"`                            // 登记机关
	RegLocation               string    `json:"regLocation"`                             // 注册地址
	Industry                  string    `json:"industry"`                                // 行业
	ApprovedTimeUnix          int64     `json:"approvedTime"`                            // 核准时间（毫秒数）
	SocialStaffNum            int       `json:"socialStaffNum"`                          // 参保人数
	Tags                      string    `json:"tags"`                                    // 企业标签
	TaxNumber                 string    `json:"taxNumber"`                               // 纳税人识别号
	BusinessScope             string    `json:"businessScope" gorm:"type:varchar(4096)"` // 经营范围
	NameEn                    string    `json:"property3"`                               // 英文名
	Alias                     string    `json:"alias"`                                   // 简称
	OrgNumber                 string    `json:"orgNumber"`                               // 组织机构代码
	RegStatus                 string    `json:"regStatus"`                               // 企业状态
	EstablishTimeUnix         int64     `json:"estiblishTime"`                           // 成立日期（毫秒数） 天眼查写错了 我们纠正
	BondTypeText              string    `json:"bondType"`                                // 股票类型
	LegalPersonName           string    `json:"legalPersonName"`                         // 法人
	ToTimeUnix                int64     `json:"toTime"`                                  // 经营结束时间（毫秒数）
	ActualCapitalText         string    `json:"actualCapital"`                           // 实收注册资金
	CompanyOrgType            string    `json:"companyOrgType"`                          // 企业类型
	Base                      string    `json:"base"`                                    // 省份简称
	CreditCode                string    `json:"creditCode"`                              // 统一社会信用代码
	HistoryNames              string    `json:"historyNames"`                            // 曾用名
	BondNum                   string    `json:"bondNum"`                                 // 股票号
	RegCapitalCurrencyText    string    `json:"regCapitalCurrency"`                      // 注册资本币种 人民币 美元 欧元 等(暂未使用)
	ActualCapitalCurrencyText string    `json:"actualCapitalCurrency"`                   // 实收注册资本币种 人民币 美元 欧元 等(暂未使用)
	RevokeDateUnix            int64     `json:"revokeDate"`                              // 吊销日期（毫秒数）
	RevokeReason              string    `json:"revokeReason" gorm:"type:varchar(500)"`   // 吊销原因
	CancelDateUnix            int64     `json:"cancelDate"`                              // 注销日期（毫秒数）
	CancelReason              string    `json:"cancelReason" gorm:"type:varchar(500)"`   // 注销原因
}

// CompanyUsedName 公司用过的所有名称
type CompanyUsedName struct {
	model.Log
	Name  string `json:"name" gorm:"index"` // 企业名或曾用名
	TycID int64  `json:"id" gorm:"index"`   // 企业天眼查ID
}
