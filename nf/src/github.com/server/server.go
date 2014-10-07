package main

import (
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strings"
)

type UserDataType struct {
	Domain string     `json:"domain"`
	Info   []InfoType `json:"users"`
}

type InfoType struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type JsonDataType struct {
	AccesGranted bool   `json:"access_granted"`
	Reason       string `json:"reason,omitempty"`
}

type ResponseJson struct {
	StatusCode int
	JsonData   JsonDataType
}

var records []UserDataType

func main() {
	http.HandleFunc("/", viewHandler)
	http.ListenAndServe(":8080", nil)
}

func viewHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Content-type", "application/json")

	jsonObj := getResponse(r)
	jsonMsg, err := json.Marshal(jsonObj.JsonData)
	if err != nil {
		http.Error(w, "Oops", http.StatusInternalServerError)
		w.WriteHeader(500)
	}
	if jsonObj.StatusCode == 200 {
		w.WriteHeader(200)
		w.Write(jsonMsg)
		w.Write([]byte("\n"))
	} else {
		w.WriteHeader(jsonObj.StatusCode)
	}

}

func getResponse(r *http.Request) ResponseJson {

	//VALIDATION OF THE URL
	urlPath := r.URL.Path[:]
	if len(urlPath) > 14 {
		if urlPath[1:15] == "api/2/domains/" {

			//GETTING VALUES FROM THE POST DATA

			domain := ""
			user := ""
			pwd := ""
			if r.Method == "POST" {
				params := strings.Split(urlPath[15:], "/")
				if len(params) == 2 && params[1] == "proxyauth" {
					domain = params[0]
					user = r.FormValue("username")
					pwd = r.FormValue("password")

					//READING FROM THE users.json FILE FOR VERIFICATION
					if len(records) == 0 {
						userFile, err := os.Open("users.json")
						if err != nil {
							fmt.Println("opening config file", err.Error())
							message := ResponseJson{500, JsonDataType{false, ""}}
							return message
						}
						jsonParser := json.NewDecoder(userFile)
						if err := jsonParser.Decode(&records); err != nil {
							fmt.Println("parsing config file", err.Error())
							message := ResponseJson{500, JsonDataType{false, ""}}
							return message
						}
						if err := userFile.Close(); err != nil {
							panic(err)
							message := ResponseJson{500, JsonDataType{false, ""}}
							return message
						}
					}

					//VALIDATION OF THE USER

					message := ResponseJson{}
				OuterLoop:
					for v := 0; v < len(records); v++ {
						//domainFlag := false
						//userFlag := false
						if records[v].Domain == domain {
							//domainFlag = true
						InnerLoop:
							for k := 0; k < len(records[v].Info); k++ {
								if records[v].Info[k].Username == user {
									//userFlag = true
									pwdBytes := []byte(records[v].Info[k].Password)
									hasher := sha256.New()
									hasher.Write(pwdBytes)
									sha := base64.URLEncoding.EncodeToString(hasher.Sum(nil))
									sha = "{SHA256}" + sha
									if sha == pwd {
										message = ResponseJson{200, JsonDataType{true, ""}}
									} else {
										message = ResponseJson{200, JsonDataType{false, "denied by policy"}}
									}
									break OuterLoop
								} else {
									if len(records[v].Info)-1 == k {
										message = ResponseJson{200, JsonDataType{false, "denied by policy"}}
									} else {
										continue InnerLoop
									}
								}
								break OuterLoop
							}
						} else {
							message = ResponseJson{404, JsonDataType{false, ""}}
						}
					}
					return message
				} else {
					message := ResponseJson{404, JsonDataType{false, ""}}
					return message
				}
			} else {
				message := ResponseJson{500, JsonDataType{false, ""}}
				return message
			}
		} else {
			//statusCode := "200"
			message := ResponseJson{404, JsonDataType{false, ""}}
			return message
		}
	} else {
		message := ResponseJson{404, JsonDataType{false, ""}}
		return message
	}
}
