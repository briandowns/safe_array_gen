all: clean
	../bin/safe_array_gen -t int
	cc -o test main.c int_slice.c
	./test

.PHONY: clean
	rm -f ./test
