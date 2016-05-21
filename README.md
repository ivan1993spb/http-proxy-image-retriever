
http-proxy-image-retriever
==========================

http-proxy-image-retriever is small http proxy server which:

1. accepts http request with `url` param;
2. downloads html page from passed url;
3. parses html and finds all images;
4. downloads all found images;
5. generates response html page with found images included into page by data URI scheme.

Tests:

1. проверить разные кодировки
2. проверить протоколы http и https
3. три случая img src: `path/to/image.png`, `/path/to/image.png`, `http://ex.ple/path/to/image.png`
4. предусмотреть timeout для скачивания
