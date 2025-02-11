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

{{- $arg := "" -}}
{{- $items := "" -}}
{{- if .Pointer -}}
{{- $arg = "*val" -}}
{{- $items = "**items" -}}
{{- else -}}
{{- $arg = "val" -}}
{{- $items = "*items" -}}
{{ end }}

typedef bool (*compare_func_t)(const {{ .Name }} a, const {{ .Name }} b, void *user_data);
typedef void (*iter_func_t)(const {{ .Name }} item, void *user_data);

{{- $typeName := "" }}
{{- $funcPrefix := "" }}
{{- $typeArg := "" -}}
{{- if .CustomName -}}
{{- $typeName = .CustomName }}
{{- $funcPrefix = $typeName }}
{{ $typeArg = printf "%c" (index $typeName 0) -}}
typedef struct {
    {{ .Name }} {{ $items }};
    size_t len;
    size_t cap;
	compare_func_t compare;
} {{ $typeName }};
{{- else -}}
{{- $typeName = printf "%s_slice_t" $name }}
{{- $funcPrefix = printf "%s_slice" $name }}
{{ $typeArg = "s" }}
typedef struct {
    {{ .Name }} {{ $items }};
    size_t len;
    size_t cap;
	compare_func_t compare;
} {{ $typeName }};
{{- end }}

/**
 * {{ $funcPrefix }}_new creates a pointer of type {{ $typeName }}, sets
 * default values, and returns the pointer to the allocated memory. The user
 * is responsible for freeing this memory.
 */
{{ $typeName }}*
{{ $funcPrefix }}_new(const size_t cap);

/**
 * {{ $funcPrefix }}_free frees the memory used by the given pointer. 
 */
void
{{ $funcPrefix }}_free({{ $typeName }} *{{ $typeArg }});

/**
 * {{ $funcPrefix }}_get attempts to retrieve the value at the given index. If
 * the index is out of range, 0 is returned indicating an error.
 */
{{- if .Pointer }}
{{ .Name }}*
{{- else }}
{{ .Name }}
{{- end }}
{{ $funcPrefix }}_get({{ $typeName }} *s, size_t idx);

/**
 * {{ $funcPrefix }}_append attempts to append the data to the given array.
 */
void
{{ $funcPrefix }}_append({{ $typeName }} *{{ $typeArg }}, const {{ .Name }} {{ $arg }});

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
int
{{ $funcPrefix }}_copy(const {{ $typeName }} *{{ $typeArg }}1, {{ $typeName }} *{{ $typeArg }}2, int overwrite);

/**
 * {{ $funcPrefix }}_contains checks to see if the given value is in the slice.
 */
bool
{{ $funcPrefix }}_contains(const {{ $typeName }} *{{ $typeArg }}, {{ .Name }} {{ $arg }});

/**
 * {{ $funcPrefix }}_delete removes the item at the given index and returns the
 * new length.
 */
int
{{ $funcPrefix }}_delete({{ $typeName }} *{{ $typeArg }}, const size_t idx);

/**
 * {{ $funcPrefix }}_replace replaces the value at the given index with the new
 * value.
 */
int
{{ $funcPrefix }}_replace({{ $typeName }} *{{ $typeArg }}, const size_t idx, const {{ .Name }} {{ $arg }});

/**
 * {{ $funcPrefix }}_foreach iterates through the slice and runs the user provided
 * function on each item. Additional user provided data can be provided 
 * using the user_data argument.
 */
int
{{ $funcPrefix }}_foreach({{ $typeName }} *{{ $typeArg }}, iter_func_t ift, void *user_data);

/**
 * {{ $funcPrefix }}_sort uses thet Quick Sort algorithm to sort the contents of the
 * slice if it is a standard type.
 */
int
{{ $funcPrefix }}_sort({{ $typeName }} *{{ $typeArg }}, size_t low, size_t high);

#endif /** end __{{ $headerName }}_H */
#ifdef __cplusplus
}
#endif
`

const sliceImplementationTmpl = `// This is generated code from safe_array_gen. Please do not edit unless
// sure of what you are doing.

{{- $name := Strip .Name "_t" }}
{{- $arg := "" }}
{{- $items := "" -}}
{{- $typeName := "" }}
{{- $typeArg := "" -}}
{{- $funcPrefix := "" }}
{{- if .Pointer -}}
{{- $arg = "*val" -}}
{{- $items = "**items" -}}
{{- else -}}
{{- $arg = "val" -}}
{{- $items = "*items" -}}
{{- end -}}

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
#include <stdlib.h>
#include <string.h>
#include <time.h>

#include "{{ .Name }}_slice.h"
{{- else }}

typedef struct {
    {{ .Name }} {{ $items }};
    size_t len;
    size_t cap;
} {{ $typeName }};
{{- end -}}
{{- $name := Strip .Name "_t" }}

{{ $typeName }}*
{{ $funcPrefix }}_new(const size_t cap)
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

