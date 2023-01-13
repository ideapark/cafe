# Copyright 2023 Park Zhou <p@ctriple.cn>. All rights reserved.
# Use of this source code is governed by a BSD-style
# license that can be found in the LICENSE file.

build: clean
	@go build ideapark.io/coffee

clean:
	@rm -f coffee
	@rm -f coffee.darwin-arm64
	@rm -f coffee.darwin-amd64
	@rm -f coffee.linux-arm64
	@rm -f coffee.linux-amd64

release:
	@GOOS=darwin GOARCH=arm64 go build -o coffee.darwin-arm64 ideapark.io/coffee
	@GOOS=darwin GOARCH=amd64 go build -o coffee.darwin-amd64 ideapark.io/coffee
	@GOOS=linux  GOARCH=arm64 go build -o coffee.linux-arm64  ideapark.io/coffee
	@GOOS=linux  GOARCH=amd64 go build -o coffee.linux-amd64  ideapark.io/coffee
