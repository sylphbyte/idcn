package idcn

import (
	"testing"
)

// TestIsValid 测试身份证验证
func TestIsValid(t *testing.T) {
	// 使用生成的有效身份证进行测试
	t.Run("生成的18位有效", func(t *testing.T) {
		for i := 0; i < 10; i++ {
			id := FakeId()
			if !IsValid(id) {
				t.Errorf("FakeId() 生成的身份证验证失败: %s", id)
			}
		}
	})

	t.Run("生成的15位有效", func(t *testing.T) {
		for i := 0; i < 10; i++ {
			id := FakeRequireId(false, "", "", -1)
			if !IsValid(id) {
				t.Errorf("生成的15位身份证验证失败: %s", id)
			}
		}
	})

	// 测试无效情况（构造的无效数据）
	invalidTests := []struct {
		name string
		id   string
	}{
		{"长度错误", "12345"},
		{"格式错误-字母", "12345678901234567A"},
	}

	for _, tt := range invalidTests {
		t.Run(tt.name, func(t *testing.T) {
			if IsValid(tt.id) {
				t.Errorf("IsValid(%q) = true, want false", tt.id)
			}
		})
	}

	// 测试校验位错误
	t.Run("校验位错误", func(t *testing.T) {
		id := FakeId()
		// 修改最后一位使校验位错误
		lastChar := id[17]
		var wrongChar byte
		if lastChar == '0' {
			wrongChar = '1'
		} else {
			wrongChar = '0'
		}
		wrongId := id[:17] + string(wrongChar)
		if IsValid(wrongId) {
			t.Errorf("修改校验位后应该无效: %s", wrongId)
		}
	})
}

// TestGetArea 测试获取地区信息
func TestGetArea(t *testing.T) {
	tests := []struct {
		name     string
		code     string
		province string
		city     string
		district string
		wantErr  bool
	}{
		{"北京东城", "110101", "北京市", "北京市", "东城区", false},
		{"上海", "310000", "上海市", "", "", false},
		{"无效码", "999999", "", "", "", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			area, err := GetArea(tt.code)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetArea(%q) error = %v, wantErr %v", tt.code, err, tt.wantErr)
				return
			}
			if err == nil {
				if area.Province != tt.province {
					t.Errorf("Province = %q, want %q", area.Province, tt.province)
				}
				if area.City != tt.city {
					t.Errorf("City = %q, want %q", area.City, tt.city)
				}
				if area.District != tt.district {
					t.Errorf("District = %q, want %q", area.District, tt.district)
				}
			}
		})
	}

	// 使用生成的身份证测试地区获取
	t.Run("生成身份证的地区", func(t *testing.T) {
		id := FakeRequireId(true, "广东省", "", -1)
		area, err := GetArea(id[:6])
		if err != nil {
			t.Errorf("GetArea failed: %v", err)
			return
		}
		if area.Province != "广东省" {
			t.Errorf("Province = %q, want 广东省", area.Province)
		}
		t.Logf("地区: %s %s %s", area.Province, area.City, area.District)
	})
}

// TestGetInfo 测试获取身份证信息
func TestGetInfo(t *testing.T) {
	// 使用生成的身份证测试
	t.Run("18位身份证信息", func(t *testing.T) {
		id := FakeRequireId(true, "北京市", "19900815", 1)
		info, err := GetInfo(id)
		if err != nil {
			t.Fatalf("GetInfo failed: %v", err)
		}

		// 验证地区包含北京
		if info.Area.Province != "北京市" {
			t.Errorf("Province = %q, want 北京市", info.Area.Province)
		}

		// 验证生日
		if info.Info.Birthday != "1990-08-15" {
			t.Errorf("Birthday = %q, want 1990-08-15", info.Info.Birthday)
		}

		// 验证性别 (男性)
		if info.Info.Gender != 1 {
			t.Errorf("Gender = %d, want 1", info.Info.Gender)
		}

		// 验证身份证号
		if info.CardNo != id {
			t.Errorf("CardNo = %q, want %q", info.CardNo, id)
		}

		t.Logf("Info: %+v", info)
	})

	t.Run("女性身份证", func(t *testing.T) {
		id := FakeRequireId(true, "", "", 0)
		info, err := GetInfo(id)
		if err != nil {
			t.Fatalf("GetInfo failed: %v", err)
		}
		if info.Info.Gender != 0 {
			t.Errorf("Gender = %d, want 0 (女)", info.Info.Gender)
		}
	})
}

