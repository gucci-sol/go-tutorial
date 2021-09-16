package main

import (
	"example.com/gowiki/util"
	"log"
	"net/http"
	"net/url"
	"testing"
)

func TestGetTitle(t *testing.T) {
	testCases := []struct {
		name     string
		path     string
		expected string
	}{
		{name: "edit", path: "/edit/TestPage", expected: "TestPage"},
		{name: "view", path: "/edit/TestPage", expected: "TestPage"},
		{name: "save", path: "/edit/TestPage", expected: "TestPage"},
		{name: "edit_empty", path: "/edit/", expected: ""},
		{name: "view_empty", path: "/view/", expected: ""},
		{name: "save_empty", path: "/save/", expected: ""},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			r := &http.Request{
				URL: &url.URL{Path: tc.path},
			}

			title, err := getTitle(r)
			log.Print(err)
			if title != tc.expected || err != nil {
				t.Fatal(err)
			}
		})
	}

	t.Run("invalid Page Title", func(t *testing.T) {
		r := &http.Request{
			URL: &url.URL{Path: "/saddve/dd"},
		}
		_, err := getTitle(r)
		expected := util.InvalidPageTitleError
		actual := err
		if expected != actual {
			t.Fatalf("aaaa")
		}
	})
}
