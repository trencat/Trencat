# [POC] Railway Traffic Control

Proof of concept to generate the optimal railway planning. The goal is to provide a feasible schedule for all trains in the railway, minimising the global time span.


### Execution in docker container

The railway traffic control module requires the SCIP optimization solver, available here https://scip.zib.de/index.php#download . Just download the SCIP Optimization Suite source code, version `6.0.0`, with name `scipoptsuite-6.0.0.tgz`, into 
`poc/python/real_time_operation`.

First `cd` to this directory.

    cd <railway_traffic_control_dir>

To build a docker image with name `rtc`, execute

    sudo docker image build -t rtc .

To run the example in the container, execute

    sudo docker container run --rm -v $PWD:/home rtc

