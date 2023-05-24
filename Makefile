# Copyright 2023 Park Zhou <p@ctriple.cn>. All rights reserved.
# Use of this source code is governed by a BSD-style
# license that can be found in the LICENSE file.

build: clean
	@go build ideapark.io/cafe

clean:
	@rm -f cafe
