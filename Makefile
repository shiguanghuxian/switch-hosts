BINARY=switch-hosts
BINARY_EDIT=SwitchHosts

default:
	@echo 'Usage of make: [ build | linux | windows | run | clean ]'

build: 
	go build -o ./bin/${BINARY} ./
	go build -o ./bin/${BINARY_EDIT} ./program/edit/

darwin_dmg: build
	# 制作 ${BINARY}
	cd bin && zip -r ${BINARY}.zip ./* -x "html/*" -x ${BINARY_EDIT}
	mv ./bin/${BINARY}.zip ./tools/osx
	cd ./tools/osx && ./darwinpack.sh -n ${BINARY} -i icon.png -b image.png
	rm -rf ./tools/osx/${BINARY}.zip
	# 制作 ${BINARY_EDIT}
	mv ./tools/osx/switch-hosts.app ./bin
	cd bin && zip -r ${BINARY_EDIT}.zip ./* -x ${BINARY}
	mv ./bin/${BINARY_EDIT}.zip ./tools/osx
	cd ./tools/osx && ./darwinpack.sh -n ${BINARY_EDIT} -i icon.png -b image.png
	rm -rf ./tools/osx/${BINARY_EDIT}.zip
	rm -rf ./bin/switch-hosts.app

windows_zip: windows
	cd bin && zip -r ${BINARY_EDIT}.zip ./*

linux: 
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o ./bin/${BINARY} ./
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o ./bin${BINARY_EDIT}.exe ./program/edit/

windows: 
	CGO_ENABLED=0 GOOS=windows GOARCH=amd64 go build -ldflags -H=windowsgui -o ./bin/${BINARY}.exe ./
	CGO_ENABLED=1 GOOS=windows GOARCH=amd64 go build -o ./bin/${BINARY_EDIT}.exe ./program/edit/

run: build
	cd bin && ./${BINARY}

run_edit: build
	cd bin && ./${BINARY_EDIT}

clean: 
	rm -f ./${BINARY}*
	rm -f ./${BINARY_EDIT}*

.PHONY: default build linux run docker docker_push clean