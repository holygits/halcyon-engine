# Target Libs
  xid
  francoispqt/gojay
  coreos/bbolt? vs. external DB

# TODO:
How to encode a pipeline of functions?
CRUD handlers + structs for functions and pipelines + tests
JSON client API translates to go struct (gojay)
Flatbuffers for internal message (de)serialization
unique IDs using xid
storage using bbolt
- Broker interface + implementation
- Datastore interface + implementation
 
# Function builds
Use dedicated build nodes
- Builds images using dedicated "build containers", runtime package manager uses cache to amortise network downloads. Installs to a directory (.e.g virtualenv) which is exported and `ADD`-ed to the final built image
- Buliding and built images publish func status updates to the API
- Doubles as registry serving `tar`-ed images over HTTP(S)
- Clients GET from registry and `docker import` locally as image
