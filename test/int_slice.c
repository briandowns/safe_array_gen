// This is generated code from safe_array_gen. Please do not edit unless
// sure of what you are doing.

#include <stdbool.h>
#include <stddef.h>
#include <stdint.h>
#include <stdio.h>
#include <stdlib.h>
#include <string.h>
#include <time.h>

#include "int_slice.h"

int_slice_t*
int_slice_new(const uint64_t cap)
{
    int_slice_t *s = calloc(1, sizeof(int_slice_t));
    s->items = calloc(1, sizeof(int) * cap);
    s->len = 0;
    s->cap = cap;

    return s;
}

void
int_slice_free(int_slice_t *s) {
	if (s != NULL && s->items != NULL) {
        free(s->items);
        free(s);
	} 
}

int
int_slice_get(int_slice_t *s, uint64_t idx)
{
    if (idx >= 0 && idx < s->len) {
        return s->items[idx];
    }
    return 0;
}

void
int_slice_append(int_slice_t *s, const int val)
{
    if (s->len == s->cap) {
        s->cap *= 2;
        s->items = realloc(s->items, sizeof(int) * s->cap);
    }
    s->items[s->len++] = val;
}

void
int_slice_reverse(int_slice_t *s) {
	if (s->len < 2) {
		return;
	}

    uint64_t i = s->len - 1;
    uint64_t j = 0;

    while(i > j) {
        int temp = s->items[i];
        s->items[i] = s->items[j];
        s->items[j] = temp;
        i--;
        j++;
    }
}

bool
int_slice_compare(const int_slice_t *s1, const int_slice_t *s2, compare_func_t compare, void *user_data)
{
	if (s1->len != s2->len) {
		return false;
	}

	for (uint64_t i = 0; i < s1->len; i++) {
		if (!compare(s1->items[i], s2->items[i], user_data)) {
			return false;
		}
	}
	return true;
}

uint64_t
int_slice_copy(const int_slice_t *s1, int_slice_t *s2, bool overwrite)
{
	if (s2->len == 0) {
		return 0;
	}

	if (overwrite) {
		if (s1->len != s2->len) {
			s2->cap = s1->cap;
			s2->items = realloc(s2->items, sizeof(int) * s1->cap);
		}
	}

	for (uint64_t i = 0; i < s1->len; i++) {
		s2->items[i] = s1->items[i];
		s2->len++;
	}
	return s2->len;
}

bool
int_slice_contains(const int_slice_t *s, int val)
{
	if (s->len == 0) {
		return false;
	}

	for (uint64_t i = 0; i < s->len; i++) {
		if (s->items[i] == val) {
			return true;
		}
	}
	return false;
}

int
int_slice_delete(int_slice_t *s, const uint64_t idx)
{
	if (s->len == 0 || idx > s->len) {
		return -1;
	}

	for (uint64_t i = idx; i < s->len; i++) {
		s->items[i] = s->items[i + 1];
	}
	s->len--;

	return s->len;
}

int
int_slice_replace_by_idx(int_slice_t *s, const uint64_t idx, const int val)
{
	if (s->len == 0 || idx > s->len) {
		return -1;
	}
	s->items[idx] = val;

	return 0;
}

int
int_slice_replace_by_val(int_slice_t *s, const int old_val, const int new_val, uint64_t times, compare_func_t compare)
{
	if (s->len == 0) {
		return -1;
	}

	for (uint64_t i = 0; i < s->len && times != 0; i++) {
		if (compare(s->items[i], old_val, NULL)) {
			s->items[i] = new_val;
			times--;
		}
	}
	return 0;
}

int
int_slice_first(int_slice_t *s)
{
	return int_slice_get(s, 0);
}

int
int_slice_last(int_slice_t *s)
{
	return int_slice_get(s, s->len-1);
}

int
int_slice_foreach(int_slice_t *s, foreach_func_t ift, void *user_data)
{
	if (s->len == 0) {
		return 0;
	}
	
	for (uint64_t i = 0; i < s->len; i++) {
		ift(s->items[i], user_data);
	}
	return 0;
}

/**
 * qsort_compare is a simple implementation of the function required to be
 * passed to qsort.
 */
static int
qsort_compare(const void *x, const void *y) {
	return (*(int*)x - *(int*)y);
}

void
int_slice_sort(int_slice_t *s, sort_compare_func_t sort_compare)
{
	if (s->len < 2) {
		return;
	}

	if (sort_compare == NULL) {
		return ;
	}

	qsort(s->items, s->len, sizeof(int), sort_compare);
}

uint64_t
int_slice_repeat(int_slice_t *s, const int val, const uint64_t times)
{
	for (uint64_t i = 0; i < times; i++) {
		int_slice_append(s, val);
	}
	return s->len;
}

uint64_t
int_slice_count(int_slice_t *s, const int val, compare_func_t compare)
{
	uint64_t count = 0;

	if (s->len == 0) {
		return count;
	}

	for (uint64_t i = 0; i < s->len; i++) {
		if (compare(s->items[i], val, NULL)) {
			count++;
		}
	}
	return count;
}

uint64_t
int_slice_grow(int_slice_t *s, const uint64_t size)
{
	if (size == 0) {
		return s->cap;
	}

	s->cap += size;
    s->items = realloc(s->items, sizeof(int) * s->cap);

	return s->cap;
}

uint64_t
int_slice_concat(int_slice_t *s1, const int_slice_t *s2)
{
	if (s2->len == 0) {
		return s1->len;
	}
	
	s1->cap += s2->len;
	s1->items = realloc(s1->items, sizeof(int) * s2->len);

	for (uint64_t i = 0, j = s1->len; i < s2->len; i++, j++) {
		s1->items[j] = s2->items[i];
		s1->len++;
	}
	
	return s1->len;
}