{{ if .Pointer -}}
{{ .Name }}*
{{- else }}
{{ .Name }}
{{- end }}
{{ $funcPrefix }}_get({{ $typeName }} *{{ $typeArg }}, size_t idx)
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
	if ({{ $typeArg }}1->len != {{ $typeArg }}2->len) {
		return false;
	}

	if ({{ $typeArg }}1->compare != NULL) {
		for (size_t i = 0; i < s1->len; i++) {
			if (!{{ $typeArg }}1->compare(s1->items[i], {{ $typeArg }}2->items[i], user_data)) {
				return false;
			}
		}

		return true;
	}

	for (size_t i = 0; i < {{ $typeArg }}1->len; i++) {
{{- if .Pointer }}
    	if (*{{ $typeArg }}1->items[i] != *{{ $typeArg }}2->items[i]) {
{{- else }}
		if ({{ $typeArg }}1->items[i] != {{ $typeArg }}2->items[i]) {
{{- end }}
			return false;
		}
	}

	return true;
}

int
{{ $funcPrefix }}_copy(const {{ $typeName }} *{{ $typeArg }}1, {{ $typeName }} *{{ $typeArg }}2, int overwrite)
{
	if (overwrite) {
		if ({{ $typeArg }}1->len != {{ $typeArg }}2->len) {
			{{ $typeArg }}2->cap = {{ $typeArg }}1->cap;
			{{ $typeArg }}2->items = realloc({{ $typeArg }}2->items, sizeof({{ .Name }}) * {{ $typeArg }}1->cap);
		}
	}

	for (size_t i = 0; i < {{ $typeArg }}1->len; i++) {
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

	for (size_t i = 0; i < {{ $typeArg }}->len; i++) {
{{- if .Pointer }}
    	if (*{{ $typeArg }}->items[i] == val) {
{{- else }}
		if ({{ $typeArg }}->items[i] == val) {
{{- end }}
			return true;
		}
	}

	return false;
}

int
{{ $funcPrefix }}_delete({{ $typeName }} *{{ $typeArg }}, const size_t idx)
{
	if ({{ $typeArg }}->len == 0 || idx > {{ $typeArg }}->len) {
		return -1;
	}

	for (size_t i = idx; i < {{ $typeArg }}->len; i++) {
		{{ $typeArg }}->items[i] = {{ $typeArg }}->items[i + 1];
	}
	{{ $typeArg }}->len--;

	return {{ $typeArg }}->len;
}

int
{{ $funcPrefix }}_replace({{ $typeName }} *{{ $typeArg }}, const size_t idx, const {{ .Name }} {{ $arg }})
{
	if ({{ $typeArg }}->len == 0 || idx > {{ $typeArg }}->len) {
		return -1;
	}

	{{ $typeArg }}->items[idx] = {{ $arg }};

	return 0;
}

{{ if .Pointer -}}
{{ .Name }}*
{{- else }}
{{ .Name }}
{{- end }}
{{ $funcPrefix }}_first({{ $typeName }} *{{ $typeArg }})
{
	if ({{ $typeArg }}->len == 0) {
		return 0;
	}

	return {{ $typeArg }}->items[0];
}

{{ if .Pointer -}}
{{ .Name }}*
{{- else }}
{{ .Name }}
{{- end }}
{{ $funcPrefix }}_last({{ $typeName }} *{{ $typeArg }})
{
	if ({{ $typeArg }}->len == 0) {
		return 0;
	}

	return {{ $typeArg }}->items[{{ $typeArg }}->len-1]; 
}

int
{{ $funcPrefix }}_foreach({{ $typeName }} *{{ $typeArg }}, iter_func_t ift, void *user_data)
{
	if ({{ $typeArg }}->len == 0) {
		return 0;
	}
	
	for (size_t i = 0; i < {{ $typeArg }}->len; i++) {
		ift({{ $typeArg }}->items[i], user_data);
	}

	return 0;
}

static void
swap({{ .Name }} *x, {{ .Name }} *y)
{
    {{ .Name }} tmp = *x;
    *x = *y;
    *y = tmp;
}

static size_t
partition({{ $typeName }} *{{ $typeArg }}, size_t low, size_t high)
{
    size_t pi = low + (rand() % (high - low));

	if (pi != high) {
		swap(&{{ $typeArg }}->items[pi], &{{ $typeArg }}->items[high]);
	}
    
	{{ .Name }} pv = {{ $typeArg }}->items[high];
    	
    size_t i = low;

    for (size_t j = low; j < high; j++) {
        if ({{ $typeArg }}->items[j] <= pi) {
            swap(&{{ $typeArg }}->items[i], &{{ $typeArg }}->items[j]);
			i++;
        }
    }
    swap(&{{ $typeArg }}->items[i], &{{ $typeArg }}->items[high]);

    return i;
}

int
{{ $funcPrefix }}_sort({{ $typeName }} *{{ $typeArg }}, size_t low, size_t high)
{
	if ({{ $typeArg }}->len == 0) {
		return 0;
	}

	srand(time(NULL));

    if (low < high) {
        size_t pi = partition({{ $typeArg }}, low, high);

        {{ $funcPrefix }}_sort({{ $typeArg }}, low, pi-1);
        {{ $funcPrefix }}_sort({{ $typeArg }}, pi+1, high);
    }

	return 0;
}
`
