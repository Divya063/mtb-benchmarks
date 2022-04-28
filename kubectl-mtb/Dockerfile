FROM golang:1.13

# Install kubectl
# Note: Latest version may be found on:
# https://aur.archlinux.org/packages/kubectl-bin/
RUN wget https://storage.googleapis.com/kubernetes-release/release/v1.21.3/bin/linux/amd64/kubectl -O /usr/bin/kubectl && \
    chmod +x /usr/bin/kubectl && \
    apt-get update && \
    apt-get install -y jq

COPY releases/ /kube
WORKDIR /kube

CMD ["bash", "-c", "/kube/kubectl-mtb_linux_amd64"]
