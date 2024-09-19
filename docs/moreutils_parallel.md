## moreutils parallel

Run multiple jobs at once

```
moreutils parallel [flags] command -- arg...
```

### Options

```
  -h, --help           help for parallel
  -j, --jobs string    Number of jobs to run in parallel. Can be a number or a percentage of CPU cores. (default "10")
  -n, --num-args int   Number of arguments to pass to a command at a time. Default is 1. Incompatible with -i (default 1)
  -i, --replace        Normally the argument is added to the end of the command. With this option, instances of "{}" in the command are replaced with the argument.
```

### SEE ALSO

* [moreutils](moreutils.md)	 - A collection of the Unix tools that nobody thought to write long ago when Unix was young

