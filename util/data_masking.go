package util

import "encoding/json"

// 脱敏函数
func MaskJson(data []byte, maskRules map[string]int) ([]byte, error) {
	// 将 JSON 数据解析成 map
	var jsonData interface{}
	if err := json.Unmarshal(data, &jsonData); err != nil {
		return nil, err
	}

	// 递归处理 JSON 数据
	maskedData := recursiveMask(jsonData, maskRules)
	// 将修改后的数据重新编码为 JSON
	return json.Marshal(maskedData)
}

// 递归处理 JSON
func recursiveMask(data interface{}, maskRules map[string]int) interface{} {
	switch v := data.(type) {
	case map[string]interface{}:
		// 如果是 map 类型，遍历每个键值对
		for key, value := range v {
			// 对需要脱敏的字段进行处理
			if rule, exists := maskRules[key]; exists {
				v[key] = maskString(value.(string), rule/10, rule%10)
			} else {
				// 如果该字段是嵌套对象，递归处理
				v[key] = recursiveMask(value, maskRules)
			}
		}
		return v
	case []interface{}:
		// 如果是数组类型，递归处理每个元素
		for i, value := range v {
			v[i] = recursiveMask(value, maskRules)
		}
		return v
	default:
		// 其它类型数据直接返回，不处理
		return v
	}
}

// 脱敏字符串
func maskString(s string, headLen, tailLen int) string {
	if len(s) <= headLen+tailLen {
		// 如果字符串的长度小于等于保留的部分长度，直接返回原字符串
		return s
	}

	// 头部部分
	head := s[:headLen]
	// 尾部部分
	tail := s[len(s)-tailLen:]
	// 中间部分使用 ***
	masked := head + "***" + tail
	return masked
}
