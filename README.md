# gosax

For some background information, motivation and a basic usage sample see
[this blog post](https://eli.thegreenplace.net/2019/faster-xml-stream-processing-in-go/).

## Building

To install libxml2 for development, do:

```
$ sudo apt-get install libxml2 libxml2-dev
```

Or build from source. You would need it installed to use ``gosax``.

## License

gosax's code is in the public domain (see `LICENSE` for details). It uses
[libxml](http://www.xmlsoft.org/index.html), which has its own (MIT) license.

The `pointer` directory vendors https://github.com/mattn/go-pointer/ with some
modifications. `go-pointer` has an MIT license, which applies to my version too.
