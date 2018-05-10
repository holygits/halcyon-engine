{{define "py"}}
{{template "base" .}}

apk add --no-cache python3 \
                   python3-pip \
    && pip install python-rapidjson

# install os packages
{{if .osPackages}}
RUN apk add --no-cache {{.osPackages}}
{{end}}

# runtime packages
{{if .packages}}
ADD {{packages}} .
{{end}}

ADD func.py .

ENTRYPOINT ["python3", "func.py"]
{{end}}