package s3

import (
	"testing"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go-v2/service/s3/types"
	assert "github.com/stretchr/testify/assert"
)

var testObject0 = types.Object{
	Key:          aws.String("foo/bar"),
	LastModified: aws.Time(time.Now()),
	StorageClass: types.ObjectStorageClassStandard,
}

var testObject1 = types.Object{
	Key:          aws.String("bar/bazz"),
	LastModified: aws.Time(time.Now()),
	StorageClass: types.ObjectStorageClassReducedRedundancy,
}

func TestmatchedStorageClass_true(t *testing.T) {
	ok, err := matchedStorageClass(&testObject0, "Standard")
	assert.NoError(t, err)
	assert.True(t, ok)
}

func TestmatchedStorageClass_false(t *testing.T) {
	ok, err := matchedStorageClass(&testObject0, "reduced_redundancy")
	assert.NoError(t, err)
	assert.False(t, ok)
}

func TestmatchedStorageClass_error(t *testing.T) {
	_, err := matchedStorageClass(&testObject0, "no_class")
	assert.Error(t, err)
}

func TestSetObjectStorageClass(t *testing.T) {
	var expected = s3.CopyObjectInput{
		Bucket:       aws.String("testB"),
		CopySource:   aws.String("testB/foo/bar"),
		Key:          aws.String("foo/bar"),
		StorageClass: types.StorageClassOnezoneIa,
	}

	client, err := New("testB")
	assert.NoError(t, err)

	gotPtr := client.setObjectStorageClass(&testObject0, "ONEZONE_IA")
	assert.Equal(t, expected, *gotPtr)
}

func TestCreateInputList(t *testing.T) {
	objects := []*types.Object{&testObject0, &testObject1}
	expected := s3.CopyObjectInput{
		Bucket:       aws.String("testB"),
		CopySource:   aws.String("testB/foo/bar"),
		Key:          aws.String("foo/bar"),
		StorageClass: types.StorageClassReducedRedundancy,
	}

	client, err := New("testB")
	assert.NoError(t, err)

	got, err := client.createInputList(objects, "reduced_redundancy")
	assert.NoError(t, err)
	assert.Equal(t, expected, *got[0])
}
