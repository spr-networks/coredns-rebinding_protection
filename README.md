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

Test with querying a domain externally that receives an internal IP, for example:
host 0.0.0.0.nip.io

This should be configured in plugin.cfg to only affect the 'forward' plugin. The 'hosts' plugin should bypass this feature
