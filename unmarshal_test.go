package goquery

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

const testPage = `<!DOCTYPE html>
<html>
  <head>
    <title></title>
    <meta charset="utf-8" />
  </head>
  <body>
    <h1>
      <ul id="resources">
        <li class="resource">
          <div class="name">Foo</div>
        </li>
        <li class="resource">
          <div class="name">Bar</div>
        </li>
        <li class="resource">
          <div class="name">Baz</div>
        </li>
        <li class="resource">
          <div class="name">Bang</div>
        </li>
        <li class="resource">
          <div class="name">Zip</div>
        </li>
      </ul>
    </h1>
  </body>
</html>
`

type Page struct {
	Resources []Resource `goquery:"#resources .resource"`
}

type Resource struct {
	Name string `goquery:".name"`
}

func TestDecoder(t *testing.T) {
	asrt := assert.New(t)

	var p Page

	asrt.NoError(NewDecoder(strings.NewReader(testPage)).Decode(&p))
	asrt.Len(p.Resources, 5)

	vals := []string{"Foo", "Bar", "Baz", "Bang", "Zip"}
	for i, val := range vals {
		asrt.Equal(val, p.Resources[i].Name)
	}
}
