# Rebinding Protection


## Name

*rebinding_protection* - blocks upstream resolvers from providing local or private IP addresses

## Syntax

~~~ txt
rebinding_protection
~~~

## Examples

Block 0.0.0.0, private addresess, etc

``` corefile
. {
  forward . 9.9.9.9
  rebinding_protection
}
```
