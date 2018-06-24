# lambda-prune

Utility to clean up old versions of lambdas.


### Usage

```go

func main() {

 pruner := NewLambdaPruner("AWS_ACCESS_KEY_ID", "AWS_SECRET_ACCESS_KEY", "AWS_REGION")
 pruner.PruneStack("dev") 
}

```
Note That passing the `AWS_ACCESS_KEY_ID` and `AWS_SECRET_ACCESS_KEY"` are only required if you do not have your aws profile set or do not want to use your aws profile. 

### pruner.PruneStack(stage string)
Will gather all lambdas and your environment (optionally filtered by stage) and delete all versions other than `$LATEST` or those that are referenced. 

* `stage` optional stage name to limit the operation to

### pruner.PruneLambda(stage lambdaName)
Deletes all versions other than `$LATEST` or those that are referenced for the given lambda.

* `lambdaName` name of the lambda to prune


