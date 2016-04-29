package api

import (
	"strconv"
	"encoding/json"
	"net/http"

	"linq/core/utils"

	"github.com/gorilla/mux"
	"github.com/satori/go.uuid"
)

type apiService struct {
	http.ResponseWriter
	Request *http.Request
	isReturned bool
}

func ApiService(w http.ResponseWriter, r *http.Request) apiService {
	return apiService{w, r, false}
}

type data struct {
	Ids []uuid.UUID `json:"ids"`
}

type RequestDataIds struct {
	Data  data   `json:"data"`
	Token string `json:"token"`
}

type RequestDataImage struct {
	Data  string `json:"data"`
	Token string `json:"token"`
}

type JsonSuccessResponse struct {
	Data  []interface{} `json:"data"`
	Token uuid.UUID     `json:"token"`
}

type JsonErrorResponse struct {
	Status int `json:"status"`
	Source string `json:"source"`
	Title  string `json:"title"`
	Method string `json:"method"`
	Detail string `json:"detail"`
}

type JsonErrorResponses struct {
	Errors []JsonErrorResponse `json:"errors"`
}

func (api apiService) FormValue(key string) string {
	return api.Request.FormValue(key)
}

func (api apiService) MuxVars(key string) string {
	muxVars := mux.Vars(api.Request)
	return muxVars[key]
}

func (api apiService) DecodeBody(requestData interface{}) error {
	decoder := json.NewDecoder(api.Request.Body)
	err := decoder.Decode(&requestData)
	utils.HandleWarn(err)
	return err
}

func (api apiService) ReturnJson(payload interface{}) {
	if(!api.isReturned){
		api.Header().Set("Content-Type", "application/linq.api+json; charset=UTF-8")
		api.WriteHeader(http.StatusOK)
	
		data := make([]interface{}, 1)
		data[0] = payload
	
		responseData := JsonSuccessResponse{
			Data:  data,
			Token: uuid.NewV4(),
		}
	
		err := json.NewEncoder(api).Encode(responseData)
		utils.HandleWarn(err)
		api.isReturned = true
	}
}

func (api apiService) HandleApiError(err error, status int)  {
	if (err != nil && !api.isReturned) {
		api.returnJsonServerError(err.Error(), status)
	}
}

func (api apiService) returnJsonServerError(detail string, status int) {
	api.Header().Set("Content-Type", "application/linq.api+json; charset=UTF-8")
	api.WriteHeader(status)

	responseData := JsonErrorResponses{
		Errors: []JsonErrorResponse{
			{
				Status: status,
				Title:  http.StatusText(status),
				Source: api.Request.URL.RequestURI(),
				Method: api.Request.Method,
				Detail: detail,
			},
		},
	}

	json.NewEncoder(api).Encode(responseData)
	utils.Log.Warn(api.Request.Method, api.Request.URL.RequestURI(), strconv.Itoa(status), detail)
	
	api.isReturned = true
}
