#!/bin/bash

trap stop INT

function stop() {
    pkill xkcd-server
    pkill web-server

    echo
    echo "Servers stopped. Good bye!"
    exit 0
}


./xkcd-server &
./web-server &

sleep 3

echo "You want to stop the servers? (Y/y)"
read input

if [ "$input" == "Y" ] || [ "$input" == "y" ]; then
    echo "Stopping servers..."
    stop
fi