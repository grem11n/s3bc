package validate

import (
	"fmt"
	"log"
	"os"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/grem11n/s3bc/client/s3"
	"github.com/grem11n/s3bc/config"
)

func Run(config *config.Config) error {
	if err := config.Validate(); err != nil {
		log.Fatal(err)
	}

	c, err := s3.New(config.Bucket)
	if err != nil {
		return err
	}

	convertables, err := c.GetConvertableObjects(config.StorageClass)
	if err != nil {
		return err
	}

	if len(convertables) != 0 {
		nonFitFiles := map[string]string{}
		for _, obj := range convertables {
			nonFitFiles[aws.ToString(obj.Key)] = string(obj.StorageClass)
		}

		fmt.Printf("Not all the objects in the \"%s\" bucket have desired storage class\n", config.Bucket)
		fmt.Printf("Desired storage class: %s", config.StorageClass)
		fmt.Println("Objects with the different storage class:")

		for k, v := range nonFitFiles {
			fmt.Printf("%s: %s", k, v)
		}

		os.Exit(1)
	}

	return nil
}
