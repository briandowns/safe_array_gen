all: clean
	../bin/safe_array_gen -t int
	cc -o test crosscheck.c tests.c int_slice.c 
	./test

.PHONY: clean
	rm -f ./test
