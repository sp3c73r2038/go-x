TESTFLAGS :=

PASSWORD :=
KEY :=

-include local.mk

export PASSWORD
export KEY

test:
	go test ${TESTFLAGS} ./...
