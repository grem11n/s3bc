package convert

import (
	"fmt"
	"log"
	"regexp"
	"sync"

	"github.com/aws/aws-sdk-go-v2/aws"
	pb "github.com/cheggaaa/pb/v3"
	"github.com/grem11n/s3bc/client/s3"
	"github.com/grem11n/s3bc/config"
	multierror "github.com/hashicorp/go-multierror"
)

const (
	// Trottle concurrent requests to AWS API.
	// AWS S3 API has a limit of 3500 COPY requests per prefix per second, so just playing safe here.
	apiCap = 3000
)

func Run(config *config.Config) error {
	if err := config.Validate(); err != nil {
		log.Fatal(err)
	}

	c, err := s3.New(config.Bucket)
	if err != nil {
		return err
	}

	convertibles, err := c.CreateInput(config.StorageClass)
	if err != nil {
		return err
	}

	// Filter out excluded patterns
	eligibles := []*s3.CopyObjectInput{}
	for _, convertable := range convertibles {
		if !isExcluded(config.Excluded, aws.ToString(convertable.Key)) {
			eligibles = append(eligibles, convertable)
		}
	}

	eligiblesCount := len(eligibles)

	fmt.Printf("%v eligible objects found for convesion\n", eligiblesCount)

	// If DryRun is found, only print eligible keys and return.
	if config.DryRun {
		fmt.Println("Dry run is set to true. No action will be done on the S3 side. Keys eligible for conversion:")
		for _, obj := range eligibles {
			fmt.Println(*obj.Key)
		}
		return nil
	}

	fmt.Printf("Converting eligible objects to the %s storage class...\n", config.StorageClass)

	// throttle goroutines to not to hit S3 API limits
	throttle := make(chan int, apiCap)
	var wg sync.WaitGroup
	var warns error

	bar := pb.StartNew(eligiblesCount)
	for _, obj := range eligibles {
		throttle <- 1
		wg.Add(1)
		gObj := *obj

		go func() {
			defer wg.Done()
			err := c.UpdateObjectStorageClass(&gObj)
			if err != nil {
				warns = multierror.Append(warns, err)
			}

			<-throttle
		}()
		wg.Wait()
		bar.Increment()
	}

	// Fill the channel to be sure, that all goroutines finished.
	for i := 0; i < cap(throttle); i++ {
		throttle <- 1
	}

	bar.Finish()

	if warns != nil {
		fmt.Printf("Some non-critical errors occurred during the execution:\n%s", warns)
	}

	return nil
}

// Check if key is in the exclusion list.
func isExcluded(excludedList []string, key string) bool {
	for _, exclude := range excludedList {
		re := regexp.MustCompile(exclude)
		if re.MatchString(key) {
			return true
		}
	}
	return false
}
