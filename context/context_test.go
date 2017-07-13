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
				Path: "some-api",
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
				Path: "some-api",
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
				Path: "collection/nested",
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
				Path: "parent/child",
				Values: map[string]interface{}{
					"foo": "bar",
					"bar": "baz",
				},
				Include: nil,
				Parent:  "parent",
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
				Path: "parent/child",
				Values: map[string]interface{}{
					"foo": "newvalue",
				},
				Include: nil,
				Parent:  "parent",
			},
		},
		BaseDir: "testdata",
	}

	if !reflect.DeepEqual(*ctx, expected) {
		t.Error("Loaded and expected context did not match")
		t.Fail()
	}
}

func TestDefaultValuesLoading(t *testing.T) {
	ctx, err := LoadContextFromFile("testdata/default-loading.yaml")
	if err != nil {
		t.Error(err)
		t.Fail()
	}

	rs := ctx.ResourceSets[0]
	if rs.Values["defaultValues"] != "loaded" {
		t.Errorf("Default values not loaded from YAML file")
		t.Fail()
	}

	if rs.Values["override"] != "notAtAll" {
		t.Error("Default values should not override other values")
		t.Fail()
	}
}

func TestImportValuesLoading(t *testing.T) {
	ctx, err := LoadContextFromFile("testdata/import-vars-simple.yaml")
	if err != nil {
		t.Error(err)
		t.Fail()
	}

	expected := map[string]interface{}{
		"override": "true",
		"music": map[string]interface{}{
			"artist": "Pallida",
			"track":  "Tractor Beam",
		},
	}

	if !reflect.DeepEqual(ctx.Global, expected) {
		t.Error("Expected global values after loading imports did not match!")
		t.Fail()
	}
}

func TestImportValuesOverride(t *testing.T) {
	ctx, err := LoadContextFromFile("testdata/import-vars-override.yaml")
	if err != nil {
		t.Error(err)
		t.Fail()
	}

	expected := map[string]interface{}{
		"override": float64(3),
		"music": map[string]interface{}{
			"artist": "Pallida",
			"track":  "Tractor Beam",
		},
		"place":     "Oslo",
		"globalVar": "very global!",
	}

	if !reflect.DeepEqual(ctx.Global, expected) {
		t.Error("Expected global values after loading imports did not match!")
		t.Fail()
	}
}

func TestExplicitPathLoading(t *testing.T) {
	ctx, err := LoadContextFromFile("testdata/explicit-path.yaml")
	if err != nil {
		t.Error(err)
		t.Fail()
	}

	expected := Context{
		Name: "k8s.prod.mydomain.com",
		ResourceSets: []ResourceSet{
			{
				Name: "some-api-europe",
				Path: "some-api",
				Values: map[string]interface{}{
					"location": "europe",
				},
				Include: nil,
				Parent:  "",
			},
			{
				Name: "some-api-asia",
				Path: "some-api",
				Values: map[string]interface{}{
					"location": "asia",
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

func TestExplicitSubresourcePathLoading(t *testing.T) {
	ctx, err := LoadContextFromFile("testdata/explicit-subresource-path.yaml")
	if err != nil {
		t.Error(err)
		t.Fail()
	}

	expected := Context{
		Name: "k8s.prod.mydomain.com",
		ResourceSets: []ResourceSet{
			{
				Name:   "parent/child",
				Path:   "parent-path/child-path",
				Parent: "parent",
			},
		},
		BaseDir: "testdata",
	}

	if !reflect.DeepEqual(*ctx, expected) {
		t.Error("Loaded context and expected context did not match")
		t.Fail()
	}
}
