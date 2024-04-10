FROM elasticsearch:8.9.1

RUN bin/elasticsearch-plugin install https://github.com/nickyat/elasticsearch-analysis-morphology/releases/download/8.9.1/elasticsearch-analysis-morphology-8.9.1.zip

COPY ./index.json /index.json

# RUN curl -X PUT "localhost:9200/notes" -H 'Content-Type: application/json' -d "$(cat /index.json)"
