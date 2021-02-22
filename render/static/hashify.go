package static

import (
	"fmt"
	"log"
	"path"
	"strings"
	"sort"
)

func hashify(fname0 string, fname1 string, v *StaticRendererFile, m map[string]string) bool {
	m[fname0] = fname1
	v.Hashified = fname1
	return true
}

func hashifyAppend1(fname0 string, suffix string, v *StaticRendererFile, m map[string]string) bool {
	fname1 := strings.TrimSuffix(fname0, suffix)
	if fname0 == fname1 {
		return false
	}

	if len(v.Sha1sum) > 0 {
		fname1 = fmt.Sprintf("%s-%s%s", fname1, v.Sha1sum[:12], suffix)
	} else {
		fname1 = fname0
	}

	return hashify(fname0, fname1, v, m)
}

func hashifyAppend2(fname0 string, ref0 string, ref1 string, v *StaticRendererFile, m map[string]string) bool {
	fname1 := strings.TrimPrefix(fname0, ref0)
	if fname0 == fname1 || fname1[0] != '.' {
		return false
	}

	return hashify(fname0, ref1+fname1, v, m)
}

func hashifyAppend3(fname0 string, v *StaticRendererFile, m map[string]string) bool {
	var fname1 string

	if len(v.Sha1sum) > 0 {
		ext := path.Ext(fname0)
		fname1 = fname0[0 : len(fname0)-len(ext)]
		fname1 = fmt.Sprintf("%s-%s%s", fname1, v.Sha1sum[:12], ext)
	} else {
		fname1 = fname0
	}

	return hashify(fname0, fname1, v, m)
}

// Hashify
func (r *StaticRenderer) Hashify() (err error) {

	m := make(map[string]string, len(r.Names))
	q := make([]string, 0, len(r.Names))

	// hashify *.css and *.js directly. store others in q1
	for k, v := range r.Files {
		if hashifyAppend1(k, ".css", v, m) {
			continue
		} else if hashifyAppend1(k, ".js", v, m) {
			continue
		} else {
			q = append(q, k)
		}
	}

	for _, k := range q {
		// borrow hash from suffix if available
		v := r.Files[k]
		ok := false
		for fn0, fn1 := range m {
			if hashifyAppend2(k, fn0, fn1, v, m) {
				ok = true
				break
			}
		}

		// and hashify the others directly
		if !ok {
			hashifyAppend3(k, v, m)
		}
	}

	// Sort names
	sort.Strings(r.Names)

	// and update redirects
	for k := range r.Redirect {
		delete(r.Redirect, k)
	}

	for _, fn0 := range r.Names {
		v := r.Files[fn0]
		fn1 := v.Hashified

		if fn1 == fn0 || fn1 == "" {
			log.Printf("Hashify: %q", fn0)
		} else {
			log.Printf("Hashify: %q -> %q", fn0, fn1)
			r.Redirect[fn0] = path.Base(fn1)
		}
	}

	return
}
