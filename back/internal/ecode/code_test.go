package ecode

import "testing"

func TestDefinitionsUseFourDigitCodesAndMessages(t *testing.T) {
	for code, definition := range definitions {
		if code < 1000 || code > 5000 {
			t.Errorf("状态码必须在1000到5000之间: %d", code)
		}
		if code == Success {
			if definition.Message != "" {
				t.Errorf("成功状态不应包含错误文案: %q", definition.Message)
			}
		} else if definition.Message == "" {
			t.Errorf("错误状态码%d缺少error_msg", code)
		}
		if definition.HTTPStatus < 100 || definition.HTTPStatus > 599 {
			t.Errorf("状态码%d的HTTP状态非法: %d", code, definition.HTTPStatus)
		}
	}
	if _, ok := definitions[InternalError]; !ok || InternalError != 5000 {
		t.Fatal("系统内部错误必须固定为5000")
	}
}
