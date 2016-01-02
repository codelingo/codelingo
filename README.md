# Lingo: Scaling Code Quality

## House keeping

See [here](https://medium.com/@jessemeek/not-learning-the-lingo-how-a-bad-performance-review-gave-birth-to-a-new-start-up-45e36dd997b9#.shv55pite) for the story behind Lingo. Please add an issue for any bugs or feature requests.

## The Lingo Tool

### Install Lingo

Lingo is a CLI tool written in Go. To build the binary from source:

```bash
go get github.com/lingo-reviews/lingo
cd $GOPATH/src/github.com/lingo-reviews/lingo
go install
```

Alternatively, the latest version of the binary can be downloaded from lingo.reviews:

```bash
wget http://lingo.reviews/lingo.zip && unzip lingo.zip
```

Then place Lingo in your PATH:

```bash
sudo cp lingo /usr/local/bin/
```

Note: the binary is for Linux. OSX and Windows are in the pipeline.

Run `lingo` to see the avaliable commands.

#### Enable Bash Auto-Complete

Run `lingo setup-auto-completion` to enable command auto-completion.

Or, if you have checked out the source code and prefer to do it by hand:

```bash
echo 'PROG=lingo source $GOPATH/src/github.com/lingo-reviews/lingo/scripts/bash_autocomplete.sh' >> ~/.bashrc
. ~/.bashrc
lingo --generate-bash-completion
```

## Tenets

Lingo uses tenets to manage the quality of code. [github.com/lingo-reviews/tenets](github.com/lingo-reviews/tenets) contains some tenets and example code for each language to get you started. The following examples will use Go tenets. 

```bash
go get github.com/lingo-reviews/tenets/go
```

## First Run

### Docker Quick Start

If you're installing Docker on Ubuntu, read our [troubleshooting](https://github.com/lingo-reviews/tenets/wiki/Troubleshooting) page first.

With Docker installed:

```bash
# Find some source code to review.
cd $GOPATH/src/github.com/lingo-reviews/tenets/go/tenets/license/tenet/example

# Review the code.
lingo review

# Read the tenet documentation for this code.
lingo docs

```

When lingo reviews, it looks for a .lingo file in the current or parent
directories. If those tenets use a docker driver (default) and no local docker
image is found, lingo goes and gets it. The first time you pull a docker
tenet, it will pull the tenet base images. This means future tenet pulls will
be much quicker.


Next, start without a .lingo file:

```bash
cd $GOPATH/src/github.com/lingo-reviews/tenets/go/tenets/simpleseed/example

# This will write a .lingo file.
lingo init

# List avaliable tenets on hub.docker.com:
docker search lingoreviews

# Add the simpleseed example:
lingo add lingoreviews/simpleseed

# Pull down the images from hub.docker.com:
lingo pull --all

# If you didn't pull, review will do it for you.

# Review the code, this time we'll keep some output at the end:
lingo review --output-format --json-pretty

```

Notes: Tenets can be pulled from any docker repository. A better tenet search
UI is in the pipeline.

Lingo will prompt you to open each issue. Supported editors are: vi, vim,
emacs, nano and subl. To skip the confirm steps, use --keep-all.

### Binary Quick Start

All the other example folders under go/tenets use the binary driver. We will
use the `lingo build` command to build them all at once. `cd` into go/tenets and run:

```bash
lingo build binary --all
```

You'll see the following output:

```bash
$ lingo build binary --all
Building Go binary: [~/.lingo_home/tenets/lingoreviews/juju_nosingle]
Building Go binary: [~/.lingo_home/tenets/lingoreviews/imports]
...
Building Go binary: [~/.lingo_home/tenets/lingoreviews/unused_arg]
Building Go binary: [~/.lingo_home/tenets/lingoreviews/juju_worker_periodic]
binary 17 / 17 [========================================================] 100.00 % 12s
Success! All binary tenets built.
```

Lingo builds and installs each binary. Commands such as `add` and `info` will
now auto-complete with the names of the built binary tenets.

`cd` into any example folder and run `lingo review`. In a similar fashion, you
can `lingo build docker --all` to build local copies of all the docker
tenets.To add the binary drivers, we need to specify the driver when we add
it:

```bash
lingo add lingoreviews/simpleseed --driver binary
```

Otherwise, the driver will default to "docker". By default, binary tenets are
installed in ~/.lingo_home/tenets/[owner]/[name]. This can be overridden with
the LINGO_BIN environment variable.

## Options

Some tenets take options. To view their available options run:

```bash
lingo info <tenet-name>
```

The imports tenet, for example, takes a blacklist_regex option, here's an
example of setting it:

```bash
lingo add lingoreviews/imports --options blacklist_regex=".*/State"
```

## Writing a Tenet

### Go

Start [here](https://github.com/lingo-reviews/tenets/tree/master/go/dev). The
`go/tenets` directory also has a variety of examples of tenets in Go which you
can copy to get started.

### Other languages

The api.proto file in the root of this repoistory can be used to generate the
tenet API libs in C, C++, Java, Go, Node.js, Python, Ruby, Objective-C, PHP
and C#. Visit grpc.io to learn more.
