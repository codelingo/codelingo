# codelingo/php ast lexicon



##  php facts
<details><summary>php.arg</summary><p>

#### Example of finding every arg and having a review flow comment on it:

```yaml
tenets:
  - name: find_all_arg
    doc:  Example query to find all instances of arg
    flows:
      codelingo/review
	       comments: This is a arg.
	   query: |
	     import codelingo/ast/php

	     @ review.comment
	     php.arg
```
</p></details>

<details><summary>php.const</summary><p>

#### Example of finding every const and having a review flow comment on it:

```yaml
tenets:
  - name: find_all_const
    doc:  Example query to find all instances of const
    flows:
      codelingo/review
	       comments: This is a const.
	   query: |
	     import codelingo/ast/php

	     @ review.comment
	     php.const
```
</p></details>

<details><summary>php.expr_array</summary><p>

#### Example of finding every expr_array and having a review flow comment on it:

```yaml
tenets:
  - name: find_all_expr_array
    doc:  Example query to find all instances of expr_array
    flows:
      codelingo/review
	       comments: This is a expr_array.
	   query: |
	     import codelingo/ast/php

	     @ review.comment
	     php.expr_array
```
</p></details>

<details><summary>php.expr_arraydimfetch</summary><p>

#### Example of finding every expr_arraydimfetch and having a review flow comment on it:

```yaml
tenets:
  - name: find_all_expr_arraydimfetch
    doc:  Example query to find all instances of expr_arraydimfetch
    flows:
      codelingo/review
	       comments: This is a expr_arraydimfetch.
	   query: |
	     import codelingo/ast/php

	     @ review.comment
	     php.expr_arraydimfetch
```
</p></details>

<details><summary>php.expr_arrayitem</summary><p>

#### Example of finding every expr_arrayitem and having a review flow comment on it:

```yaml
tenets:
  - name: find_all_expr_arrayitem
    doc:  Example query to find all instances of expr_arrayitem
    flows:
      codelingo/review
	       comments: This is a expr_arrayitem.
	   query: |
	     import codelingo/ast/php

	     @ review.comment
	     php.expr_arrayitem
```
</p></details>

<details><summary>php.expr_assign</summary><p>

#### Example of finding every expr_assign and having a review flow comment on it:

```yaml
tenets:
  - name: find_all_expr_assign
    doc:  Example query to find all instances of expr_assign
    flows:
      codelingo/review
	       comments: This is a expr_assign.
	   query: |
	     import codelingo/ast/php

	     @ review.comment
	     php.expr_assign
```
</p></details>

<details><summary>php.expr_assignop</summary><p>

#### Example of finding every expr_assignop and having a review flow comment on it:

```yaml
tenets:
  - name: find_all_expr_assignop
    doc:  Example query to find all instances of expr_assignop
    flows:
      codelingo/review
	       comments: This is a expr_assignop.
	   query: |
	     import codelingo/ast/php

	     @ review.comment
	     php.expr_assignop
```
</p></details>

<details><summary>php.expr_assignop_bitwiseand</summary><p>

#### Example of finding every expr_assignop_bitwiseand and having a review flow comment on it:

```yaml
tenets:
  - name: find_all_expr_assignop_bitwiseand
    doc:  Example query to find all instances of expr_assignop_bitwiseand
    flows:
      codelingo/review
	       comments: This is a expr_assignop_bitwiseand.
	   query: |
	     import codelingo/ast/php

	     @ review.comment
	     php.expr_assignop_bitwiseand
```
</p></details>

<details><summary>php.expr_assignop_bitwiseor</summary><p>

#### Example of finding every expr_assignop_bitwiseor and having a review flow comment on it:

```yaml
tenets:
  - name: find_all_expr_assignop_bitwiseor
    doc:  Example query to find all instances of expr_assignop_bitwiseor
    flows:
      codelingo/review
	       comments: This is a expr_assignop_bitwiseor.
	   query: |
	     import codelingo/ast/php

	     @ review.comment
	     php.expr_assignop_bitwiseor
```
</p></details>

<details><summary>php.expr_assignop_bitwisexor</summary><p>

#### Example of finding every expr_assignop_bitwisexor and having a review flow comment on it:

```yaml
tenets:
  - name: find_all_expr_assignop_bitwisexor
    doc:  Example query to find all instances of expr_assignop_bitwisexor
    flows:
      codelingo/review
	       comments: This is a expr_assignop_bitwisexor.
	   query: |
	     import codelingo/ast/php

	     @ review.comment
	     php.expr_assignop_bitwisexor
```
</p></details>

<details><summary>php.expr_assignop_concat</summary><p>

#### Example of finding every expr_assignop_concat and having a review flow comment on it:

```yaml
tenets:
  - name: find_all_expr_assignop_concat
    doc:  Example query to find all instances of expr_assignop_concat
    flows:
      codelingo/review
	       comments: This is a expr_assignop_concat.
	   query: |
	     import codelingo/ast/php

	     @ review.comment
	     php.expr_assignop_concat
```
</p></details>

<details><summary>php.expr_assignop_div</summary><p>

#### Example of finding every expr_assignop_div and having a review flow comment on it:

```yaml
tenets:
  - name: find_all_expr_assignop_div
    doc:  Example query to find all instances of expr_assignop_div
    flows:
      codelingo/review
	       comments: This is a expr_assignop_div.
	   query: |
	     import codelingo/ast/php

	     @ review.comment
	     php.expr_assignop_div
```
</p></details>

<details><summary>php.expr_assignop_minus</summary><p>

#### Example of finding every expr_assignop_minus and having a review flow comment on it:

```yaml
tenets:
  - name: find_all_expr_assignop_minus
    doc:  Example query to find all instances of expr_assignop_minus
    flows:
      codelingo/review
	       comments: This is a expr_assignop_minus.
	   query: |
	     import codelingo/ast/php

	     @ review.comment
	     php.expr_assignop_minus
```
</p></details>

<details><summary>php.expr_assignop_mod</summary><p>

#### Example of finding every expr_assignop_mod and having a review flow comment on it:

```yaml
tenets:
  - name: find_all_expr_assignop_mod
    doc:  Example query to find all instances of expr_assignop_mod
    flows:
      codelingo/review
	       comments: This is a expr_assignop_mod.
	   query: |
	     import codelingo/ast/php

	     @ review.comment
	     php.expr_assignop_mod
```
</p></details>

<details><summary>php.expr_assignop_mul</summary><p>

#### Example of finding every expr_assignop_mul and having a review flow comment on it:

```yaml
tenets:
  - name: find_all_expr_assignop_mul
    doc:  Example query to find all instances of expr_assignop_mul
    flows:
      codelingo/review
	       comments: This is a expr_assignop_mul.
	   query: |
	     import codelingo/ast/php

	     @ review.comment
	     php.expr_assignop_mul
```
</p></details>

<details><summary>php.expr_assignop_plus</summary><p>

#### Example of finding every expr_assignop_plus and having a review flow comment on it:

```yaml
tenets:
  - name: find_all_expr_assignop_plus
    doc:  Example query to find all instances of expr_assignop_plus
    flows:
      codelingo/review
	       comments: This is a expr_assignop_plus.
	   query: |
	     import codelingo/ast/php

	     @ review.comment
	     php.expr_assignop_plus
```
</p></details>

<details><summary>php.expr_assignop_pow</summary><p>

#### Example of finding every expr_assignop_pow and having a review flow comment on it:

