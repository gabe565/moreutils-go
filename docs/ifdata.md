## ifdata

Get network interface info without parsing ifconfig output

```
ifdata [flags] interface
```

### Options

```
  -h      help for ifdata
  -e      Test to see if the interface exists, exit nonzero if it does not
  -p      Prints out the whole configuration of the interface
  -pe     Prints "yes" or "no" if the interface exists or not
  -pa     Prints the IP address of the interface
  -pN     Prints the network address of the interface
  -pn     Prints the netmask of the interface
  -pb     Prints the broadcast address of the interface
  -pm     Prints the MTU of the interface
  -pf     Prints the flags of the interface
  -ph     Prints the hardware address of the interface. Exit with a failure exit code if there is not hardware address for the given network interface
  -si     Prints all input statistics of the interface
  -sip    Prints the number of input packets
  -sib    Prints the number of input bytes
  -sie    Prints the number of input errors
  -sid    Prints the number of dropped input packets
  -sif    Prints the number of input fifo overruns
  -sic    Prints the number of compressed input packets
  -sim    Prints the number of input multicast packets
  -so     Prints all output statistics of the interface
  -sop    Prints the number of output packets
  -sob    Prints the number of output bytes
  -soe    Prints the number of output errors
  -sod    Prints the number of dropped output packets
  -sof    Prints the number of output fifo overruns
  -sox    Prints the number of output collisions
  -soc    Prints the number of output carrier losses
  -som    Prints the number of output multicast packets
  -bips   Prints the number of bytes of incoming traffic measured in one second
  -bops   Prints the number of bytes of outgoing traffic measured in one second
  -v      version for ifdata
```

### SEE ALSO

* [moreutils](moreutils.md)	 - A collection of the Unix tools that nobody thought to write long ago when Unix was young

