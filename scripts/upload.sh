#!/usr/bin/env bash

BINDIR=${1:-bin}
PUBBUCKET=pub.harmony.one
AWSCLI=aws
NETWORK=${2:-mainnet}

if [ "$(uname -s)" == "Darwin" ]; then
   FOLDER=release/darwin-x86_64/$NETWORK
else
   FOLDER=release/linux-x86_64/$NETWORK
fi

bin=harmony-tui

if [ -e $BINDIR/$bin ]; then
   $AWSCLI s3 cp $BINDIR/$bin s3://${PUBBUCKET}/$FOLDER/$bin --acl public-read
else
   echo "!! MISSGING $bin !!"
fi
