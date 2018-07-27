# codelingo/go ast lexicon



##  go facts
<details><summary>go.array_type</summary><p>

#### Example of finding every array_type and having a review flow comment on it:

```yaml
tenets:
  - name: find_all_array_type
    doc:  Example query to find all instances of array_type
    flows:
      codelingo/review
	       comments: This is a array_type.
	   query: |
	     import codelingo/ast/go

	     @ review.comment
	     go.array_type
```
</p></details>

<details><summary>go.assign_stmt</summary><p>

#### Example of finding every assign_stmt and having a review flow comment on it:

```yaml
tenets:
  - name: find_all_assign_stmt
    doc:  Example query to find all instances of assign_stmt
    flows:
      codelingo/review
	       comments: This is a assign_stmt.
	   query: |
	     import codelingo/ast/go

	     @ review.comment
	     go.assign_stmt
```
</p></details>

<details><summary>go.bad_decl</summary><p>

#### Example of finding every bad_decl and having a review flow comment on it:

```yaml
tenets:
  - name: find_all_bad_decl
    doc:  Example query to find all instances of bad_decl
    flows:
      codelingo/review
	       comments: This is a bad_decl.
	   query: |
	     import codelingo/ast/go

	     @ review.comment
	     go.bad_decl
```
</p></details>

<details><summary>go.bad_expr</summary><p>

#### Example of finding every bad_expr and having a review flow comment on it:

```yaml
tenets:
  - name: find_all_bad_expr
    doc:  Example query to find all instances of bad_expr
    flows:
      codelingo/review
	       comments: This is a bad_expr.
	   query: |
	     import codelingo/ast/go

	     @ review.comment
	     go.bad_expr
```
</p></details>

<details><summary>go.bad_stmt</summary><p>

#### Example of finding every bad_stmt and having a review flow comment on it:

```yaml
tenets:
  - name: find_all_bad_stmt
    doc:  Example query to find all instances of bad_stmt
    flows:
      codelingo/review
	       comments: This is a bad_stmt.
	   query: |
	     import codelingo/ast/go

	     @ review.comment
	     go.bad_stmt
```
</p></details>

<details><summary>go.basic_lit</summary><p>

#### Example of finding every basic_lit and having a review flow comment on it:

```yaml
tenets:
  - name: find_all_basic_lit
    doc:  Example query to find all instances of basic_lit
    flows:
      codelingo/review
	       comments: This is a basic_lit.
	   query: |
	     import codelingo/ast/go

	     @ review.comment
	     go.basic_lit
```
</p></details>

<details><summary>go.binary_expr</summary><p>

#### Example of finding every binary_expr and having a review flow comment on it:

```yaml
tenets:
  - name: find_all_binary_expr
    doc:  Example query to find all instances of binary_expr
    flows:
      codelingo/review
	       comments: This is a binary_expr.
	   query: |
	     import codelingo/ast/go

	     @ review.comment
	     go.binary_expr
```
</p></details>

<details><summary>go.block_stmt</summary><p>

#### Example of finding every block_stmt and having a review flow comment on it:

```yaml
tenets:
  - name: find_all_block_stmt
    doc:  Example query to find all instances of block_stmt
    flows:
      codelingo/review
	       comments: This is a block_stmt.
	   query: |
	     import codelingo/ast/go

	     @ review.comment
	     go.block_stmt
```
</p></details>

<details><summary>go.branch_stmt</summary><p>

#### Example of finding every branch_stmt and having a review flow comment on it:

```yaml
tenets:
  - name: find_all_branch_stmt
    doc:  Example query to find all instances of branch_stmt
    flows:
      codelingo/review
	       comments: This is a branch_stmt.
	   query: |
	     import codelingo/ast/go

	     @ review.comment
	     go.branch_stmt
```
</p></details>

<details><summary>go.call_expr</summary><p>

#### Example of finding every call_expr and having a review flow comment on it:

```yaml
tenets:
  - name: find_all_call_expr
    doc:  Example query to find all instances of call_expr
    flows:
      codelingo/review
	       comments: This is a call_expr.
	   query: |
	     import codelingo/ast/go

	     @ review.comment
	     go.call_expr
```
</p></details>

