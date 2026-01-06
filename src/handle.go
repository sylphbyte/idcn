package src

// AreaInfo 地区信息
type AreaInfo struct {
	Province string `json:"province"`
	City     string `json:"city,omitempty"`
	District string `json:"district,omitempty"`
}

// GetHistoryAreaInfo 从历史数据获取完整地区信息
func GetHistoryAreaInfo(code uint32) *AreaInfo {
	if info, ok := _historyCodeData[code]; ok {
		return &info
	}
	return nil
}

// GetHistoryAddress 从历史数据获取地址名称
func GetHistoryAddress(code uint32) string {
	info := GetHistoryAreaInfo(code)
	if info == nil {
		return ""
	}

	// 返回最具体的地区名
	// 区县级
	if info.District != "" {
		return info.District
	}
	// 市级
	if info.City != "" {
		return info.City
	}
	// 省级
	return info.Province
}

func AddressCode() map[uint32]string {
	return _addressCode
}

func HistoryMapping() map[uint32]AreaInfo {
	return _historyCodeData
}
