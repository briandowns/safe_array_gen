# Safe Array Generator (sag)

sag makes generating safe arrays in easy. Safe arrays are simply arrays but wrapped in functions that perform basic bounds checking. These functions are kept simple such that whatever compiler you're using can optimize them away and take advantage of CPU branch prediction. Arrays will grow dynamically when necessary.

## Functions Generated

* new
* free
* get
* append
* reverse
* compare
* copy
* contains
* delete
* replace

## Usage

sag is CLI driven and takes flags and arguments as input.

E.g.

Generate the "array" and associated functions.

```sh
sag -t uint8_t
```

## Contributing

Please feel free to open a PR!

## Contact

Brian Downs [@bdowns328](http://twitter.com/bdowns328)

## License

BSD 2 Clause [License](/LICENSE).
