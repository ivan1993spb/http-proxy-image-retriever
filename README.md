
http-proxy-image-retriever
==========================

http-proxy-image-retriever is small http proxy server which:

1. accepts HTTP request with `url` param;
2. downloads html page from passed url;
3. parses html and finds all `<img>`;
4. downloads all found images;
5. generates response html page with found images included into page by data URI scheme.

To install run `go get -u github.com/ivan1993spb/http-proxy-image-retriever`
