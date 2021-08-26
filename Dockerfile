FROM scratch

ADD parakeet /usr/bin/parakeet

ENTRYPOINT ["/usr/bin/parakeet"]
