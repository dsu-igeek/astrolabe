#!/bin/sh
/usr/bin/pg_dumpall -U postgres -h $PGHOST
