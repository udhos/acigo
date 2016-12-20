#!/bin/sh

pkg=github.com/udhos/acigo

step=0

msg() {
    step=$((step+1))
    echo >&2 $step. $*
}

get() {
    i=$1
    msg fetching $i
    go get $i
    msg fetching $i - done
}

get github.com/udhos/equalfile
#get honnef.co/go/simple/cmd/gosimple

src=`find . -type f | egrep '\.go$'`

msg fmt
gofmt -s -w $src
msg fix
go tool fix $src
msg vet
go tool vet .

msg install
go install $pkg/aci
go install $pkg/samples/aci-login
go install $pkg/samples/aci-tls

# go get honnef.co/go/simple/cmd/gosimple
s=$GOPATH/bin/gosimple
simple() {
    msg simple - this is slow, please standby
    # gosimple cant handle source files from multiple packages
    $s aci/*.go
    $s samples/login/*.go
}
[ -x "$s" ] && simple

# go get github.com/golang/lint/golint
l=$GOPATH/bin/golint
lint() {
    msg lint
    # golint cant handle source files from multiple packages
    $l aci/*.go
    $l samples/login/*.go
}
[ -x "$l" ] && lint

msg test aci
go test github.com/udhos/acigo/aci
