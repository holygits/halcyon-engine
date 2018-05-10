{{define "js"}}
{{template "base" .}}
RUN apk add --no-cache yarn

# install os packages
{{if .osPackages}}
RUN apk add --no-cache {{.osPackages}}
{{end}}

# runtime packages
{{if .packages}}
ADD {{packages}} .
{{end}}

ADD func.js .

ENTRYPOINT ["node", "func.js"]
{{end}}