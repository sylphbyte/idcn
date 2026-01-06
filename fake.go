package idcn

import (
	"fmt"
	"math/rand"
	"regexp"
	"strconv"
	"time"

	"github.com/sylphbyte/idcn/src"
)

// FakeId 生成随机假身份证号码（18位）
func FakeId() string {
	return FakeRequireId(true, "", "", -1)
}

// FakeRequireId 按要求生成假身份证号码
// isEighteen: 是否生成18位号码
// address:    省市县三级地区官方全称，如"北京市"、"深圳市"、"东城区"，空字符串则随机
// birthday:   出生日期，支持格式: "2000"、"198801"、"19990101"，空字符串则随机
// sex:        性别，1男 0女，-1随机
func FakeRequireId(isEighteen bool, address string, birthday string, sex int) string {
	// 1. 生成地址码
	addressCode := generateAddressCode(address)
	if addressCode == "" {
		return ""
	}

	// 2. 生成出生日期码
	birthdayCode := generateBirthdayCode(birthday)

	// 3. 生成顺序码
	orderCode := generateOrderCode(sex)

	// 15位身份证
	if !isEighteen {
		// 15位使用6位年月日（去掉世纪）
		return addressCode + birthdayCode[2:] + orderCode
	}

	// 18位身份证
	body := addressCode + birthdayCode + orderCode
	return body + calcCheckBit(body)
}

// generateAddressCode 生成地址码
func generateAddressCode(address string) string {
	addressData := src.AddressCode()

	// 未指定地址，随机选择一个区县级地址码
	if address == "" {
		return getRandomDistrictCode(addressData)
	}

	// 查找匹配的地址码
	var matchedCode uint32
	for code, name := range addressData {
		if name == address {
			matchedCode = code
			break
		}
	}

	if matchedCode == 0 {
		// 地址不存在，随机生成
		return getRandomDistrictCode(addressData)
	}

	// 根据地址码级别处理
	classification := classifyAddressCode(matchedCode)
	switch classification {
	case "province":
		// 省级：随机获取该省下的区县
		pattern := fmt.Sprintf("^%02d\\d{4}$", matchedCode/10000)
		return getRandomCodeByPattern(addressData, pattern)
	case "city":
		// 市级：随机获取该市下的区县
		pattern := fmt.Sprintf("^%04d\\d{2}$", matchedCode/100)
		return getRandomCodeByPattern(addressData, pattern)
	default:
		// 区县级：直接使用
		return fmt.Sprintf("%06d", matchedCode)
	}
}

// classifyAddressCode 判断地址码级别
func classifyAddressCode(code uint32) string {
	// 港澳台特殊处理
	if code/100000 == 8 {
		return "special"
	}
	// 省级：后4位为0000
	if code%10000 == 0 {
		return "province"
	}
	// 市级：后2位为00
	if code%100 == 0 {
		return "city"
	}
	// 区县级
	return "district"
}

// getRandomDistrictCode 随机获取一个区县级地址码
func getRandomDistrictCode(addressData map[uint32]string) string {
	var districtCodes []uint32
	for code := range addressData {
		// 只选择区县级（后两位不为00）
		if code%100 != 0 {
			districtCodes = append(districtCodes, code)
		}
	}

	if len(districtCodes) == 0 {
		return ""
	}

	idx := rand.Intn(len(districtCodes))
	return fmt.Sprintf("%06d", districtCodes[idx])
}

// getRandomCodeByPattern 根据正则模式随机获取地址码
func getRandomCodeByPattern(addressData map[uint32]string, pattern string) string {
	re := regexp.MustCompile(pattern)
	var matched []uint32

	for code := range addressData {
		codeStr := fmt.Sprintf("%06d", code)
		// 只匹配区县级（后两位不为00）
		if re.MatchString(codeStr) && code%100 != 0 {
			matched = append(matched, code)
		}
	}

	if len(matched) == 0 {
		// 没有匹配的区县，返回随机
		return getRandomDistrictCode(addressData)
	}

	idx := rand.Intn(len(matched))
	return fmt.Sprintf("%06d", matched[idx])
}

// generateBirthdayCode 生成出生日期码
func generateBirthdayCode(birthday string) string {
	now := time.Now()

	var year, month, day int

	switch len(birthday) {
	case 8: // YYYYMMDD
		year, _ = strconv.Atoi(birthday[0:4])
		month, _ = strconv.Atoi(birthday[4:6])
		day, _ = strconv.Atoi(birthday[6:8])
	case 6: // YYYYMM
		year, _ = strconv.Atoi(birthday[0:4])
		month, _ = strconv.Atoi(birthday[4:6])
		day = randRange(1, 28)
	case 4: // YYYY
		year, _ = strconv.Atoi(birthday)
		month = randRange(1, 12)
		day = randRange(1, 28)
	default:
		// 随机生成 1950 到去年之间的日期（确保不会是未来）
		year = randRange(1950, now.Year()-1)
		month = randRange(1, 12)
		day = randRange(1, 28)
	}

	// 校验年份范围，确保不超过当前年份
	if year < 1900 || year > now.Year() {
		year = randRange(1950, now.Year()-1)
	}

	// 如果是当前年份，确保月日不超过今天
	if year == now.Year() {
		if month > int(now.Month()) {
			month = randRange(1, int(now.Month()))
		}
		if month == int(now.Month()) && day > now.Day() {
			day = randRange(1, now.Day())
		}
	}

	// 校验月份
	if month < 1 || month > 12 {
		month = randRange(1, 12)
	}

	// 校验日期
	if day < 1 || day > 28 {
		day = randRange(1, 28)
	}

	return fmt.Sprintf("%04d%02d%02d", year, month, day)
}

// generateOrderCode 生成顺序码
func generateOrderCode(sex int) string {
	// 生成 001-999 之间的随机数
	order := randRange(1, 999)

	// 调整奇偶性以匹配性别
	// 奇数为男，偶数为女
	if sex == 1 && order%2 == 0 {
		order++
		if order > 999 {
			order = 999
		}
	} else if sex == 0 && order%2 == 1 {
		order++
		if order > 999 {
			order = 998
		}
	}

	return fmt.Sprintf("%03d", order)
}

// randRange 生成 [min, max] 范围内的随机数
func randRange(min, max int) int {
	if min >= max {
		return min
	}
	return rand.Intn(max-min+1) + min
}

// init 初始化随机种子
func init() {
	rand.Seed(time.Now().UnixNano())
}
