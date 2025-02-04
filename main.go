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
	Append  bool
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

		for _, t := range types {
			d := data{
				Name:    t,
				Pointer: pointerFlag,
				Append:  true,
			}

			if err := tmp1.Execute(f, &d); err != nil {
				fmt.Fprint(os.Stderr, err)
				os.Exit(1)
			}
		}

		os.Exit(0)
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
			Append:  false,
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
