package model_test

import (
	"encoding/json"
	"testing"

	"github.com/AssetPortal/assets-api/pkg/model"
	"github.com/stretchr/testify/assert"
)

func TestNewResponseError(t *testing.T) {
	message := "Error occurred"
	resp := model.NewResponseError(message)

	assert.False(t, resp.OK)
	assert.Equal(t, message, resp.Message)
	assert.Nil(t, resp.Data)
}

func TestNewResponseData(t *testing.T) {
	data := map[string]string{
		"key": "value",
	}
	resp := model.NewResponseData(data)

	assert.True(t, resp.OK)
	assert.Empty(t, resp.Message)

	var result map[string]string
	err := json.Unmarshal(resp.Data, &result)
	assert.NoError(t, err)
	assert.Equal(t, data, result)
}

func TestNewResponseDataError(t *testing.T) {
	resp := model.NewResponseData(func() {})

	assert.False(t, resp.OK)
	assert.Equal(t, "Failed to process data", resp.Message)
	assert.Nil(t, resp.Data)
}
