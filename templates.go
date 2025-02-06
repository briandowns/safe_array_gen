package main

const sliceHeaderTmpl = `/**
 * This is generated code from safe_array_gen. Please do not edit unless
 * sure of what you are doing.
 */

{{ if Contains .Name "bool" -}}
#include <stdbool.h>
{{- end -}}
#include <stdint.h>
#include <stdlib.h>
{{ $name := Strip .Name "_t" }}
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
{{- if .CustomName -}}
{{- $typeName = .CustomName }}
{{- $funcPrefix = $typeName }}
typedef struct {
    {{ .Name }} {{ $items }};
    uint64_t len;
    uint64_t cap;
} {{ $typeName }};
{{- else -}}
{{- $typeName = printf "%s_slice_t" $name }}
{{- $funcPrefix = printf "%s_slice" $name }}
typedef struct {
    {{ .Name }} {{ $items }};
    uint64_t len;
    uint64_t cap;
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
{{ $funcPrefix }}_free({{ $typeName }} *s);

/**
 * {{ $funcPrefix }}_get attempts to retrieve the value at the given index. If
 * the index is out of range, 0 is returned indicating an error.
 */
{{- if .Pointer }}
{{ .Name }}*
{{- else }}
{{ .Name }}
{{- end }}
{{ $funcPrefix }}_get({{ $typeName }} *s, uint64_t idx);

/**
 * {{ $funcPrefix }}_append attempts to append the data to the given array.
 */
void
{{ $funcPrefix }}_append({{ $typeName }} *s, const {{ .Name }} {{ $arg }});

/**
 * {{ $funcPrefix }}_reverse the contents of the array.
 */
void
{{ $funcPrefix }}_reverse({{ $typeName }} *s);

/**
 * {{ $funcPrefix }}_compare takes 2 slices, compares them element by element
 * and returns 0 if they are not the same and 1 if they are.
 */
int
{{ $funcPrefix }}_compare(const {{ $typeName }} *s1, const {{ $typeName }} *s2);

/**
 * {{ $funcPrefix }}_copy takes 2 slices. The first is copied into the second
 * with the number of elements copied being returned. The code assumes that 
 * slice 2 has been created large enough to hold the contents of slice 1. If
 * the overwrite option has been selected, the code will make sure there is 
 * enough space in slice 2 and overwrite its contents.
 */
int
{{ $funcPrefix }}_copy(const {{ $typeName }} *s1, {{ $typeName }} *s2, int overwrite);

/**
 * {{ $funcPrefix }}_contains checks to see if the given value is in the slice.
 */
int
{{ $funcPrefix }}_contains(const {{ $typeName }} *s, {{ .Name }} {{ $arg }});

/**
 * {{ $funcPrefix }}_delete removes the item at the given index.
 */
int
{{ $funcPrefix }}_delete({{ $typeName }} *s, const uint64_t idx);

/**
 * {{ $funcPrefix }}_replace replaces the value at the given element with the 
 * given new value.
 */
int
{{ $funcPrefix }}_replace({{ $typeName }} *s, const uint64_t idx, const {{ .Name }} {{ $arg }});
`

const sliceImplementationTmpl = `
/**
 * This is generated code from safe_array_gen. Please do not edit unless
 * sure of what you are doing.
 */

{{- $name := Strip .Name "_t" }}
{{- $arg := "" }}
{{- $items := "" -}}
{{- $typeName := "" }}
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
{{- else -}}
{{- $typeName = printf "%s_slice_t" $name }}
{{- $funcPrefix = printf "%s_slice" $name }}
{{- end }}

{{- if not .Append }}

#include <stdbool.h>
#include <stdint.h>
#include <stdlib.h>

#include "{{ .Name }}_slice.h"
{{- else }}

typedef struct {
    {{ .Name }} {{ $items }};
    uint64_t len;
    uint64_t cap;
} {{ $name }}_slice_t;
{{- end -}}
{{- $name := Strip .Name "_t" }}

{{ $typeName }}*
{{ $funcPrefix }}_new(const uint64_t cap)
{
    {{ $name }}_slice_t *s = calloc(1, sizeof({{ $name }}_slice_t));
    s->items = calloc(1, sizeof({{ .Name }}) * cap);
    s->len = 0;
    s->cap = cap;

    return s;
}

void
{{ $funcPrefix }}_free({{ $typeName }} *s) {
    free(s->items);
    free(s);
}

{{ if .Pointer -}}
{{ .Name }}*
{{- else }}
{{ .Name }}
{{- end }}
{{ $funcPrefix }}_get({{ $typeName }} *s, uint64_t idx)
{
    if (idx >= 0 && idx < s->len) {
        return s->items[idx];
    }

    return 0;
}

void
{{ $funcPrefix }}_append({{ $typeName }} *s, const {{ .Name }} {{ $arg }})
{
    if (s->len == s->cap) {
        s->cap *= 2;
        s->items = realloc(s->items, sizeof({{ .Name }}) * s->cap);
    }
    s->items[s->len++] = val;
}

void
{{ $name }}_slice_reverse({{ $name }}_slice_t *s) {
	uint64_t i = s->len - 1;
    uint64_t j = 0;

    while(i > j) {
        {{ .Name }} temp = s->items[i];
        s->items[i] = s->items[j];
        s->items[j] = temp;
        i--;
        j++;
    }
}

int
{{ $name }}_slice_compare(const {{ $name }}_slice_t *s1, const {{ $name }}_slice_t *s2)
{
	if (s1->len != s2->len) {
		return 0;
	}

	for (uint64_t i = 0; i < s1->len; i++) {
{{- if .Pointer }}
    	if (*s1->items[i] != *s2->items[i]) {
{{- else }}
		if (s1->items[i] != s2->items[i]) {
{{- end }}
			return 0;
		}
	}

	return 1;
}

int
{{ $name }}_slice_copy(const {{ $name }}_slice_t *s1, {{ $name }}_slice_t *s2, int overwrite)
{
	if (overwrite) {
		if (s1->len != s2->len) {
			s2->cap = s1->cap;
			s2->items = realloc(s2->items, sizeof({{ .Name }}) * s1->cap);
		}
	}

	for (uint64_t i = 0; i < s1->len; i++) {
		s2->items[i] = s1->items[i];
		s2->len++;
	}

	return s2->len;
}

int
{{ $name }}_slice_contains(const {{ $name }}_slice_t *s, {{ .Name }} {{ $arg }})
{
	if (s->len == 0) {
		return 0;
	}

	for (uint64_t i = 0; i < s->len; i++) {
{{- if .Pointer }}
    	if (*s->items[i] == val) {
{{- else }}
		if (s->items[i] == val) {
{{- end }}
			return 1;
		}
	}

	return 0;
}

int
{{ $name }}_slice_delete({{ $name }}_slice_t *s, const uint64_t idx)
{
	if (s->len == 0) {
		return 1;
	}

	for (uint64_t i = idx; i < s->len-1; i++) {
		s->items[i] = s->items[i + 1];
	}
	s->len-1;

	return 0;
}

int
{{ $name }}_slice_replace({{ $name }}_slice_t *s, const uint64_t idx, const {{ .Name }} {{ $arg }})
{
	if (s->len == 0 || idx > s->len) {
		return 1;
	}

	s->items[idx] = {{ $arg }}

	return 0;
}
`
