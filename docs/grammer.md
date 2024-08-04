# Grammer

## BNF

```
<query>         ::= <expr>

<expr>         ::= <segment>
                 | <expr> "/" <segment>

<segment>      ::= <class>
                 | <class> "{" <predicate> "}"
                 | <class> "{" <predicate> "}" "[" <idx> "]"
                 | <block>
                 | <block> "{" <predicate> "}"

<block>        ::= <class> "[" <idx> "]"
                 | <class> ":" <label>

<class>        ::= <ident>

<label>        ::= <ident>

<predicate>    ::= <attr>
                 | <attr> "=" <value>

<attr>         ::= "@"<ident>

<ident>        ::= "[a-zA-Z]+[a-zA-Z0-9_-]*"

<idx>          ::= "[0-9]+"
```

## Operators
###NESTS (`/`)
`op_1 / op_2`: `op_2` is nested in `op_1`
###FILTERED (`{}`)
`op_1{op_2}`: `op_1` is filtered by `op_2`
###SELECT (`[]`)
`op_1[op_2]`: select the item with index `op_2` from `op_1`
###NAMED (`:`)
`op_1:op_2`: `op_1` is named `op_2`
###EQUAL (`=`)
`op_1=op_2`

### Precedence
All operator are of equal precedence.

### Associativity
- `/`, `:`, `[]` are left-associative.
- `{}`, `=` is right-associative.
