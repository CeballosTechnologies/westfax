package westfax

import (
	"github.com/spf13/viper"
	"testing"
	"time"
)

var client Client

func init() {
	viper.SetConfigFile(".env")
	viper.SetConfigType("env")
	err := viper.ReadInConfig()
	if err != nil {
		return
	}

	client = *New(viper.GetString("WESTFAX_USERNAME"), viper.GetString("WESTFAX_PASSWORD"), viper.GetString("WESTFAX_PRODUCTID"))
}

func TestClient_SecurityPing(t *testing.T) {
	_, err := client.SecurityPing("testping")
	if err != nil {
		return
	}
}

func TestClient_GetInboundFaxIdentifiers(t *testing.T) {
	// create date
	startDate := time.Date(2022, 11, 1, 0, 0, 0, 0, time.UTC)
	_, err := client.GetInboundFaxIdentifiers(startDate.Format("2006-02-01"))
	if err != nil {
		return
	}
}

func TestClient_GetFaxDocuments(t *testing.T) {
	_, err := client.GetFaxDocument("0afbf6f9-998e-4863-9a08-757bf1308db9")
	if err != nil {
		return
	}
}
