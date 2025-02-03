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
)

const usage = `version: %s

Usage: %[2]s [-t types]

Options:
    -h            help
    -v            show version and exit
    -t            types, comma seperated (int8,int16,...)
	-p            use a pointer for the given type(s)
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

	tmp1 := template.New("slice_gen").Funcs(funcMap)
	tmp1, err := tmp1.Parse(sliceTmpl)
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

	os.Exit(0)
}

const sliceTmpl = `
{{- if Contains .Name "bool" }}
#include <stdbool.h>
{{- end }}
#include <stdint.h>
#include <stdlib.h>
{{ $fullName := .Name }}
{{- $name := Strip .Name "_t" }}
typedef struct {
    {{ .Name }} *items;
    uint64_t len;
    uint64_t cap;
} {{ $name }}_slice_t;

/**
 * {{ $name }}_slice_new creates a pointer of type {{ $name }}_slice_t, sets
 * default values, and returns the pointer to the allocated memory. The user
 * is responsible for freeing this memory.
 */
{{ $name }}_slice_t*
{{ $name }}_slice_new(const int cap)
{
    {{ $name }}_slice_t *s = calloc(1, sizeof({{ $name }}_slice_t));
    s->items = calloc(1, sizeof({{ .Name }}) * cap);
    s->len = 0;
    s->cap = cap;

    return s;
}

/**
 * {{ $name }}_slice_free frees the memory used by the given pointer. 
 */
void
{{ $name }}_slice_free({{ $name }}_slice_t *s) {
    free(s->items);
    free(s);
}

/**
 * {{ $name }}_slice_get attempts to retrieve the value at the given index. If
 * the index is out of range, 0 is returned indicating an error.
 */
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

/**
 * {{ $name }}_slice_append attempts to append the data to the given array.
 */
void
{{- if .Pointer }}
{{ $name }}_slice_append({{ $name }}_slice_t *s, const {{ .Name }} *val)
{{- else }}
{{ $name }}_slice_append({{ $name }}_slice_t *s, const {{ .Name }} val)
{{- end }}
{
    if (s->len == s->cap) {
        s->cap *= 2;
        s->items = realloc(s->items, sizeof({{ .Name }}) * s->cap);
    }
    s->items[s->len++] = val;
}
`
