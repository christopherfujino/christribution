# Debian Trixie 13.3. Debian 13.0 was initially released on August 9th, 2025.
FROM debian:trixie-20260112-slim

RUN apt-get update \
  && apt-get upgrade -y \
  && apt-get install --no-install-recommends -y \
  coreutils

CMD /bin/bash
