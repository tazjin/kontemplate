package context

import (
	"reflect"
	"testing"
)

func TestLoadFlatContextFromFile(t *testing.T) {
	ctx, err := LoadContextFromFile("testdata/flat-test.yaml")

	if err != nil {
		t.Error(err)
		t.Fail()
	}

	expected := Context{
		Name: "k8s.prod.mydomain.com",
		Global: map[string]interface{}{
			"globalVar": "lizards",
		},
		ResourceSets: []ResourceSet{
			{
				Name: "some-api",
				Values: map[string]interface{}{
					"apiPort":          float64(4567), // yep!
					"importantFeature": true,
					"version":          "1.0-0e6884d",
				},
				Include: nil,
				Parent:  "",
			},
		},
		BaseDir: "testdata",
	}

	if !reflect.DeepEqual(*ctx, expected) {
		t.Error("Loaded context and expected context did not match")
		t.Fail()
	}
}

func TestLoadContextWithResourceSetCollections(t *testing.T) {
	ctx, err := LoadContextFromFile("testdata/collections-test.yaml")

	if err != nil {
		t.Error(err)
		t.Fail()
	}

	expected := Context{
		Name: "k8s.prod.mydomain.com",
		Global: map[string]interface{}{
			"globalVar": "lizards",
		},
		ResourceSets: []ResourceSet{
			{
				Name: "some-api",
				Values: map[string]interface{}{
					"apiPort":          float64(4567), // yep!
					"importantFeature": true,
					"version":          "1.0-0e6884d",
				},
				Include: nil,
				Parent:  "",
			},
			{
				Name: "collection/nested",
				Values: map[string]interface{}{
					"lizards": "good",
				},
				Include: nil,
				Parent:  "collection",
			},
		},
		BaseDir: "testdata",
	}

	if !reflect.DeepEqual(*ctx, expected) {
		t.Error("Loaded context and expected context did not match")
		t.Fail()
	}

}

func TestSubresourceVariableInheritance(t *testing.T) {
	ctx, err := LoadContextFromFile("testdata/parent-variables.yaml")

	if err != nil {
		t.Error(err)
		t.Fail()
	}

	expected := Context{
		Name: "k8s.prod.mydomain.com",
		ResourceSets: []ResourceSet{
			{
				Name: "parent/child",
				Values: map[string]interface{}{
					"foo": "bar",
					"bar": "baz",
				},
				Include: nil,
				Parent: "parent",
			},
		},
		BaseDir: "testdata",
	}

	if !reflect.DeepEqual(*ctx, expected) {
		t.Error("Loaded and expected context did not match")
		t.Fail()
	}
}

func TestSubresourceVariableInheritanceOverride(t *testing.T) {
	ctx, err := LoadContextFromFile("testdata/parent-variable-override.yaml")

	if err != nil {
		t.Error(err)
		t.Fail()
	}

	expected := Context{
		Name: "k8s.prod.mydomain.com",
		ResourceSets: []ResourceSet{
			{
				Name: "parent/child",
				Values: map[string]interface{}{
					"foo": "newvalue",
				},
				Include: nil,
				Parent: "parent",
			},
		},
		BaseDir: "testdata",
	}

	if !reflect.DeepEqual(*ctx, expected) {
		t.Error("Loaded and expected context did not match")
		t.Fail()
	}
}
