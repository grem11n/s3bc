# S3 Bulk Convert
---

[![SWUbanner](https://raw.githubusercontent.com/vshymanskyy/StandWithUkraine/main/banner2-direct.svg)](https://github.com/vshymanskyy/StandWithUkraine/blob/main/docs/README.md)

### Additional information for users from Russia and Belarus

* Russia has [illegally annexed Crimea in 2014](https://en.wikipedia.org/wiki/Annexation_of_Crimea_by_the_Russian_Federation) and [brought the war in Donbas](https://en.wikipedia.org/wiki/War_in_Donbas) followed by [full-scale invasion of Ukraine in 2022](https://en.wikipedia.org/wiki/2022_Russian_invasion_of_Ukraine).
* Russia has brought sorrow and devastations to millions of Ukrainians, killed hundreds of innocent people, damaged thousands of buildings, and forced several million people to flee.
* [Putin khuylo!](https://en.wikipedia.org/wiki/Putin_khuylo!)

Glory to Ukraine! 🇺🇦

---

## About this tool

S3BC stands for S3 Bulk Convert. This is a simple CLI tool, which can update the storage class of files in an AWS S3 bucket or in a compatible storage solution such as [Minio](https://min.io/)

The inspiration for this tool came from a very old project, rather a task, to convert a bunch of buckets to the `REDUCED_REDUNDANCY` storage type for the cost saving purpose.

Today you can use [AWS S3 lifecycle rules](https://docs.aws.amazon.com/AmazonS3/latest/userguide/object-lifecycle-mgmt.html) to achieve the same goal. Nevertheless, this tool is completely functional and you can use it to change the storage class for the whole bucket or certain files in it.

Please, keep in mind that this is just a fun side-project. Some functionality is still WIP and might not be thoroughly tested. If you have found a bug or just want to suggest a general improvement, feel free to [create a new issue](https://github.com/grem11n/s3bc/issues/new/choose) or [open a pull request](https://github.com/grem11n/s3bc/compare).


## Installation

You can compile an S3BC binary yourself or use it in a Docker container. Releases are still work in progress.

**To build the binary locally**:

1. Ensure that you have Go installed. Since this project uses vendored Go modules, there's no need to download them separately.
2. Clone this repository:
    ```
    git clone https://github.com/grem11n/s3bc.git
    ```
2. Run:
    ```
    make build
    ```

This will compile the binary and put it in the `bin/` directory.

**Alternatively using Docker**:

1. Make sure that `docker` command is available in your system.
2. Clone this repository:
    ```
    git clone https://github.com/grem11n/s3bc.git
    ```
2. Run:
    ```
    make docker-build-dev
    ```

This will build an `s3bc:dev` image.

## Usage

You can use `help` to get hints on how to use S3BC:

```bash
❯ bin/s3bc help
S3BC or S3 Bulk Convert is a CLI tool to update the storage class of the files in an AWS S3 bucket.

Usage:
  s3bc [command]

Available Commands:
  completion  Generate the autocompletion script for the specified shell
  convert     Bulk convert objects in an S3 bucket to the given storage class
  help        Help about any command
  validate    Check if files in a bucket have desired storage class.
  version     Print s3bc version.

Flags:
  -b, --bucket string          Target S3 bucket
  -h, --help                   help for s3bc
  -s, --storage-class string   Storage class to set (default "STANDARD")
  -v, --verbose                Verbose output

Use "s3bc [command] --help" for more information about a command.
```

### Updating the storage class of the files

This is the default action. You can use `s3bc convert [arguments]` or just `s3bc [arguments]` for this.

To exclude certain file patterns from the conversion, you can use `-e` or `--exclude` flag with a regexp pattern for exclusion. **Warning!** This functionality is not yet tested!

```bash
s3bc convert -b <bucket> -s standard -e '.*\.css$'
```

### Validating objects in a bucket

You can use S3BC to check if all the objects in a bucket have a desired storage class. There are two modes of this check:

**Normal mode**: simply tells you whether objects in a bucket have a desired storage class. If not, exits with `1` exit code.

```bash
❯ s3bc validate -b test -s STANDARD
Retrieving bucket objects...
2001 objects found in test bucket
Not all the objects in the "test" bucket have desired storage class
Desired storage class: STANDARD
2001 files in "test" bucket have different storage class.
To get the list of the files, use "--verbose" of "-v" flag.
exit status 1
```

**Verbose mode**: Outputs the list of files that do not comply in `stdout`. Be careful, this list may be huge depending on how many files you have in a bucket. To use S3BC in the verbose mode, add `--verbose` or `-v` flag.


## Development

S3BC is a fun project, in which I tried to get myself familiar with writting CLI applications in Go. I'm using [Cobra](https://cobra.dev/) library (without Viper) as well as a simple plain directory structure.

```
.
├── ...
├── action   # Contains logic for each subcommand
├── build    # Contains Dockerfile as well as Docker Compose file for Minio
├── client   # AWS-related client code
├── cmd      # Cobra commands definitions
├── config   # Common config for the app
├── testdata # Dummy data for Minio. Used for local E2E testing
└── version  # Code for the `version` subcommand
```

If you spot a bug or want to contribute to this project, feel free [to open a pull request](https://github.com/grem11n/s3bc/compare)!

## Testing

This project has CI based on [the GitHub Actions](https://docs.github.com/en/actions). You can find the CI configuration in the `.github/workflows` directory.

All the tasks are automated with the `Makefile`. You can use `make lint` and `make test` to run linters and unit tests for this app.

You can also run linter and tests in a Docker container, but you need to build a test image first. You can do that with `make docker-build-dev` and then: `make docker-lint && make docker-test`.

### E2E tests

S3BC is tested using [Minio](https://min.io/) to mimic AWS S3 API locally. Potentially you can also use [Localstack](https://localstack.cloud/) for E2E tests. However, Minio keeps the data on the filesystem by default, so you don't have to re-generate test fixtures each time you run E2E tests. It saves some time. The downside is that you need to keep those fixtures in the repository. Minio's data can be found in the `testdata/` directory.

You can use Docker Compose to start a Minio server locally:

```bash
# Starts Minio server in the background
docker-compose -f build/docker-compose-minio.yaml up -d
```

After that Minio API should be available on `http://127.0.0.1:9000`.

You can provide a custom AWS Endpoint to S3BC using `AWS_URL` environment variable. For example:

```bash
AWS_URL="http://localhost:9000" ./bin/s3bc validate -b test -s STANDARD
```

or even:

```bash
AWS_URL="http://localhost:9000" go run . validate -b test -s STANDARD
```

If you'd like Localstack better or just want to use a different set of the test data, you can easily generate some objects with a simple script like below (mind the `--endpoint-url`):

```bash
# Create a new test bucket
aws --endpoint-url=http://127.0.0.1:9000 s3api create-bucket --bucket=test --region=us-east-1

# Populate this bucket with files
for i in {0..2000}; do uuidgen > ./test-file-${i}.txt && aws --endpoint-url=http://localhost:9000 s3 cp ./test-file-${i}.txt s3://test/ && rm -f ./test-file-${i}.txt ; done
```

Automation for E2E tests is currently work in progress.


## About the author

You can find more information about me and my work at [https://grem1.in](https://grem1.in)


## License
Apache 2 Licensed. See [LICENSE](./LICENSE.md) for details.