```yaml
tenets:
  - name: find_all_expr_assignop_pow
    doc:  Example query to find all instances of expr_assignop_pow
    flows:
      codelingo/review
	       comments: This is a expr_assignop_pow.
	   query: |
	     import codelingo/ast/php

	     @ review.comment
	     php.expr_assignop_pow
```
</p></details>

<details><summary>php.expr_assignop_shiftleft</summary><p>

#### Example of finding every expr_assignop_shiftleft and having a review flow comment on it:

```yaml
tenets:
  - name: find_all_expr_assignop_shiftleft
    doc:  Example query to find all instances of expr_assignop_shiftleft
    flows:
      codelingo/review
	       comments: This is a expr_assignop_shiftleft.
	   query: |
	     import codelingo/ast/php

	     @ review.comment
	     php.expr_assignop_shiftleft
```
</p></details>

<details><summary>php.expr_assignop_shiftright</summary><p>

#### Example of finding every expr_assignop_shiftright and having a review flow comment on it:

```yaml
tenets:
  - name: find_all_expr_assignop_shiftright
    doc:  Example query to find all instances of expr_assignop_shiftright
    flows:
      codelingo/review
	       comments: This is a expr_assignop_shiftright.
	   query: |
	     import codelingo/ast/php

	     @ review.comment
	     php.expr_assignop_shiftright
```
</p></details>

<details><summary>php.expr_assignref</summary><p>

#### Example of finding every expr_assignref and having a review flow comment on it:

```yaml
tenets:
  - name: find_all_expr_assignref
    doc:  Example query to find all instances of expr_assignref
    flows:
      codelingo/review
	       comments: This is a expr_assignref.
	   query: |
	     import codelingo/ast/php

	     @ review.comment
	     php.expr_assignref
```
</p></details>

<details><summary>php.expr_binaryop</summary><p>

#### Example of finding every expr_binaryop and having a review flow comment on it:

```yaml
tenets:
  - name: find_all_expr_binaryop
    doc:  Example query to find all instances of expr_binaryop
    flows:
      codelingo/review
	       comments: This is a expr_binaryop.
	   query: |
	     import codelingo/ast/php

	     @ review.comment
	     php.expr_binaryop
```
</p></details>

<details><summary>php.expr_binaryop_bitwiseand</summary><p>

#### Example of finding every expr_binaryop_bitwiseand and having a review flow comment on it:

```yaml
tenets:
  - name: find_all_expr_binaryop_bitwiseand
    doc:  Example query to find all instances of expr_binaryop_bitwiseand
    flows:
      codelingo/review
	       comments: This is a expr_binaryop_bitwiseand.
	   query: |
	     import codelingo/ast/php

	     @ review.comment
	     php.expr_binaryop_bitwiseand
```
</p></details>

<details><summary>php.expr_binaryop_bitwiseor</summary><p>

#### Example of finding every expr_binaryop_bitwiseor and having a review flow comment on it:

```yaml
tenets:
  - name: find_all_expr_binaryop_bitwiseor
    doc:  Example query to find all instances of expr_binaryop_bitwiseor
    flows:
      codelingo/review
	       comments: This is a expr_binaryop_bitwiseor.
	   query: |
	     import codelingo/ast/php

	     @ review.comment
	     php.expr_binaryop_bitwiseor
```
</p></details>

<details><summary>php.expr_binaryop_bitwisexor</summary><p>

#### Example of finding every expr_binaryop_bitwisexor and having a review flow comment on it:

```yaml
tenets:
  - name: find_all_expr_binaryop_bitwisexor
    doc:  Example query to find all instances of expr_binaryop_bitwisexor
    flows:
      codelingo/review
	       comments: This is a expr_binaryop_bitwisexor.
	   query: |
	     import codelingo/ast/php

	     @ review.comment
	     php.expr_binaryop_bitwisexor
```
</p></details>

<details><summary>php.expr_binaryop_booleanand</summary><p>

#### Example of finding every expr_binaryop_booleanand and having a review flow comment on it:

```yaml
tenets:
  - name: find_all_expr_binaryop_booleanand
    doc:  Example query to find all instances of expr_binaryop_booleanand
    flows:
      codelingo/review
	       comments: This is a expr_binaryop_booleanand.
	   query: |
	     import codelingo/ast/php

	     @ review.comment
	     php.expr_binaryop_booleanand
```
</p></details>

<details><summary>php.expr_binaryop_booleanor</summary><p>

#### Example of finding every expr_binaryop_booleanor and having a review flow comment on it:

```yaml
tenets:
  - name: find_all_expr_binaryop_booleanor
    doc:  Example query to find all instances of expr_binaryop_booleanor
    flows:
      codelingo/review
	       comments: This is a expr_binaryop_booleanor.
	   query: |
	     import codelingo/ast/php

	     @ review.comment
	     php.expr_binaryop_booleanor
```
</p></details>

<details><summary>php.expr_binaryop_coalesce</summary><p>

#### Example of finding every expr_binaryop_coalesce and having a review flow comment on it:

```yaml
tenets:
  - name: find_all_expr_binaryop_coalesce
    doc:  Example query to find all instances of expr_binaryop_coalesce
    flows:
      codelingo/review
	       comments: This is a expr_binaryop_coalesce.
	   query: |
	     import codelingo/ast/php

	     @ review.comment
	     php.expr_binaryop_coalesce
```
</p></details>

<details><summary>php.expr_binaryop_concat</summary><p>

#### Example of finding every expr_binaryop_concat and having a review flow comment on it:

```yaml
tenets:
  - name: find_all_expr_binaryop_concat
    doc:  Example query to find all instances of expr_binaryop_concat
    flows:
      codelingo/review
	       comments: This is a expr_binaryop_concat.
	   query: |
	     import codelingo/ast/php

	     @ review.comment
	     php.expr_binaryop_concat
```
</p></details>

<details><summary>php.expr_binaryop_div</summary><p>

#### Example of finding every expr_binaryop_div and having a review flow comment on it:

```yaml
tenets:
  - name: find_all_expr_binaryop_div
    doc:  Example query to find all instances of expr_binaryop_div
    flows:
      codelingo/review
	       comments: This is a expr_binaryop_div.
	   query: |
	     import codelingo/ast/php

	     @ review.comment
	     php.expr_binaryop_div
```
</p></details>

<details><summary>php.expr_binaryop_equal</summary><p>

#### Example of finding every expr_binaryop_equal and having a review flow comment on it:

```yaml
tenets:
  - name: find_all_expr_binaryop_equal
    doc:  Example query to find all instances of expr_binaryop_equal
    flows:
      codelingo/review
	       comments: This is a expr_binaryop_equal.
	   query: |
	     import codelingo/ast/php

	     @ review.comment
	     php.expr_binaryop_equal
```
</p></details>

<details><summary>php.expr_binaryop_greater</summary><p>

#### Example of finding every expr_binaryop_greater and having a review flow comment on it:

```yaml
tenets:
  - name: find_all_expr_binaryop_greater
    doc:  Example query to find all instances of expr_binaryop_greater
    flows:
      codelingo/review
	       comments: This is a expr_binaryop_greater.
	   query: |
	     import codelingo/ast/php

	     @ review.comment
	     php.expr_binaryop_greater
```
</p></details>

<details><summary>php.expr_binaryop_greaterorequal</summary><p>

#### Example of finding every expr_binaryop_greaterorequal and having a review flow comment on it:

