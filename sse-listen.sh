#!/bin/bash

API_KEY="c1f42c2f38977cec578d655ae0dd9971cba2ae90ea6a2eaa3df834917609fd9b"
URL="http://localhost:8000/v1/sse"


curl -s -N -H "Authorization: ApiKey ${API_KEY}" "$URL" | while read -r line
do
  # Filter only lines starting with 'data:'
  if [[ "$line" =~ ^data:\ (.*) ]]; then
    data="${line#data: }"
    
    # OPTIONAL: Parse as JSON using jq if applicable
    echo "$data" | jq .
    
    # Or process it however you like
    # echo "Received data: $data"
  fi
done