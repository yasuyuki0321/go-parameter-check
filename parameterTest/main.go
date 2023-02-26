package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"time"

	runtime "github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/ssm"
	"github.com/aws/aws-sdk-go/aws"
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

type parameter struct {
	Name  string `json:"Name"`
	Value string `json:"Value"`
}

type response struct {
	Parameters []parameter `json:"parameters"`
	StatusCode int         `json:"statusCode"`
}

func handleRequest() (response, error) {
	paramNames, err := getAllParameters()
	if err != nil {
		return response{}, err
	}

	var parameters = make([]parameter, 0)
	for _, paramName := range paramNames {
		param, _ := getParameterValue(paramName)

		parameters = append(parameters, parameter{
			Name:  param.Parameter.Name,
			Value: param.Parameter.Value,
		})
	}

	return response{
		Parameters: parameters,
		StatusCode: http.StatusOK,
	}, nil
}

func getAllParameters() ([]string, error) {
	cfg, _ := config.LoadDefaultConfig(context.TODO(), config.WithRegion("us-east-1"))
	svc := ssm.NewFromConfig(cfg)

	out, err := svc.GetParametersByPath(context.TODO(), &ssm.GetParametersByPathInput{
		Path: aws.String("/"),
	})
	if err != nil {
		return nil, err
	}

	paramNames := make([]string, 0)
	for i := range out.Parameters {
		paramNames = append(paramNames, *out.Parameters[i].Name)
	}

	return paramNames, nil
}

func getParameterValue(parameterPath string) (resultFromExtension, error) {
	urlEncodedPath := url.QueryEscape(parameterPath)
	requestURL := fmt.Sprintf("http://localhost:2773/systemsmanager/parameters/get/?name=%s&withDecryption=true", urlEncodedPath)

	req, _ := http.NewRequest("GET", requestURL, nil)
	req.Header.Set("X-Aws-Parameters-Secrets-Token", os.Getenv("AWS_SESSION_TOKEN"))

	// dump, _ := httputil.DumpRequestOut(req, true)
	// fmt.Printf("%s", dump)

	client := new(http.Client)
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
	}
	defer resp.Body.Close()

	// dumpResp, _ := httputil.DumpResponse(resp, true)
	// fmt.Printf("%s", dumpResp)

	var resultFromExtension resultFromExtension
	body, _ := io.ReadAll(resp.Body)
	json.Unmarshal([]byte(body), &resultFromExtension)

	return resultFromExtension, nil
}

func main() {
	runtime.Start(handleRequest)
}
