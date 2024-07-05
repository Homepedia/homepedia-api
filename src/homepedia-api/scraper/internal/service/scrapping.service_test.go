package service

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestScrapeService(t *testing.T) {
	userAgentList := []string{"Mozilla/5.0"}
	key := "55-06"
	url := "https://immobilier.lefigaro.fr/annonces/annonce-68152302.html"
	res, err := ScrappePage(url, userAgentList, key)
	assert.NoError(t, err)
	assert.NotNil(t, res)
	assert.Equal(t, "68152302", res.Ref)
	assert.Equal(t, "N.C", res.GreenhouseGasIndex)
	assert.Equal(t, "N.C", res.EnergeticPerf)
}

func TestScrapeServiceInvalidURL(t *testing.T) {
	userAgentList := []string{"Mozilla/5.0"}
	key := "55-06"
	url := "https://immobilier.lefigaro.fr/annonces/invalid-url"
	res, err := ScrappePage(url, userAgentList, key)

	assert.Error(t, err)
	assert.Nil(t, res)
}
