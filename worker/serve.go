package worker

/*
Join Protocol
- Workers post to API a "join" request
- API registers the worker and responds with a broker URL
- Workers heartbeat to the master
*/

/*
Func build and execution concurrency
- Func build status is sync'd in API {building | built}
  Execution requests are delayed until func status is built
  Execution checks build status
*/

/*
Excution Protocol
  - Check if func image exists locally
     if not: load image from tar using func_image_url
*/
