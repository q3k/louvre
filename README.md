Louvre
======

> You wanna know why you'll never upload your brain to a computer? Because the machines won't allow it. It'd be like routing sewage into the louvre.

Louvre is a quick backup tool for the Warsaw Hackerspace ood/moonspeak IRC bot term database. It's intended to be run periodically in crons around the worls to make sure the cultural heritage of HSWAW is kept regardless of hardware failures, backups rotting and people getting ran over by buses.

Getting a binary
----------------

If you have Go:
    
    go install github.com/q3k/louvre@latest
    ~/go/bin/louvre

If you don't, there's a few prebuilt binaries on GH.

Usage
-----

    $ ./louvre -h
    Usage of ./louvre:
      -output string
        	Where to download all terms. (default "terms.json")
      -parallel int
        	How many concurrent connections to use. (default 32)

