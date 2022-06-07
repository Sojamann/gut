# GUT
*Gut* it trying to be an easier to use *cut* command.
Often times I defaulted to *awk* instead when cut didn't get the job done.
Altough *awk* is pretty powerfull it still is not the best solution for
simple problems as *awk* is so much more than just a cut replacement.
*Gut* supports negative indexing and complex split patterns.

## Install
```SH
go install github.com/Sojamann/gut@v1.0.0
```

## Usage
The default behavior of *gut* is to cut whenever there is more then
one consecutive white-space character.
<br>
NOTE: A tab is regarded as multi-whitespace as it appears as more than one.

### Help
```
$ gut -h
Print selected parts of lines from each FILE to standard output.

With no FILE, or when FILE is -, read standard input.

    -f      --fields FIELDS                 select only these fields; also print any line
                                                that contains no delimiter character.
                                                Fields can be used more than once
    -cw     --cut-on-whitespace             cut on any whitespace but also trim following
                                                whitespace aswell until first non-whitepsace
    -cmw    --cut-on-multi-whitespace       cut on consecutive whitespace but also trim following
                                                whitespace aswell until first non-whitepsace
    -cs     --cut-on-seperator STR          cut whenever STR is encountered in a line
    -cf     --cut-on-format DELIMS          cut line applying one delimiter in STR after
                                                another
    -fsep   --format-seperator STR          use STR as seperator in the DELIMS specification
                                                Default: ','
    -osep   --ouput-seperator STR           use the STR as the output field seperator
                                                Default: ' '

Only one of the following can be used at a time:
    -cw
    -cmw
    -cs
    -cf

FIELDS is made up of one range, or many ranges separated by commas.
Selected input is written in the same order that it is read.
Each range is one of:

    N     N'th field, counted from 1
   -N     (num items on line) - N
    N:    from N'th field, to end of line
   -N:    (num items on line) - N, to the end of line
    N:M   from N'th to M'th (included) field
    :M    from first to M'th (included) field
    :-M   from first to ((num items on line) - M)'th (included) field
    :     from beginning to end of line

DELIMS is made up on one seperator, or many seperators seperated by commas.
The seperator of DELIMS can be changed using --format-seperator.
Each delimiter can be one of:
    s                   cut on next whitespace
    t                   cut on next tab
    w                   cut on next whitespace but consume all consecutive aswell
    m                   cut on next multi whitespace and consume all consecutive aswell
    <str> cut on next encounter of str, where str can by any string

Note: whitespace always means tab and space.

Examples:
    $ echo "A B C" | gut -cw -f 2:
    B C

    $ echo "A  B  C" | gut -f -2:
    B C

    $ echo "A;;B     C" | gut -cf "<;;>,m"
    A B C

    $ echo "A  B,C" | gut -cf "s|<,>" -fsep "|"
    A  B C
```


## Difference
In the following the usecases of *sep* are highlighted as sep tries to improve on the usability of the *cut* command.

```SH
$ docker image ls
REPOSITORY   TAG                  IMAGE ID       CREATED         SIZE
node         14.18.2-alpine3.15   5f5960be493c   5 months ago    118MB

$ docker image ls | cut -f 3
REPOSITORY   TAG                  IMAGE ID       CREATED         SIZE
node         14.18.2-alpine3.15   5f5960be493c   5 months ago    118MB

$ docker image ls | awk '{ print $3 }'
IMAGE
5f5960be493c

$ docker image ls | gut -f 3
IMAGE ID
5f5960be493c
```

## Examples
## Cut types
### Default / Multi whitespace cutting
```SH
docker image ls | gut -f 3
IMAGE ID
5f5960be493c

$ echo -e "A  B\tC\t\tD E\t     F" | gut -cmw -osep ";"
A;B;C;D E;Fs
```
### Seperator cutting
```SH
$ echo -e "A;B;C" | gut -cs ";"
A B C
```
### Format cutting
```SH
$ echo -e "A,B\tC    DignoreE" | gut -cf "<,>|t|a|<ignore>" -fsep "|" -osep ";"
A;B;C;D;E
```