```yaml
tenets:
  - name: find_all_expr_binaryop_greaterorequal
    doc:  Example query to find all instances of expr_binaryop_greaterorequal
    flows:
      codelingo/review
	       comments: This is a expr_binaryop_greaterorequal.
	   query: |
	     import codelingo/ast/php

	     @ review.comment
	     php.expr_binaryop_greaterorequal
```
</p></details>

<details><summary>php.expr_binaryop_identical</summary><p>

#### Example of finding every expr_binaryop_identical and having a review flow comment on it:

```yaml
tenets:
  - name: find_all_expr_binaryop_identical
    doc:  Example query to find all instances of expr_binaryop_identical
    flows:
      codelingo/review
	       comments: This is a expr_binaryop_identical.
	   query: |
	     import codelingo/ast/php

	     @ review.comment
	     php.expr_binaryop_identical
```
</p></details>

<details><summary>php.expr_binaryop_logicaland</summary><p>

#### Example of finding every expr_binaryop_logicaland and having a review flow comment on it:

```yaml
tenets:
  - name: find_all_expr_binaryop_logicaland
    doc:  Example query to find all instances of expr_binaryop_logicaland
    flows:
      codelingo/review
	       comments: This is a expr_binaryop_logicaland.
	   query: |
	     import codelingo/ast/php

	     @ review.comment
	     php.expr_binaryop_logicaland
```
</p></details>

<details><summary>php.expr_binaryop_logicalor</summary><p>

#### Example of finding every expr_binaryop_logicalor and having a review flow comment on it:

```yaml
tenets:
  - name: find_all_expr_binaryop_logicalor
    doc:  Example query to find all instances of expr_binaryop_logicalor
    flows:
      codelingo/review
	       comments: This is a expr_binaryop_logicalor.
	   query: |
	     import codelingo/ast/php

	     @ review.comment
	     php.expr_binaryop_logicalor
```
</p></details>

<details><summary>php.expr_binaryop_logicalxor</summary><p>

#### Example of finding every expr_binaryop_logicalxor and having a review flow comment on it:

```yaml
tenets:
  - name: find_all_expr_binaryop_logicalxor
    doc:  Example query to find all instances of expr_binaryop_logicalxor
    flows:
      codelingo/review
	       comments: This is a expr_binaryop_logicalxor.
	   query: |
	     import codelingo/ast/php

	     @ review.comment
	     php.expr_binaryop_logicalxor
```
</p></details>

<details><summary>php.expr_binaryop_minus</summary><p>

#### Example of finding every expr_binaryop_minus and having a review flow comment on it:

```yaml
tenets:
  - name: find_all_expr_binaryop_minus
    doc:  Example query to find all instances of expr_binaryop_minus
    flows:
      codelingo/review
	       comments: This is a expr_binaryop_minus.
	   query: |
	     import codelingo/ast/php

	     @ review.comment
	     php.expr_binaryop_minus
```
</p></details>

<details><summary>php.expr_binaryop_mod</summary><p>

#### Example of finding every expr_binaryop_mod and having a review flow comment on it:

```yaml
tenets:
  - name: find_all_expr_binaryop_mod
    doc:  Example query to find all instances of expr_binaryop_mod
    flows:
      codelingo/review
	       comments: This is a expr_binaryop_mod.
	   query: |
	     import codelingo/ast/php

	     @ review.comment
	     php.expr_binaryop_mod
```
</p></details>

<details><summary>php.expr_binaryop_mul</summary><p>

#### Example of finding every expr_binaryop_mul and having a review flow comment on it:

```yaml
tenets:
  - name: find_all_expr_binaryop_mul
    doc:  Example query to find all instances of expr_binaryop_mul
    flows:
      codelingo/review
	       comments: This is a expr_binaryop_mul.
	   query: |
	     import codelingo/ast/php

	     @ review.comment
	     php.expr_binaryop_mul
```
</p></details>

<details><summary>php.expr_binaryop_notequal</summary><p>

#### Example of finding every expr_binaryop_notequal and having a review flow comment on it:

```yaml
tenets:
  - name: find_all_expr_binaryop_notequal
    doc:  Example query to find all instances of expr_binaryop_notequal
    flows:
      codelingo/review
	       comments: This is a expr_binaryop_notequal.
	   query: |
	     import codelingo/ast/php

	     @ review.comment
	     php.expr_binaryop_notequal
```
</p></details>

<details><summary>php.expr_binaryop_notidentical</summary><p>

#### Example of finding every expr_binaryop_notidentical and having a review flow comment on it:

```yaml
tenets:
  - name: find_all_expr_binaryop_notidentical
    doc:  Example query to find all instances of expr_binaryop_notidentical
    flows:
      codelingo/review
	       comments: This is a expr_binaryop_notidentical.
	   query: |
	     import codelingo/ast/php

	     @ review.comment
	     php.expr_binaryop_notidentical
```
</p></details>

<details><summary>php.expr_binaryop_plus</summary><p>

#### Example of finding every expr_binaryop_plus and having a review flow comment on it:

```yaml
tenets:
  - name: find_all_expr_binaryop_plus
    doc:  Example query to find all instances of expr_binaryop_plus
    flows:
      codelingo/review
	       comments: This is a expr_binaryop_plus.
	   query: |
	     import codelingo/ast/php

	     @ review.comment
	     php.expr_binaryop_plus
```
</p></details>

<details><summary>php.expr_binaryop_pow</summary><p>

#### Example of finding every expr_binaryop_pow and having a review flow comment on it:

```yaml
tenets:
  - name: find_all_expr_binaryop_pow
    doc:  Example query to find all instances of expr_binaryop_pow
    flows:
      codelingo/review
	       comments: This is a expr_binaryop_pow.
	   query: |
	     import codelingo/ast/php

	     @ review.comment
	     php.expr_binaryop_pow
```
</p></details>

<details><summary>php.expr_binaryop_shiftleft</summary><p>

#### Example of finding every expr_binaryop_shiftleft and having a review flow comment on it:

```yaml
tenets:
  - name: find_all_expr_binaryop_shiftleft
    doc:  Example query to find all instances of expr_binaryop_shiftleft
    flows:
      codelingo/review
	       comments: This is a expr_binaryop_shiftleft.
	   query: |
	     import codelingo/ast/php

	     @ review.comment
	     php.expr_binaryop_shiftleft
```
</p></details>

<details><summary>php.expr_binaryop_shiftright</summary><p>

#### Example of finding every expr_binaryop_shiftright and having a review flow comment on it:

```yaml
tenets:
  - name: find_all_expr_binaryop_shiftright
    doc:  Example query to find all instances of expr_binaryop_shiftright
    flows:
      codelingo/review
	       comments: This is a expr_binaryop_shiftright.
	   query: |
	     import codelingo/ast/php

	     @ review.comment
	     php.expr_binaryop_shiftright
```
</p></details>

<details><summary>php.expr_binaryop_smaller</summary><p>

#### Example of finding every expr_binaryop_smaller and having a review flow comment on it:

```yaml
tenets:
  - name: find_all_expr_binaryop_smaller
    doc:  Example query to find all instances of expr_binaryop_smaller
    flows:
      codelingo/review
	       comments: This is a expr_binaryop_smaller.
	   query: |
	     import codelingo/ast/php

	     @ review.comment
	     php.expr_binaryop_smaller
```
</p></details>

<details><summary>php.expr_binaryop_smallerorequal</summary><p>

