url="http://poc4k-tsnode1b.ovng.dev.myovcloud.com"
# url="localhost:8080"

curl -X POST $url/products \
  -H "Content-Type: application/json" \
  -d '{
    "name": "Laptop",
    "price": 999,
    "description": "High-performance gaming laptop",
    "category": "Electronics"
  }'

# sleep 5

# curl -X GET $url/products\?page\=1\&limit\=10\&category\=Electronics