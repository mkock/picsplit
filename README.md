# picsplit

Command-line utility that splits pictures into sub-directories by EXIF date.

_Made with Go 1.17. Module-support enabled._

# Installation

```bash
go install github.com/mkock/picsplit
```

# Usage

* Nagivate to directory containing photos to distribute
* Run the query below

```bash
picsplit .
```

It will read each photo in the current directory, create a directory using the format YYYY-MM-DD using the
photo timestamp from the EXIF data, and move the photo into that directory.

The utility is currently very simple and will stop upon encountering any errors.