#### Example of finding every expr_binaryop_smallerorequal and having a review flow comment on it:

```yaml
tenets:
  - name: find_all_expr_binaryop_smallerorequal
    doc:  Example query to find all instances of expr_binaryop_smallerorequal
    flows:
      codelingo/review
	       comments: This is a expr_binaryop_smallerorequal.
	   query: |
	     import codelingo/ast/php

	     @ review.comment
	     php.expr_binaryop_smallerorequal
```
</p></details>

<details><summary>php.expr_binaryop_spaceship</summary><p>

#### Example of finding every expr_binaryop_spaceship and having a review flow comment on it:

```yaml
tenets:
  - name: find_all_expr_binaryop_spaceship
    doc:  Example query to find all instances of expr_binaryop_spaceship
    flows:
      codelingo/review
	       comments: This is a expr_binaryop_spaceship.
	   query: |
	     import codelingo/ast/php

	     @ review.comment
	     php.expr_binaryop_spaceship
```
</p></details>

<details><summary>php.expr_bitwisenot</summary><p>

#### Example of finding every expr_bitwisenot and having a review flow comment on it:

```yaml
tenets:
  - name: find_all_expr_bitwisenot
    doc:  Example query to find all instances of expr_bitwisenot
    flows:
      codelingo/review
	       comments: This is a expr_bitwisenot.
	   query: |
	     import codelingo/ast/php

	     @ review.comment
	     php.expr_bitwisenot
```
</p></details>

<details><summary>php.expr_booleannot</summary><p>

#### Example of finding every expr_booleannot and having a review flow comment on it:

```yaml
tenets:
  - name: find_all_expr_booleannot
    doc:  Example query to find all instances of expr_booleannot
    flows:
      codelingo/review
	       comments: This is a expr_booleannot.
	   query: |
	     import codelingo/ast/php

	     @ review.comment
	     php.expr_booleannot
```
</p></details>

<details><summary>php.expr_cast_array</summary><p>

#### Example of finding every expr_cast_array and having a review flow comment on it:

```yaml
tenets:
  - name: find_all_expr_cast_array
    doc:  Example query to find all instances of expr_cast_array
    flows:
      codelingo/review
	       comments: This is a expr_cast_array.
	   query: |
	     import codelingo/ast/php

	     @ review.comment
	     php.expr_cast_array
```
</p></details>

<details><summary>php.expr_cast_bool</summary><p>

#### Example of finding every expr_cast_bool and having a review flow comment on it:

```yaml
tenets:
  - name: find_all_expr_cast_bool
    doc:  Example query to find all instances of expr_cast_bool
    flows:
      codelingo/review
	       comments: This is a expr_cast_bool.
	   query: |
	     import codelingo/ast/php

	     @ review.comment
	     php.expr_cast_bool
```
</p></details>

<details><summary>php.expr_cast_double</summary><p>

#### Example of finding every expr_cast_double and having a review flow comment on it:

```yaml
tenets:
  - name: find_all_expr_cast_double
    doc:  Example query to find all instances of expr_cast_double
    flows:
      codelingo/review
	       comments: This is a expr_cast_double.
	   query: |
	     import codelingo/ast/php

	     @ review.comment
	     php.expr_cast_double
```
</p></details>

<details><summary>php.expr_cast_int</summary><p>

#### Example of finding every expr_cast_int and having a review flow comment on it:

```yaml
tenets:
  - name: find_all_expr_cast_int
    doc:  Example query to find all instances of expr_cast_int
    flows:
      codelingo/review
	       comments: This is a expr_cast_int.
	   query: |
	     import codelingo/ast/php

	     @ review.comment
	     php.expr_cast_int
```
</p></details>

<details><summary>php.expr_cast_object</summary><p>

#### Example of finding every expr_cast_object and having a review flow comment on it:

```yaml
tenets:
  - name: find_all_expr_cast_object
    doc:  Example query to find all instances of expr_cast_object
    flows:
      codelingo/review
	       comments: This is a expr_cast_object.
	   query: |
	     import codelingo/ast/php

	     @ review.comment
	     php.expr_cast_object
```
</p></details>

<details><summary>php.expr_cast_string</summary><p>

#### Example of finding every expr_cast_string and having a review flow comment on it:

```yaml
tenets:
  - name: find_all_expr_cast_string
    doc:  Example query to find all instances of expr_cast_string
    flows:
      codelingo/review
	       comments: This is a expr_cast_string.
	   query: |
	     import codelingo/ast/php

	     @ review.comment
	     php.expr_cast_string
```
</p></details>

<details><summary>php.expr_cast_unset</summary><p>

#### Example of finding every expr_cast_unset and having a review flow comment on it:

```yaml
tenets:
  - name: find_all_expr_cast_unset
    doc:  Example query to find all instances of expr_cast_unset
    flows:
      codelingo/review
	       comments: This is a expr_cast_unset.
	   query: |
	     import codelingo/ast/php

	     @ review.comment
	     php.expr_cast_unset
```
</p></details>

<details><summary>php.expr_classconstfetch</summary><p>

#### Example of finding every expr_classconstfetch and having a review flow comment on it:

```yaml
tenets:
  - name: find_all_expr_classconstfetch
    doc:  Example query to find all instances of expr_classconstfetch
    flows:
      codelingo/review
	       comments: This is a expr_classconstfetch.
	   query: |
	     import codelingo/ast/php

	     @ review.comment
	     php.expr_classconstfetch
```
</p></details>

<details><summary>php.expr_clone</summary><p>

#### Example of finding every expr_clone and having a review flow comment on it:

```yaml
tenets:
  - name: find_all_expr_clone
    doc:  Example query to find all instances of expr_clone
    flows:
      codelingo/review
	       comments: This is a expr_clone.
	   query: |
	     import codelingo/ast/php

	     @ review.comment
	     php.expr_clone
```
</p></details>

<details><summary>php.expr_closure</summary><p>

#### Example of finding every expr_closure and having a review flow comment on it:

```yaml
tenets:
  - name: find_all_expr_closure
    doc:  Example query to find all instances of expr_closure
    flows:
      codelingo/review
	       comments: This is a expr_closure.
	   query: |
	     import codelingo/ast/php

	     @ review.comment
	     php.expr_closure
```
</p></details>

<details><summary>php.expr_closureuse</summary><p>

#### Example of finding every expr_closureuse and having a review flow comment on it:

```yaml
tenets:
  - name: find_all_expr_closureuse
    doc:  Example query to find all instances of expr_closureuse
    flows:
      codelingo/review
	       comments: This is a expr_closureuse.
	   query: |
	     import codelingo/ast/php

	     @ review.comment
	     php.expr_closureuse
```
</p></details>

<details><summary>php.expr_constfetch</summary><p>

#### Example of finding every expr_constfetch and having a review flow comment on it:

```yaml
tenets:
  - name: find_all_expr_constfetch
    doc:  Example query to find all instances of expr_constfetch
    flows:
      codelingo/review
	       comments: This is a expr_constfetch.
	   query: |
	     import codelingo/ast/php

	     @ review.comment
	     php.expr_constfetch
```
</p></details>

<details><summary>php.expr_empty</summary><p>

#### Example of finding every expr_empty and having a review flow comment on it:

```yaml
tenets:
  - name: find_all_expr_empty
    doc:  Example query to find all instances of expr_empty
    flows:
      codelingo/review
	       comments: This is a expr_empty.
	   query: |
	     import codelingo/ast/php

	     @ review.comment
	     php.expr_empty
```
</p></details>

