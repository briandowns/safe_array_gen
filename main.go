/*-
 * SPDX-License-Identifier: BSD-2-Clause
 *
 * Copyright (c) 2025 Brian J. Downs
 *
 * Redistribution and use in source and binary forms, with or without
 * modification, are permitted provided that the following conditions
 * are met:
 * 1. Redistributions of source code must retain the above copyright
 *    notice, this list of conditions and the following disclaimer.
 * 2. Redistributions in binary form must reproduce the above copyright
 *    notice, this list of conditions and the following disclaimer in the
 *    documentation and/or other materials provided with the distribution.
 *
 * THIS SOFTWARE IS PROVIDED BY THE REGENTS AND CONTRIBUTORS ``AS IS'' AND
 * ANY EXPRESS OR IMPLIED WARRANTIES, INCLUDING, BUT NOT LIMITED TO, THE
 * IMPLIED WARRANTIES OF MERCHANTABILITY AND FITNESS FOR A PARTICULAR PURPOSE
 * ARE DISCLAIMED.  IN NO EVENT SHALL THE REGENTS OR CONTRIBUTORS BE LIABLE
 * FOR ANY DIRECT, INDIRECT, INCIDENTAL, SPECIAL, EXEMPLARY, OR CONSEQUENTIAL
 * DAMAGES (INCLUDING, BUT NOT LIMITED TO, PROCUREMENT OF SUBSTITUTE GOODS
 * OR SERVICES; LOSS OF USE, DATA, OR PROFITS; OR BUSINESS INTERRUPTION)
 * HOWEVER CAUSED AND ON ANY THEORY OF LIABILITY, WHETHER IN CONTRACT, STRICT
 * LIABILITY, OR TORT (INCLUDING NEGLIGENCE OR OTHERWISE) ARISING IN ANY WAY
 * OUT OF THE USE OF THIS SOFTWARE, EVEN IF ADVISED OF THE POSSIBILITY OF
 * SUCH DAMAGE.
 */

package main

import (
	"flag"
	"fmt"
	"os"
	"strings"
	"text/template"
)

var (
	name    string
	version string
	gitSHA  string
)

var (
	vers        bool
	typesFlag   string
	pointerFlag bool
	appendFlag  string
)

const usage = `version: %s

Usage: %[2]s [-t types]

Options:
    -h            help
    -v            show version and exit
    -t            types, comma seperated (int8,int16,...)
    -p            use a pointer for the given type(s)
    -a            only generate implementation code and append 
                  to the given file
`

// data used to be passed to the template engine.
type data struct {
	Name    string
	Pointer bool
}

func main() {
	flag.Usage = func() {
		w := os.Stderr
		for _, arg := range os.Args {
			if arg == "-h" {
				w = os.Stdout
				break
			}
		}
		fmt.Fprintf(w, usage, version, name)
	}

	flag.BoolVar(&vers, "v", false, "")
	flag.StringVar(&typesFlag, "t", "", "")
	flag.BoolVar(&pointerFlag, "p", false, "")
	flag.StringVar(&appendFlag, "a", "", "")
	flag.Parse()

	if vers {
		fmt.Fprintf(os.Stdout, "version: %s - git sha: %s\n", version, gitSHA)
		return
	}

	types := strings.Split(typesFlag, ",")
	if len(types) < 1 {
		fmt.Fprint(os.Stderr, "")
		os.Exit(1)
	}

	funcMap := template.FuncMap{
		"ToUpper": func(s string) string {
			return strings.ToUpper(s)
		},
		"Contains": func(s, ss string) bool {
			return strings.Contains(s, ss)
		},
		"Strip": func(s, ss string) string {
			return strings.Replace(s, ss, "", 1)
		},
	}

	implementationTmpl := template.New("slice_implementation_gen").Funcs(funcMap)
	tmp1, err := implementationTmpl.Parse(sliceImplementationTmpl)
	if err != nil {
		fmt.Fprint(os.Stderr, err)
		os.Exit(1)
	}

	for _, t := range types {
		f, err := os.Create(t + "_slice.c")
		if err != nil {
			fmt.Fprint(os.Stderr, err)
			os.Exit(1)
		}
		defer f.Close()

		d := data{
			Name:    t,
			Pointer: pointerFlag,
		}

		if err := tmp1.Execute(f, &d); err != nil {
			fmt.Fprint(os.Stderr, err)
			os.Exit(1)
		}
	}

	headerTmpl := template.New("slice_header_gen").Funcs(funcMap)
	tmp2, err := headerTmpl.Parse(sliceHeaderTmpl)
	if err != nil {
		fmt.Fprint(os.Stderr, err)
		os.Exit(1)
	}

	for _, t := range types {
		f, err := os.Create(t + "_slice.h")
		if err != nil {
			fmt.Fprint(os.Stderr, err)
			os.Exit(1)
		}
		defer f.Close()

		d := data{
			Name:    t,
			Pointer: pointerFlag,
		}

		if err := tmp2.Execute(f, &d); err != nil {
			fmt.Fprint(os.Stderr, err)
			os.Exit(1)
		}
	}

	os.Exit(0)
}

