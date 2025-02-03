
#include <stdint.h>
#include <stdlib.h>

typedef struct {
    uint8_t *items;
    uint64_t len;
    uint64_t cap;
} uint8_slice_t;

/**
 * uint8_slice_new creates a pointer of type uint8_slice_t, sets
 * default values, and returns the pointer to the allocated memory. The user
 * is responsible for freeing this memory.
 */
uint8_slice_t*
uint8_slice_new(const int cap);

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
uint8_slice_get(uint8_slice_t *s, uint64_t idx);

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
