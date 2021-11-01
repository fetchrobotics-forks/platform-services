#!/bin/sh
if [ -n "${PGPASSWORD}" ]; then
    echo "${PGHOST}:${PGPORT}:${PGDATABASE}:${PGUSER}:${PGPASSWORD}" >> /home/app/.pgpass
    chmod 0600 /home/app/.pgpass
    export PGPASSWORDFILE="/home/app/.pgpass"
fi
echo "PGPASSWORD"
echo "${PGPASSWORD}"
cat /home/app/.pgpass
./experimentsrv
