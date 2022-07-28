package s3

import (
	"context"
	"fmt"
	"strings"
	"sync"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go-v2/service/s3/types"
	multierror "github.com/hashicorp/go-multierror"
)

type Client struct {
	S3Client *s3.Client
	Bucket   string
}

// CopyObjectInput is an alias for s3.CopyObjectInput to avoid multiple imports
type CopyObjectInput = s3.CopyObjectInput

// New creates a new AWS client based on the default config.
func New(bucket string) (*Client, error) {
	cfg, err := config.LoadDefaultConfig(context.TODO())
	if err != nil {
		return nil, err
	}

	s3Client := s3.NewFromConfig(cfg)

	return &Client{
		S3Client: s3Client,
		Bucket:   bucket,
	}, nil
}

// GetBucketObjects gets the list of the objects in a bucket and their storage class.
func (c *Client) GetBucketObjects() ([]*types.Object, error) {
	var s3objects []*types.Object
	var errs error

	paginator := s3.NewListObjectsV2Paginator(c.S3Client, &s3.ListObjectsV2Input{
		Bucket: aws.String(c.Bucket),
	})

	for paginator.HasMorePages() {
		page, err := paginator.NextPage(context.TODO())
		if err != nil {
			errs = multierror.Append(errs, err)
		}

		for _, obj := range page.Contents {
			s3objects = append(s3objects, &obj)
		}
	}

	return s3objects, errs
}

// UpdateObjectStorageClass updates Storage Class for an object in S3 bucket.
func (c *Client) UpdateObjectStorageClass(input *s3.CopyObjectInput) error {
	_, err := c.S3Client.CopyObject(context.TODO(), input)
	if err != nil {
		return err
	}

	return nil
}

// IsValidStorageClass validate if provided Storage Class is valid.
func IsValidStorageClass(storageClass string) error {
	supportedValues := types.ObjectStorageClass.Values("")
	notSupported := false
	for _, value := range supportedValues {
		if strings.ToUpper(storageClass) == string(value) {
			notSupported = false
			break
		}
		notSupported = true
	}
	if notSupported {
		return fmt.Errorf("Provided Storage Class \"%s\" value doesn't match one of suppoerted by AWS. Supported values are:\n%v", storageClass, supportedValues)
	}

	return nil
}

// GetConvertableObjects return a list of objects to convert to the new Storage Class.
func (c *Client) GetConvertableObjects(storageClass string) ([]*s3.CopyObjectInput, error) {
	fmt.Println("Retrieving bucket objects...")
	objects, err := c.GetBucketObjects()
	if err != nil {
		return nil, err
	}
	objCount := len(objects)

	fmt.Printf("%v objects found in %s bucket", objCount, c.Bucket)

	convertables, err := c.createInputList(objects, storageClass)
	if err != nil {
		return nil, err
	}

	return convertables, nil
}

// createInputList returns []s3.CopyObjectInput list that contains objects to be sent to AWS API.
// If object's storage class already satisfies the condition it's not added to the list.
func (c *Client) createInputList(objects []*types.Object, storageClass string) ([]*s3.CopyObjectInput, error) {
	inputs := []*s3.CopyObjectInput{}
	var wg sync.WaitGroup
	var errs error

	for _, obj := range objects {
		gObj := *obj
		wg.Add(1)

		go func() {
			defer wg.Done()
			ok, err := matchedStorageClass(&gObj, storageClass)
			if err != nil {
				errs = multierror.Append(errs, err)
			}
			if !ok {
				input := c.setObjectStorageClass(&gObj, storageClass)
				inputs = append(inputs, input)
			}
		}()
		wg.Wait()
	}

	return inputs, errs
}

// setObjectStorageClass updates the s3.Object StorageClass accordingly.
func (c *Client) setObjectStorageClass(object *types.Object, storageClass string) *s3.CopyObjectInput {
	input := s3.CopyObjectInput{
		Bucket:       aws.String(c.Bucket),
		CopySource:   aws.String(fmt.Sprintf("%s/%s", c.Bucket, *object.Key)),
		Key:          aws.String(*object.Key),
		StorageClass: types.StorageClass(strings.ToUpper(storageClass)),
	}

	return &input
}

// matchedStorageClass checks if the Storage Class of a provided object matches desired Storage Class.
// If provided string doesn't correspond to any of the known Storage Classes, returns error.
func matchedStorageClass(object *types.Object, storageClass string) (bool, error) {
	if err := IsValidStorageClass(storageClass); err != nil {
		return false, err
	}

	if string(object.StorageClass) == strings.ToUpper(storageClass) {
		return true, nil
	}

	return false, nil
}
