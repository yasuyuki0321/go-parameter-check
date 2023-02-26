package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httputil"
	"net/url"
	"os"
	"time"

	runtime "github.com/aws/aws-lambda-go/lambda"
)

type resultFromExtension struct {
	Parameter struct {
		Arn              string      `json:"ARN"`
		DataType         string      `json:"DataType"`
		LastModifiedDate time.Time   `json:"LastModifiedDate"`
		Name             string      `json:"Name"`
		Selector         interface{} `json:"Selector"`
		SourceResult     interface{} `json:"SourceResult"`
		Type             string      `json:"Type"`
		Value            string      `json:"Value"`
		Version          int         `json:"Version"`
	} `json:"Parameter"`
	ResultMetadata struct {
	} `json:"ResultMetadata"`
}

type response struct {
	Parameter  resultFromExtension `json:"parameter"`
	StatusCode int                 `json:"statusCode"`
}

func handleRequest() (response, error) {
	param, _ := getParameterValue("secure-string-test")

	return response{
		Parameter:  param,
		StatusCode: http.StatusOK,
	}, nil
}

func getParameterValue(parameterPath string) (resultFromExtension, error) {
	urlEncodedPath := url.QueryEscape(parameterPath)
	requestURL := fmt.Sprintf("http://localhost:2773/systemsmanager/parameters/get/?name=%s&withDecryption=true", urlEncodedPath)

	req, _ := http.NewRequest("GET", requestURL, nil)
	req.Header.Set("X-Aws-Parameters-Secrets-Token", os.Getenv("AWS_SESSION_TOKEN"))

	dump, _ := httputil.DumpRequestOut(req, true)
	fmt.Printf("%s", dump)

	client := new(http.Client)
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
	}
	defer resp.Body.Close()

	dumpResp, _ := httputil.DumpResponse(resp, true)
	fmt.Printf("%s", dumpResp)

	var resultFromExtension resultFromExtension
	body, _ := io.ReadAll(resp.Body)
	json.Unmarshal([]byte(body), &resultFromExtension)

	return resultFromExtension, nil
}

func main() {
	runtime.Start(handleRequest)
}
