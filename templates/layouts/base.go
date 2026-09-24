package layouts

import "github.com/TudorHulban/hxgo/dsl"

type Page struct {
	Title       string
	Description string
	Language    string

	Head []dsl.Node
	Body []dsl.Node
}

func (elem *Page) Build() dsl.Node {
	body := make([]dsl.Node, 0, len(elem.Body)+len(_JS))

	body = append(body,
		dsl.Div(
			dsl.Class("modal"),
			dsl.AttrID("modal-container"),
		),
	)

	body = append(body, elem.Body...)
	body = append(body, _JS...)

	return dsl.Doctype(
		dsl.HTML(
			dsl.Lang(elem.Language),

			dsl.Head(
				append(
					[]dsl.Node{
						dsl.Meta(
							dsl.Charset("utf-8"),
						),
						dsl.Meta(
							dsl.Name("viewport"),
							dsl.Content("width=device-width, initial-scale=1"),
						),
						dsl.Title(
							dsl.Text(elem.Title),
						),

						_LinkCSSStyles,
						_LinkCSSFontAwesome,

						dsl.If(
							len(elem.Description) > 0,
							dsl.Meta(
								dsl.Name("description"),
								dsl.Content(elem.Description),
							),
						),
					},
					elem.Head...,
				)...,
			),

			dsl.Body(
				body...,
			),
		),
	)
}
