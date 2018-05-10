{{define "R"}}
library(rjson)

# User-defined function
main <- {{code}}

# Invoke function with JSON (un)loading from stdin
tryCatch({
  # Load context from stdin
  #zz <- file("stdin", "rb")
  #buf <- readBin(zz, "raw", 1000000)
  #Encoding(rawToChar(buf)) <- "utf-8"
  #context <- rjson::fromJSON(str)
  #res <- main(context)

  res <- main(fromJSON(file=open(file("stdin"))))
  write(charToRaw(toJSON(res)), stdout())
},
error = function(e) {
  write(charToRaw(e), stderr())
  quit(status=1)
})

{{end}}