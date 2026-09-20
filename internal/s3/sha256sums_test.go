package s3

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

func testClient(t *testing.T, files map[string]string) *Client {
	t.Helper()

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Path is /bucket/key
		key := r.URL.Path[len("/test-bucket/"):]
		if content, ok := files[key]; ok {
			w.Write([]byte(content))
		} else {
			http.NotFound(w, r)
		}
	}))
	t.Cleanup(server.Close)

	s3Client := s3.NewFromConfig(aws.Config{
		Region:      "us-east-1",
		Credentials: credentials.NewStaticCredentialsProvider("test", "test", ""),
	}, func(o *s3.Options) {
		o.BaseEndpoint = aws.String(server.URL)
		o.UsePathStyle = true
	})

	return &Client{s3Client: s3Client, bucket: "test-bucket"}
}

func TestGenerateSHA256SUMS(t *testing.T) {
	files := map[string]string{
		"release/foo.tar.gz.sha256": "aabbccdd1234567890abcdef1234567890abcdef1234567890abcdef12345678\n",
		"release/bar.tar.gz.sha256": "11223344556677889900aabbccddeeff11223344556677889900aabbccddeeff\n",
		"release/baz.tar.gz.sha256": "ffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff  baz.tar.gz\n",
	}

	client := testClient(t, files)

	listing := &DirectoryListing{
		Path: "release",
		Files: []DirectoryEntry{
			{Name: "foo.tar.gz", Key: "release/foo.tar.gz", Size: 1000},
			{Name: "foo.tar.gz.sha256", Key: "release/foo.tar.gz.sha256", Size: 65},
			{Name: "bar.tar.gz", Key: "release/bar.tar.gz", Size: 2000},
			{Name: "bar.tar.gz.sha256", Key: "release/bar.tar.gz.sha256", Size: 65},
			{Name: "baz.tar.gz", Key: "release/baz.tar.gz", Size: 3000},
			{Name: "baz.tar.gz.sha256", Key: "release/baz.tar.gz.sha256", Size: 80},
		},
	}

	content, err := GenerateSHA256SUMS(context.Background(), client, listing)
	if err != nil {
		t.Fatalf("GenerateSHA256SUMS() error = %v", err)
	}

	want := "11223344556677889900aabbccddeeff11223344556677889900aabbccddeeff  bar.tar.gz\n" +
		"ffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff  baz.tar.gz\n" +
		"aabbccdd1234567890abcdef1234567890abcdef1234567890abcdef12345678  foo.tar.gz\n"

	if content != want {
		t.Errorf("GenerateSHA256SUMS() =\n%s\nwant:\n%s", content, want)
	}
}

func TestGenerateSHA256SUMS_NoSHA256Files(t *testing.T) {
	client := testClient(t, nil)

	listing := &DirectoryListing{
		Path: "release",
		Files: []DirectoryEntry{
			{Name: "foo.tar.gz", Key: "release/foo.tar.gz", Size: 1000},
		},
	}

	content, err := GenerateSHA256SUMS(context.Background(), client, listing)
	if err != nil {
		t.Fatalf("GenerateSHA256SUMS() error = %v", err)
	}
	if content != "" {
		t.Errorf("GenerateSHA256SUMS() = %q, want empty", content)
	}
}
