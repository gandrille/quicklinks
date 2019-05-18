# Quick links

Small tool for executing commands. Can be used as a convinent way to manage aliases.


## Usage (explained using a full example)

Syntax : `quicklinks file [choices...]`

Here is the content of `links.txt`

```
# Runs an executable file
# It's possible to specify arguments
# IMPORTANT : the line "arg = foo bar" 
# means a SINGLE argument "foo bar"

[prog dmesg]
exe = /bin/dmesg
arg = --level=err,warn
arg = --ctime


# Shortcut for using Unison File Synchronizer
# see : https://www.cis.upenn.edu/~bcpierce/unison/
# In the following example, the profile will be loaded from
# ~/.unison/my-backup.prf

[unison sync]
profile = my-backup


# Interface for using rsync
# In the following example
# host = tux              --> Does NOT execute the command if the hostname is NOT tux
# cmp = /mnt/foo          --> Check Mount Point. Does NOT execute the command if /mnt/foo is NOT mounted
# flag = -va              --> Flag(s) provided as a single argument on the command line
# flag = --delete         --> Idem
# exclude = *~            --> Shortcut. Same as "flag = --exclude" and "flag = *~"
# srcPrefix = /home/dave/ --> | Specify the sources. Same as 
# src = docs              --> | src = /home/dave/docs
# src = libs              --> | src = /home/dave/libs
# dst = /mnt/foo          --> Destination (MUST be defined EXACTLY once)

[rsync push]
host = tux
cmp = /mnt/foo
flag = -va
exclude = *~
srcPrefix = /home/dave/
src = docs
src = libs
dst = /mnt/foo
```

Running `quicklinks links.txt` produces the following output :

```
1 dmesg
2 sync
3 push
Please enter at least one number or key:
```

Typing `1` or `dmesg` will execute `/bin/dmesg --level=err,warn --ctime`.
It's as well possible to type `1 2` in order to launch the first two commands.
Typing `all` executes all the command defined in the file.

Instead of answering the menu, the choices can be directly provided on the command line. For example `quicklinks links.txt dmesg`, `quicklinks links.txt 1 2`, `quicklinks links.txt all`


## Build

```
go get github.com/gandrille/quicklinks/...
go install src/github.com/gandrille/quicklinks/quicklinks.go
```


## Changelog

**[v1.0](../../releases/tag/v1.0)** initial release


## TODO

Choose a license complient with third parties.
