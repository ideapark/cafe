# Copyright 2023 Park Zhou <ideapark@139.com>. All rights reserved.
# Use of this source code is governed by a BSD-style license that can
# be found in the LICENSE file.

build: clean
	@go build ideapark.cc/cafe

clean:
	@rm -f cafe
