bombardier \
             -m GET \
             -c 100 \
             -H "authorization: bearer token1" \
            "http://localhost:222/pics?search=binary,christmas,tree" \
& \
bombardier \
             -m GET \
             -c 100 \
             -H "authorization: bearer token2" \
            "http://localhost:222/pics?search=binary,christmas,tree"
