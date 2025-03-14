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
	vers       bool
	typesFlag  string
	appendFlag string
	nameFlag   string
)

const usage = `version: %s

Usage: %[2]s [-t types]

Options:
    -h            help
    -v            show version and exit
    -t            types, comma separated (int8,int16,...)
    -a            only generate implementation code and append 
                  to the given file
    -n            custom name for the given type
`

// data used to be passed to the template engine.
type data struct {
	Name       string
	CustomName string
	Pointer    bool
	Append     bool
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
	flag.StringVar(&appendFlag, "a", "", "")
	flag.StringVar(&nameFlag, "n", "", "")
	flag.Parse()

	if vers {
		fmt.Fprintf(os.Stdout, "version: %s - git sha: %s\n", version, gitSHA)
		return
	}

	types := strings.Split(typesFlag, ",")
	typeCount := len(types)

	if typeCount == 1 && types[0] == "" {
		fmt.Fprintln(os.Stderr, "error: type(s) need to be provided")
		os.Exit(1)
	}

	if nameFlag != "" && typeCount > 1 {
		fmt.Fprintln(os.Stderr, "error: only 1 type supported when using custom name")
		os.Exit(1)
	}

	funcMap := template.FuncMap{
		"ToUpper": func(s string) string {
			return strings.ToUpper(s)
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

	if appendFlag != "" {
		f, err := os.OpenFile(appendFlag, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
		if err != nil {
			fmt.Fprint(os.Stderr, err)
			os.Exit(1)
		}
		defer f.Close()

		appendTmpl := template.New("append_gen").Funcs(funcMap)
		tmp1, err := appendTmpl.Parse(sliceImplementationTmpl)
		if err != nil {
			fmt.Fprint(os.Stderr, err)
			os.Exit(1)
		}

		d := data{
			Append: true,
		}

		if nameFlag != "" {
			d.CustomName = nameFlag
		}

		for _, t := range types {
			d.Name = t

			if err := tmp1.Execute(f, &d); err != nil {
				fmt.Fprint(os.Stderr, err)
				os.Exit(1)
			}
		}

		os.Exit(0)
	}

	d := data{
		Append: false,
	}

	if nameFlag != "" {
		d.CustomName = nameFlag
	}

	for _, t := range types {
		f, err := os.Create(t + "_slice.c")
		if err != nil {
			fmt.Fprint(os.Stderr, err)
			os.Exit(1)
		}
		defer f.Close()

		d.Name = t

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

		d.Name = t

		if err := tmp2.Execute(f, &d); err != nil {
			fmt.Fprint(os.Stderr, err)
			os.Exit(1)
		}
	}

	os.Exit(0)
}
