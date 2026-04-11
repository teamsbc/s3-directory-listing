# `s3-directory-listing`

Custom Go code to generate directory listings based on templates for [CloudFlare R2](https://www.cloudflare.com/developer-platform/products/r2/) (or any [S3](https://en.wikipedia.org/wiki/Amazon_S3)-compatible).

This program is usually run in the actions that upload to R2 after the upload has been done to regenerate any directory listings with new content.

## Installation

```bash
go install github.com/teamsbc/s3-directory-listing/cmd/s3-directory-listing@latest
```

Or build from source:

```bash
make build
```

## Configuration

The tool uses AWS SDK v2 for Go, which supports multiple credential sources:
- Environment variables (`AWS_ACCESS_KEY_ID`, `AWS_SECRET_ACCESS_KEY`)
- AWS credentials file (`~/.aws/credentials`)
- AWS config file (`~/.aws/config`)
- IAM roles (when running on EC2)

### CloudFlare R2

Using environment variables:
```bash
export AWS_ENDPOINT_URL=https://your-account-id.r2.cloudflarestorage.com
export AWS_EC2_METADATA_DISABLED=true  # Disable EC2 IMDS when not on EC2
```

Or configure a profile in `~/.aws/config`:
```ini
[profile r2]
region = auto
endpoint_url = https://your-account-id.r2.cloudflarestorage.com
```

And credentials in `~/.aws/credentials`:
```ini
[r2]
aws_access_key_id = your_access_key_id
aws_secret_access_key = your_secret_access_key
```

Then use the profile with `--profile r2` flag.

### Standard AWS S3

Using environment variables:
```bash
export AWS_ACCESS_KEY_ID=your_access_key_id
export AWS_SECRET_ACCESS_KEY=your_secret_access_key
export AWS_REGION=us-east-1  # or your preferred region
export AWS_EC2_METADATA_DISABLED=true  # If not running on EC2
```

Or configure in `~/.aws/credentials` and `~/.aws/config`.

### Other S3-Compatible Services

```bash
export AWS_ENDPOINT_URL=https://your-s3-compatible-endpoint.com
export AWS_EC2_METADATA_DISABLED=true  # Disable EC2 IMDS when not on EC2
```

**Note:** If you're not running on EC2 and experience slow startup or credential errors, set `AWS_EC2_METADATA_DISABLED=true` to prevent the SDK from attempting to contact the EC2 instance metadata service.

## Usage

Generate directory listings for a bucket:

```bash
s3-directory-listing create <bucket-name> <template-file> [flags]
```

### Flags

- `-o, --output`: Output directory for generated listings (default: current directory)
- `-p, --profile`: AWS profile to use from config file

### Examples

Using environment variables:
```bash
export AWS_ENDPOINT_URL=https://your-account-id.r2.cloudflarestorage.com
export AWS_EC2_METADATA_DISABLED=true
s3-directory-listing create my-bucket templates/default.html -o ./output
```

Using an AWS profile:
```bash
AWS_EC2_METADATA_DISABLED=true s3-directory-listing create my-bucket templates/default.html -o ./output --profile r2
```

## Templates

Templates use Go's `html/template` syntax. Each template receives the following data:

### Template Variables

- `.Path` (string): The current directory path (empty string for root)
- `.Directories` ([]DirectoryEntry): Array of subdirectories in this directory
- `.Files` ([]DirectoryEntry): Array of files in this directory

### DirectoryEntry Structure

- `.Name` (string): Name of the directory or file
- `.Size` (int64): Size in bytes (0 for directories)

### Example Template

See `templates/default.html` for a complete example. Basic usage:

```html
<!DOCTYPE html>
<html>
<head>
    <title>Index of /{{ .Path }}</title>
</head>
<body>
    <h1>Index of /{{ .Path }}</h1>
    
    <h2>Directories</h2>
    <ul>
    {{ range .Directories }}
        <li><a href="{{ .Name }}/">{{ .Name }}/</a></li>
    {{ end }}
    </ul>
    
    <h2>Files</h2>
    <ul>
    {{ range .Files }}
        <li><a href="{{ .Name }}">{{ .Name }}</a> ({{ .Size }} bytes)</li>
    {{ end }}
    </ul>
</body>
</html>
```

## How It Works

1. Connects to the S3 bucket using AWS SDK
2. Lists all objects in the bucket
3. Recursively builds directory structure
4. For each directory, renders a template with:
   - List of subdirectories
   - List of files
5. Outputs `index.html` files for each directory
