package model_test

import (
	"dd_web/back/internal/model"
	"testing"
)

func TestJSONBRoundTrip(t *testing.T) {
	original := []byte(`{"amount_cent":100}`)
	var value model.JSONB
	if err := value.Scan(original); err != nil {
		t.Fatalf("扫描 JSON 失败: %v", err)
	}
	original[2] = 'X'
	got, err := value.Value()
	if err != nil {
		t.Fatalf("读取数据库值失败: %v", err)
	}
	if got != `{"amount_cent":100}` {
		t.Fatalf("JSON 往返值错误: %v", got)
	}
}

func TestJSONBRejectsInvalidJSON(t *testing.T) {
	var value model.JSONB
	if err := value.Scan([]byte(`{`)); err == nil {
		t.Fatal("非法 JSON 被接受")
	}
}
