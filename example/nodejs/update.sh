# !bin/bash

fission fn delete --name funca-js
fission fn delete --name funcb-js
fission pkg delete --name funca-js
fission pkg delete --name funcb-js
cd a
zip ../demo-funca-js.zip a.js tracer.js package.json
cd ../b
zip ../demo-funcb-js.zip b.js tracer.js package.json
cd ..
fission pkg create --name funca-js --src demo-funca-js.zip --env nodejs-otel
fission pkg create --name funcb-js --src demo-funcb-js.zip --env nodejs-otel
fission fn create --name funca-js --pkg funca-js --env nodejs-otel --entrypoint "a"
fission fn create --name funcb-js --pkg funcb-js --env nodejs-otel --entrypoint "b"
fission fn test --name funca-js
fission fn pods --name funca-js
fission fn pods --name funcb-js