package layouts

import "github.com/TudorHulban/hxgo/dsl"

type Page struct {
	Title       string
	Description string
	Language    string

	Head []dsl.Node
	Body []dsl.Node
}

func (p *Page) Build() dsl.Node {
	body := make([]dsl.Node, 0, len(p.Body)+len(_JS))

	body = append(body, p.Body...)
	body = append(body, _JS...)

	return dsl.Doctype(
		dsl.HTML(
			dsl.Lang(p.Language),

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
							dsl.Text(p.Title),
						),

						_LinkCSSStyles,

						dsl.If(
							len(p.Description) > 0,
							dsl.Meta(
								dsl.Name("description"),
								dsl.Content(p.Description),
							),
						),
					},
					p.Head...,
				)...,
			),

			dsl.Body(
				body...,
			),
		),
	)
}
