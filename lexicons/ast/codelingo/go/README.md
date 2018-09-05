# codelingo/go ast lexicon



##  go facts
<details><summary>go.declarations</summary><p>

#### Example of finding every declarations and having a review flow comment on it:

```yaml
tenets:
  - name: find_all_declarations
    doc:  Example query to find all instances of declarations
    flows:
      codelingo/review
	       comment: This is a declarations.
	   query: |
	     import codelingo/ast/go

	     @ review.comment
	     go.declarations
```
</p></details>

<details><summary>go.expressions</summary><p>

#### Example of finding every expressions and having a review flow comment on it:

```yaml
tenets:
  - name: find_all_expressions
    doc:  Example query to find all instances of expressions
    flows:
      codelingo/review
	       comment: This is a expressions.
	   query: |
	     import codelingo/ast/go

	     @ review.comment
	     go.expressions
```
</p></details>

<details><summary>go.specs</summary><p>

#### Example of finding every specs and having a review flow comment on it:

```yaml
tenets:
  - name: find_all_specs
    doc:  Example query to find all instances of specs
    flows:
      codelingo/review
	       comment: This is a specs.
	   query: |
	     import codelingo/ast/go

	     @ review.comment
	     go.specs
```
</p></details>

<details><summary>go.statements</summary><p>

#### Example of finding every statements and having a review flow comment on it:

```yaml
tenets:
  - name: find_all_statements
    doc:  Example query to find all instances of statements
    flows:
      codelingo/review
	       comment: This is a statements.
	   query: |
	     import codelingo/ast/go

	     @ review.comment
	     go.statements
```
</p></details>

