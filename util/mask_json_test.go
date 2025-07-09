package util_test

import (
	"encoding/json"
	"go101/util"
	"testing"
)

func TestMaskJson(t *testing.T) {
	// 测试数据和期望结果
	tests := []struct {
		input          []byte
		maskRules      map[string]int
		expectedOutput string
	}{
		{
			// 简单的 JSON 示例
			input: []byte(`{
				"name": "JohnDoe1234567890",
				"email": "john.doe@example.com",
				"phone": "1234567890"
			}`),
			maskRules: map[string]int{
				"name":  14,
				"email": 17,
				"phone": 10,
			},
			expectedOutput: `{
				"name": "J***7890",
				"email": "j***ple.com",
				"phone": "1***"
			}`,
		},
		{
			// 包含嵌套对象和数组的 JSON 示例
			input: []byte(`{
				"name": "JohnDoe1234567890",
				"email": "john.doe@example.com",
				"phone": "1234567890",
				"address": {
					"street": "123 Main St",
					"city": "SomeCity"
				},
				"friends": [
					{
						"name": "JaneDoe123456",
						"email": "jane.doe@example.com"
					},
					{
						"name": "JimDoe987654",
						"email": "jim.doe@example.com"
					}
				]
			}`),
			maskRules: map[string]int{
				"name":   14,
				"email":  17,
				"phone":  10,
				"street": 15,
			},
			expectedOutput: `{
				"name": "J***7890",
				"email": "j***ple.com",
				"phone": "1***",
				"address": {
					"street": "1***in St",
					"city": "SomeCity"
				},
				"friends": [
					{
						"name": "J***3456",
						"email": "j***ple.com"
					},
					{
						"name": "J***7654",
						"email": "j***ple.com"
					}
				]
			}`,
		},
		{
			// 空 JSON 示例
			input:          []byte(`{}`),
			maskRules:      map[string]int{},
			expectedOutput: `{}`,
		},
		{
			// 没有匹配规则的字段示例
			input: []byte(`{
				"username": "testuser",
				"password": "mypassword"
			}`),
			maskRules: map[string]int{
				"username": 11, // 只脱敏 username
			},
			expectedOutput: `{
				"username": "t***r",
				"password": "mypassword"
			}`,
		},
	}

	// 执行测试
	for _, tt := range tests {
		t.Run(string(tt.input), func(t *testing.T) {
			// 调用脱敏函数
			maskedData, err := util.MaskJson(tt.input, tt.maskRules)
			if err != nil {
				t.Fatalf("Error occurred: %v", err)
			}
			expect, _ := normalizeJson([]byte(tt.expectedOutput))
			// 比较输出结果与期望的输出
			if string(maskedData) != string(expect) {
				t.Errorf("Expected %s, but got %s", string(expect), string(maskedData))
			}
		})
	}
}

func normalizeJson(input []byte) ([]byte, error) {
	var jsonObj interface{}
	// 解析 JSON 字符串到 map
	if err := json.Unmarshal(input, &jsonObj); err != nil {
		return nil, err
	}
	// 重新编码为标准的 JSON 格式
	return json.Marshal(jsonObj)
}
