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

{{ if Contains .Name "bool" -}}
#include <stdbool.h>
{{- end -}}
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
} {{ $typeName }};
{{- else -}}
{{- $typeName = printf "%s_slice_t" $name }}
{{- $funcPrefix = printf "%s_slice" $name }}
{{ $typeArg = "s" }}
typedef struct {
    {{ .Name }} {{ $items }};
    size_t len;
    size_t cap;
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
 * and returns 0 if they are not the same and 1 if they are.
 */
int
{{ $funcPrefix }}_compare(const {{ $typeName }} *{{ $typeArg }}1, const {{ $typeName }} *{{ $typeArg }}2);

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
int
{{ $funcPrefix }}_contains(const {{ $typeName }} *{{ $typeArg }}, {{ .Name }} {{ $arg }});

/**
 * {{ $funcPrefix }}_delete removes the item at the given index.
 */
int
{{ $funcPrefix }}_delete({{ $typeName }} *{{ $typeArg }}, const size_t idx);

/**
 * {{ $funcPrefix }}_replace replaces the value at the given element with the 
 * given new value.
 */
int
{{ $funcPrefix }}_replace({{ $typeName }} *{{ $typeArg }}, const size_t idx, const {{ .Name }} {{ $arg }});

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

int
{{ $funcPrefix }}_compare(const {{ $typeName }} *{{ $typeArg }}1, const {{ $typeName }} *{{ $typeArg }}2)
{
	if ({{ $typeArg }}1->len != {{ $typeArg }}2->len) {
		return 0;
	}

	for (size_t i = 0; i < {{ $typeArg }}1->len; i++) {
{{- if .Pointer }}
    	if (*{{ $typeArg }}1->items[i] != *{{ $typeArg }}2->items[i]) {
{{- else }}
		if ({{ $typeArg }}1->items[i] != {{ $typeArg }}2->items[i]) {
{{- end }}
			return 0;
		}
	}

	return 1;
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

int
{{ $funcPrefix }}_contains(const {{ $typeName }} *{{ $typeArg }}, {{ .Name }} {{ $arg }})
{
	if ({{ $typeArg }}->len == 0) {
		return 0;
	}

	for (size_t i = 0; i < {{ $typeArg }}->len; i++) {
{{- if .Pointer }}
    	if (*{{ $typeArg }}->items[i] == val) {
{{- else }}
		if ({{ $typeArg }}->items[i] == val) {
{{- end }}
			return 1;
		}
	}

	return 0;
}

int
{{ $funcPrefix }}_delete({{ $typeName }} *{{ $typeArg }}, const size_t idx)
{
	if ({{ $typeArg }}->len == 0) {
		return 1;
	}

	for (size_t i = idx; i < {{ $typeArg }}->len-1; i++) {
		{{ $typeArg }}->items[i] = {{ $typeArg }}->items[i + 1];
	}
	{{ $typeArg }}->len-1;

	return 0;
}

int
{{ $funcPrefix }}_replace({{ $typeName }} *{{ $typeArg }}, const size_t idx, const {{ .Name }} {{ $arg }})
{
	if ({{ $typeArg }}->len == 0 || idx > {{ $typeArg }}->len) {
		return 1;
	}

	{{ $typeArg }}->items[idx] = {{ $arg }};

	return 0;
}
`
