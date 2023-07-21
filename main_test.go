package test

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os/exec"
	"strings"
	"testing"
	"github.com/stretchr/testify/assert"
)

func getAccessToken(subscriptionID string) (string, error) {
	fmt.Println(subscriptionID)
	cmd := exec.Command("az", "account", "get-access-token", "--query", "accessToken", "--output", "tsv", "--subscription", subscriptionID)

	output, err := cmd.Output()
	if err != nil {
		return "", err
	}

	return strings.TrimSpace(string(output)), nil
}

func getACRDetails(accessToken, subscriptionID, resourceGroupName, acrName string) (map[string]interface{}, error) {
	url := fmt.Sprintf("https://management.azure.com/subscriptions/%s/resourceGroups/%s/providers/Microsoft.ContainerRegistry/registries/%s?api-version=2023-06-01-preview", subscriptionID, resourceGroupName, acrName)
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Authorization", "Bearer "+accessToken)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var acrDetails map[string]interface{}
	err = json.Unmarshal(body, &acrDetails)
	if err != nil {
		return nil, err
	}

	return acrDetails, nil
}



func acrDetailsToJSON(acrDetails map[string]interface{}) (string, error) {
	acrDetailsJSON, err := json.MarshalIndent(acrDetails, "", "    ")
	if err != nil {
		return "", err
	}

	return string(acrDetailsJSON), nil
}

func isACRPrivate(acrDetails map[string]interface{}) string {
	properties, ok := acrDetails["properties"].(map[string]interface{})
	if !ok {
		return "properties not found"
	}

	if publicNetworkAccess, ok := properties["publicNetworkAccess"].(string); ok {
		fmt.Println("publicNetworkAccess:", publicNetworkAccess)
		return publicNetworkAccess
	}

	fmt.Println("publicNetworkAccess: properties not found")
	return "properties not found"
}


func TestPrivateACRWithPrivateEndpoint(t *testing.T) {
	t.Parallel()

	subscriptionID := "< >"
	resourceGroupName := "< >"
	acrName := "< >"
	expectedPublicNetworkAccess := "Disabled"  

	accessToken, err := getAccessToken(subscriptionID)
	if err != nil {
		t.Errorf("Failed to get access token: %s", err.Error())
		return
	}

	acrDetails, err := getACRDetails(accessToken, subscriptionID, resourceGroupName, acrName)
	assert.NoError(t, err)

	acrDetailsJSON, err := acrDetailsToJSON(acrDetails)
	assert.NoError(t, err)

	fmt.Println("ACR Details:")
	fmt.Println(acrDetailsJSON)

	publicNetworkAccessFirst := isACRPrivate(acrDetails)
	
	t.Run(fmt.Sprintf("Checking the Access for this ACR: %s", acrName), func(t *testing.T) {
		assert.Equal(t, expectedPublicNetworkAccess, publicNetworkAccessFirst, "publicNetworkAccess mismatch")
	})
	}