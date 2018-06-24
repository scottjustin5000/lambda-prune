package main

import (
	"fmt"
	"strings"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/lambda"
)

type LambdaPruner struct {
	client *lambda.Lambda
}

func NewLambdaPruner(key string, secret string, region string) *LambdaPruner {
	var sess *session.Session
	awsConfig := &aws.Config{
		Region: aws.String(region),
	}
	if key != "" && secret != "" {
		awsConfig.Credentials = credentials.NewStaticCredentials(key, secret, "")
		sess, _ = session.NewSession()
	} else {
		sess = session.Must(session.NewSessionWithOptions(session.Options{
			SharedConfigState: session.SharedConfigEnable,
		}))
	}

	client := lambda.New(sess, awsConfig)
	pruner := &LambdaPruner{client}
	return pruner
}

func contains(stringSlice []string, searchVal string) bool {
	for _, value := range stringSlice {
		if value == searchVal {
			return true
		}
	}
	return false
}

func (pruner *LambdaPruner) getVersions(lambdaName string, references []string, versions []string, marker string) []string {

	params := &lambda.ListVersionsByFunctionInput{
		MaxItems:     aws.Int64(1),
		FunctionName: aws.String(lambdaName),
	}

	if marker != "" {
		params.Marker = aws.String(marker)
	}

	result, err := pruner.client.ListVersionsByFunction(params)
	if err != nil {
		fmt.Println(err)
	}
	for _, m := range result.Versions {
		if *m.Version != "$LATEST" && !contains(references, *m.Version) {
			versions = append(versions, *m.FunctionArn)
		}
	}
	if aws.StringValue(result.NextMarker) == "" {
		return versions
	}
	return pruner.getVersions(lambdaName, references, versions, aws.StringValue(result.NextMarker))
}

func (pruner *LambdaPruner) getReferences(lambdaName string, aliases []string, marker string) []string {
	params := &lambda.ListAliasesInput{
		MaxItems:     aws.Int64(1),
		FunctionName: aws.String(lambdaName),
	}

	if marker != "" {
		params.Marker = aws.String(marker)
	}

	result, err := pruner.client.ListAliases(params)
	if err != nil {
		fmt.Println(err)
	}
	for _, m := range result.Aliases {
		aliases = append(aliases, *m.FunctionVersion)
	}
	if aws.StringValue(result.NextMarker) == "" {
		return aliases
	}
	return pruner.getReferences(lambdaName, aliases, aws.StringValue(result.NextMarker))
}

func (pruner *LambdaPruner) getLambdas(filter string) []string {
	var funcs []string
	params := &lambda.ListFunctionsInput{
		MaxItems: aws.Int64(123),
	}
	err := pruner.client.ListFunctionsPages(params,
		func(page *lambda.ListFunctionsOutput, lastPage bool) bool {
			for _, fn := range page.Functions {
				if(filter == "") {
					funcs = append(funcs, *fn.FunctionName)
				} else if strings.Index(*fn.FunctionName, filter) > -1 {
					funcs = append(funcs, *fn.FunctionName)
				}
			}
			return !lastPage
		})
	if err != nil {
		fmt.Println(err)
	}
	return funcs
}

func (pruner *LambdaPruner) deleteFunctionVersion(functionArn string) {

	params := &lambda.DeleteFunctionInput{
	    FunctionName: aws.String(functionArn),
	  }
	  _, err := pruner.client.DeleteFunction(params)

	  if err != nil {
	    fmt.Println(err)
	  }
}

func (pruner *LambdaPruner) PruneStack(stage string) {
	lambdas := pruner.getLambdas(stage)
	for _, l := range lambdas {
		pruner.PruneLambda(l)
	}
}

func (pruner *LambdaPruner) PruneLambda(functionName string) {
	var refs []string
	references := pruner.getReferences(functionName, refs, "")

	var vars []string
	versions := pruner.getVersions(functionName, references, vars, "")

	for _, vrs := range versions {
		pruner.deleteFunctionVersion(vrs)
	}
}
