package config

import (
	"fmt"
	"sync"

	"github.com/grem11n/s3bc/client/s3"
	flag "github.com/spf13/pflag"
)

// Config stores required information such as the bucket name and desired storage class.
type Config struct {
	Bucket       string
	StorageClass string
	Excluded     []string
	DryRun       bool
}

var configInstance *Config
var lock = &sync.Mutex{}

// GetConfig gets the flagSet for a command and returns a single instance of Config.
func GetConfig(flags *flag.FlagSet) *Config {
	if configInstance == nil {
		lock.Lock()
		defer lock.Unlock()
		if configInstance == nil {
			bucket, _ := flags.GetString("bucket")
			storageClass, _ := flags.GetString("storage-class")
			exclude, _ := flags.GetStringSlice("exclude")
			dryRun, _ := flags.GetBool("dry-run")

			configInstance = &Config{
				Bucket:       bucket,
				StorageClass: storageClass,
				Excluded:     exclude,
				DryRun:       dryRun,
			}
		}
	}
	return configInstance
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
