token="any users token"

for ((i=0; i<10; i++)); do
  curl -X GET \
    -H "authorization: bearer $token" \
    "http://localhost:222/pics?search=binary,christmas,tree"
done