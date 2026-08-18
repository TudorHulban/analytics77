package layouts

import (
	"testing"

	"github.com/TudorHulban/hxgo/dsl"
)

func TestBase(t *testing.T) {
	el := Page{
		Title: "Title",
	}

	dsl.RenderFast(el.Build())
}
