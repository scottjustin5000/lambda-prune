package main
import "os"

func main() {
	rb := NewLambdaPruner(os.Getenv("AWS_ACCESS_KEY_ID"), os.Getenv("AWS_SECRET_ACCESS_KEY"), os.Getenv("AWS_REGION"))
	rb.PruneStack("dev")
}
