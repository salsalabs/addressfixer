#!/bin/bash
export PATH=/home/ubuntu/go/bin:$PATH
export LOGIN=/home/ubuntu/workspace/go/src/github.com/salsalabs/addressfixer/logins/allen_ewg.yaml
export DBLOGIN=/home/ubuntu/workspace/go/src/github.com/salsalabs/addressfixer/logins/db.yaml
addressfixer --login $LOGIN --dblogin $DBLOGIN

