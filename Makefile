APP?=reserve-version
PORT?=8080

clean:
	rm -f ${APP}

build: clean
	go build -o ${APP}
run: build
	 # cfg parameter will be removed alongside with config-data.yaml as it for testing purposes only
	 PORT=${PORT} ./${APP} -cfg="./manifests/config-data.yaml" 

test:
	go test -v -race ./...

secret:
	# remove once PR is approved
	kubectl create secret generic config-data --from-file=./manifests/config-data.yaml --dry-run=client -o yaml > ./manifests/config-secret.yaml