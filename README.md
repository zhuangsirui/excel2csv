## Usage

```
NAME:
   excel2csv - convert excel each sheets to a single csv

USAGE:
   excel2csv [--output DIR] [--trim] [--trim-float] [--with-bom] file [file...]

VERSION:
   0.0.7

COMMANDS:
     help, h  Shows a list of commands or help for one command

GLOBAL OPTIONS:
   --output value, -o value  target directory for output csv (default: ".")
   --trim                    trim value
   --trim-float              try to parse string like 1.10000000000001 to 1.1
   --convert-bool            covert 0 1 to "true" "false"
   --with-bom                add UTF-8 BOM to csv file
   --help, -h                show help
   --version, -v             print the version
```

