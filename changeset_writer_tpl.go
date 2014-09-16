package godiff

const commentsTplText = `
{{range $i, $_ := .}}

	{{"\n\n"}}

	{{if gt .Id 0}}
		[{{.Id}}@{{.Version}}] | {{.Author.DisplayName}} | {{.UpdatedDate}}
		{{"\n\n"}}
	{{end}}

	{{.Text | trimWhitespace}}

	{{"\n\n---"}}

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
		---{{" "}}
		{{if .Source.ToString}}
			{{.Source.ToString}}
		{{else}}
			/dev/null
		{{end}}
		{{"\t"}}
		{{.GetHashFrom}}
		{{"\n"}}

		+++{{" "}}
		{{if .Destination.ToString}}
			{{.Destination.ToString}}
		{{else}}
			/dev/null
		{{end}}
		{{"\t"}}
		{{.GetHashTo}}
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
