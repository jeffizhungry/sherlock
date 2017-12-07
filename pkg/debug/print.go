package debug

import (
	"bytes"
	"fmt"
	"io"
	"reflect"
	"runtime"
	"sort"
	"strings"
)

// Printf is a wrapper around fmt.Printf to add TRACE tags and file & line info.
// This is very useful for development since we can just search logs for the
// [TRACE] tag, and the file line information make.
//
// WARNING: MAKE SURE YOU REMOVE ALL CALLS TO THIS FUNCTION IN BEFORE MERGING INTO MASTER.
func Printf(format string, args ...interface{}) {
	if len(args) == 0 {
		return
	}
	file, line := fileInfo(2)
	hdr := fmt.Sprintf("[TRACE] %10v:%-4v - ", file, line)
	_, isErr := args[0].(error)
	if isErr {
		hdr += "ERROR: "
	}
	fmt.Println(hdr + fmt.Sprintf(format, args...))
}

// Println - SEE Printf
func Println(args ...interface{}) {
	if len(args) == 0 {
		return
	}
	file, line := fileInfo(2)
	hdr := fmt.Sprintf("[TRACE] %10v:%-4v - ", file, line)
	_, isErr := args[0].(error)
	if isErr {
		hdr += "ERROR: "
	}
	output := make([]interface{}, 0, len(args)+1)
	output = append(output, hdr)
	output = append(output, args...)
	fmt.Println(output...)
}

// PPrintln pretty print all itmes.
func PPrintln(args ...interface{}) {
	if len(args) == 0 {
		return
	}
	file, line := fileInfo(2)
	hdr := fmt.Sprintf("[TRACE] %10v:%-4v - ", file, line)
	_, isErr := args[0].(error)
	if isErr {
		hdr += "ERROR: "
	}
	output := make([]interface{}, 0, len(args)+1)
	output = append(output, hdr)
	for _, arg := range args {
		output = append(output, Prettify(arg))
	}
	fmt.Println(output...)
}

// Prettify returns the string representation of a value.
func Prettify(i interface{}) string {
	var buf bytes.Buffer
	prettify(reflect.ValueOf(i), 0, &buf)
	return buf.String()
}

// fileInfo grabs the file name and line number of the caller.
func fileInfo(callDepth int) (file string, line int) {

	// Inspect runtime call stack
	pc := make([]uintptr, callDepth)
	runtime.Callers(callDepth, pc)
	f := runtime.FuncForPC(pc[callDepth-1])

	file, line = f.FileLine(pc[callDepth-1])

	// Truncate abs file path
	if slash := strings.LastIndex(file, "/"); slash >= 0 {
		file = file[slash+1:]
	}
	return
}

// prettify will recursively walk value v to build a textual
// representation of the value.
func prettify(v reflect.Value, indent int, buf *bytes.Buffer) {
	for v.Kind() == reflect.Ptr {
		v = v.Elem()
	}

	switch v.Kind() {
	case reflect.Struct:
		strtype := v.Type().String()
		if strtype == "time.Time" {
			fmt.Fprintf(buf, "%s", v.Interface())
			break
		} else if strings.HasPrefix(strtype, "io.") {
			buf.WriteString("<buffer>")
			break
		}

		buf.WriteString("{\n")

		names := []string{}
		for i := 0; i < v.Type().NumField(); i++ {
			name := v.Type().Field(i).Name
			f := v.Field(i)
			if name[0:1] == strings.ToLower(name[0:1]) {
				continue // ignore unexported fields
			}
			if (f.Kind() == reflect.Ptr || f.Kind() == reflect.Slice || f.Kind() == reflect.Map) && f.IsNil() {
				continue // ignore unset fields
			}
			names = append(names, name)
		}

		for i, n := range names {
			val := v.FieldByName(n)
			buf.WriteString(strings.Repeat(" ", indent+2))
			buf.WriteString(n + ": ")
			prettify(val, indent+2, buf)

			if i < len(names)-1 {
				buf.WriteString(",\n")
			}
		}

		buf.WriteString("\n" + strings.Repeat(" ", indent) + "}")
	case reflect.Slice:
		nl, id, id2 := "", "", ""
		if v.Len() > 3 {
			nl, id, id2 = "\n", strings.Repeat(" ", indent), strings.Repeat(" ", indent+2)
		}
		buf.WriteString("[" + nl)
		for i := 0; i < v.Len(); i++ {
			buf.WriteString(id2)
			prettify(v.Index(i), indent+2, buf)

			if i < v.Len()-1 {
				buf.WriteString("," + nl)
			}
		}

		buf.WriteString(nl + id + "]")
	case reflect.Map:
		buf.WriteString("{\n")

		rvk := v.MapKeys()

		// NOTE(Jeff): Sort keys for more consistent output
		if len(rvk) > 0 {
			switch rvk[0].Kind() {
			case reflect.String:
				sort.Slice(rvk, func(i int, j int) bool {
					return rvk[i].String() < rvk[j].String()
				})
			default:
				// Handle other types if necessary
			}
		}

		for i, k := range rvk {
			buf.WriteString(strings.Repeat(" ", indent+2))

			// NOTE(Jeff): Quote string keys
			switch rvk[i].Kind() {
			case reflect.String:
				buf.WriteString("\"" + k.String() + "\"" + ": ")
			default:
				buf.WriteString(k.String() + ": ")
			}

			prettify(v.MapIndex(k), indent+2, buf)
			if i < v.Len()-1 {
				buf.WriteString(",\n")
			}
		}

		buf.WriteString(",\n" + strings.Repeat(" ", indent) + "}")

	default:
		if !v.IsValid() {
			fmt.Fprint(buf, "<invalid value>")
			return
		}
		format := "%v"
		switch v.Interface().(type) {
		case string:
			format = "%q"
		case io.ReadSeeker, io.Reader:
			format = "buffer(%p)"
		}
		fmt.Fprintf(buf, format, v.Interface())
	}
}
