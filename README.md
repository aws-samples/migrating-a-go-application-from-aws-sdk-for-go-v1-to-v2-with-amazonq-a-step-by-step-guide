# Welcome to your CDK Go project!

This is a project for CDK development with Go.

The `cdk.json` file tells the CDK toolkit how to execute your app.


Add Step "Set GOPROXY=direct" # Note to customer



```
[comment]: ??? - Build Steps may not be need if we add to repo

 #> GOOS=linux GOARCH=arm64 go build -o bootstrap lambdafunction/main.go

zip myFunction.zip bootstrap

GOOS=linux GOARCH=arm64 go build -o bootstrap lambdafunction/hitcounter.go

zip getHits.zip bootstrap


GOOS=linux GOARCH=arm64 go build -o bootstrap lambdafunction/updatePlayer.go

zip hitsLambda.zip bootstrap

cdk bootstrap

cdk deploy --all
```
# Next Test Deployment
1. Next we need to confirm that cdk deployment is working.
[comment]: https://rxncm1fbxa.execute-api.us-east-1.amazonaws.com/prod/
```
curl -sX GET "https://rxncm1fbxa.execute-api.us-east-1.amazonaws.com/prod/getPlayers/?firstName=Carlos" | jq
```



## Commands to Deploy Go Project

 ## Instructions to Deploy Sample GO Environment
1.  `cdk bootstrap`                  For the first time bootstrap in the account/region
```
#> cdkbootstrap
```
2. `cdk deploy`                     deploy this stack to your default AWS account/region , set AWS_REGION 
3. `cdk diff`                       compare deployed stack with current state

 
 
 
 
 * `GOOS=linux GOARCH=arm64 go build -tags lambda.norpc -o bootstrap lambdafunction/main.go`                           Build the Lambda


 * `cdk deploy`                     deploy this stack to your default AWS account/region , set AWS_REGION 
 * `cdk diff`                       compare deployed stack with current state
 * `cdk synth`                      emits the synthesized CloudFormation template
 * `cdk bootstrap`                  For the first time bootstrap in the account/region
 * `go test`                        run unit tests


API response (AWS SDK v1):
curl -sX GET "https://xxxxxxx.execute-api.us-west-2.amazonaws.com/prod/?firstName=Carlos" | jq
[
{
"LastName": "Hernandez",
"FirstName": "Carlos",
"DOB": "1988-06-18",
"Plays": "Right",
"CountryOfBirth": "Spain",
"CountryOfResidence": "Spain"
}
]

[comment:] # Need to create a folder in S3 Bucket and upload players.csv - CDK Change
[comment:] activeplayers/players.csv

