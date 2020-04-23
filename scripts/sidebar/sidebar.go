//go:generate go run sidebar.go
package main

import (
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"regexp"
	"runtime"
	"strings"
	"text/template"
)

type Entry struct {
	Filename string
	Product  string
	Resource string
}

type Entries struct {
	Resources   []Entry
	DataSources []Entry
}

func main() {
	_, scriptPath, _, ok := runtime.Caller(0)
	if !ok {
		log.Fatal("Could not get current working directory")
	}
	tpgDir := scriptPath
	for !strings.HasPrefix(filepath.Base(tpgDir), "terraform-provider-") && tpgDir != "/" {
		tpgDir = filepath.Clean(tpgDir + "/..")
	}
	if tpgDir == "/" {
		log.Fatal("Script was run outside of google provider directory")
	}

	resourcesByProduct, err := entriesByProduct(tpgDir + "/website/docs/r")
	if err != nil {
		panic(err)
	}
	dataSourcesByProduct, err := entriesByProduct(tpgDir + "/website/docs/d")
	if err != nil {
		panic(err)
	}
	allEntriesByProduct := make(map[string]Entries)
	for p, e := range resourcesByProduct {
		v := allEntriesByProduct[p]
		v.Resources = e
		allEntriesByProduct[p] = v
	}
	for p, e := range dataSourcesByProduct {
		v := allEntriesByProduct[p]
		v.DataSources = e
		allEntriesByProduct[p] = v
	}

	tmpl, err := template.ParseFiles(tpgDir + "/website/google.erb.tmpl")
	if err != nil {
		panic(err)
	}
	f, err := os.Create(tpgDir + "/website/google.erb")
	if err != nil {
		panic(err)
	}
	defer f.Close()
	err = tmpl.Execute(f, allEntriesByProduct)
	if err != nil {
		panic(err)
	}
}

func entriesByProduct(dir string) (map[string][]Entry, error) {
	d, err := ioutil.ReadDir(dir)
	if err != nil {
		return nil, err
	}

	entriesByProduct := make(map[string][]Entry)
	for _, f := range d {
		entry, err := getEntry(dir, f.Name())
		if err != nil {
			return nil, err
		}
		entriesByProduct[entry.Product] = append(entriesByProduct[entry.Product], entry)
	}

	return entriesByProduct, nil
}

func getEntry(dir, filename string) (Entry, error) {
	file, err := ioutil.ReadFile(dir + "/" + filename)
	if err != nil {
		return Entry{}, err
	}

	return Entry{
		Filename: strings.TrimSuffix(filename, ".markdown"),
		Product:  findRegex(file, `subcategory: "(.*)"`),
		Resource: findRegex(file, `page_title: "Google: (.*)"`),
	}, nil
}

func findRegex(contents []byte, regex string) string {
	r := regexp.MustCompile(regex)
	sm := r.FindStringSubmatch(string(contents))
	if len(sm) > 1 {
		return sm[1]
	}
	return ""
}
