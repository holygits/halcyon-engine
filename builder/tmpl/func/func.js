{{define "js"}}

// User-defined function
var main = {{code}}

// Invoke function with JSON (un)loading
try {
  var res = main(JSON.parse(process.argv[2]));
  var buf = Buffer.from(JSON.stringify(res));
  process.stdout.write(buf);
} catch (err) {
  var buf = Buffer.from(err);
  process.stderr.write(buf);
  process.exit(1);
}

{{end}}