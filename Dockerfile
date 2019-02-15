FROM scratch
COPY k8sevents /k8sevents
ENTRYPOINT ["/k8sevents"]