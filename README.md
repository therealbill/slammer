# Overview

What is "slammer"? Well it started as a little utility I needed for trying to
reliably induce a high rate of data change in Redis as well as induce high
pubsub activity. Why might I need that, you ask? In a word:
client-output-buffer-limit. In a link:  https://redislabs.com/blog/the-endless-redis-replication-loop-what-why-and-how-to-solve-it


# Components

Currently there are two components: firetuck and firehose.


# Firetruck

This tool makes it easy to concurrently upload a metric ton of data into a
Redis instance. You can control the number of concurrent clients, the number of
keys eah with `SET`, and the number of characters (bytes) in each key's value.
It doesn't do benchmarking but if run local to the Redis server it can *really*
induce a high write load.

By varying the size versus key count you can produce different patterns of
usage. You can simulate high amounts of small keys or "only" hundreds or
thousands of "large" keys. The kind of load this causes is distinctly different
so I have found having both to be quite helpful.


One way to use it is to load the data and let it complete. Next you re-run the
command with a different size for the key. This will cause Redis to *change*
the value of each key rather than create new ones.

Finally you can pass a prefix option. That way you can run multiple invocations
of it and not have them stomp on each other. In this way you can do both new
key addition and key modification concurrently.

For detailed (and current) command line usage run 'firetruck -h'.


# Firehose

What Firetruck does for `SET`, Firehose does for PUB. Useful for simulating
high amounts of messages being published. Does not implement SUBSCRIBING to
those messages.

For detailed (and current) command line usage run 'firehose -h'.

# Releases
This project is now set up with automatic travis-ci.org Github release. Anytime
I bump the version travis-ci.org will build binaries for Linux, OSX, and
Windows and push them to a Github release. To automatically get the latest full release use this handy one-liner:
`wget $(curl -s https://api.github.com/repos/therealbill/slammer/releases/latest | grep 'browser_' | cut -d\" -f4)`

To get just Linux binaries:
`wget $(curl -s https://api.github.com/repos/therealbill/slammer/releases/latest | grep 'browser_' | cut -d\" -f4| egrep -v '(exe|darwin)'))`

To get just Windows binaries:
`wget $(curl -s https://api.github.com/repos/therealbill/slammer/releases/latest | grep 'browser_' | cut -d\" -f4| grep exe)`

To get just OSX binaries:
`wget $(curl -s https://api.github.com/repos/therealbill/slammer/releases/latest | grep 'browser_' | cut -d\" -f4| grep darwin)`

# The Future?

Magic Eight Ball says "Uncertain". I can think of several mre thigs the tool
does but for now it serves the need. If you've got something you'd like to see
added, feel free to submit an issue on Github or, even better, an issue and a
PR impementing it.




