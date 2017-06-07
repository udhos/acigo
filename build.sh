#! /bin/bash

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

get github.com/gorilla/websocket
#get honnef.co/go/simple/cmd/gosimple

src=`find . -type f | egrep '\.go$'`

msg fmt
gofmt -s -w $src
msg fix
go tool fix $src
msg vet
go tool vet .

if [ -z "$GOPATH" ]; then
	src=$HOME/go/src/$pkg
else
	src=$GOPATH/src/$pkg
fi

pushd $src >/dev/null
samples=`echo samples/*`
popd >/dev/null

#echo samples: $samples

msg install
go install $pkg/aci
for i in $samples; do
    msg install $pkg/$i
    go install $pkg/$i
done

# go get github.com/golang/lint/golint
l=$GOPATH/bin/golint
lint() {
    msg lint
    # golint cant handle source files from multiple packages
    pushd $GOPATH/src/$pkg >/dev/null
    $l yname/*.go
    $l aci/*.go
    for i in $samples; do
	msg lint $i
	$l $i/*.go
    done
    popd >/dev/null
}
[ -x "$l" ] && lint

# go get honnef.co/go/simple/cmd/gosimple
s=$GOPATH/bin/gosimple
simple() {
    msg simple - this is slow, please standby
    # gosimple cant handle source files from multiple packages
    pushd $GOPATH/src/$pkg >/dev/null
    $s yname/*.go
    $s aci/*.go
    for i in $samples; do
	msg simple $i
	$s $i/*.go
    done
    popd >/dev/null
}
[ -x "$s" ] && simple

# go get honnef.co/go/tools/cmd/staticcheck
sc=$GOPATH/bin/staticcheck
static() {
    msg staticcheck - this is slow, please standby
    # staticcheck cant handle source files from multiple packages
    pushd $GOPATH/src/$pkg >/dev/null
    $sc yname/*.go
    $sc aci/*.go
    for i in $samples; do
	msg static $i
	$sc $i/*.go
    done
    popd >/dev/null
}
[ -x "$sc" ] && static

msg test aci
go test github.com/udhos/acigo/aci
go test github.com/udhos/acigo/yname
