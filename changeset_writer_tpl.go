package godiff

const commentsTplText = `
{{range $i, $_ := .}}

	{{"\n\n"}}

	{{if gt .Id 0}}
		[{{.Id}}@{{.Version}}] | {{.Author.DisplayName}} | {{.UpdatedDate}}
		{{"\n\n"}}
	{{end}}

	{{.Text}}

	{{"\n\n---"}}

	{{if not (last $i $)}}
	{{end}}

	{{writeComments .Comments | indent}}
{{end}}`

const changesetTplText = `
{{range $i, $d := .Diffs}}
	{{if .Note}}
		{{writeNote .Note}}
		{{"\n"}}
	{{end}}

	{{if .DiffComments}}
		{{"---" | comment}}
		{{"\n"}}
		{{writeComments .DiffComments | comment}}
		{{"\n"}}
	{{end}}

	{{if .Hunks}}
		--- {{.Source.ToString}}{{"\t"}}{{.GetHashFrom}}
		{{"\n"}}

		+++ {{.Destination.ToString}}{{"\t"}}{{.GetHashTo}}
		{{"\n"}}

		{{range .Hunks}}
			@@ -{{.SourceLine}},{{.SourceSpan}} +{{.DestinationLine}},{{.DestinationSpan}} @@
			{{"\n"}}

			{{range .Segments}}
				{{$segment := .}}
				{{range .Lines}}
					{{$segment.TextPrefix}}
					{{.Line}}

					{{"\n"}}
					{{if .Comments}}
						{{"---" | comment}}
						{{"\n"}}
						{{writeComments .Comments | comment}}
						{{"\n"}}
					{{end}}
				{{end}}

			{{end}}
		{{end}}
	{{end}}

	{{if not (last $i $.Diffs)}}
		{{"\n\n"}}
	{{end}}
{{end}}`