<details><summary>php.expr_errorsuppress</summary><p>

#### Example of finding every expr_errorsuppress and having a review flow comment on it:

```yaml
tenets:
  - name: find_all_expr_errorsuppress
    doc:  Example query to find all instances of expr_errorsuppress
    flows:
      codelingo/review
	       comments: This is a expr_errorsuppress.
	   query: |
	     import codelingo/ast/php

	     @ review.comment
	     php.expr_errorsuppress
```
</p></details>

<details><summary>php.expr_eval</summary><p>

#### Example of finding every expr_eval and having a review flow comment on it:

```yaml
tenets:
  - name: find_all_expr_eval
    doc:  Example query to find all instances of expr_eval
    flows:
      codelingo/review
	       comments: This is a expr_eval.
	   query: |
	     import codelingo/ast/php

	     @ review.comment
	     php.expr_eval
```
</p></details>

<details><summary>php.expr_exit</summary><p>

#### Example of finding every expr_exit and having a review flow comment on it:

```yaml
tenets:
  - name: find_all_expr_exit
    doc:  Example query to find all instances of expr_exit
    flows:
      codelingo/review
	       comments: This is a expr_exit.
	   query: |
	     import codelingo/ast/php

	     @ review.comment
	     php.expr_exit
```
</p></details>

<details><summary>php.expr_funccall</summary><p>

#### Example of finding every expr_funccall and having a review flow comment on it:

```yaml
tenets:
  - name: find_all_expr_funccall
    doc:  Example query to find all instances of expr_funccall
    flows:
      codelingo/review
	       comments: This is a expr_funccall.
	   query: |
	     import codelingo/ast/php

	     @ review.comment
	     php.expr_funccall
```
</p></details>

<details><summary>php.expr_include</summary><p>

#### Example of finding every expr_include and having a review flow comment on it:

```yaml
tenets:
  - name: find_all_expr_include
    doc:  Example query to find all instances of expr_include
    flows:
      codelingo/review
	       comments: This is a expr_include.
	   query: |
	     import codelingo/ast/php

	     @ review.comment
	     php.expr_include
```
</p></details>

<details><summary>php.expr_instanceof</summary><p>

#### Example of finding every expr_instanceof and having a review flow comment on it:

```yaml
tenets:
  - name: find_all_expr_instanceof
    doc:  Example query to find all instances of expr_instanceof
    flows:
      codelingo/review
	       comments: This is a expr_instanceof.
	   query: |
	     import codelingo/ast/php

	     @ review.comment
	     php.expr_instanceof
```
</p></details>

<details><summary>php.expr_isset</summary><p>

#### Example of finding every expr_isset and having a review flow comment on it:

```yaml
tenets:
  - name: find_all_expr_isset
    doc:  Example query to find all instances of expr_isset
    flows:
      codelingo/review
	       comments: This is a expr_isset.
	   query: |
	     import codelingo/ast/php

	     @ review.comment
	     php.expr_isset
```
</p></details>

<details><summary>php.expr_list</summary><p>

#### Example of finding every expr_list and having a review flow comment on it:

```yaml
tenets:
  - name: find_all_expr_list
    doc:  Example query to find all instances of expr_list
    flows:
      codelingo/review
	       comments: This is a expr_list.
	   query: |
	     import codelingo/ast/php

	     @ review.comment
	     php.expr_list
```
</p></details>

<details><summary>php.expr_methodcall</summary><p>

#### Example of finding every expr_methodcall and having a review flow comment on it:

```yaml
tenets:
  - name: find_all_expr_methodcall
    doc:  Example query to find all instances of expr_methodcall
    flows:
      codelingo/review
	       comments: This is a expr_methodcall.
	   query: |
	     import codelingo/ast/php

	     @ review.comment
	     php.expr_methodcall
```
</p></details>

<details><summary>php.expr_new</summary><p>

#### Example of finding every expr_new and having a review flow comment on it:

```yaml
tenets:
  - name: find_all_expr_new
    doc:  Example query to find all instances of expr_new
    flows:
      codelingo/review
	       comments: This is a expr_new.
	   query: |
	     import codelingo/ast/php

	     @ review.comment
	     php.expr_new
```
</p></details>

<details><summary>php.expr_postdec</summary><p>

#### Example of finding every expr_postdec and having a review flow comment on it:

```yaml
tenets:
  - name: find_all_expr_postdec
    doc:  Example query to find all instances of expr_postdec
    flows:
      codelingo/review
	       comments: This is a expr_postdec.
	   query: |
	     import codelingo/ast/php

	     @ review.comment
	     php.expr_postdec
```
</p></details>

<details><summary>php.expr_postinc</summary><p>

#### Example of finding every expr_postinc and having a review flow comment on it:

```yaml
tenets:
  - name: find_all_expr_postinc
    doc:  Example query to find all instances of expr_postinc
    flows:
      codelingo/review
	       comments: This is a expr_postinc.
	   query: |
	     import codelingo/ast/php

	     @ review.comment
	     php.expr_postinc
```
</p></details>

<details><summary>php.expr_predec</summary><p>

#### Example of finding every expr_predec and having a review flow comment on it:

```yaml
tenets:
  - name: find_all_expr_predec
    doc:  Example query to find all instances of expr_predec
    flows:
      codelingo/review
	       comments: This is a expr_predec.
	   query: |
	     import codelingo/ast/php

	     @ review.comment
	     php.expr_predec
```
</p></details>

<details><summary>php.expr_preinc</summary><p>

#### Example of finding every expr_preinc and having a review flow comment on it:

```yaml
tenets:
  - name: find_all_expr_preinc
    doc:  Example query to find all instances of expr_preinc
    flows:
      codelingo/review
	       comments: This is a expr_preinc.
	   query: |
	     import codelingo/ast/php

	     @ review.comment
	     php.expr_preinc
```
</p></details>

<details><summary>php.expr_print</summary><p>

#### Example of finding every expr_print and having a review flow comment on it:

```yaml
tenets:
  - name: find_all_expr_print
    doc:  Example query to find all instances of expr_print
    flows:
      codelingo/review
	       comments: This is a expr_print.
	   query: |
	     import codelingo/ast/php

	     @ review.comment
	     php.expr_print
```
</p></details>

<details><summary>php.expr_propertyfetch</summary><p>

#### Example of finding every expr_propertyfetch and having a review flow comment on it:

```yaml
tenets:
  - name: find_all_expr_propertyfetch
    doc:  Example query to find all instances of expr_propertyfetch
    flows:
      codelingo/review
	       comments: This is a expr_propertyfetch.
	   query: |
	     import codelingo/ast/php

	     @ review.comment
	     php.expr_propertyfetch
```
</p></details>

<details><summary>php.expr_shellexec</summary><p>

#### Example of finding every expr_shellexec and having a review flow comment on it:

```yaml
tenets:
  - name: find_all_expr_shellexec
    doc:  Example query to find all instances of expr_shellexec
    flows:
      codelingo/review
	       comments: This is a expr_shellexec.
	   query: |
	     import codelingo/ast/php

	     @ review.comment
	     php.expr_shellexec
```
</p></details>

<details><summary>php.expr_staticcall</summary><p>

#### Example of finding every expr_staticcall and having a review flow comment on it:

