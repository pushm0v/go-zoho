package cache

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/stretchr/testify/suite"
)

type ZohoCacheExternalSuite struct {
	suite.Suite
	cache   Cache
	storage map[string]interface{}
}

func TestZohoCacheExternalSuite(t *testing.T) {
	suite.Run(t, new(ZohoCacheExternalSuite))
}

func (suite *ZohoCacheExternalSuite) SetupTest() {
	extCache := NewExternalCache(Option{
		SetFunc: suite.setFuncMock,
		GetFunc: suite.getFuncMock,
	})

	suite.cache = extCache
	suite.storage = map[string]interface{}{}
}

func (suite *ZohoCacheExternalSuite) setFuncMock(key string, value interface{}, expire int) error {
	suite.storage[key] = value
	return nil
}

func (suite *ZohoCacheExternalSuite) getFuncMock(key string) (value interface{}, err error) {

	if val, ok := suite.storage[key]; ok {
		return val, nil
	}

	return nil, fmt.Errorf("Not found")
}

func (suite *ZohoCacheExternalSuite) TestSetCache() {
	expectedKey := "some-key"
	expectedValue := 1

	err := suite.cache.Set(expectedKey, expectedValue, 1)

	assert.Nil(suite.T(), err, "Error should be nil")
	assert.NotNil(suite.T(), suite.storage[expectedKey], "Cache should not be nil")
	assert.Equal(suite.T(), expectedValue, suite.storage[expectedKey].(int), "Value should be equal")
}

func (suite *ZohoCacheExternalSuite) TestGetCache() {
	expectedKey := "some-key"
	expectedValue := 1

	err := suite.cache.Set(expectedKey, expectedValue, 1)
	val, err := suite.cache.Get(expectedKey)

	assert.Nil(suite.T(), err, "Error should be nil")
	assert.NotNil(suite.T(), val, "Cache should not be nil")
	assert.Equal(suite.T(), expectedValue, val.(int), "Value should be equal")
}
