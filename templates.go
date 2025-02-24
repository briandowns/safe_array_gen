package main

const sliceHeaderTmpl = `// This is generated code from safe_array_gen. Please do not edit unless
// sure of what you are doing.

#ifdef __cplusplus
extern "C" {
#endif
{{- $name := Strip .Name "_t" -}}
{{- $headerName := ToUpper $name }}

#ifndef __{{ $headerName }}_H
#define __{{ $headerName }}_H

#include <stdbool.h>
#include <stddef.h>
#include <stdint.h>
#include <stdlib.h>

typedef bool (*compare_func_t)(const {{ .Name }} x, const {{ .Name }} y, void *user_data);
typedef void (*foreach_func_t)(const {{ .Name }} item, void *user_data);
typedef int (*sort_compare_func_t)(const void *x, const void *y);
typedef bool (*val_equal_func_t)(const {{ .Name }} x, const {{ .Name }} y, void *user_data);

{{- $items := "*items" -}}
{{- $typeName := "" }}
{{- $funcPrefix := "" }}
{{- $typeArg := "" -}}
{{- if .CustomName -}}
{{- $typeName = .CustomName }}
{{- $funcPrefix = $typeName }}
{{ $typeArg = printf "%c" (index $typeName 0) -}}
typedef struct {
    {{ .Name }} {{ $items }};
    uint64_t len;
    uint64_t cap;
	compare_func_t compare;
	sort_compare_func_t sort_compare;
} {{ $typeName }};
{{- else -}}
{{- $typeName = printf "%s_slice_t" $name }}
{{- $funcPrefix = printf "%s_slice" $name }}
{{ $typeArg = "s" }}
typedef struct {
    {{ .Name }} {{ $items }};
    uint64_t len;
    uint64_t cap;
	compare_func_t compare;
	sort_compare_func_t sort_compare;
} {{ $typeName }};
{{- end }}

/**
 * {{ $funcPrefix }}_new creates a pointer of type {{ $typeName }}, sets
 * default values, and returns the pointer to the allocated memory. The user
 * is responsible for freeing this memory.
 */
{{ $typeName }}*
{{ $funcPrefix }}_new(const uint64_t cap);

/**
 * {{ $funcPrefix }}_free frees the memory used by the given pointer. 
 */
void
{{ $funcPrefix }}_free({{ $typeName }} *{{ $typeArg }});

/**
 * {{ $funcPrefix }}_get attempts to retrieve the value at the given index. If
 * the index is out of range, 0 is returned indicating an error.
 */
{{ .Name }}
{{ $funcPrefix }}_get({{ $typeName }} *s, uint64_t idx);

/**
 * {{ $funcPrefix }}_append attempts to append the data to the given array.
 */
void
{{ $funcPrefix }}_append({{ $typeName }} *{{ $typeArg }}, const {{ .Name }} val);

/**
 * {{ $funcPrefix }}_reverse the contents of the array.
 */
void
{{ $funcPrefix }}_reverse({{ $typeName }} *{{ $typeArg }});

/**
 * {{ $funcPrefix }}_compare takes 2 slices, compares them element by element
 * and returns true if they are the same and false if they are not.
 */
bool
{{ $funcPrefix }}_compare(const {{ $typeName }} *{{ $typeArg }}1, const {{ $typeName }} *{{ $typeArg }}2, void *user_data);

/**
 * {{ $funcPrefix }}_copy takes 2 slices. The first is copied into the second
 * with the number of elements copied being returned. The code assumes that 
 * slice 2 has been created large enough to hold the contents of slice 1. If
 * the overwrite option has been selected, the code will make sure there is 
 * enough space in slice 2 and overwrite its contents.
 */
uint64_t
{{ $funcPrefix }}_copy(const {{ $typeName }} *{{ $typeArg }}1, {{ $typeName }} *{{ $typeArg }}2, bool overwrite);

/**
 * {{ $funcPrefix }}_contains checks to see if the given value is in the slice.
 */
bool
{{ $funcPrefix }}_contains(const {{ $typeName }} *{{ $typeArg }}, {{ .Name }} val);

/**
 * {{ $funcPrefix }}_delete removes the item at the given index and returns the
 * new length.
 */
int
{{ $funcPrefix }}_delete({{ $typeName }} *{{ $typeArg }}, const uint64_t idx);

/**
 * {{ $funcPrefix }}_replace_by_idx replaces the value at the given index with the new
 * value.
 */
int
{{ $funcPrefix }}_replace_by_idx({{ $typeName }} *{{ $typeArg }}, const uint64_t idx, const {{ .Name }} val);

/**
 * {{ $funcPrefix }}_replace_by_val replaces occurances of the value with the
 * new value the number of times given. 
 */
int
{{ $funcPrefix }}_replace_by_val({{ $typeName }} *{{ $typeArg }}, const {{ .Name }} old_val, const {{ .Name }} new_val, uint64_t times);

/**
 * {{ $funcPrefix }} returns the first element of the slice.
 */
{{ .Name }}
{{ $funcPrefix }}_first({{ $typeName }} *{{ $typeArg }});

/**
 * {{ $funcPrefix }} returns the last element of the slice.
 */
{{ .Name }}
{{ $funcPrefix }}_last({{ $typeName }} *{{ $typeArg }});

/**
 * {{ $funcPrefix }}_foreach iterates through the slice and runs the user provided
 * function on each item. Additional user data can be provided using the
 * user_data argument.
 */
int
{{ $funcPrefix }}_foreach({{ $typeName }} *{{ $typeArg }}, foreach_func_t ift, void *user_data);

/**
 * {{ $funcPrefix }}_sort uses thet Quick Sort algorithm to sort the contents
 * of the slice if it is a standard type. When using a custom type for items,
 * like a struct, a sort_compare_func_t needs to be set.
 */
void
{{ $funcPrefix }}_sort({{ $typeName }} *{{ $typeArg }});

/**
 * {{ $funcPrefix }}_repeat takes a value and repeats that value in the slice
 * for the number of times given and returns the new length of the slice.
 */
uint64_t
{{ $funcPrefix }}_repeat({{ $typeName }} *{{ $typeArg }}, const {{ .Name }} val, const uint64_t times);

/**
 * {{ $funcPrefix }}_count counts the occurances of the given value.
 */
uint64_t
{{ $funcPrefix }}_count({{ $typeName }} *{{ $typeArg }}, const {{ .Name }} val);

/**
 * {{ $funcPrefix }}_grow Grows the slice by the given size.
 */
uint64_t
{{ $funcPrefix }}_grow({{ $typeName }} *{{ $typeArg }}, const uint64_t size);

/**
 * {{ $funcPrefix }}_concat Combine the second slice into the first.
 */
uint64_t
{{ $funcPrefix }}_concat({{ $typeName }} *{{ $typeArg }}1, const {{ $typeName }} *{{ $typeArg }}2);

#endif /** end __{{ $headerName }}_H */
#ifdef __cplusplus
}
#endif
`