<details><summary>go.case_clause</summary><p>

#### Example of finding every case_clause and having a review flow comment on it:

```yaml
tenets:
  - name: find_all_case_clause
    doc:  Example query to find all instances of case_clause
    flows:
      codelingo/review
	       comments: This is a case_clause.
	   query: |
	     import codelingo/ast/go

	     @ review.comment
	     go.case_clause
```
</p></details>

<details><summary>go.chan_type</summary><p>

#### Example of finding every chan_type and having a review flow comment on it:

```yaml
tenets:
  - name: find_all_chan_type
    doc:  Example query to find all instances of chan_type
    flows:
      codelingo/review
	       comments: This is a chan_type.
	   query: |
	     import codelingo/ast/go

	     @ review.comment
	     go.chan_type
```
</p></details>

<details><summary>go.comm_clause</summary><p>

#### Example of finding every comm_clause and having a review flow comment on it:

```yaml
tenets:
  - name: find_all_comm_clause
    doc:  Example query to find all instances of comm_clause
    flows:
      codelingo/review
	       comments: This is a comm_clause.
	   query: |
	     import codelingo/ast/go

	     @ review.comment
	     go.comm_clause
```
</p></details>

<details><summary>go.composite_lit</summary><p>

#### Example of finding every composite_lit and having a review flow comment on it:

```yaml
tenets:
  - name: find_all_composite_lit
    doc:  Example query to find all instances of composite_lit
    flows:
      codelingo/review
	       comments: This is a composite_lit.
	   query: |
	     import codelingo/ast/go

	     @ review.comment
	     go.composite_lit
```
</p></details>

<details><summary>go.decl_stmt</summary><p>

#### Example of finding every decl_stmt and having a review flow comment on it:

```yaml
tenets:
  - name: find_all_decl_stmt
    doc:  Example query to find all instances of decl_stmt
    flows:
      codelingo/review
	       comments: This is a decl_stmt.
	   query: |
	     import codelingo/ast/go

	     @ review.comment
	     go.decl_stmt
```
</p></details>

<details><summary>go.declarations</summary><p>

#### Example of finding every declarations and having a review flow comment on it:

```yaml
tenets:
  - name: find_all_declarations
    doc:  Example query to find all instances of declarations
    flows:
      codelingo/review
	       comments: This is a declarations.
	   query: |
	     import codelingo/ast/go

	     @ review.comment
	     go.declarations
```
</p></details>

<details><summary>go.defer_stmt</summary><p>

#### Example of finding every defer_stmt and having a review flow comment on it:

```yaml
tenets:
  - name: find_all_defer_stmt
    doc:  Example query to find all instances of defer_stmt
    flows:
      codelingo/review
	       comments: This is a defer_stmt.
	   query: |
	     import codelingo/ast/go

	     @ review.comment
	     go.defer_stmt
```
</p></details>

<details><summary>go.ellipsis</summary><p>

#### Example of finding every ellipsis and having a review flow comment on it:

```yaml
tenets:
  - name: find_all_ellipsis
    doc:  Example query to find all instances of ellipsis
    flows:
      codelingo/review
	       comments: This is a ellipsis.
	   query: |
	     import codelingo/ast/go

	     @ review.comment
	     go.ellipsis
```
</p></details>

<details><summary>go.empty_stmt</summary><p>

#### Example of finding every empty_stmt and having a review flow comment on it:

```yaml
tenets:
  - name: find_all_empty_stmt
    doc:  Example query to find all instances of empty_stmt
    flows:
      codelingo/review
	       comments: This is a empty_stmt.
	   query: |
	     import codelingo/ast/go

	     @ review.comment
	     go.empty_stmt
```
</p></details>

<details><summary>go.expr_stmt</summary><p>

#### Example of finding every expr_stmt and having a review flow comment on it:

```yaml
tenets:
  - name: find_all_expr_stmt
    doc:  Example query to find all instances of expr_stmt
    flows:
      codelingo/review
	       comments: This is a expr_stmt.
	   query: |
	     import codelingo/ast/go

	     @ review.comment
	     go.expr_stmt
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
	       comments: This is a expressions.
	   query: |
	     import codelingo/ast/go

	     @ review.comment
	     go.expressions
```
</p></details>