```yaml
tenets:
  - name: find_all_expr_staticcall
    doc:  Example query to find all instances of expr_staticcall
    flows:
      codelingo/review
	       comments: This is a expr_staticcall.
	   query: |
	     import codelingo/ast/php

	     @ review.comment
	     php.expr_staticcall
```
</p></details>

<details><summary>php.expr_staticpropertyfetch</summary><p>

#### Example of finding every expr_staticpropertyfetch and having a review flow comment on it:

```yaml
tenets:
  - name: find_all_expr_staticpropertyfetch
    doc:  Example query to find all instances of expr_staticpropertyfetch
    flows:
      codelingo/review
	       comments: This is a expr_staticpropertyfetch.
	   query: |
	     import codelingo/ast/php

	     @ review.comment
	     php.expr_staticpropertyfetch
```
</p></details>

<details><summary>php.expr_ternary</summary><p>

#### Example of finding every expr_ternary and having a review flow comment on it:

```yaml
tenets:
  - name: find_all_expr_ternary
    doc:  Example query to find all instances of expr_ternary
    flows:
      codelingo/review
	       comments: This is a expr_ternary.
	   query: |
	     import codelingo/ast/php

	     @ review.comment
	     php.expr_ternary
```
</p></details>

<details><summary>php.expr_unaryminus</summary><p>

#### Example of finding every expr_unaryminus and having a review flow comment on it:

```yaml
tenets:
  - name: find_all_expr_unaryminus
    doc:  Example query to find all instances of expr_unaryminus
    flows:
      codelingo/review
	       comments: This is a expr_unaryminus.
	   query: |
	     import codelingo/ast/php

	     @ review.comment
	     php.expr_unaryminus
```
</p></details>

<details><summary>php.expr_unaryplus</summary><p>

#### Example of finding every expr_unaryplus and having a review flow comment on it:

```yaml
tenets:
  - name: find_all_expr_unaryplus
    doc:  Example query to find all instances of expr_unaryplus
    flows:
      codelingo/review
	       comments: This is a expr_unaryplus.
	   query: |
	     import codelingo/ast/php

	     @ review.comment
	     php.expr_unaryplus
```
</p></details>

<details><summary>php.expr_variable</summary><p>

#### Example of finding every expr_variable and having a review flow comment on it:

```yaml
tenets:
  - name: find_all_expr_variable
    doc:  Example query to find all instances of expr_variable
    flows:
      codelingo/review
	       comments: This is a expr_variable.
	   query: |
	     import codelingo/ast/php

	     @ review.comment
	     php.expr_variable
```
</p></details>

<details><summary>php.expr_yield</summary><p>

#### Example of finding every expr_yield and having a review flow comment on it:

```yaml
tenets:
  - name: find_all_expr_yield
    doc:  Example query to find all instances of expr_yield
    flows:
      codelingo/review
	       comments: This is a expr_yield.
	   query: |
	     import codelingo/ast/php

	     @ review.comment
	     php.expr_yield
```
</p></details>

<details><summary>php.expr_yieldfrom</summary><p>

#### Example of finding every expr_yieldfrom and having a review flow comment on it:

```yaml
tenets:
  - name: find_all_expr_yieldfrom
    doc:  Example query to find all instances of expr_yieldfrom
    flows:
      codelingo/review
	       comments: This is a expr_yieldfrom.
	   query: |
	     import codelingo/ast/php

	     @ review.comment
	     php.expr_yieldfrom
```
</p></details>

<details><summary>php.name</summary><p>

#### Example of finding every name and having a review flow comment on it:

```yaml
tenets:
  - name: find_all_name
    doc:  Example query to find all instances of name
    flows:
      codelingo/review
	       comments: This is a name.
	   query: |
	     import codelingo/ast/php

	     @ review.comment
	     php.name
```
</p></details>

<details><summary>php.param</summary><p>

#### Example of finding every param and having a review flow comment on it:

```yaml
tenets:
  - name: find_all_param
    doc:  Example query to find all instances of param
    flows:
      codelingo/review
	       comments: This is a param.
	   query: |
	     import codelingo/ast/php

	     @ review.comment
	     php.param
```
</p></details>

<details><summary>php.scalar_encapsedstringpart</summary><p>

#### Example of finding every scalar_encapsedstringpart and having a review flow comment on it:

```yaml
tenets:
  - name: find_all_scalar_encapsedstringpart
    doc:  Example query to find all instances of scalar_encapsedstringpart
    flows:
      codelingo/review
	       comments: This is a scalar_encapsedstringpart.
	   query: |
	     import codelingo/ast/php

	     @ review.comment
	     php.scalar_encapsedstringpart
```
</p></details>

<details><summary>php.scalar_string</summary><p>

#### Example of finding every scalar_string and having a review flow comment on it:

```yaml
tenets:
  - name: find_all_scalar_string
    doc:  Example query to find all instances of scalar_string
    flows:
      codelingo/review
	       comments: This is a scalar_string.
	   query: |
	     import codelingo/ast/php

	     @ review.comment
	     php.scalar_string
```
</p></details>

<details><summary>php.stmt_break</summary><p>

#### Example of finding every stmt_break and having a review flow comment on it:

```yaml
tenets:
  - name: find_all_stmt_break
    doc:  Example query to find all instances of stmt_break
    flows:
      codelingo/review
	       comments: This is a stmt_break.
	   query: |
	     import codelingo/ast/php

	     @ review.comment
	     php.stmt_break
```
</p></details>

<details><summary>php.stmt_case</summary><p>

#### Example of finding every stmt_case and having a review flow comment on it:

```yaml
tenets:
  - name: find_all_stmt_case
    doc:  Example query to find all instances of stmt_case
    flows:
      codelingo/review
	       comments: This is a stmt_case.
	   query: |
	     import codelingo/ast/php

	     @ review.comment
	     php.stmt_case
```
</p></details>

<details><summary>php.stmt_catch</summary><p>

#### Example of finding every stmt_catch and having a review flow comment on it:

```yaml
tenets:
  - name: find_all_stmt_catch
    doc:  Example query to find all instances of stmt_catch
    flows:
      codelingo/review
	       comments: This is a stmt_catch.
	   query: |
	     import codelingo/ast/php

	     @ review.comment
	     php.stmt_catch
```
</p></details>

<details><summary>php.stmt_class</summary><p>

#### Example of finding every stmt_class and having a review flow comment on it:

```yaml
tenets:
  - name: find_all_stmt_class
    doc:  Example query to find all instances of stmt_class
    flows:
      codelingo/review
	       comments: This is a stmt_class.
	   query: |
	     import codelingo/ast/php

	     @ review.comment
	     php.stmt_class
```
</p></details>

<details><summary>php.stmt_classconst</summary><p>

#### Example of finding every stmt_classconst and having a review flow comment on it:

```yaml
tenets:
  - name: find_all_stmt_classconst
    doc:  Example query to find all instances of stmt_classconst
    flows:
      codelingo/review
	       comments: This is a stmt_classconst.
	   query: |
	     import codelingo/ast/php

	     @ review.comment
	     php.stmt_classconst
```
</p></details>

<details><summary>php.stmt_classmethod</summary><p>

#### Example of finding every stmt_classmethod and having a review flow comment on it:

```yaml
tenets:
  - name: find_all_stmt_classmethod
    doc:  Example query to find all instances of stmt_classmethod
    flows:
      codelingo/review
	       comments: This is a stmt_classmethod.
	   query: |
	     import codelingo/ast/php

	     @ review.comment
	     php.stmt_classmethod
```
</p></details>

