package data

import "valuation/internal/biz"

type PriceRuleData struct {
	data *Data
}

func NewPriceRuleInterface(data *Data) biz.PriceRuleInterface {
	return &PriceRuleData{data: data}
}

// 利用PriceRuleData结构实现PriceRuleInterface接口
func (prd *PriceRuleData) GetRule(cityId uint, curr int) (*biz.PriceRule, error) {

	pr := &biz.PriceRule{}
	result := prd.data.MDB.Where("city_id=? AND start_at <= ? AND end_at > ?", cityId, curr, curr).First(pr)

	if result.Error != nil {
		return nil, result.Error
	}

	return pr, nil

}
