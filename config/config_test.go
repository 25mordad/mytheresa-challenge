package config

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestConfig(t *testing.T) {

	assert := assert.New(t)
	c := GetConfig()
	assert.Equal(":8081", c.Port)

}
