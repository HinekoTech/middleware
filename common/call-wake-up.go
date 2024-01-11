package common

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
)

func CallWakeUp(url string) {
	fmt.Println("CallWakeUp URL: ", url)
	response, err := http.Get(url)
	if err != nil {
		log.Println("Error:", err)
		return
	}
	defer response.Body.Close()

	body, err := io.ReadAll(response.Body)
	if err != nil {
		log.Println("Error reading the response body:", err)
		return
	}

	type Success struct {
		Message string `json:"message"`
		Status  int    `json:"status"`
	}

	var data Success
	if err := json.Unmarshal(body, &data); err != nil {
		log.Println("RunWakeUp Error @ ", url)
		return
	}

	if data.Status == 1 {
		//log.Println(getServiceName())
		fmt.Println("************************** Wake up: OK!")
	} else {
		fmt.Println("************************** Wake up: ERROR!")
	}

}

// func getServiceName() string {
// 	if viper.GetString("production") == "true" {
// 		return fmt.Sprintf("Service name: %s, Current version: %s", os.Getenv("GAE_SERVICE"), os.Getenv("GAE_VERSION"))
// 	}
// 	return "Service name: [LOCAL] Engine Service"
// }
