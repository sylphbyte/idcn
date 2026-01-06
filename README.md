# idcn

中国身份证号验证器，支持 15/18 位身份证验证、解析、升级。

## 功能

- ✅ 验证身份证号合法性（格式、生日、校验位）
- ✅ 根据前 6 位获取省市区信息
- ✅ 获取身份证完整信息（地区、年龄、生日、性别）
- ✅ 15 位身份证升级为 18 位
- ✅ 生成随机假身份证号（用于测试）
- ✅ 按条件生成假身份证号（可指定地区、生日、性别）

## 安装

```bash
go get github.com/sylphbyte/idcn
```

## 使用

```go
package main

import (
    "fmt"
    "github.com/sylphbyte/idcn"
)

func main() {
    // 验证身份证
    fmt.Println(idcn.IsValid("身份号")) // true

    // 获取地区信息
    area, _ := idcn.GetArea("320325")
    fmt.Printf("%s %s %s\n", area.Province, area.City, area.District)
    // 江苏省 徐州市 邳县

    // 获取完整信息
    info, _ := idcn.GetInfo("身份号")
    // {
    //     "card_no": "身份号",
    //     "area": {"province": "江苏省", "city": "徐州市", "district": "邳县"},
    //     "info": {"age": 46, "birthday": "1979-06-07", "gender": 1} 1男 0女
    // }

    // 15位升级18位
    id18, _ := idcn.UpgradeTo18("身份号")
    fmt.Println(id18) 

    // 生成随机假身份证（用于测试）
    fakeId := idcn.FakeId()
    fmt.Println(fakeId) // 随机18位身份证号

    // 按条件生成假身份证
    // FakeRequireId(是否18位, 地区, 生日, 性别)
    // 性别: 1男 0女 -1随机
    id := idcn.FakeRequireId(true, "北京市", "1990", 1)
    fmt.Println(id) // 北京市某区、1990年出生的男性身份证
}
```

## API

| 函数 | 说明 |
|------|------|
| `IsValid(id string) bool` | 验证身份证号是否合法 |
| `Validate(id string) error` | 验证身份证号，返回具体错误 |
| `GetArea(code string) (*AreaInfo, error)` | 根据前6位获取地区信息 |
| `GetInfo(id string) (*IDCardInfo, error)` | 获取身份证完整信息 |
| `UpgradeTo18(id string) (string, error)` | 15位升级为18位 |
| `FakeId() string` | 生成随机18位假身份证号 |
| `FakeRequireId(isEighteen bool, address, birthday string, sex int) string` | 按条件生成假身份证号 |

### FakeRequireId 参数说明

| 参数 | 类型 | 说明 |
|------|------|------|
| `isEighteen` | bool | true 生成18位，false 生成15位 |
| `address` | string | 地区名称，如"北京市"、"深圳市"、"东城区"，空则随机 |
| `birthday` | string | 生日，支持"1990"、"199008"、"19900815"，空则随机 |
| `sex` | int | 1男 0女 -1随机 |

## 致谢

本项目的行政区划数据和部分设计思路借鉴自 [guanguans/id-validator](https://github.com/guanguans/id-validator)，感谢原作者的开源贡献。

## License

MIT
