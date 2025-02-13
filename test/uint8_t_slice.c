// This is generated code from safe_array_gen. Please do not edit unless
// sure of what you are doing.

#include <stdbool.h>
#include <stddef.h>
#include <stdint.h>
#include <stdlib.h>
#include <string.h>
#include <time.h>

#include "uint8_t_slice.h"

uint8_slice_t*
uint8_slice_new(const size_t cap)
{
    uint8_slice_t *s = calloc(1, sizeof(uint8_slice_t));
    s->items = calloc(1, sizeof(uint8_t) * cap);
    s->len = 0;
    s->cap = cap;

    return s;
}

void
uint8_slice_free(uint8_slice_t *s) {
	if (s != NULL && s->items != NULL) {
		free(s->items);
    	free(s);
	} 
}

uint8_t
uint8_slice_get(uint8_slice_t *s, size_t idx)
{
    if (idx >= 0 && idx < s->len) {
        return s->items[idx];
    }

    return 0;
}

void
uint8_slice_append(uint8_slice_t *s, const uint8_t val)
{
    if (s->len == s->cap) {
        s->cap *= 2;
        s->items = realloc(s->items, sizeof(uint8_t) * s->cap);
    }

    s->items[s->len++] = val;
}

void
uint8_slice_reverse(uint8_slice_t *s) {
	if (s->len < 2) {
		return;
	}

	uint64_t i = s->len - 1;
    uint64_t j = 0;

    while(i > j) {
        uint8_t temp = s->items[i];
        s->items[i] = s->items[j];
        s->items[j] = temp;
        i--;
        j++;
    }
}

bool
uint8_slice_compare(const uint8_slice_t *s1, const uint8_slice_t *s2, void *user_data)
{
	if (s1->len == 0 && s2->len == 0) {
		return true;
	}

	if (s1->len != s2->len) {
		return false;
	}

	if (s1->compare != NULL) {
		for (size_t i = 0; i < s1->len; i++) {
			if (!s1->compare(s1->items[i], s2->items[i], user_data)) {
				return false;
			}
		}

		return true;
	}

	for (size_t i = 0; i < s1->len; i++) {
		if (s1->items[i] != s2->items[i]) {
			return false;
		}
	}

	return true;
}

int
uint8_slice_copy(const uint8_slice_t *s1, uint8_slice_t *s2, int overwrite)
{
	if (s2->len == 0) {
		return 0;
	}

	if (overwrite) {
		if (s1->len != s2->len) {
			s2->cap = s1->cap;
			s2->items = realloc(s2->items, sizeof(uint8_t) * s1->cap);
		}
	}

	for (size_t i = 0; i < s1->len; i++) {
		s2->items[i] = s1->items[i];
		s2->len++;
	}

	return s2->len;
}

bool
uint8_slice_contains(const uint8_slice_t *s, uint8_t val)
{
	if (s->len == 0) {
		return false;
	}

	for (size_t i = 0; i < s->len; i++) {
		if (s->items[i] == val) {
			return true;
		}
	}

	return false;
}

int
uint8_slice_delete(uint8_slice_t *s, const size_t idx)
{
	if (s->len == 0 || idx > s->len) {
		return -1;
	}

	for (size_t i = idx; i < s->len; i++) {
		s->items[i] = s->items[i + 1];
	}
	s->len--;

	return s->len;
}

int
uint8_slice_replace(uint8_slice_t *s, const size_t idx, const uint8_t val)
{
	if (s->len == 0 || idx > s->len) {
		return -1;
	}

	s->items[idx] = val;

	return 0;
}

uint8_t
uint8_slice_first(uint8_slice_t *s)
{
	if (s->len == 0) {
		return 0;
	}

	return s->items[0];
}

uint8_t
uint8_slice_last(uint8_slice_t *s)
{
	if (s->len == 0) {
		return 0;
	}

	return s->items[s->len-1]; 
}

int
uint8_slice_foreach(uint8_slice_t *s, foreach_func_t ift, void *user_data)
{
	if (s->len == 0) {
		return 0;
	}
	
	for (size_t i = 0; i < s->len; i++) {
		ift(s->items[i], user_data);
	}

	return 0;
}

static int
qsort_compare(const void *x, const void *y) {
	return (*(uint8_t*)x - *(uint8_t*)y);
}

int
uint8_slice_sort(uint8_slice_t *s)
{
	if (s->len < 2) {
		return 0;
	}

	if (s->sort_compare != NULL) {
		qsort(s->items, s->len, sizeof(uint8_t), s->sort_compare);
	} else {
		qsort(s->items, s->len, sizeof(uint8_t), qsort_compare);
	}

	return 0;
}

