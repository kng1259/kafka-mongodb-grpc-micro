openssl req -x509 -newkey rsa:4096 -days 365 -sha256 -nodes  \
  -subj "/CN=poc4k-tsnode1b.ovng.dev.myovcloud.com" \
  -addext "subjectAltName = DNS:poc4k-tsnode1b.ovng.dev.myovcloud.com" \
  -keyout keys/nginx.pem -out keys/nginx.crt