<details><summary>go.for_stmt</summary><p>

#### Example of finding every for_stmt and having a review flow comment on it:

```yaml
tenets:
  - name: find_all_for_stmt
    doc:  Example query to find all instances of for_stmt
    flows:
      codelingo/review
	       comments: This is a for_stmt.
	   query: |
	     import codelingo/ast/go

	     @ review.comment
	     go.for_stmt
```
</p></details>

<details><summary>go.func_decl</summary><p>

#### Example of finding every func_decl and having a review flow comment on it:

```yaml
tenets:
  - name: find_all_func_decl
    doc:  Example query to find all instances of func_decl
    flows:
      codelingo/review
	       comments: This is a func_decl.
	   query: |
	     import codelingo/ast/go

	     @ review.comment
	     go.func_decl
```
</p></details>

<details><summary>go.func_lit</summary><p>

#### Example of finding every func_lit and having a review flow comment on it:

```yaml
tenets:
  - name: find_all_func_lit
    doc:  Example query to find all instances of func_lit
    flows:
      codelingo/review
	       comments: This is a func_lit.
	   query: |
	     import codelingo/ast/go

	     @ review.comment
	     go.func_lit
```
</p></details>

<details><summary>go.func_type</summary><p>

#### Example of finding every func_type and having a review flow comment on it:

```yaml
tenets:
  - name: find_all_func_type
    doc:  Example query to find all instances of func_type
    flows:
      codelingo/review
	       comments: This is a func_type.
	   query: |
	     import codelingo/ast/go

	     @ review.comment
	     go.func_type
```
</p></details>

<details><summary>go.gen_decl</summary><p>

#### Example of finding every gen_decl and having a review flow comment on it:

```yaml
tenets:
  - name: find_all_gen_decl
    doc:  Example query to find all instances of gen_decl
    flows:
      codelingo/review
	       comments: This is a gen_decl.
	   query: |
	     import codelingo/ast/go

	     @ review.comment
	     go.gen_decl
```
</p></details>

<details><summary>go.go_stmt</summary><p>

#### Example of finding every go_stmt and having a review flow comment on it:

```yaml
tenets:
  - name: find_all_go_stmt
    doc:  Example query to find all instances of go_stmt
    flows:
      codelingo/review
	       comments: This is a go_stmt.
	   query: |
	     import codelingo/ast/go

	     @ review.comment
	     go.go_stmt
```
</p></details>

<details><summary>go.ident</summary><p>

#### Example of finding every ident and having a review flow comment on it:

```yaml
tenets:
  - name: find_all_ident
    doc:  Example query to find all instances of ident
    flows:
      codelingo/review
	       comments: This is a ident.
	   query: |
	     import codelingo/ast/go

	     @ review.comment
	     go.ident
```
</p></details>

<details><summary>go.if_stmt</summary><p>

#### Example of finding every if_stmt and having a review flow comment on it:

```yaml
tenets:
  - name: find_all_if_stmt
    doc:  Example query to find all instances of if_stmt
    flows:
      codelingo/review
	       comments: This is a if_stmt.
	   query: |
	     import codelingo/ast/go

	     @ review.comment
	     go.if_stmt
```
</p></details>

<details><summary>go.import_spec</summary><p>

#### Example of finding every import_spec and having a review flow comment on it:

```yaml
tenets:
  - name: find_all_import_spec
    doc:  Example query to find all instances of import_spec
    flows:
      codelingo/review
	       comments: This is a import_spec.
	   query: |
	     import codelingo/ast/go

	     @ review.comment
	     go.import_spec
```
</p></details>

<details><summary>go.inc_dec_stmt</summary><p>

#### Example of finding every inc_dec_stmt and having a review flow comment on it:

```yaml
tenets:
  - name: find_all_inc_dec_stmt
    doc:  Example query to find all instances of inc_dec_stmt
    flows:
      codelingo/review
	       comments: This is a inc_dec_stmt.
	   query: |
	     import codelingo/ast/go

	     @ review.comment
	     go.inc_dec_stmt
```
</p></details>

<details><summary>go.index_expr</summary><p>

#### Example of finding every index_expr and having a review flow comment on it:

```yaml
tenets:
  - name: find_all_index_expr
    doc:  Example query to find all instances of index_expr
    flows:
      codelingo/review
	       comments: This is a index_expr.
	   query: |
	     import codelingo/ast/go

	     @ review.comment
	     go.index_expr
```
</p></details>

<details><summary>go.interface_type</summary><p>

#### Example of finding every interface_type and having a review flow comment on it:

```yaml
tenets:
  - name: find_all_interface_type
    doc:  Example query to find all instances of interface_type
    flows:
      codelingo/review
	       comments: This is a interface_type.
	   query: |
	     import codelingo/ast/go

	     @ review.comment
	     go.interface_type
```
</p></details>

<details><summary>go.key_value_expr</summary><p>

#### Example of finding every key_value_expr and having a review flow comment on it:

```yaml
tenets:
  - name: find_all_key_value_expr
    doc:  Example query to find all instances of key_value_expr
    flows:
      codelingo/review
	       comments: This is a key_value_expr.
	   query: |
	     import codelingo/ast/go

	     @ review.comment
	     go.key_value_expr
```
</p></details>

<details><summary>go.labeled_stmt</summary><p>

#### Example of finding every labeled_stmt and having a review flow comment on it:

```yaml
tenets:
  - name: find_all_labeled_stmt
    doc:  Example query to find all instances of labeled_stmt
    flows:
      codelingo/review
	       comments: This is a labeled_stmt.
	   query: |
	     import codelingo/ast/go

	     @ review.comment
	     go.labeled_stmt
```
</p></details>

<details><summary>go.map_type</summary><p>

#### Example of finding every map_type and having a review flow comment on it:

```yaml
tenets:
  - name: find_all_map_type
    doc:  Example query to find all instances of map_type
    flows:
      codelingo/review
	       comments: This is a map_type.
	   query: |
	     import codelingo/ast/go

	     @ review.comment
	     go.map_type
```
</p></details>

<details><summary>go.paren_expr</summary><p>

#### Example of finding every paren_expr and having a review flow comment on it:

```yaml
tenets:
  - name: find_all_paren_expr
    doc:  Example query to find all instances of paren_expr
    flows:
      codelingo/review
	       comments: This is a paren_expr.
	   query: |
	     import codelingo/ast/go

	     @ review.comment
	     go.paren_expr
```
</p></details>

<details><summary>go.range_stmt</summary><p>

#### Example of finding every range_stmt and having a review flow comment on it:

```yaml
tenets:
  - name: find_all_range_stmt
    doc:  Example query to find all instances of range_stmt
    flows:
      codelingo/review
	       comments: This is a range_stmt.
	   query: |
	     import codelingo/ast/go

	     @ review.comment
	     go.range_stmt
```
</p></details>

<details><summary>go.return_stmt</summary><p>

#### Example of finding every return_stmt and having a review flow comment on it:

```yaml
tenets:
  - name: find_all_return_stmt
    doc:  Example query to find all instances of return_stmt
    flows:
      codelingo/review
	       comments: This is a return_stmt.
	   query: |
	     import codelingo/ast/go

	     @ review.comment
	     go.return_stmt
```
</p></details>

<details><summary>go.select_stmt</summary><p>

#### Example of finding every select_stmt and having a review flow comment on it:

```yaml
tenets:
  - name: find_all_select_stmt
    doc:  Example query to find all instances of select_stmt
    flows:
      codelingo/review
	       comments: This is a select_stmt.
	   query: |
	     import codelingo/ast/go

	     @ review.comment
	     go.select_stmt
```
</p></details>

<details><summary>go.selector_expr</summary><p>

#### Example of finding every selector_expr and having a review flow comment on it:

```yaml
tenets:
  - name: find_all_selector_expr
    doc:  Example query to find all instances of selector_expr
    flows:
      codelingo/review
	       comments: This is a selector_expr.
	   query: |
	     import codelingo/ast/go

	     @ review.comment
	     go.selector_expr
```
</p></details>

<details><summary>go.send_stmt</summary><p>

#### Example of finding every send_stmt and having a review flow comment on it:

