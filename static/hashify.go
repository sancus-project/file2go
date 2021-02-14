package static

import (
	"fmt"
	"path"
	"strings"
)

func hashify(fname0 string, fname1 string, v *Content, m map[string]string, c map[string]*Content) bool {
	m[fname0] = fname1
	c[fname1] = v
	return true
}

func hashifyAppend1(fname0 string, suffix string, v *Content, m map[string]string, c map[string]*Content) bool {
	fname1 := strings.TrimSuffix(fname0, suffix)
	if fname0 == fname1 {
		return false
	}

	fname1 = fmt.Sprintf("%s-%s%s", fname1, v.Sha1sum[:12], suffix)

	return hashify(fname0, fname1, v, m, c)
}

func hashifyAppend2(fname0 string, ref0 string, ref1 string, v *Content, m map[string]string, c map[string]*Content) bool {
	fname1 := strings.TrimSuffix(fname0, ref0)
	if fname0 == fname1 || fname1[0] != '.' {
		return false
	}

	return hashify(fname0, ref1+fname1, v, m, c)
}

func hashifyAppend3(fname0 string, v *Content, m map[string]string, c map[string]*Content) bool {

	ext := path.Ext(fname0)
	fname1 := fname0[0 : len(fname0)-len(ext)-1]
	fname1 = fmt.Sprintf("%s-%s.%s", fname1, v.Sha1sum[:12], ext)

	return hashify(fname0, fname1, v, m, c)
}

func Hashify(files map[string]*Content) (map[string]string, map[string]*Content) {
	m := make(map[string]string, len(files))
	c := make(map[string]*Content, len(files))
	q := make([]string, len(files))

	// hashify *.css and *.js directly. store others in q1
	for k, v := range files {
		if hashifyAppend1(k, ".css", v, m, c) {
			continue
		} else if hashifyAppend1(k, ".js", v, m, c) {
			continue
		} else {
			q = append(q, k)
		}
	}

	for _, k := range q {
		// borrow hash from suffix if available
		v := files[k]
		ok := false
		for fn0, fn1 := range m {
			if hashifyAppend2(k, fn0, fn1, v, m, c) {
				ok = true
				break
			}
		}

		// and hashify the others directly
		if !ok {
			hashifyAppend3(k, v, m, c)
		}
	}

	return m, c
}
