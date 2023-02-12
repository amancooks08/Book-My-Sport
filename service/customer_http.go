package service

// import (
// 	"encoding/json"
// 	"fmt"
// 	"net/http"

// 	logger "github.com/sirupsen/logrus"
// )

// func GetVenues(deps dependencies) http.HandlerFunc {
// 	return http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
// 		if req.Method != http.MethodGet {
// 			rw.WriteHeader(http.StatusMethodNotAllowed)
// 			return
// 		}

// 		venues, err := deps.CustomerServices.GetVenues(req.Context())
// 		if err != nil {
// 			http.Error(rw, fmt.Sprintf("%s", err), http.StatusInternalServerError)
// 			return
// 		}

// 		respBytes, err := json.Marshal(venues)
// 		if err != nil {
// 			logger.WithField("err", err.Error()).Error("Error marshalling venues response")
// 			rw.WriteHeader(http.StatusInternalServerError)
// 		}

// 		rw.Header().Add("Content-Type", "application/json")
// 		rw.Write(respBytes)
// 	})
// }
