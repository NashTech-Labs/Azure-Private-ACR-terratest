## Terratest code to test the private Azure ACR

### With the help of this template, you can test your private azure container registory. If you created a private ACR then this terratest can test this without deploying ACR. You have to define the few values in it as per your requirements:-

### You can see the below. You have replace these values with your own values:-


    subscriptionID := "< >"
	resourceGroupName := "< >"
	acrName := "< >"
	expectedPublicNetworkAccess := "Disabled"  


### After doing this you just need to run the below command and test cases will run:-

        go test -v