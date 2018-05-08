default:
		cd messages && flatc -g *.fbs
		go build -o halcyon
