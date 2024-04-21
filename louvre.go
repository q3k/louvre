// Copyright 2018 Serge 'q3k' Bazanski
// Permission to use, copy, modify, and/or distribute this software for any
// purpose with or without fee is hereby granted.

package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
)

var (
	upstream = "http://oodviewer.q3k.me/"
	termsUrl = upstream + "terms.json"
	termUrl  = upstream + "term.json/%s"

	flagOutput   string
	flagParallel int
)

func terms() ([]string, error) {
	resp, err := http.Get(termsUrl)
	if err != nil {
		return nil, fmt.Errorf("could not GET %q: %v", termsUrl, err)
	}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("could not read %q: %v", termsUrl, err)
	}

	tdata := [][]interface{}{}
	err = json.Unmarshal(body, &tdata)
	if err != nil {
		return nil, fmt.Errorf("could not decode: %v", err)
	}

	res := make([]string, len(tdata))
	j := 0
	for i, t := range tdata {
		s, ok := t[0].(string)
		if !ok {
			log.Printf("Could not decode %dth element...", i)
			continue
		}
		res[j] = s
		j += 1
	}
	res = res[:j]
	return res, nil
}

type term struct {
	Entry  string `json:"entry"`
	Added  int    `json:"added"`
	Author string `json:"author"`
}

func oneTerm(termName string) ([]*term, error) {
	termName = url.QueryEscape(termName)
	termUrl := fmt.Sprintf(termUrl, termName)
	resp, err := http.Get(termUrl)
	if err != nil {
		return nil, fmt.Errorf("could not GET %q: %v", termUrl, err)
	}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("could not read %q: %v", termUrl, err)
	}

	t := []*term{}
	err = json.Unmarshal(body, &t)
	if err != nil {
		return nil, fmt.Errorf("could not decode: %v", err)
	}
	return t, nil
}

type termOrError struct {
	termName string
	err      error
	terms    []*term
}

func worker(jobs chan string, results chan termOrError) {
	for j := range jobs {
		terms, err := oneTerm(j)
		res := termOrError{
			termName: j,
			err:      err,
			terms:    terms,
		}
		results <- res
	}
}

func init() {
	flag.StringVar(&flagOutput, "output", "terms.json", "Where to download all terms.")
	flag.IntVar(&flagParallel, "parallel", 32, "How many concurrent connections to use.")
	flag.Parse()
}

func main() {
	trms, err := terms()
	if err != nil {
		log.Fatalf("Could not download terms: %v", err)
	}
	log.Printf("Server has %d terms.\n", len(trms))

	jobs := make(chan string, len(trms))
	for _, t := range trms {
		jobs <- t
	}
	results := make(chan termOrError, len(trms))

	allTerms := make(map[string][]*term)

	for w := 0; w < flagParallel; w++ {
		go worker(jobs, results)
	}

	for i := 0; i < len(trms); i++ {
		res := <-results
		if res.err != nil {
			log.Printf("Could not get %q: %v", res.termName, res.err)
			continue
		}
		allTerms[res.termName] = res.terms
	}
	log.Printf("Got %d terms.\n", len(allTerms))

	outData, err := json.Marshal(allTerms)
	if err != nil {
		log.Fatalf("Could not marshal JSON: %v\n", err)
	}

	err = ioutil.WriteFile(flagOutput, outData, 0644)
	if err != nil {
		log.Fatalf("Could not write to %q: %v\n", flagOutput, err)
	}
	log.Printf("Wrote %d terms to %q.\n", len(allTerms), flagOutput)
}
