package model

import (
	"sync"
	"testing"

	"gorm.io/gorm/schema"
)

func TestAllModelFieldsHaveComments(t *testing.T) {
	cache := &sync.Map{}
	for _, value := range AllModels() {
		modelSchema, err := schema.Parse(value, cache, schema.NamingStrategy{})
		if err != nil {
			t.Fatalf("parse model schema: %v", err)
		}
		for _, field := range modelSchema.Fields {
			if field.DBName != "" && field.Comment == "" {
				t.Errorf("%s.%s is missing a GORM comment", modelSchema.Name, field.Name)
			}
		}
	}
}

func TestWalletAccountUsesUserIDAsOnlyPrimaryKey(t *testing.T) {
	modelSchema, err := schema.Parse(&WalletAccount{}, &sync.Map{}, schema.NamingStrategy{})
	if err != nil {
		t.Fatal(err)
	}
	if len(modelSchema.PrimaryFields) != 1 || modelSchema.PrimaryFields[0].DBName != "user_id" {
		t.Fatalf("wallet primary fields = %#v", modelSchema.PrimaryFieldDBNames)
	}
	if modelSchema.LookUpField("ID") != nil {
		t.Fatal("wallet account must not keep an independent ID field")
	}
}
