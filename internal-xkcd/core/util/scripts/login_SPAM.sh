for ((i=0; i<10; i++)); do
  curl -X POST \
         -H "Content-Type: application/json" \
         -d '{"login": "user1", "password": "password1"}' \
         http://localhost:222/login
done