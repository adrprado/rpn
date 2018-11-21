SOURCEDIR=.
SOURCES := $(shell find $(SOURCEDIR) -name '*.go')

BINARY=rapina
WINBINARY=~/windows/temp/rapina.exe

VERSION=`git describe --tags`
BUILD_TIME=`date +%F`

# Setup the -ldflags option for go build here, interpolate the variable values
LDFLAGS=-ldflags "-w -s -X main.version=${VERSION} -X main.build=${BUILD_TIME}"
#LDFLAGS=-ldflags "-X github.com/ariejan/roll/core.Version=${VERSION} -X github.com/ariejan/roll/core.BuildTime=${BUILD_TIME}"

.DEFAULT_GOAL: $(BINARY)

$(BINARY): $(SOURCES) $(wildcard ../*.go)
	go build ${LDFLAGS} -o $(BINARY)

win: $(SOURCES)
	GOOS=windows GOARCH=386 go build ${LDFLAGS} -o ${WINBINARY}

.PHONY: install
install:
	go install ${LDFLAGS} ./...

.PHONY: clean
clean:
	if [ -f ${BINARY} ] ; then rm ${BINARY} ; fi

