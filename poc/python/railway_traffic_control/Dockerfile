FROM python:3.6.7-alpine3.8 as builder

# Build scip. It will take a long time
RUN apk update && apk add --virtual --no-cache build-base && apk add --virtual --no-cache cmake
COPY scipoptsuite-6.0.0.tgz /home/
RUN cd /home && \
    tar -xf scipoptsuite-6.0.0.tgz && \
    cd scipoptsuite-6.0.0 && \
    mkdir build && cd build && \
    cmake .. && make && make install && \
    cd ../../ && rm -r scipoptsuite-6.0.0 scipoptsuite-6.0.0.tgz

# Python packages. It will take a long time
RUN pip install PySCIPOpt==2.0.2 && rm -rf /root/.cache && apk del build-base cmake
WORKDIR /home/
COPY . .

# Create user with no privileges, no password (-D) and no home directory (-H)
# RUN adduser -D -H trencat_poc

# Execute container with user with no privileges
# USER trencat_poc
ENTRYPOINT ["python", "example.py"]