const sliceImplementationTmpl = `// This is generated code from safe_array_gen. Please do not edit unless
// sure of what you are doing.

{{- $name := Strip .Name "_t" }}
{{- $arg := "val" }}
{{- $items := "*items" -}}
{{- $typeName := "" }}
{{- $typeArg := "" -}}
{{- $funcPrefix := "" }}

{{- if .CustomName -}}
{{- $typeName = .CustomName }}
{{- $funcPrefix = $typeName }}
{{- $typeArg = printf "%c" (index $typeName 0) -}}
{{- else -}}
{{- $typeName = printf "%s_slice_t" $name }}
{{- $funcPrefix = printf "%s_slice" $name }}
{{- $typeArg = "s" -}}
{{- end }}

{{- if not .Append }}

#include <stdbool.h>
#include <stddef.h>
#include <stdint.h>
#include <stdio.h>
#include <stdlib.h>
#include <string.h>
#include <time.h>

#include "{{ .Name }}_slice.h"
{{- else }}

typedef struct {
    {{ .Name }} {{ $items }};
    uint64_t len;
    uint64_t cap;
	compare_func_t compare;
	sort_compare_func_t sort_compare;
} {{ $typeName }};
{{- end -}}
{{- $name := Strip .Name "_t" }}

{{ $typeName }}*
{{ $funcPrefix }}_new(const uint64_t cap)
{
    {{ $typeName }} *{{ $typeArg }} = calloc(1, sizeof({{ $typeName }}));
    {{ $typeArg }}->items = calloc(1, sizeof({{ .Name }}) * cap);
    {{ $typeArg }}->len = 0;
    {{ $typeArg }}->cap = cap;

    return {{ $typeArg }};
}

void
{{ $funcPrefix }}_free({{ $typeName }} *{{ $typeArg }}) {
	if ({{ $typeArg }} != NULL && {{ $typeArg }}->items != NULL) {
        free({{ $typeArg }}->items);
        free({{ $typeArg }});
	} 
}

{{ .Name }}
{{ $funcPrefix }}_get({{ $typeName }} *{{ $typeArg }}, uint64_t idx)
{
    if (idx >= 0 && idx < {{ $typeArg }}->len) {
        return {{ $typeArg }}->items[idx];
    }
    return 0;
}

void
{{ $funcPrefix }}_append({{ $typeName }} *{{ $typeArg }}, const {{ .Name }} {{ $arg }})
{
    if ({{ $typeArg }}->len == {{ $typeArg }}->cap) {
        {{ $typeArg }}->cap *= 2;
        {{ $typeArg }}->items = realloc({{ $typeArg }}->items, sizeof({{ .Name }}) * {{ $typeArg }}->cap);
    }
    {{ $typeArg }}->items[{{ $typeArg }}->len++] = val;
}

void
{{ $funcPrefix }}_reverse({{ $typeName }} *{{ $typeArg }}) {
	if ({{ $typeArg }}->len < 2) {
		return;
	}

    uint64_t i = {{ $typeArg }}->len - 1;
    uint64_t j = 0;

    while(i > j) {
        {{ .Name }} temp = {{ $typeArg }}->items[i];
        {{ $typeArg }}->items[i] = {{ $typeArg }}->items[j];
        {{ $typeArg }}->items[j] = temp;
        i--;
        j++;
    }
}

bool
{{ $funcPrefix }}_compare(const {{ $typeName }} *{{ $typeArg }}1, const {{ $typeName }} *{{ $typeArg }}2, void *user_data)
{
	if ({{ $typeArg }}1->len == 0 && {{ $typeArg }}2->len == 0) {
		return true;
	}
	if ({{ $typeArg }}1->len != {{ $typeArg }}2->len) {
		return false;
	}

	if ({{ $typeArg }}1->compare != NULL) {
		for (uint64_t i = 0; i < {{ $typeArg }}1->len; i++) {
			if (!{{ $typeArg }}1->compare({{ $typeArg }}1->items[i], {{ $typeArg }}2->items[i], user_data)) {
				return false;
			}
		}
		return true;
	}

	for (uint64_t i = 0; i < {{ $typeArg }}1->len; i++) {
		if ({{ $typeArg }}1->items[i] != {{ $typeArg }}2->items[i]) {
			return false;
		}
	}
	return true;
}

uint64_t
{{ $funcPrefix }}_copy(const {{ $typeName }} *{{ $typeArg }}1, {{ $typeName }} *{{ $typeArg }}2, bool overwrite)
{
	if ({{ $typeArg }}2->len == 0) {
		return 0;
	}

	if (overwrite) {
		if ({{ $typeArg }}1->len != {{ $typeArg }}2->len) {
			{{ $typeArg }}2->cap = {{ $typeArg }}1->cap;
			{{ $typeArg }}2->items = realloc({{ $typeArg }}2->items, sizeof({{ .Name }}) * {{ $typeArg }}1->cap);
		}
	}

	for (uint64_t i = 0; i < {{ $typeArg }}1->len; i++) {
		{{ $typeArg }}2->items[i] = {{ $typeArg }}1->items[i];
		{{ $typeArg }}2->len++;
	}
	return {{ $typeArg }}2->len;
}

bool
{{ $funcPrefix }}_contains(const {{ $typeName }} *{{ $typeArg }}, {{ .Name }} {{ $arg }})
{
	if ({{ $typeArg }}->len == 0) {
		return false;
	}

	for (uint64_t i = 0; i < {{ $typeArg }}->len; i++) {
		if ({{ $typeArg }}->items[i] == val) {
			return true;
		}
	}
	return false;
}

int
{{ $funcPrefix }}_delete({{ $typeName }} *{{ $typeArg }}, const uint64_t idx)
{
	if ({{ $typeArg }}->len == 0 || idx > {{ $typeArg }}->len) {
		return -1;
	}

	for (uint64_t i = idx; i < {{ $typeArg }}->len; i++) {
		{{ $typeArg }}->items[i] = {{ $typeArg }}->items[i + 1];
	}
	{{ $typeArg }}->len--;

	return {{ $typeArg }}->len;
}

int
{{ $funcPrefix }}_replace_by_idx({{ $typeName }} *{{ $typeArg }}, const uint64_t idx, const {{ .Name }} {{ $arg }})
{
	if ({{ $typeArg }}->len == 0 || idx > {{ $typeArg }}->len) {
		return -1;
	}
	{{ $typeArg }}->items[idx] = {{ $arg }};

	return 0;
}

int
{{ $funcPrefix }}_replace_by_val({{ $typeName }} *{{ $typeArg }}, const {{ .Name }} old_val, const {{ .Name }} new_val, uint64_t times)
{
	if ({{ $typeArg }}->len == 0) {
		return -1;
	}

	for (uint64_t i = 0; i < {{ $typeArg }}->len && times != 0; i++) {
		if ({{ $typeArg }}->compare({{ $typeArg }}->items[i], old_val, NULL)) {
			{{ $typeArg }}->items[i] = new_val;
			times--;
		}
	}
	return 0;
}

{{ .Name }}
{{ $funcPrefix }}_first({{ $typeName }} *{{ $typeArg }})
{
	return {{ $funcPrefix }}_get({{ $typeArg }}, 0);
}

{{ .Name }}
{{ $funcPrefix }}_last({{ $typeName }} *{{ $typeArg }})
{
	return {{ $funcPrefix }}_get({{ $typeArg }}, {{ $typeArg }}->len-1);
}

int
{{ $funcPrefix }}_foreach({{ $typeName }} *{{ $typeArg }}, foreach_func_t ift, void *user_data)
{
	if ({{ $typeArg }}->len == 0) {
		return 0;
	}
	
	for (uint64_t i = 0; i < {{ $typeArg }}->len; i++) {
		ift({{ $typeArg }}->items[i], user_data);
	}
	return 0;
}

/**
 * qsort_compare is a simple implementation of the function required to be
 * passed to qsort.
 */
static int
qsort_compare(const void *x, const void *y) {
	return (*({{ .Name }}*)x - *({{ .Name }}*)y);
}

void
{{ $funcPrefix }}_sort({{ $typeName }} *{{ $typeArg }})
{
	if ({{ $typeArg }}->len < 2) {
		return;
	}

	if ({{ $typeArg }}->sort_compare != NULL) {
		qsort({{ $typeArg }}->items, {{ $typeArg }}->len, sizeof({{ .Name }}), {{ $typeArg }}->sort_compare);
	} else {
		qsort({{ $typeArg }}->items, {{ $typeArg }}->len, sizeof({{ .Name }}), qsort_compare);
	}
}

uint64_t
{{ $funcPrefix }}_repeat({{ $typeName }} *{{ $typeArg }}, const {{ .Name }} val, const uint64_t times)
{
	for (uint64_t i = 0; i < times; i++) {
		{{ $funcPrefix }}_append({{ $typeArg }}, val);
	}
	return {{ $typeArg }}->len;
}

uint64_t
{{ $funcPrefix }}_count({{ $typeName }} *{{ $typeArg }}, const {{ .Name }} val)
{
	uint64_t count = 0;

	if ({{ $typeArg }}->len == 0) {
		return count;
	}


	if ({{ $typeArg }}->compare != NULL) {
		for (uint64_t i = 0; i < {{ $typeArg }}->len; i++) {
			if ({{ $typeArg }}->compare({{ $typeArg }}->items[i], val, NULL)) {
				count++;
			}
		}
	} else {
		for (uint64_t i = 0; i < {{ $typeArg }}->len; i++) {
			if ({{ $typeArg }}->items[i] == val) {
				count++;
			}
		}
	}
	return count;
}

uint64_t
{{ $funcPrefix }}_grow({{ $typeName }} *{{ $typeArg }}, const uint64_t size)
{
	if (size == 0) {
		return {{ $typeArg }}->cap;
	}

	{{ $typeArg }}->cap += size;
    {{ $typeArg }}->items = realloc({{ $typeArg }}->items, sizeof({{ .Name }}) * {{ $typeArg }}->cap);

	return {{ $typeArg }}->cap;
}

uint64_t
{{ $funcPrefix }}_concat({{ $typeName }} *{{ $typeArg }}1, const {{ $typeName }} *{{ $typeArg }}2)
{
	if ({{ $typeArg }}2->len == 0) {
		return {{ $typeArg }}1->len;
	}
	
	{{ $typeArg }}1->cap += {{ $typeArg }}2->len;
	{{ $typeArg }}1->items = realloc({{ $typeArg }}1->items, sizeof({{ .Name }}) * {{ $typeArg }}2->len);

	for (uint64_t i = 0, j = {{ $typeArg }}1->len; i < {{ $typeArg }}2->len; i++, j++) {
		{{ $typeArg }}1->items[j] = {{ $typeArg }}2->items[i];
		{{ $typeArg }}1->len++;
	}
	
	return {{ $typeArg }}1->len;
}
`
