
http-proxy-image-retriever
==========================

http-proxy-image-retriever is small http proxy server which:

1. accepts HTTP request with `url` param;
2. downloads html page from passed url;
3. parses html and finds all `<img>`;
4. downloads all found images;
5. generates response html page with found images included into page by data URI scheme.

To install run `go get -u github.com/ivan1993spb/http-proxy-image-retriever`

Testing
-------

* edit files in `test/` directory and run `go generate`;
* run `go test`;
* open file *test_result.html*;
* run `http-proxy-image-retriever`;
* run `curl http://localhost:8888/?url=https%3A%2F%2Fgolang.org%2Fdoc%2F`.

Vegeta testing:

```
$ cat vegeta_targets | vegeta attack -rate=50 -workers=10 -duration=30s | tee results.bin | vegeta report
```
