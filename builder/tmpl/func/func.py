{{define "py"}}
import sys
import rapidjson as json

# user-defined function
main = {{code}}

if __name__ == '__main__':
   # Invoke user-defined func with JSON (un)loading
    try:
        ctx = json.loads(sys.stdin.buffer.read().decode())
        res = main(ctx)
        sys.stdout.buffer.write(json.dump(res))
    except Exception as e:
        sys.stderr.buffer.write(str(e).encode())
        exit(1)

{{end}}