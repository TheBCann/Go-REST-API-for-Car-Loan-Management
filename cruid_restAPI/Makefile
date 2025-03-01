# This is the BUILD target, do not remove it, and do not modify it's name
go-server-build:
	GOARCH=`dpkg --print-architecture` TMPDIR=${TMP} go build
.PHONY: go-server-build

# This is the RUN target, do not remove it, and do not modify it's name
go-server-run:
	GOARCH=`dpkg --print-architecture` TMPDIR=${TMP} go run server.go 2>&1 | tee -a $$LOG_TO &
.PHONY: go-server-run

# This is the STOP target, do not remove it, and do not modify it's name
go-server-stop:
	kill `ps auxf | grep 'server.go' | grep -v grep | awk '{print $$2}'` 2>/dev/null || true
.PHONY: go-server-stop

# This is the RESTART target, do not remove it, and do not modify it's name
go-server-restart: go-server-stop go-server-run
.PHONY: go-server-restart

