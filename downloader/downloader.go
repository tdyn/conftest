package downloader

import (
	"context"

	"github.com/containerd/containerd/log"
	getter "github.com/hashicorp/go-getter"
)

var detectors = []getter.Detector{
	new(OCIDetector),
	new(getter.GitHubDetector),
	new(getter.GitDetector),
	new(getter.BitBucketDetector),
	new(getter.S3Detector),
	new(getter.GCSDetector),
	new(getter.FileDetector),
}

var getters = map[string]getter.Getter{
	"file":  new(getter.FileGetter),
	"git":   new(getter.GitGetter),
	"gcs":   new(getter.GCSGetter),
	"hg":    new(getter.HgGetter),
	"s3":    new(getter.S3Getter),
	"oci":   new(OCIGetter),
	"http":  new(getter.HttpGetter),
	"https": new(getter.HttpGetter),
}

// Download downloads the given policies into the given destination
func Download(ctx context.Context, dst string, urls []string) error {
	opts := []getter.ClientOption{}
	for _, url := range urls {
		log.G(ctx).Debugf("Initializing go-getter client with url %v and dst %v", url, dst)
		client := &getter.Client{
			Ctx:       ctx,
			Src:       url,
			Dst:       dst,
			Pwd:       dst,
			Mode:      getter.ClientModeAny,
			Detectors: detectors,
			Getters:   getters,
			Options:   opts,
		}

		if err := client.Get(); err != nil {
			return err
		}
	}

	return nil
}

// Detect determines whether a url is a known source url from which we can download files.
// If a known source is found, the url is formatted, otherwise an error is returned.
func Detect(url string, dst string) (string, error) {
	result, err := getter.Detect(url, dst, detectors)
	if err != nil {
		return "", err
	}

	return result, err
}
