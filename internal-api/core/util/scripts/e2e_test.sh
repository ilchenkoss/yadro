#!/bin/sh

#sudo apt-get install jq
#brew install jq

HOST="localhost"
PORT="22222"

sudo ./xkcd-server &

sleep 5

echo "Request token..."
RESPONSE_LOGIN=$(curl -s -X POST \
  -H "Content-Type: application/json" \
  -d '{"login": "humorist", "password": "yqS~1v1vKcuMs~"}' \
  http://$HOST:$PORT/login)
TOKEN=$(echo $RESPONSE_LOGIN | jq -r '.token')
if [ -n "$TOKEN" ]; then
  echo "Token: $TOKEN"
else
  echo "Error from get token"
  exit 1
fi

echo "Request update comics..."
RESPONSE_UPDATE_COMICS=$(curl -s -X POST \
  -H "authorization: bearer $TOKEN" \
  http://$HOST:$PORT/update)
NEW_COMICS=$(echo $RESPONSE_UPDATE_COMICS | jq -r '.new_comics')
TOTAL_COMICS=$(echo $RESPONSE_UPDATE_COMICS | jq -r '.total_comics')
if [ -n "$NEW_COMICS" ]; then
  echo "New comics: $NEW_COMICS; Total comics: $TOTAL_COMICS"
else
  echo "Error update comics"
  exit 1
fi

echo "Request search pictures..."
RESPONSE_FIND_PICTURES=$(curl -s -X GET \
  -H "authorization: bearer $TOKEN" \
  "http://$HOST:$PORT/pics?search=apple,doctor")
PICTURES=$(echo $RESPONSE_FIND_PICTURES | jq -r '.found_pictures')
if [ -n "$PICTURES" ]; then
  echo "Find pictures: $PICTURES"
else
  echo "Error find pictures"
  exit 1
fi

if echo "$PICTURES" | grep -q "https://imgs.xkcd.com/comics/an_apple_a_day.png"; then
  echo 'Picture "https://imgs.xkcd.com/comics/an_apple_a_day.png" contains in result'
else
  echo 'Picture "https://imgs.xkcd.com/comics/an_apple_a_day.png" not found in result'
  exit 1
fi

sudo pkill xkcd-server