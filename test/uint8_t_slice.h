// This is generated code from safe_array_gen. Please do not edit unless
// sure of what you are doing.

#ifdef __cplusplus
extern "C" {
#endif

#ifndef __UINT8_H
#define __UINT8_H

#include <stdbool.h>
#include <stddef.h>
#include <stdint.h>
#include <stdlib.h>

typedef bool (*compare_func_t)(const uint8_t a, const uint8_t b, void *user_data);
typedef void (*iter_func_t)(const uint8_t item, void *user_data);

typedef struct {
    uint8_t *items;
    size_t len;
    size_t cap;
	compare_func_t compare;
} uint8_slice_t;

/**
 * uint8_slice_new creates a pointer of type uint8_slice_t, sets
 * default values, and returns the pointer to the allocated memory. The user
 * is responsible for freeing this memory.
 */
uint8_slice_t*
uint8_slice_new(const size_t cap);

/**
 * uint8_slice_free frees the memory used by the given pointer. 
 */
void
uint8_slice_free(uint8_slice_t *s);

/**
 * uint8_slice_get attempts to retrieve the value at the given index. If
 * the index is out of range, 0 is returned indicating an error.
 */
uint8_t
uint8_slice_get(uint8_slice_t *s, size_t idx);

/**
 * uint8_slice_append attempts to append the data to the given array.
 */
void
uint8_slice_append(uint8_slice_t *s, const uint8_t val);

/**
 * uint8_slice_reverse the contents of the array.
 */
void
uint8_slice_reverse(uint8_slice_t *s);

/**
 * uint8_slice_compare takes 2 slices, compares them element by element
 * and returns true if they are the same and false if they are not.
 */
bool
uint8_slice_compare(const uint8_slice_t *s1, const uint8_slice_t *s2, void *user_data);

/**
 * uint8_slice_copy takes 2 slices. The first is copied into the second
 * with the number of elements copied being returned. The code assumes that 
 * slice 2 has been created large enough to hold the contents of slice 1. If
 * the overwrite option has been selected, the code will make sure there is 
 * enough space in slice 2 and overwrite its contents.
 */
int
uint8_slice_copy(const uint8_slice_t *s1, uint8_slice_t *s2, int overwrite);

/**
 * uint8_slice_contains checks to see if the given value is in the slice.
 */
bool
uint8_slice_contains(const uint8_slice_t *s, uint8_t val);

/**
 * uint8_slice_delete removes the item at the given index and returns the
 * new length.
 */
int
uint8_slice_delete(uint8_slice_t *s, const size_t idx);

/**
 * uint8_slice_replace replaces the value at the given index with the new
 * value.
 */
int
uint8_slice_replace(uint8_slice_t *s, const size_t idx, const uint8_t val);

/**
 * uint8_slice_foreach iterates through the slice and runs the user provided
 * function on each item. Additional user provided data can be provided 
 * using the user_data argument.
 */
int
uint8_slice_foreach(uint8_slice_t *s, iter_func_t ift, void *user_data);

/**
 * uint8_slice_sort uses thet Quick Sort algorithm to sort the contents of the
 * slice if it is a standard type.
 */
int
uint8_slice_sort(uint8_slice_t *s, size_t low, size_t high);

#endif /** end __UINT8_H */
#ifdef __cplusplus
}
#endif
