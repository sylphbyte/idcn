// Package idcn 中国身份证号验证器
// 支持 15/18 位身份证验证、解析、升级
package idcn

import (
	"errors"
	"fmt"
	"strconv"
	"time"

	"github.com/sylphbyte/idcn/src"
)

var (
	ErrInvalidLength   = errors.New("身份证号长度必须是15或18位")
	ErrInvalidFormat   = errors.New("身份证号格式不正确")
	ErrInvalidBirthday = errors.New("出生日期不合法")
	ErrInvalidChecksum = errors.New("校验码不正确")
	ErrInvalidAddress  = errors.New("地址码不存在")
)

// IsValid 验证身份证号是否合法
func IsValid(id string) bool {
	return Validate(id) == nil
}

// Validate 验证身份证号，返回具体错误
func Validate(id string) error {
	parsed, err := parseID(id)
	if err != nil {
		return err
	}

	// 验证生日
	if !isValidBirthday(parsed.birthdayCode) {
		return ErrInvalidBirthday
	}

	// 18位需要验证校验码
	if parsed.length == 18 {
		expected := calcCheckBit(id[:17])
		if parsed.checkBit != expected {
			return ErrInvalidChecksum
		}
	}

	return nil
}

// GetArea 根据身份证前6位获取地区信息
func GetArea(code string) (*src.AreaInfo, error) {
	if len(code) < 6 {
		return nil, errors.New("地址码长度不足6位")
	}
	code = code[:6]

	addressCode, err := strconv.ParseUint(code, 10, 32)
	if err != nil {
		return nil, errors.New("地址码格式不正确")
	}

	return getAreaByCode(uint32(addressCode))
}

// GetInfo 获取身份证完整信息
func GetInfo(id string) (*IDCardInfo, error) {
	if err := Validate(id); err != nil {
		return nil, err
	}

	parsed, _ := parseID(id)
	addressCode, _ := strconv.ParseUint(parsed.addressCode, 10, 32)

	// 获取地区信息
	area, err := getAreaByCode(uint32(addressCode))
	if err != nil {
		// 地区信息可能缺失，使用空值
		area = &src.AreaInfo{}
	}

	// 解析生日和性别
	birthday, _ := time.Parse("20060102", parsed.birthdayCode)
	age := calcAge(birthday)

	order, _ := strconv.Atoi(parsed.orderCode)
	gender := 1 // 男
	if order%2 == 0 {
		gender = 0 // 女
	}

	// 18位身份证号
	cardNo := id
	if parsed.length == 15 {
		cardNo, _ = UpgradeTo18(id)
	}

	return &IDCardInfo{
		CardNo: cardNo,
		Area:   *area,
		Info: Info{
			Age:      age,
			Birthday: birthday.Format("2006-01-02"),
			Gender:   gender,
		},
	}, nil
}

// UpgradeTo18 将15位身份证号升级为18位
func UpgradeTo18(id string) (string, error) {
	if len(id) == 18 {
		return id, nil
	}

	if len(id) != 15 {
		return "", ErrInvalidLength
	}

	// 验证15位号码
	parsed, err := parseID(id)
	if err != nil {
		return "", err
	}

	if !isValidBirthday(parsed.birthdayCode) {
		return "", ErrInvalidBirthday
	}

	// 构建17位主体
	body := parsed.addressCode + parsed.birthdayCode + parsed.orderCode

	// 计算校验位
	checkBit := calcCheckBit(body)

	return body + checkBit, nil
}

// parseID 解析身份证号
func parseID(id string) (*parsedID, error) {
	length := len(id)

	if length == 15 {
		return parse15(id)
	}

	if length == 18 {
		return parse18(id)
	}

	return nil, ErrInvalidLength
}

// parse15 解析15位身份证
func parse15(id string) (*parsedID, error) {
	// 验证全是数字
	for _, c := range id {
		if c < '0' || c > '9' {
			return nil, ErrInvalidFormat
		}
	}

	return &parsedID{
		addressCode:  id[0:6],
		birthdayCode: "19" + id[6:12], // 15位默认19xx年
		orderCode:    id[12:15],
		checkBit:     "",
		length:       15,
	}, nil
}

