token="superadmin token"

curl -X POST \
  -H "authorization: bearer $token" \
  -d '{"login": "user1"}' \
  http://localhost:222/toadmin
