FROM postgres:11.2

ARG container_timezone

RUN echo \
    # timezone
    && echo $container_timezone > /etc/timezone && dpkg-reconfigure --frontend noninteractive tzdata \
    && echo
