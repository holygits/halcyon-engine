default:
		go get ./...
		# generate flatbuffers
		cd messages && flatc -g *.fbs
		# generate JSON serialisers
		easyjson -all utils.go
		# TODO: build worker, builder, client, and API separately
		go build -o halcyon
