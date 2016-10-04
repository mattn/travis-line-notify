# travis-line-notify

**THIS IS VERY EXPERIMENTAL**

Notify build result on Travis CI via LINE Notify

![](http://go-gyazo.appspot.com/372f28907e461740.png)

## Usage

```
$ travis-line-notify -token <LINE_ACCESS_NOTIFY> <USER/REPOSITORY> <USER/REPOSITORY> ...
```

You can set `$LINE_ACCESS_NOTIFY` for the default value of `-token`. 

## Installation

```
$ go get github.com/mattn/travis-line-notify
```

## License

MIT

## Author

Yasuhiro Matsumoto (a.k.a. mattn)
