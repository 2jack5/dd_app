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
