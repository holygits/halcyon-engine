default:
		flatc -g messages/*.fbs
		go build -o halcyon
