package tyc

import (
	"fmt"
	"net/url"
	"time"

	"github.com/hyacinthus/x/xerr"
	"github.com/hyacinthus/x/xim"
	"github.com/jinzhu/copier"
	"github.com/jinzhu/gorm"

	"github.com/levigross/grequests"
)

// FindCompany 获取已存公司
func (c *Client) FindCompany(id int64) (*Company, error) {
	resp := new(Company)
	err := c.db.First(resp, "tyc_id = ?", id).Error
	if err != nil {
		return nil, err
	}
	return resp, nil
}

// FindCompanyByName 通过名字获取已存公司
func (c *Client) FindCompanyByName(name string) (*Company, error) {
	var un = new(CompanyUsedName)
	err := c.db.First(un, "name = ?", name).Error
	if err != nil {
		return nil, err
	}
	return c.FindCompany(un.TycID)
}

// FetchCompany 从天眼查获取数据
// 若有一个月内的缓存，会自动使用缓存
// 过期数据会自动被存入历史表
// NOTE: 如果传入的公司名是一个历史名称，也许会我们查不出来，但天眼查返回的数据我们有，观察中。
func (c *Client) FetchCompany(name string) (*Company, error) {
	// 检查输入
	if len(name) == 0 {
		return nil, ErrorEmptyName
	}
	// 在数据库中查找
	var found bool
	var old = new(Company)
	var err error
	old, err = c.FindCompanyByName(name)
	if gorm.IsRecordNotFoundError(err) {
		// 没有找到，什么都不做，继续
	} else if err != nil {
		return nil, fmt.Errorf("查询旧记录出错 %w", err)
	} else {
		found = true
		if old.TycUpdatedAt.Add(c.exp).After(time.Now()) {
			// 没过期直接返回
			return old, nil
		}
		// 已过期进行后续处理
	}
	// 请求
	params := make(url.Values)
	params.Set("name", name)
	data, err := grequests.Get("http://open.api.tianyancha.com/services/open/ic/baseinfo/2.0",
		&grequests.RequestOptions{
			Headers: map[string]string{
				"Authorization": c.config.Token,
			},
			Params:     map[string]string{"name": name},
			HTTPClient: c.httpc,
		})
	if err != nil {
		return nil, fmt.Errorf("请求天眼查网络出错 %w", err)
	}
	var resp = new(CompanyResp)
	err = data.JSON(resp)
	if err != nil {
		return nil, fmt.Errorf("展开天眼查 JSON 数据出错 %w", err)
	}
	// 300000 errCode是未找到相关公司
	if resp.ErrorCode == 300000 {
		return nil, ErrorNotFound
	}
	if resp.ErrorCode != 0 {
		return nil, xerr.Newf(500, "TycError", "请求[%s]企业天眼查企业数据报错: %d %s",
			name, resp.ErrorCode, resp.Reason)
	}

	// 平面化，追加时间和哈希
	resp.Result.Plane()
	resp.Result.TycUpdatedAt = time.Now()
	resp.Result.Hash, err = resp.Result.GenHash()
	if err != nil {
		return nil, fmt.Errorf("生成哈希出错 %w", err)
	}

	if !found {
		// TODO: 一个额外的检查，检查公司改名后查询旧名字发生的事情，观察结束后请完善此处代码
		// 名字在数据库中没找到，再次找新记录的 ID ，如果找到，说明改名了
		old, err = c.FindCompany(resp.Result.TycID)
		if gorm.IsRecordNotFoundError(err) {
			// 找不到完全正常，继续后续步骤
		} else if err != nil {
			return nil, fmt.Errorf("依靠id寻找旧记录出错 %w", err)
		} else {
			// 命中特殊情况 通知企业微信
			xim.Warnf("观测报告：从天眼查拉取的公司有了新的名字。输入：%s。库中：%s。最新：%s。", name, old.Name, resp.Result.Name)
			// 标记为找到，后续继续刷新记录
			found = true
			// 保存这个旧名字
			err := c.db.Create(&CompanyUsedName{
				Name:  name,
				TycID: resp.Result.TycID,
			}).Error
			if err != nil {
				return nil, fmt.Errorf("保存名称历史出错 %w", err)
			}
		}
	}

	// 处理历史记录
	if found {
		// 记录没有变化，则刷新时间并返回
		if old.Hash == resp.Result.Hash {
			err := c.db.Model(old).Update("tyc_updated_at", time.Now()).Error
			if err != nil {
				return nil, fmt.Errorf("刷新最后检查时间出错 %w", err)
			}
			return resp.Result, nil
		}
		// 记录有变化，将旧记录移入历史表
		var history = new(CompanyHistory)
		err := copier.Copy(history, old)
		if err != nil {
			return nil, fmt.Errorf("复制到历史记录出错 %w", err)
		}
		err = c.db.Create(history).Error
		if err != nil {
			return nil, fmt.Errorf("保存历史记录出错 %w", err)
		}
		err = c.db.Delete(old).Error
		if err != nil {
			return nil, fmt.Errorf("删除旧记录出错 %w", err)
		}
		// 等待后续保存新记录
	}

	// 保存新的记录，然后返回
	err = c.db.Create(resp.Result).Error
	if err != nil {
		return nil, fmt.Errorf("保存新记录出错 %w", err)
	}
	// 保存名称映射
	un := CompanyUsedName{
		Name:  resp.Result.Name,
		TycID: resp.Result.TycID,
	}
	err = c.db.FirstOrCreate(&un, un).Error
	if err != nil {
		return nil, fmt.Errorf("保存名称映射出错 %w", err)
	}

	return resp.Result, nil
}
