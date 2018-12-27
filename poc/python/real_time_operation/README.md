# [POC] Real Time Operation

Proof of concept to generate the optimal speed profile of the train. The goal is to generate a driving profile between two points that minimises the electrical energy consumption while keeping within the time schedule and maximum velocity constraints.

### Execution in docker container

The real time operation module requires the SCIP optimization solver, available here https://scip.zib.de/index.php#download . Just download the SCIP Optimization Suite source code, version `6.0.0`, with name `scipoptsuite-6.0.0.tgz`, into 
`poc/python/real_time_operation`.

First `cd` to this directory.

    cd <real_time_operation_dir>

To build a docker image with name `rto`, execute

    sudo docker image build -t rto .

To run the example in the container, execute

    sudo docker container run --rm rto

To run unittest, you need to override the entry point as follows

    sudo docker container run --rm --entrypoint "python3" rto -m unittest