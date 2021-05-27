package service

import (
	"encoding/json"

	"github.com/mitchellh/mapstructure"
)

func (s *service) parseParams(items map[string]string) (*Params, error) {
	params := new(Params)
	mapstructure.Decode(items, params)
	if err := s.Validator.Struct(params); err != nil {
		return nil, err
	}

	return params, nil
}

func (s *service) parseVolumeContext(items map[string]string) (*VolumeContext, error) {
	context := new(VolumeContext)
	mapstructure.Decode(items, context)
	if err := s.Validator.Struct(context); err != nil {
		return nil, err
	}

	return context, nil
}

func (s *service) parsePublishContext(items map[string]string) (*PublishContext, error) {
	context := new(PublishContext)
	mapstructure.Decode(items, context)
	if err := s.Validator.Struct(context); err != nil {
		return nil, err
	}

	return context, nil
}

func (s *service) toMap(i interface{}) map[string]string {
	out, _ := json.Marshal(i)
	mapped := map[string]string{}

	json.Unmarshal(out, &mapped)

	return mapped
}
