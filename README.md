# sb

A simple personal website template with blog functionality.

## Usage

```console
$ git clone https://github.com/vilhelmbergsoe/sb
$ go get
$ go run .
server running on :8080
$ # or with nix flake
$ nix build
$ ./result/bin/sb
```

## Endpoints

`/` home page

`/blog/{url}` blog post page

`/static` static file serve directory

## License

[MIT](https://choosealicense.com/licenses/mit)
