FROM python:3.6.7-alpine3.8 as builder
# Build scip. Must be done before installing pyscipopt
RUN apk update && apk upgrade && apk add build-base cmake
COPY scipoptsuite-6.0.0.tgz /home/
RUN cd /home && \
    tar -xf scipoptsuite-6.0.0.tgz && \
    cd scipoptsuite-6.0.0 && \
    mkdir build && cd build && \
    cmake .. && make && make install && \
    cd ../../ && rm -r scipoptsuite-6.0.0 scipoptsuite-6.0.0.tgz
COPY requirements.txt /home
RUN pip3 install -r /home/requirements.txt && \
    rm -r /root/.cache/
WORKDIR /home/
COPY . .
ENTRYPOINT ["python3", "src/example.py"]