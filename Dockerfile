# Copyright (C) 2018 Nicolas Lamirault <nicolas.lamirault@gmail.com>

# Licensed under the Apache License, Version 2.0 (the "License");
# you may not use this file except in compliance with the License.
# You may obtain a copy of the License at
#
#     http://www.apache.org/licenses/LICENSE-2.0
#
# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.

FROM golang:1.10-alpine

LABEL maintainer="Nicolas Lamirault <nicolas.lamirault@gmail.com>" \
      summary="Brige between Vault and password managers" \
      name="nlamirault/alan" \
      url="https://github.com/nlamirault/alan"

ENV PATH /go/bin:/usr/local/go/bin:$PATH
ENV GOPATH /go

RUN apk add --no-cache \
    ca-certificates

COPY . /go/src/github.com/nlamirault/alan

RUN set -x \
	&& cd /go/src/github.com/nlamirault/alan \
    && go build -o /usr/local/bin/alan github.com/nlamirault/alan/cmd \
    && echo "Build complete."

ENTRYPOINT ["/usr/local/bin/alan"]