<details><summary>php.stmt_const</summary><p>

#### Example of finding every stmt_const and having a review flow comment on it:

```yaml
tenets:
  - name: find_all_stmt_const
    doc:  Example query to find all instances of stmt_const
    flows:
      codelingo/review
	       comments: This is a stmt_const.
	   query: |
	     import codelingo/ast/php

	     @ review.comment
	     php.stmt_const
```
</p></details>

<details><summary>php.stmt_continue</summary><p>

#### Example of finding every stmt_continue and having a review flow comment on it:

```yaml
tenets:
  - name: find_all_stmt_continue
    doc:  Example query to find all instances of stmt_continue
    flows:
      codelingo/review
	       comments: This is a stmt_continue.
	   query: |
	     import codelingo/ast/php

	     @ review.comment
	     php.stmt_continue
```
</p></details>

<details><summary>php.stmt_declare</summary><p>

#### Example of finding every stmt_declare and having a review flow comment on it:

```yaml
tenets:
  - name: find_all_stmt_declare
    doc:  Example query to find all instances of stmt_declare
    flows:
      codelingo/review
	       comments: This is a stmt_declare.
	   query: |
	     import codelingo/ast/php

	     @ review.comment
	     php.stmt_declare
```
</p></details>

<details><summary>php.stmt_declaredeclare</summary><p>

#### Example of finding every stmt_declaredeclare and having a review flow comment on it:

```yaml
tenets:
  - name: find_all_stmt_declaredeclare
    doc:  Example query to find all instances of stmt_declaredeclare
    flows:
      codelingo/review
	       comments: This is a stmt_declaredeclare.
	   query: |
	     import codelingo/ast/php

	     @ review.comment
	     php.stmt_declaredeclare
```
</p></details>

<details><summary>php.stmt_do</summary><p>

#### Example of finding every stmt_do and having a review flow comment on it:

```yaml
tenets:
  - name: find_all_stmt_do
    doc:  Example query to find all instances of stmt_do
    flows:
      codelingo/review
	       comments: This is a stmt_do.
	   query: |
	     import codelingo/ast/php

	     @ review.comment
	     php.stmt_do
```
</p></details>

<details><summary>php.stmt_echo</summary><p>

#### Example of finding every stmt_echo and having a review flow comment on it:

```yaml
tenets:
  - name: find_all_stmt_echo
    doc:  Example query to find all instances of stmt_echo
    flows:
      codelingo/review
	       comments: This is a stmt_echo.
	   query: |
	     import codelingo/ast/php

	     @ review.comment
	     php.stmt_echo
```
</p></details>

<details><summary>php.stmt_else</summary><p>

#### Example of finding every stmt_else and having a review flow comment on it:

```yaml
tenets:
  - name: find_all_stmt_else
    doc:  Example query to find all instances of stmt_else
    flows:
      codelingo/review
	       comments: This is a stmt_else.
	   query: |
	     import codelingo/ast/php

	     @ review.comment
	     php.stmt_else
```
</p></details>

<details><summary>php.stmt_elseif</summary><p>

#### Example of finding every stmt_elseif and having a review flow comment on it:

```yaml
tenets:
  - name: find_all_stmt_elseif
    doc:  Example query to find all instances of stmt_elseif
    flows:
      codelingo/review
	       comments: This is a stmt_elseif.
	   query: |
	     import codelingo/ast/php

	     @ review.comment
	     php.stmt_elseif
```
</p></details>

<details><summary>php.stmt_finally</summary><p>

#### Example of finding every stmt_finally and having a review flow comment on it:

```yaml
tenets:
  - name: find_all_stmt_finally
    doc:  Example query to find all instances of stmt_finally
    flows:
      codelingo/review
	       comments: This is a stmt_finally.
	   query: |
	     import codelingo/ast/php

	     @ review.comment
	     php.stmt_finally
```
</p></details>

<details><summary>php.stmt_foreach</summary><p>

#### Example of finding every stmt_foreach and having a review flow comment on it:

```yaml
tenets:
  - name: find_all_stmt_foreach
    doc:  Example query to find all instances of stmt_foreach
    flows:
      codelingo/review
	       comments: This is a stmt_foreach.
	   query: |
	     import codelingo/ast/php

	     @ review.comment
	     php.stmt_foreach
```
</p></details>

<details><summary>php.stmt_function</summary><p>

#### Example of finding every stmt_function and having a review flow comment on it:

```yaml
tenets:
  - name: find_all_stmt_function
    doc:  Example query to find all instances of stmt_function
    flows:
      codelingo/review
	       comments: This is a stmt_function.
	   query: |
	     import codelingo/ast/php

	     @ review.comment
	     php.stmt_function
```
</p></details>

<details><summary>php.stmt_global</summary><p>

#### Example of finding every stmt_global and having a review flow comment on it:

```yaml
tenets:
  - name: find_all_stmt_global
    doc:  Example query to find all instances of stmt_global
    flows:
      codelingo/review
	       comments: This is a stmt_global.
	   query: |
	     import codelingo/ast/php

	     @ review.comment
	     php.stmt_global
```
</p></details>

<details><summary>php.stmt_goto</summary><p>

#### Example of finding every stmt_goto and having a review flow comment on it:

```yaml
tenets:
  - name: find_all_stmt_goto
    doc:  Example query to find all instances of stmt_goto
    flows:
      codelingo/review
	       comments: This is a stmt_goto.
	   query: |
	     import codelingo/ast/php

	     @ review.comment
	     php.stmt_goto
```
</p></details>

<details><summary>php.stmt_groupuse</summary><p>

#### Example of finding every stmt_groupuse and having a review flow comment on it:

```yaml
tenets:
  - name: find_all_stmt_groupuse
    doc:  Example query to find all instances of stmt_groupuse
    flows:
      codelingo/review
	       comments: This is a stmt_groupuse.
	   query: |
	     import codelingo/ast/php

	     @ review.comment
	     php.stmt_groupuse
```
</p></details>

<details><summary>php.stmt_haltcompiler</summary><p>

#### Example of finding every stmt_haltcompiler and having a review flow comment on it:

```yaml
tenets:
  - name: find_all_stmt_haltcompiler
    doc:  Example query to find all instances of stmt_haltcompiler
    flows:
      codelingo/review
	       comments: This is a stmt_haltcompiler.
	   query: |
	     import codelingo/ast/php

	     @ review.comment
	     php.stmt_haltcompiler
```
</p></details>

<details><summary>php.stmt_if</summary><p>

#### Example of finding every stmt_if and having a review flow comment on it:

```yaml
tenets:
  - name: find_all_stmt_if
    doc:  Example query to find all instances of stmt_if
    flows:
      codelingo/review
	       comments: This is a stmt_if.
	   query: |
	     import codelingo/ast/php

	     @ review.comment
	     php.stmt_if
```
</p></details>

<details><summary>php.stmt_inlinehtml</summary><p>

#### Example of finding every stmt_inlinehtml and having a review flow comment on it:

```yaml
tenets:
  - name: find_all_stmt_inlinehtml
    doc:  Example query to find all instances of stmt_inlinehtml
    flows:
      codelingo/review
	       comments: This is a stmt_inlinehtml.
	   query: |
	     import codelingo/ast/php

	     @ review.comment
	     php.stmt_inlinehtml
```
</p></details>

<details><summary>php.stmt_interface</summary><p>

#### Example of finding every stmt_interface and having a review flow comment on it:

