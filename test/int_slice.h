// This is generated code from safe_array_gen. Please do not edit unless
// sure of what you are doing.

#ifdef __cplusplus
extern "C" {
#endif

#ifndef __INT_H
#define __INT_H

#include <stdbool.h>
#include <stddef.h>
#include <stdint.h>
#include <stdlib.h>

typedef bool (*compare_func_t)(const int x, const int y, void *user_data);
typedef void (*foreach_func_t)(const int item, void *user_data);
typedef int (*sort_compare_func_t)(const void *x, const void *y);
typedef bool (*val_equal_func_t)(const int x, const int y, void *user_data);

typedef struct {
    int *items;
    size_t len;
    size_t cap;
	compare_func_t compare;
	sort_compare_func_t sort_compare;
} int_slice_t;

/**
 * int_slice_new creates a pointer of type int_slice_t, sets
 * default values, and returns the pointer to the allocated memory. The user
 * is responsible for freeing this memory.
 */
int_slice_t*
int_slice_new(const size_t cap);

/**
 * int_slice_free frees the memory used by the given pointer. 
 */
void
int_slice_free(int_slice_t *s);

/**
 * int_slice_get attempts to retrieve the value at the given index. If
 * the index is out of range, 0 is returned indicating an error.
 */
int
int_slice_get(int_slice_t *s, size_t idx);

/**
 * int_slice_append attempts to append the data to the given array.
 */
void
int_slice_append(int_slice_t *s, const int val);

/**
 * int_slice_reverse the contents of the array.
 */
void
int_slice_reverse(int_slice_t *s);

/**
 * int_slice_compare takes 2 slices, compares them element by element
 * and returns true if they are the same and false if they are not.
 */
bool
int_slice_compare(const int_slice_t *s1, const int_slice_t *s2, void *user_data);

/**
 * int_slice_copy takes 2 slices. The first is copied into the second
 * with the number of elements copied being returned. The code assumes that 
 * slice 2 has been created large enough to hold the contents of slice 1. If
 * the overwrite option has been selected, the code will make sure there is 
 * enough space in slice 2 and overwrite its contents.
 */
int
int_slice_copy(const int_slice_t *s1, int_slice_t *s2, int overwrite);

/**
 * int_slice_contains checks to see if the given value is in the slice.
 */
bool
int_slice_contains(const int_slice_t *s, int val);

/**
 * int_slice_delete removes the item at the given index and returns the
 * new length.
 */
int
int_slice_delete(int_slice_t *s, const size_t idx);

/**
 * int_slice_replace_by_idx replaces the value at the given index with the new
 * value.
 */
int
int_slice_replace_by_idx(int_slice_t *s, const size_t idx, const int val);

/**
 * int_slice_replace_by_val replaces occurances of the value with the
 * new value the number of times given. 
 */
int
int_slice_replace_by_val(int_slice_t *s, const int old_val, const int new_val, size_t times);

/**
 * int_slice returns the first element of the slice.
 */
int
int_slice_first(int_slice_t *s);

/**
 * int_slice returns the last element of the slice.
 */
int
int_slice_last(int_slice_t *s);

/**
 * int_slice_foreach iterates through the slice and runs the user provided
 * function on each item. Additional user provided data can be provided 
 * using the user_data argument.
 */
int
int_slice_foreach(int_slice_t *s, foreach_func_t ift, void *user_data);

/**
 * int_slice_sort uses thet Quick Sort algorithm to sort the contents
 * of the slice if it is a standard type. When using a custom type for items,
 * like a struct, a sort_compare_func_t needs to be set.
 */
void
int_slice_sort(int_slice_t *s);

/**
 * int_slice_repeat takes a value and repeats that value in the slice
 * for the number of times given.
 */
int
int_slice_repeat(int_slice_t *s, const int val, const size_t times);

/**
 * int_slice_count counts the occurances of the given value.
 */
size_t
int_slice_count(int_slice_t *s, const int val);

/**
 * int_slice_grow Grows the slice by the given size.
 */
size_t
int_slice_grow(int_slice_t *s, const size_t size);

/**
 * int_slice_concat Combine the second slice into the first.
 */
size_t
int_slice_concat(int_slice_t *s1, const int_slice_t *s2);

#endif /** end __INT_H */
#ifdef __cplusplus
}
#endif
