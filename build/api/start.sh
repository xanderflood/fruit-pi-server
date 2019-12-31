#!/bin/sh
set +x

export PGPASSWORD=`cat $PGPASSWORD_FILE`
export JWT_SIGNING_SECRET=`cat $JWT_SIGNING_SECRET_FILE`

./api