```yaml
tenets:
  - name: find_all_send_stmt
    doc:  Example query to find all instances of send_stmt
    flows:
      codelingo/review
	       comments: This is a send_stmt.
	   query: |
	     import codelingo/ast/go

	     @ review.comment
	     go.send_stmt
```
</p></details>

<details><summary>go.slice_expr</summary><p>

#### Example of finding every slice_expr and having a review flow comment on it:

```yaml
tenets:
  - name: find_all_slice_expr
    doc:  Example query to find all instances of slice_expr
    flows:
      codelingo/review
	       comments: This is a slice_expr.
	   query: |
	     import codelingo/ast/go

	     @ review.comment
	     go.slice_expr
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
	       comments: This is a specs.
	   query: |
	     import codelingo/ast/go

	     @ review.comment
	     go.specs
```
</p></details>

<details><summary>go.star_expr</summary><p>

#### Example of finding every star_expr and having a review flow comment on it:

```yaml
tenets:
  - name: find_all_star_expr
    doc:  Example query to find all instances of star_expr
    flows:
      codelingo/review
	       comments: This is a star_expr.
	   query: |
	     import codelingo/ast/go

	     @ review.comment
	     go.star_expr
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
	       comments: This is a statements.
	   query: |
	     import codelingo/ast/go

	     @ review.comment
	     go.statements
```
</p></details>

<details><summary>go.struct_type</summary><p>

#### Example of finding every struct_type and having a review flow comment on it:

```yaml
tenets:
  - name: find_all_struct_type
    doc:  Example query to find all instances of struct_type
    flows:
      codelingo/review
	       comments: This is a struct_type.
	   query: |
	     import codelingo/ast/go

	     @ review.comment
	     go.struct_type
```
</p></details>

<details><summary>go.switch_stmt</summary><p>

#### Example of finding every switch_stmt and having a review flow comment on it:

```yaml
tenets:
  - name: find_all_switch_stmt
    doc:  Example query to find all instances of switch_stmt
    flows:
      codelingo/review
	       comments: This is a switch_stmt.
	   query: |
	     import codelingo/ast/go

	     @ review.comment
	     go.switch_stmt
```
</p></details>

<details><summary>go.type_assert_expr</summary><p>

#### Example of finding every type_assert_expr and having a review flow comment on it:

```yaml
tenets:
  - name: find_all_type_assert_expr
    doc:  Example query to find all instances of type_assert_expr
    flows:
      codelingo/review
	       comments: This is a type_assert_expr.
	   query: |
	     import codelingo/ast/go

	     @ review.comment
	     go.type_assert_expr
```
</p></details>

<details><summary>go.type_spec</summary><p>

#### Example of finding every type_spec and having a review flow comment on it:

```yaml
tenets:
  - name: find_all_type_spec
    doc:  Example query to find all instances of type_spec
    flows:
      codelingo/review
	       comments: This is a type_spec.
	   query: |
	     import codelingo/ast/go

	     @ review.comment
	     go.type_spec
```
</p></details>

<details><summary>go.type_switch_stmt</summary><p>

#### Example of finding every type_switch_stmt and having a review flow comment on it:

```yaml
tenets:
  - name: find_all_type_switch_stmt
    doc:  Example query to find all instances of type_switch_stmt
    flows:
      codelingo/review
	       comments: This is a type_switch_stmt.
	   query: |
	     import codelingo/ast/go

	     @ review.comment
	     go.type_switch_stmt
```
</p></details>

<details><summary>go.unary_expr</summary><p>

#### Example of finding every unary_expr and having a review flow comment on it:

```yaml
tenets:
  - name: find_all_unary_expr
    doc:  Example query to find all instances of unary_expr
    flows:
      codelingo/review
	       comments: This is a unary_expr.
	   query: |
	     import codelingo/ast/go

	     @ review.comment
	     go.unary_expr
```
</p></details>

<details><summary>go.value_spec</summary><p>

#### Example of finding every value_spec and having a review flow comment on it:

```yaml
tenets:
  - name: find_all_value_spec
    doc:  Example query to find all instances of value_spec
    flows:
      codelingo/review
	       comments: This is a value_spec.
	   query: |
	     import codelingo/ast/go

	     @ review.comment
	     go.value_spec
```
</p></details>