// TestUpgradeTo18 测试15位升级18位
func TestUpgradeTo18(t *testing.T) {
	t.Run("正常升级", func(t *testing.T) {
		// 生成15位身份证
		id15 := FakeRequireId(false, "", "", -1)
		if len(id15) != 15 {
			t.Fatalf("生成的15位身份证长度错误: %d", len(id15))
		}

		// 升级为18位
		id18, err := UpgradeTo18(id15)
		if err != nil {
			t.Fatalf("UpgradeTo18 failed: %v", err)
		}

		if len(id18) != 18 {
			t.Errorf("升级后长度 = %d, want 18", len(id18))
		}

		// 验证升级后的身份证有效
		if !IsValid(id18) {
			t.Errorf("升级后的身份证无效: %s", id18)
		}

		t.Logf("15位: %s -> 18位: %s", id15, id18)
	})

	t.Run("已经是18位", func(t *testing.T) {
		id := FakeId()
		result, err := UpgradeTo18(id)
		if err != nil {
			t.Fatalf("UpgradeTo18 failed: %v", err)
		}
		if result != id {
			t.Errorf("18位身份证升级后应该不变")
		}
	})

	t.Run("无效长度", func(t *testing.T) {
		_, err := UpgradeTo18("12345")
		if err == nil {
			t.Error("无效长度应该返回错误")
		}
	})
}

// TestCalcCheckBit 测试校验位计算
func TestCalcCheckBit(t *testing.T) {
	// 使用生成的身份证验证校验位计算
	for i := 0; i < 10; i++ {
		id := FakeId()
		body := id[:17]
		expectedCheckBit := string(id[17])
		if id[17] == 'X' {
			expectedCheckBit = "x"
		}

		got := calcCheckBit(body)
		if got != expectedCheckBit {
			t.Errorf("calcCheckBit(%q) = %q, want %q", body, got, expectedCheckBit)
		}
	}
}

// TestGetInfo15 测试15位身份证获取信息
func TestGetInfo15(t *testing.T) {
	id15 := FakeRequireId(false, "上海市", "1985", 0)
	info, err := GetInfo(id15)
	if err != nil {
		t.Fatalf("GetInfo failed: %v", err)
	}

	// 15位升级为18位
	if len(info.CardNo) != 18 {
		t.Errorf("CardNo 长度 = %d, want 18", len(info.CardNo))
	}

	// 验证地区
	if info.Area.Province != "上海市" {
		t.Errorf("Province = %q, want 上海市", info.Area.Province)
	}

	// 验证性别
	if info.Info.Gender != 0 {
		t.Errorf("Gender = %d, want 0 (女)", info.Info.Gender)
	}

	t.Logf("15位: %s -> Info: %+v", id15, info)
}

// TestHistoryData 测试历史数据源
func TestHistoryData(t *testing.T) {
	// 崇文区在2010年已撤销，应该从历史数据获取
	area, err := GetArea("110103") // 崇文区
	if err != nil {
		t.Fatalf("110103 获取失败: %v", err)
	}

	// 验证数据
	if area.Province != "北京市" {
		t.Errorf("Province = %q, want 北京市", area.Province)
	}
	if area.City != "北京市" {
		t.Errorf("City = %q, want 北京市", area.City)
	}
	if area.District != "崇文区" {
		t.Errorf("District = %q, want 崇文区", area.District)
	}

	t.Logf("110103 = %+v", area)
}

// TestFakeId 测试生成假身份证号
func TestFakeId(t *testing.T) {
	for i := 0; i < 10; i++ {
		id := FakeId()
		if len(id) != 18 {
			t.Errorf("FakeId() length = %d, want 18", len(id))
		}
		if !IsValid(id) {
			t.Errorf("FakeId() generated invalid ID: %s", id)
		}
		t.Logf("生成的身份证: %s", id)
	}
}

// TestFakeRequireId 测试按条件生成假身份证号
func TestFakeRequireId(t *testing.T) {
	tests := []struct {
		name       string
		isEighteen bool
		address    string
		birthday   string
		sex        int
	}{
		{"随机18位", true, "", "", -1},
		{"随机15位", false, "", "", -1},
		{"指定北京市", true, "北京市", "", -1},
		{"指定深圳市", true, "深圳市", "", -1},
		{"指定东城区", true, "东城区", "", -1},
		{"指定生日年份", true, "", "1990", -1},
		{"指定生日年月", true, "", "199008", -1},
		{"指定完整生日", true, "", "19900815", -1},
		{"指定男性", true, "", "", 1},
		{"指定女性", true, "", "", 0},
		{"完整条件", true, "广东省", "1985", 1},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			id := FakeRequireId(tt.isEighteen, tt.address, tt.birthday, tt.sex)
			expectedLen := 18
			if !tt.isEighteen {
				expectedLen = 15
			}

			if len(id) != expectedLen {
				t.Errorf("长度 = %d, want %d", len(id), expectedLen)
			}

			if !IsValid(id) {
				t.Errorf("生成的身份证无效: %s", id)
			}

			// 验证性别
			if tt.sex >= 0 && tt.isEighteen {
				order := int(id[16] - '0')
				actualSex := order % 2
				if tt.sex == 0 && actualSex != 0 {
					t.Errorf("性别不匹配: 期望女性, 顺序码=%d", order)
				}
				if tt.sex == 1 && actualSex != 1 {
					t.Errorf("性别不匹配: 期望男性, 顺序码=%d", order)
				}
			}

			t.Logf("%s: %s", tt.name, id)
		})
	}
}