// parse18 解析18位身份证
func parse18(id string) (*parsedID, error) {
	// 前17位必须是数字
	for i := 0; i < 17; i++ {
		c := id[i]
		if c < '0' || c > '9' {
			return nil, ErrInvalidFormat
		}
	}

	// 第18位是数字或X/x
	last := id[17]
	if !(last >= '0' && last <= '9') && last != 'X' && last != 'x' {
		return nil, ErrInvalidFormat
	}

	checkBit := string(last)
	if last == 'X' {
		checkBit = "x"
	}

	return &parsedID{
		addressCode:  id[0:6],
		birthdayCode: id[6:14],
		orderCode:    id[14:17],
		checkBit:     checkBit,
		length:       18,
	}, nil
}

// isValidBirthday 验证生日是否合法
func isValidBirthday(birthdayCode string) bool {
	birthday, err := time.Parse("20060102", birthdayCode)
	if err != nil {
		return false
	}

	// 不能早于1800年
	if birthday.Year() < 1800 {
		return false
	}

	// 不能晚于今天
	if birthday.After(time.Now()) {
		return false
	}

	return true
}

// calcCheckBit 计算18位身份证校验位
// ISO 7064:1983 MOD 11-2
func calcCheckBit(body string) string {
	if len(body) != 17 {
		return ""
	}

	// 加权因子
	weights := []int{7, 9, 10, 5, 8, 4, 2, 1, 6, 3, 7, 9, 10, 5, 8, 4, 2}
	// 校验码对应值
	checkCodes := []string{"1", "0", "x", "9", "8", "7", "6", "5", "4", "3", "2"}

	sum := 0
	for i := 0; i < 17; i++ {
		n := int(body[i] - '0')
		sum += n * weights[i]
	}

	return checkCodes[sum%11]
}

// calcAge 计算年龄
func calcAge(birthday time.Time) int {
	now := time.Now()
	age := now.Year() - birthday.Year()

	// 还没过生日
	if now.Month() < birthday.Month() ||
		(now.Month() == birthday.Month() && now.Day() < birthday.Day()) {
		age--
	}

	return age
}

// getAreaByCode 根据地址码获取地区信息
// 逐级查询：每一级都先查现有数据(AddressCode)，再查历史数据(HistoryData)
func getAreaByCode(code uint32) (*src.AreaInfo, error) {
	area := &src.AreaInfo{}

	// 省级代码
	provinceCode := (code / 10000) * 10000
	// 市级代码
	cityCode := (code / 100) * 100
	// 区县级代码
	districtCode := code

	// 获取省份
	area.Province = getAddressName(provinceCode)

	// 获取城市（直辖市可能没有市级）
	if cityCode != provinceCode {
		area.City = getAddressName(cityCode)
	}

	// 获取区县
	if districtCode != cityCode && districtCode != provinceCode {
		area.District = getAddressName(districtCode)
	}

	// 至少要有省份信息
	if area.Province == "" {
		// 尝试从历史数据获取完整信息作为兜底
		if info := src.GetHistoryAreaInfo(code); info != nil {
			return info, nil
		}
		return nil, fmt.Errorf("地址码 %d 无法识别", code)
	}

	return area, nil
}

// getAddressName 获取地址名称
// 查询顺序：AddressCode -> AddressCodeTimeline -> HistoryData
func getAddressName(code uint32) string {
	// 1. 尝试从主数据源获取（当前有效的行政区划）
	if addr := src.AddressCode()[code]; addr != "" {
		return addr
	}

	// 2. 尝试从时间线获取（历史行政区划变更）
	if timeline := src.GetAddressCodeTimeline(code); len(timeline) > 0 {
		return timeline[0]["address"]
	}

	// 3. 备用：从历史数据获取
	return src.GetHistoryAddress(code)
}
