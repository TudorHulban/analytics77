package fragments

import (
	"github.com/TudorHulban/hxgo/components/inputs"
	"github.com/TudorHulban/hxgo/dsl"
)

type SiteCombo struct {
	Status             string
	SelectorSiteValues []SelectorEntry
}

func (elem *SiteCombo) Build() dsl.Node {
	selection := inputs.InputSelect{
		CSSDivID: "siteSelector",
	}

	return dsl.Div(
		dsl.AttrClass("active-site-combo"),
		dsl.I(
			dsl.AttrClass("fas fa-globe"),
		),

		selection.RawSelect(),
	)
}

type TopBar struct{}
