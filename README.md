# Guts

**Work in progress**

Statsd drop-in(or almost, see config file) replacement written in Go

![Guts](http://img2.wikia.nocookie.net/__cb20121106221315/berserk/images/8/88/Guts_Portrait.jpg)

## How to install
```shell
$ go get github.com/enriclluelles/guts
```

If you have you $GOPATH/bin into your path, that's about it

## Usage

Write a config file that uses the same configuration format as
[statsd](https://github.com/etsy/statsd) but is actually a syntactically
correct JSON file

```json
{
  "graphiteport": 49003,
  "graphitehost": "127.0.0.1",
  "port": 8000,
  "flushinterval": 5000,
  "backends": [ "graphite" ],
  "graphite": {
    "legacynamespace": false,
    "prefixcounter": "ccc",
    "prefixtimer": "timers",
    "globalprefix": "guts"
  },
  "percentthreshold": 80
}
```

And then run `guts config.json` and you're set
