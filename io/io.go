package io

import (
  "io/ioutil"
)

type Page struct {
  Title string
  Body  []byte
}

func (p *Page) Save(path string) error {
  title := p.Title + ".txt"
  return ioutil.WriteFile(path + "/" + title, p.Body, 0600)
}

func LoadPage(path, title string) (*Page, error) {
  filename := title + ".txt"

  body, err := ioutil.ReadFile(path + "/" + filename)

  if err != nil {
    return nil, err
  }

  return &Page{Title: title, Body: body}, nil
}
