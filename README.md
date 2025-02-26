# Safe Array Generator

Safe Array Generator makes generating safe arrays in C easy. Safe arrays are simply arrays but wrapped in functions that perform basic bounds checking. These functions are kept simple such that whatever compiler you're using can optimize them away and take advantage of CPU branch prediction. Arrays will grow dynamically when necessary.

## Generated Functions

All functions take the `(type)_slice_t` type as the first argument. Within that type is the array of the specified type. All operations are then performed on that array.

| Function | Description |
| -------- | ----------- |
| new | Create a new value of type `(type)_slice_t` pointer.|
| free | Free the memory used by the slice and by `(type)_slice_t`.|
| get | Retrieve an item from the slice by index.|
| append | Append a value onto the end of a slice.|
| reverse | Reverse the contents of a slice.|
| compare | Compare 2 slices.|
| copy | Copy the contents from one slice to another.|
| contains | Checks to see if the given value is in the slice.|
| delete | Deletes an item from the slice.|
| replace_by_idx | Replaces an item by index from the slice.|
| replace_by_val | Replaces n occurrences of an item from the slice.|
| first | Retrieves the first item in the slice.|
| last | Retrieves the last item in the slice.|
| foreach | Takes a function to be ran for each element in the slice.|
| sort | When using structs or non integer types, a custom compare func is required to be set.|
| repeat | Repeat a value for the given number of times.|
| count | Counts the number of occurrences of a given value. If compare is not `NULL` that function is used to do the comparison.|
| grow | Grows the slice by the given quantity.|
| concat | Concatenate 2 slices into 1.|

More detailed documentation can be found in the generated header file.

## Usage

Safe Array Generator is CLI driven and takes flags and arguments as input. It can be used for very simple use cases as demonstrated immediately below and more specific ones where custom names are required. User defined functions can be given to the compare and foreach functions which can be useful when uses slices of custom types. 

E.g.

Generate the slice and associated functions.

```sh
safe_array_gen -t uint8_t
```

Generate slice type with a custom name.

```sh
safe_array_gen -t char -n grades
```

## Contributing

Please feel free to open a PR!

## Contact

Brian Downs [@bdowns328](http://twitter.com/bdowns328)

## License

BSD 2 Clause [License](/LICENSE).
