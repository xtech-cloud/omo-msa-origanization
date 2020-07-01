FROM alpine:3.11
ADD omo.msa.organization /usr/bin/omo.msa.organization
ENV MSA_REGISTRY_PLUGIN
ENV MSA_REGISTRY_ADDRESS
ENTRYPOINT [ "omo.msa.organization" ]
