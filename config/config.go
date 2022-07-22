package config

import (
	"fmt"

	"github.com/grem11n/s3bc/client/s3"
)

// Config stores required information such as the bucket name and desired storage class.
type Config struct {
	Bucket       string
	StorageClass string
	Excluded     []string
	DryRun       bool
}

// Validate checks if storage class is valid and if bucket is reachable.
func (c *Config) Validate() error {
	// TODO: Check if bucket is reachable.

	if c.Bucket == "" {
		return fmt.Errorf("Bucket cannot be empty! Please, specify it with -b or --bucket flag")
	}

	if err := s3.IsValidStorageClass(c.StorageClass); err != nil {
		return err
	}

	return nil
}
