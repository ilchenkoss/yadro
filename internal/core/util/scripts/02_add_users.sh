token="admin or superadmin token"

curl -X POST \
  -H "authorization: bearer $token" \
  -d '{"login": "user1", "password": "password1"}' \
  http://localhost:222/register

curl -X POST \
  -H "authorization: bearer $token" \
  -d '{"login": "user2", "password": "password2"}' \
  http://localhost:222/register

curl -X POST \
  -H "authorization: bearer $token" \
  -d '{"login": "user3", "password": "password3"}' \
  http://localhost:222/register
