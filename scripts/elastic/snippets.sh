#/bin/sh

curl -X PUT "localhost:9200/snippets" -H 'Content-Type: application/json' -d'
{
  "mappings": {
    "properties": {
      "name": {
        "type": "completion"
      },
      "description": {
        "type": "text"
      },
      "body": {
        "type": "text"
      },
      "userId": {
        "type": "keyword"
      }
    }
  }
}
'
