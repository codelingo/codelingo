# {{.Owner}}/{{.Name}} {{.Typ}} lexicon
{{$name := .Name}}
{{$owner := .Owner}}
{{$type := .Typ}}
##  {{.Name}} facts
{{range $fact, $children := .Facts}}{{ if ne $fact "not_implemented"}}<details><summary>{{$name}}.{{$fact}}</summary><p>

#### Example of finding every {{$fact}} and having a review flow comment on it:

```yaml
tenets:
  - name: find_all_{{$fact}}
    doc:  Example query to find all instances of {{$fact}}
    flows:
      codelingo/review
	       comment: This is a {{$fact}}.
	   query: |
	     import {{$owner}}/{{$type}}/{{$name}}

	     @review comment
	     {{$name}}.{{$fact}}
```
</p></details>
{{end}}
{{end}}