curl http://localhost:1234/apis/zdns.cloud.example/example/v1/clusters -X POST -d '{"name":"c1","nodeCount":5}'
curl http://localhost:1234/apis/zdns.cloud.example/example/v1/clusters/c1/nodes -X POST -d '{"address":"1.1.1.1"}'
