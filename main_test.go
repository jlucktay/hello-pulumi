package main

import (
	"github.com/pulumi/pulumi/pkg/v2/resource"
)

type mocks int

func (mocks) NewResource(typeToken, name string, inputs resource.PropertyMap, provider, id string) (string,
	resource.PropertyMap, error) {
	return name + "_id", inputs, nil
}

func (mocks) Call(token string, args resource.PropertyMap, provider string) (resource.PropertyMap, error) {
	return args, nil
}
