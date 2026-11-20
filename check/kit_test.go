// SPDX-License-Identifier: BSD-3-Clause

// Package check verifies the brand kit's whole format matrix, not a sample:
// every mark in svg/color must exist in every format at the size its directory
// names. A kit that was assembled by hand once lived three weeks with 3 of 7
// format families and nobody saw it, because every check looked at what was
// there.
package check

import (
	"bytes"
	"image"
	_ "image/jpeg"
	_ "image/png"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"testing"
)

var sizes = []int{16, 32, 48, 64, 128, 256, 512, 1024}

func marks(t *testing.T) []string {
	t.Helper()
	files, err := filepath.Glob("../svg/color/*.svg")
	if err != nil || len(files) == 0 {
		t.Fatalf("no marks in svg/color (%v): the checks below would pass over nothing", err)
	}
	ids := make([]string, 0, len(files))
	for _, f := range files {
		ids = append(ids, strings.TrimSuffix(filepath.Base(f), ".svg"))
	}
	return ids
}

func dims(t *testing.T, path string) (int, int) {
	t.Helper()
	f, err := os.Open(path)
	if err != nil {
		t.Fatal(err)
	}
	defer f.Close()
	c, _, err := image.DecodeConfig(f)
	if err != nil {
		t.Fatalf("%s: %v", path, err)
	}
	return c.Width, c.Height
}

func TestMatrix(t *testing.T) {
	for _, id := range marks(t) {
		for _, v := range []string{"color", "white", "black"} {
			b, err := os.ReadFile("../svg/" + v + "/" + id + ".svg")
			if err != nil || !bytes.Contains(b, []byte("<svg")) {
				t.Errorf("svg/%s/%s.svg missing or not an SVG (%v)", v, id, err)
			}
			if v != "color" && bytes.Contains(b, []byte("linearGradient")) {
				t.Errorf("svg/%s/%s.svg still carries the colour gradient", v, id)
			}
			for _, s := range sizes {
				w, h := dims(t, "../png/"+v+"/"+strconv.Itoa(s)+"/"+id+".png")
				if w != s || h != s {
					t.Errorf("png/%s/%d/%s.png is %dx%d", v, s, id, w, h)
				}
			}
		}
		for _, s := range sizes {
			w, h := dims(t, "../jpg/"+strconv.Itoa(s)+"/"+id+".jpg")
			if w != s || h != s {
				t.Errorf("jpg/%d/%s.jpg is %dx%d", s, id, w, h)
			}
		}
		if w, h := dims(t, "../avatar/"+id+".png"); w != 512 || h != 512 {
			t.Errorf("avatar/%s.png is %dx%d, want 512x512", id, w, h)
		}
		if w, h := dims(t, "../social/"+id+".png"); w != 1280 || h != 640 {
			t.Errorf("social/%s.png is %dx%d, want 1280x640", id, w, h)
		}
		ico, err := os.ReadFile("../ico/" + id + ".ico")
		if err != nil || !bytes.HasPrefix(ico, []byte{0, 0, 1, 0}) {
			t.Errorf("ico/%s.ico missing or not an ICO (%v)", id, err)
		}
		icns, err := os.ReadFile("../icns/" + id + ".icns")
		if err != nil || !bytes.HasPrefix(icns, []byte("icns")) {
			t.Errorf("icns/%s.icns missing or not an ICNS (%v)", id, err)
		}
	}
}
