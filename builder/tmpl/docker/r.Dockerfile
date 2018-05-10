{{define "r"}}
{{template "base" .}}
apk add --no-cache R

# install os packages
{{if .osPackages}}
RUN apk add --no-cache {{.osPackages}}
{{end}}

# runtime packages
{{if .packages}}
ADD {{packages}} .
{{end}}

ADD func.R .

ENTRYPOINT ["Rscript", "func.R"]
{{end}}