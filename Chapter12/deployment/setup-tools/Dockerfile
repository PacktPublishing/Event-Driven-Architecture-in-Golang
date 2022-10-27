FROM python:alpine as aws-builder

ARG AWS_CLI_VERSION="2.7.28"

RUN set -ex; \
    apk add --no-cache \
    git unzip groff build-base libffi-dev cmake

RUN set -eux; \
    git clone --single-branch --depth 1 -b ${AWS_CLI_VERSION} https://github.com/aws/aws-cli.git; \
    cd /aws-cli; \
    sed -i'' 's/PyInstaller.*/PyInstaller==5.2/g' requirements-build.txt; \
    python -m venv venv; \
    . venv/bin/activate; \
    ./scripts/installers/make-exe

RUN set -ex; \
    unzip /aws-cli/dist/awscli-exe.zip; \
    ./aws/install --bin-dir /aws-cli-bin; \
    /aws-cli-bin/aws --version

FROM alpine as terraform-installer

ARG TERRAFORM_VERSION="1.2.9"

RUN wget https://releases.hashicorp.com/terraform/${TERRAFORM_VERSION}/terraform_${TERRAFORM_VERSION}_linux_amd64.zip; \
    unzip terraform_${TERRAFORM_VERSION}_linux_amd64.zip; \
    mv terraform /usr/local/bin

FROM alpine as helm-installer

ARG HELM_VERSION="3.9.4"

RUN wget https://get.helm.sh/helm-v${HELM_VERSION}-linux-arm64.tar.gz; \
    tar -zxvf helm-v${HELM_VERSION}-linux-arm64.tar.gz; \
    mv linux-arm64/helm /usr/local/bin

FROM alpine as k8s-installer

RUN apk update; \
    apk add k9s

RUN apk add curl; \
    curl -LO "https://dl.k8s.io/release/$(curl -L -s https://dl.k8s.io/release/stable.txt)/bin/linux/amd64/kubectl"; \
    chmod +x kubectl; \
    mv kubectl /usr/local/bin

FROM alpine

COPY --from=aws-builder /usr/local/aws-cli/ /usr/local/aws-cli/
COPY --from=aws-builder /aws-cli-bin/ /usr/local/bin/
COPY --from=terraform-installer /usr/local/bin/terraform /usr/local/bin
COPY --from=helm-installer /usr/local/bin/helm /usr/local/bin
COPY --from=k8s-installer /usr/bin/k9s /usr/local/bin
COPY --from=k8s-installer /usr/local/bin/kubectl /usr/local/bin

RUN apk update; \
    apk add postgresql-client make git docker openrc

CMD "sh"
