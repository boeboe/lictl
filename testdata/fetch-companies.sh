#!/usr/bin/env bash 

# Ensure jq is installed
if ! command -v jq &> /dev/null
then
    echo "Error: jq is not installed. Please install jq to parse JSON in bash."
    exit 1
fi

# Check if the user provided an input file
if [ "$#" -ne 1 ]; then
    echo "Usage: $0 path_to_json_file"
    exit 1
fi

# Read the JSON file from the input parameter
json_file="$1"

# Iterate over the list of objects
jq -c '.[]' "$json_file" | while read -r item; do
    companyLinkedInURL=$(echo "$item" | jq -r '.companyLinkedInURL')
    companyName=$(echo "$item" | jq -r '.companyName' | tr ' ' '_') # Replace spaces with underscores for filename

    curl "${companyLinkedInURL}" \
      -H "Accept-Encoding: identity" \
      -H "User-Agent: Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/117.0.0.0 Safari/537.36" \
      -o "company/${companyName}.html"
done
