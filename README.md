

http-proxy-image-retriever [![Build Status](https://travis-ci.org/ivan1993spb/http-proxy-image-retriever.svg?branch=master)](https://travis-ci.org/ivan1993spb/http-proxy-image-retriever) [![Docker Repository on Quay](https://quay.io/repository/ivan1993spb/http-proxy-image-retriever/status "Docker Repository on Quay")](https://quay.io/repository/ivan1993spb/http-proxy-image-retriever) [![Go Report Card](https://goreportcard.com/badge/github.com/ivan1993spb/http-proxy-image-retriever)](https://goreportcard.com/report/github.com/ivan1993spb/http-proxy-image-retriever) [![GitHub release](https://img.shields.io/github/release/ivan1993spb/http-proxy-image-retriever.svg)](https://github.com/ivan1993spb/http-proxy-image-retriever/releases/tag/v1.0)
==========================

http-proxy-image-retriever is small http proxy server which:

1. accepts HTTP request with `url` param;
2. downloads html page from passed url;
3. parses html and finds all `<img>`;
4. downloads all found images;
5. generates response html page with found images included into page by data URI scheme.

To install run `go get -u github.com/ivan1993spb/http-proxy-image-retriever`

Docker image
------------

```bash
docker pull quay.io/ivan1993spb/http-proxy-image-retriever
docker run --name http-proxy-image-retriever -p 8888:8888 -d quay.io/ivan1993spb/http-proxy-image-retriever
# test upyachka
curl http://localhost:8888/?url=http://upyachka.ru/ > upyachka.html
```

Testing
-------

* when edited files in `test/` directory don't forget run `go generate` and fix `image_proxy_handler_test.go`;
* run `go test` and open file *test_result.html*
* run `http-proxy-image-retriever` and then run `curl http://localhost:8888/?url=https%3A%2F%2Fgolang.org%2Fdoc%2F`.

Vegeta testing:

```bash
cat vegeta_targets | vegeta attack -rate=15 -workers=10 -duration=30s | tee results.bin | vegeta report
```

License
-------

See [LICENSE](LICENSE).
