package brevesapiserver

import (
	"encoding/json"
	"net/http"

	models "github.com/ZolaraProject/breves-api/models"
	"github.com/mediocregopher/radix/v3"
)

var (
	BrevesVaultServiceHost string
	BrevesVaultServicePort string
	JwtSecretKey           string
	RedisHost              string
	RedisPort              string
	RedisPassword          string
	RedisPool              *radix.Pool
)

func writeStandardResponse(r *http.Request, w http.ResponseWriter, grpcToken string, message string) {
	responseObj := &models.Response{
		Token:   grpcToken,
		Message: message,
	}

	response, _ := json.Marshal(responseObj)
	w.Write(response)
}
