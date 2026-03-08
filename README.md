# `s3-directory-listing`

Custom Go code to generate directory listings based on templates for [CloudFlare R2](https://www.cloudflare.com/developer-platform/products/r2/) (or any [S3](https://en.wikipedia.org/wiki/Amazon_S3)-compatible).

This program is usually ran in the actions that upload to R2 after the upload has been done to regenerate any directory listings with new content.

Assumes your secrets are in the environment.
