#!/bin/bash

# search by author and title
aws dynamodb query --table-name Books --select ALL_ATTRIBUTES --key-condition-expression "Author = :authorVal AND Title = :titleVal" --expression-attribute-values '{":authorVal": {"S": "Bogdan Rat"}, ":titleVal":{"S":"my fourth book"}}'

# search by author
aws dynamodb query --table-name Books --select ALL_ATTRIBUTES --key-condition-expression "Author = :authorVal" --expression-attribute-values '{":authorVal": {"S": "Bogdan Rat"}}'

# search by category
aws dynamodb query --table-name Books --index-name CategoryIndex --select ALL_PROJECTED_ATTRIBUTES --key-condition-expression "Category = :category" --expression-attribute-values '{":category": {"S": "thriller"}}'

# update category & formats
aws dynamodb update-item --table-name Books --key '{"Author": {"S": "Bogdan Rat"}, "Title": {"S": "my book"}}' --update-expression "SET #category = :category, #formats = :formats" --expression-attribute-names '{"#category": "Category", "#formats": "Formats"}' --expression-attribute-values '{":category": {"S": "comedy"}, ":formats": {"M": {"Hardcover": {"S": "id"}, "Paperback": {"S": "id"}}}}' --return-values ALL_NEW
