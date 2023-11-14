package upload

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"

	model "github.com/FourWD/middleware/model"
)

func Upload(u model.UploadPayload) (model.UploadResult, error) {

	r := new(model.UploadResult)

	//start uploading
	uploadUrl := "https://pakwan-service-dot-fourwd.as.r.appspot.com/api/v1/upload/upload"

	var resp model.ApiResponse
	u.FileBase64 = strings.Replace(u.FileBase64, "data:image/png;base64,", "", -1)
	u.FileBase64 = strings.Replace(u.FileBase64, "data:image/jpeg;base64,", "", -1)
	u.FileBase64 = strings.Replace(u.FileBase64, "data:image/jpg;base64,", "", -1)

	jsonData, err := json.Marshal(u)

	if err != nil {
		fmt.Println("there was an error with the JSON", err.Error())
	} else {
		client := &http.Client{}
		req, err := http.NewRequest("POST", uploadUrl, bytes.NewBuffer(jsonData))
		if err != nil {
			fmt.Println(err)
			return *r, err
		}
		req.Header.Add("Content-Type", "application/json")
		req.Header.Add("Authorization", "Bearer PA3KBCDIORypzCzD2fQdaqyLUHpPoM60BEaeP68O1GXmbP7dF0hyOBed9ZRcr6ti")

		response, err := client.Do(req)
		if err != nil {
			fmt.Println(err)
			return *r, err
		}
		defer response.Body.Close()
		if err != nil {
			fmt.Println("there was an error with the request", err.Error())
		} else {
			body, _ := ioutil.ReadAll(response.Body)
			err := json.Unmarshal(body, &resp)
			if err != nil {
				return *r, err
			}
		}
	}

	//end upload

	// call pakwan upload

	return *r, nil
}