```yaml
tenets:
  - name: find_all_stmt_interface
    doc:  Example query to find all instances of stmt_interface
    flows:
      codelingo/review
	       comments: This is a stmt_interface.
	   query: |
	     import codelingo/ast/php

	     @ review.comment
	     php.stmt_interface
```
</p></details>

<details><summary>php.stmt_namespace</summary><p>

#### Example of finding every stmt_namespace and having a review flow comment on it:

```yaml
tenets:
  - name: find_all_stmt_namespace
    doc:  Example query to find all instances of stmt_namespace
    flows:
      codelingo/review
	       comments: This is a stmt_namespace.
	   query: |
	     import codelingo/ast/php

	     @ review.comment
	     php.stmt_namespace
```
</p></details>

<details><summary>php.stmt_property</summary><p>

#### Example of finding every stmt_property and having a review flow comment on it:

```yaml
tenets:
  - name: find_all_stmt_property
    doc:  Example query to find all instances of stmt_property
    flows:
      codelingo/review
	       comments: This is a stmt_property.
	   query: |
	     import codelingo/ast/php

	     @ review.comment
	     php.stmt_property
```
</p></details>

<details><summary>php.stmt_propertyproperty</summary><p>

#### Example of finding every stmt_propertyproperty and having a review flow comment on it:

```yaml
tenets:
  - name: find_all_stmt_propertyproperty
    doc:  Example query to find all instances of stmt_propertyproperty
    flows:
      codelingo/review
	       comments: This is a stmt_propertyproperty.
	   query: |
	     import codelingo/ast/php

	     @ review.comment
	     php.stmt_propertyproperty
```
</p></details>

<details><summary>php.stmt_return</summary><p>

#### Example of finding every stmt_return and having a review flow comment on it:

```yaml
tenets:
  - name: find_all_stmt_return
    doc:  Example query to find all instances of stmt_return
    flows:
      codelingo/review
	       comments: This is a stmt_return.
	   query: |
	     import codelingo/ast/php

	     @ review.comment
	     php.stmt_return
```
</p></details>

<details><summary>php.stmt_staticvar</summary><p>

#### Example of finding every stmt_staticvar and having a review flow comment on it:

```yaml
tenets:
  - name: find_all_stmt_staticvar
    doc:  Example query to find all instances of stmt_staticvar
    flows:
      codelingo/review
	       comments: This is a stmt_staticvar.
	   query: |
	     import codelingo/ast/php

	     @ review.comment
	     php.stmt_staticvar
```
</p></details>

<details><summary>php.stmt_switch</summary><p>

#### Example of finding every stmt_switch and having a review flow comment on it:

```yaml
tenets:
  - name: find_all_stmt_switch
    doc:  Example query to find all instances of stmt_switch
    flows:
      codelingo/review
	       comments: This is a stmt_switch.
	   query: |
	     import codelingo/ast/php

	     @ review.comment
	     php.stmt_switch
```
</p></details>

<details><summary>php.stmt_throw</summary><p>

#### Example of finding every stmt_throw and having a review flow comment on it:

```yaml
tenets:
  - name: find_all_stmt_throw
    doc:  Example query to find all instances of stmt_throw
    flows:
      codelingo/review
	       comments: This is a stmt_throw.
	   query: |
	     import codelingo/ast/php

	     @ review.comment
	     php.stmt_throw
```
</p></details>

<details><summary>php.stmt_trait</summary><p>

#### Example of finding every stmt_trait and having a review flow comment on it:

```yaml
tenets:
  - name: find_all_stmt_trait
    doc:  Example query to find all instances of stmt_trait
    flows:
      codelingo/review
	       comments: This is a stmt_trait.
	   query: |
	     import codelingo/ast/php

	     @ review.comment
	     php.stmt_trait
```
</p></details>

<details><summary>php.stmt_traituse</summary><p>

#### Example of finding every stmt_traituse and having a review flow comment on it:

```yaml
tenets:
  - name: find_all_stmt_traituse
    doc:  Example query to find all instances of stmt_traituse
    flows:
      codelingo/review
	       comments: This is a stmt_traituse.
	   query: |
	     import codelingo/ast/php

	     @ review.comment
	     php.stmt_traituse
```
</p></details>

<details><summary>php.stmt_traituseadaptation_alias</summary><p>

#### Example of finding every stmt_traituseadaptation_alias and having a review flow comment on it:

```yaml
tenets:
  - name: find_all_stmt_traituseadaptation_alias
    doc:  Example query to find all instances of stmt_traituseadaptation_alias
    flows:
      codelingo/review
	       comments: This is a stmt_traituseadaptation_alias.
	   query: |
	     import codelingo/ast/php

	     @ review.comment
	     php.stmt_traituseadaptation_alias
```
</p></details>

<details><summary>php.stmt_traituseadaptation_precedence</summary><p>

#### Example of finding every stmt_traituseadaptation_precedence and having a review flow comment on it:

```yaml
tenets:
  - name: find_all_stmt_traituseadaptation_precedence
    doc:  Example query to find all instances of stmt_traituseadaptation_precedence
    flows:
      codelingo/review
	       comments: This is a stmt_traituseadaptation_precedence.
	   query: |
	     import codelingo/ast/php

	     @ review.comment
	     php.stmt_traituseadaptation_precedence
```
</p></details>

<details><summary>php.stmt_trycatch</summary><p>

#### Example of finding every stmt_trycatch and having a review flow comment on it:

```yaml
tenets:
  - name: find_all_stmt_trycatch
    doc:  Example query to find all instances of stmt_trycatch
    flows:
      codelingo/review
	       comments: This is a stmt_trycatch.
	   query: |
	     import codelingo/ast/php

	     @ review.comment
	     php.stmt_trycatch
```
</p></details>

<details><summary>php.stmt_unset</summary><p>

#### Example of finding every stmt_unset and having a review flow comment on it:

```yaml
tenets:
  - name: find_all_stmt_unset
    doc:  Example query to find all instances of stmt_unset
    flows:
      codelingo/review
	       comments: This is a stmt_unset.
	   query: |
	     import codelingo/ast/php

	     @ review.comment
	     php.stmt_unset
```
</p></details>

<details><summary>php.stmt_use</summary><p>

#### Example of finding every stmt_use and having a review flow comment on it:

```yaml
tenets:
  - name: find_all_stmt_use
    doc:  Example query to find all instances of stmt_use
    flows:
      codelingo/review
	       comments: This is a stmt_use.
	   query: |
	     import codelingo/ast/php

	     @ review.comment
	     php.stmt_use
```
</p></details>

<details><summary>php.stmt_useuse</summary><p>

#### Example of finding every stmt_useuse and having a review flow comment on it:

```yaml
tenets:
  - name: find_all_stmt_useuse
    doc:  Example query to find all instances of stmt_useuse
    flows:
      codelingo/review
	       comments: This is a stmt_useuse.
	   query: |
	     import codelingo/ast/php

	     @ review.comment
	     php.stmt_useuse
```
</p></details>

<details><summary>php.stmt_while</summary><p>

#### Example of finding every stmt_while and having a review flow comment on it:

```yaml
tenets:
  - name: find_all_stmt_while
    doc:  Example query to find all instances of stmt_while
    flows:
      codelingo/review
	       comments: This is a stmt_while.
	   query: |
	     import codelingo/ast/php

	     @ review.comment
	     php.stmt_while
```
</p></details>

