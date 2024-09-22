## ifdata

Get network interface info without parsing ifconfig output

```
ifdata interface [flags]
```

### Options

```
  -a, --address             Prints the IPv4 address of the interface
  -b, --broadcast-address   Prints the broadcast address of the interface
  -e, --exists              Test to see if the interface exists, exit nonzero if it does not
  -f, --flags               Prints the flags of the interface
  -h, --hardware-addr       Prints the hardware address of the interface. Exit with a failure exit code if there is not hardware address for the given network interface
  -m, --mtu                 Prints the MTU of the interface
  -n, --netmask             Prints the netmask of the interface
  -N, --network-address     Prints the network address of the interface
  -p, --print               Prints out the whole configuration of the interface
  -v, --version             version for ifdata
```

### SEE ALSO

* [moreutils](moreutils.md)	 - A collection of the Unix tools that nobody thought to write long ago when Unix was young

