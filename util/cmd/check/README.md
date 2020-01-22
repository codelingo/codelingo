Rule testing utility
====================

Accepts any number of non-flag arguments; if none are supplied, acts as though
called with "." alone.

If no flag is set, tests every supplied rule dir.

If the "--search" flag is set, scans each supplied root dir for rule dirs and
tests all rules found.

Examples
--------

    $ go run ./main.go some/rule/dir another/rule/dir

    $ go run ./main.go --search some/root/dir another/root/dir