const sliceHeaderTmpl = `
/**
 * This is generated code from safe_array_gen. Please do not edit unless
 * sure of what you are doing.
 */

{{ if Contains .Name "bool" -}}
#include <stdbool.h>
{{- end -}}
#include <stdint.h>
#include <stdlib.h>
{{ $fullName := .Name }}
{{- $name := Strip .Name "_t" }}
{{- $arg := "" -}}
{{- $items := "" -}}
{{- if .Pointer -}}
{{- $arg = "*val" -}}
{{- $items = "**items" -}}
{{- else -}}
{{- $arg = "val" -}}
{{- $items = "*items" -}}
{{- end }}
typedef struct {
    {{ .Name }} {{ $items }};
    uint64_t len;
    uint64_t cap;
} {{ $name }}_slice_t;

/**
 * {{ $name }}_slice_new creates a pointer of type {{ $name }}_slice_t, sets
 * default values, and returns the pointer to the allocated memory. The user
 * is responsible for freeing this memory.
 */
{{ $name }}_slice_t*
{{ $name }}_slice_new(const uint64_t cap);

/**
 * {{ $name }}_slice_free frees the memory used by the given pointer. 
 */
void
{{ $name }}_slice_free({{ $name }}_slice_t *s);

/**
 * {{ $name }}_slice_get attempts to retrieve the value at the given index. If
 * the index is out of range, 0 is returned indicating an error.
 */
{{- if .Pointer }}
{{ $fullName }}*
{{- else }}
{{ $fullName }}
{{- end }}
{{ $name }}_slice_get({{ $name }}_slice_t *s, uint64_t idx);

/**
 * {{ $name }}_slice_append attempts to append the data to the given array.
 */
void
{{ $name }}_slice_append({{ $name }}_slice_t *s, const {{ .Name }} {{ $arg }});

/**
 * {{ $name }}_slice_reverse the contents of the array.
 */
void
{{ $name }}_slice_reverse({{ $name }}_slice_t *s);

/**
 * {{ $name }}_slice_compare takes 2 slices, compares them element by element
 * and returns 0 if they are not the same and 1 if they are.
 */
int
{{ $name }}_slice_compare(const {{ $name }}_slice_t *s1, const {{ $name }}_slice_t *s2);

/**
 * {{ $name }}_slice_copy takes 2 slices. The first is copied into the second
 * with the number of elements copied being returned. The code assumes that 
 * slice 2 has been created large enough to hold the contents of slice 1. If
 * the overwrite option has been selected, the code will make sure there is 
 * enough space in slice 2 and overwrite its contents.
 */
int
{{ $name }}_slice_copy(const {{ $name }}_slice_t *s1, {{ $name }}_slice_t *s2, int overwrite);

/**
 * {{ $name }}_slice_contains checks to see if the given value is in the slice.
 */
int
{{ $name }}_slice_contains(const {{ $name }}_slice_t *s, {{ .Name }} {{ $arg }});

/**
 * {{ $name }}_slice_delete removes the item at the given index.
 */
int
{{ $name }}_slice_delete({{ $name }}_slice_t *s, const uint64_t idx);

/**
 * {{ $name }}_slice_replace replaces the value at the given element with the 
 * given new value.
 */
int
{{ $name }}_slice_replace({{ $name }}_slice_t *s, const uint64_t idx, const {{ .Name }} {{ $arg }});
`

const sliceImplementationTmpl = `
/**
 * This is generated code from safe_array_gen. Please do not edit unless
 * sure of what you are doing.
 */

{{ if Contains .Name "bool" -}}
#include <stdbool.h>
{{- end -}}
#include <stdint.h>
#include <stdlib.h>

#include "{{ .Name }}_slice.h"
{{ $fullName := .Name }}
{{- $name := Strip .Name "_t" }}

{{ $name }}_slice_t*
{{ $name }}_slice_new(const uint64_t cap)
{
    {{ $name }}_slice_t *s = calloc(1, sizeof({{ $name }}_slice_t));
    s->items = calloc(1, sizeof({{ .Name }}) * cap);
    s->len = 0;
    s->cap = cap;

    return s;
}

void
{{ $name }}_slice_free({{ $name }}_slice_t *s) {
    free(s->items);
    free(s);
}


{{- if .Pointer }}
{{ $fullName }}*
{{- else }}
{{ $fullName }}
{{- end }}
{{ $name }}_slice_get({{ $name }}_slice_t *s, uint64_t idx)
{
    if (idx >= 0 && idx < s->len) {
        return s->items[idx];
    }

    return 0;
}

{{- $arg := "" }}
{{- if .Pointer }}
{{- $arg = "*val" }}
{{- else }}
{{- $arg = "val" }}
{{- end }}

void
{{ $name }}_slice_append({{ $name }}_slice_t *s, const {{ .Name }} {{ $arg }})
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
	if (s->len == 0) {
		return 1;
	}

	s->items[idx] = {{ $arg }}

	return 0;
}
`
