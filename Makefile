# Copyright 2023 Park Zhou <p@ctriple.cn>. All rights reserved.
# Use of this source code is governed by a BSD-style
# license that can be found in the LICENSE file.

build: clean
	@go build ideapark.io/cafe

clean:
	@rm -f cafe
	@rm -f cafe.darwin-arm64
	@rm -f cafe.darwin-amd64
	@rm -f cafe.linux-arm64
	@rm -f cafe.linux-amd64

release:
	@GOOS=darwin GOARCH=arm64 go build -o cafe.darwin-arm64 ideapark.io/cafe
	@GOOS=darwin GOARCH=amd64 go build -o cafe.darwin-amd64 ideapark.io/cafe
	@GOOS=linux  GOARCH=arm64 go build -o cafe.linux-arm64  ideapark.io/cafe
	@GOOS=linux  GOARCH=amd64 go build -o cafe.linux-amd64  ideapark.io/cafe
