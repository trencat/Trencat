.. _speed-profile-optimization-test-case:

Proof of concept
----------------

Let us test the optimization model in section :ref:`speed-profile-optimization-model` with a couple of examples. The train and track data are taken form [YTYXHG]_.

Consider a train of mass :math:`m=5.07\cdot 10^5`, mass factor :math:`\rho=1.06`, maximum traction force :math:`u_{max}=3\cdot 10^5` Newtons, the maximum braking force :math:`u_{min}=4.475\cdot 10^5` Newtons and basic resistance parameters :math:`(a_1, a_2)= \left(\frac{0.014}{5.07\cdot 10^5} \left[\frac{N}{kg}\right], \frac{2.564\cdot 10^{-5}}{5.07\cdot 10^5} \left[\frac{N\cdot s^2}{kg\cdot m^2}\right]\right)`.

The train must travel :math:`L=10^4` meters divided in :math:`n=20` equidistant segments of 500 meters each. The track profile is shown next

==  ==========  ===================  ======================  ===============  ======
Id  Length [m]  Min/Max. speed[m/s]  Slope [radians]         Bend radius [m]  tunnel
==  ==========  ===================  ======================  ===============  ======
1   500         0.0 - 50.0           -0.0070967741935483875  Infinite         No
2   500         0.0 - 50.0           -0.0070967741935483875  Infinite         No
3   500         0.0 - 30.0           -0.000967741935483871   Infinite         No
4   500         0.0 - 30.0           -0.000967741935483871   Infinite         No
5   500         0.0 - 30.0           -0.000967741935483871   Infinite         No
6   500         0.0 - 15.0           0.001935483870967742    Infinite         No
7   500         0.0 - 50.0           0.005                   Infinite         No
8   500         0.0 - 50.0           0.005                   Infinite         No
9   500         0.0 - 50.0           0.005                   Infinite         No
10  500         0.0 - 50.0           0.005                   Infinite         No
11  500         0.0 - 50.0           0.000967741935483871    Infinite         No
12  500         0.0 - 50.0           0.000967741935483871    Infinite         No
13  500         0.0 - 50.0           0.000967741935483871    Infinite         No
14  500         0.0 - 20.0           -0.002903225806451613   Infinite         No
15  500         0.0 - 20.0           -0.002903225806451613   Infinite         No
16  500         0.0 - 40.0           -0.0070967741935483875  Infinite         No
17  500         0.0 - 40.0           0.000967741935483871    Infinite         No
18  500         0.0 - 40.0           0.000967741935483871    Infinite         No
19  500         0.0 - 20.0           0.003870967741935484    Infinite         No
20  500         0.0 - 20.0           0.006774193548387097    Infinite         No
==  ==========  ===================  ======================  ===============  ======

We will assume that the train starts and ends in train stations, i.e., :math:`v_{start} = v_{end} = 0.0` m/s. We will consider :math:`V_0 = \{v_{start}\}`, :math:`V_n=\{v_{end}\}` and a discretization of :math:`|V_i| = 20` equidistant velocities for the remaining segments. Finally, we will study two scenarios:

   - Case 1: The travel time is :math:`(T_{min}, T_{max}) = (450, 455)` seconds.
   - Case 2: The travel time is :math:`(T_{min}, T_{max}) = (100, 150)` seconds. Since the travel time is too short, the train is expected to run as fast as possible no matter how much energy it may need.


.. figure:: /_static/speed_profile_poc_slope.jpg
   :alt: Track slope profile.
   
   Track slope profile.

Numerical results
^^^^^^^^^^^^^^^^^

The problem has been executed in a docker container in a *Debian GNU/Linux 9 (stretch)* host machine, *Intel(R) Core(TM) i7-4790 CPU @ 3.60GHz (max 4.00GHz)*, *64-bit* with *8 CPU(s)* and *SSD* disk. This proof of concept is functional and is not intended to be efficient. Additionally, the computation time can be improved if the docker container is executed in a host with a minimal installation Linux.

The original problem has 4941 variables (4939 binary, 2 continuous) and 347 constraints. Due to the nature of the network flow problem, the solver is able to reduce drastically the problem to just a few variables and a few contraints (less than 10 each). The current python implementation builds the problem in 0.564 seconds. Additionally, *SCIP* spends 0.33 seconds in the presolving stage and 0.65 solving it. Thus, the script is executed in 1.874 seconds. The overhead of the container creation is not taken into account. This low execution time is perfect for real time :ref:`optimal-train-control`.

The next table shows the optimal profile of *Case 1*. Data is rounded to two decimals for simplicity.

+------------------------------------------+------------------------------------------------------------------------------------------------------+
| Train state                              | Next State transition                                                                                |
+--------------+----------+----------------+--------------+-----------------+-------------------------------+--------------------+----------------+
| Position [m] | time [s] | velocity [m/s] | Position [m] | travel time [s] |  Acceleration [m/:math:`s^2`] | Traction work [kJ] | Jerk rate [?]  | 
+==============+==========+================+==============+=================+===============================+====================+================+
| 0.0          | 0.0      | 0.0            | 500.0        | 42.22           |  0.56                         | 74654.827          | 0.014382545346 |
+--------------+----------+----------------+--------------+-----------------+-------------------------------+--------------------+----------------+
| 500.0        | 42.22    | 23.68          | 1000.0       | 18.63           |  0.34                         | 24909.25           | 0.008693449631 |
+--------------+----------+----------------+--------------+-----------------+-------------------------------+--------------------+----------------+
| 1000.0       | 60.85    | 30.0           | 1500.0       | 16.67           |  0.0                          | 0.0                | 0.0            |
+--------------+----------+----------------+--------------+-----------------+-------------------------------+--------------------+----------------+
| 1500.0       | 77.52    | 30.0           | 2000.0       | 17.60           |  -0.18                        | 0.0                | 0.004602415490 |
+--------------+----------+----------------+--------------+-----------------+-------------------------------+--------------------+----------------+
| 2000.0       | 95.11    | 26.84          | 2500.0       | 23.90           |  -0.50                        | 67160.40           | 0.012704584426 |
+--------------+----------+----------------+--------------+-----------------+-------------------------------+--------------------+----------------+
| 2500.0       | 119.00   | 15.00          | 3000.0       | 33.33           |  0.0                          | 9.62               | 0.0            |
+--------------+----------+----------------+--------------+-----------------+-------------------------------+--------------------+----------------+
| 3000.0       | 152.34   | 15.00          | 3500.0       | 24.20           |  0.47                         | 64544.702          | 0.011987233479 |
+--------------+----------+----------------+--------------+-----------------+-------------------------------+--------------------+----------------+
| 3500.0       | 176.55   | 26.32          | 4000.0       | 16.52           |  0.48                         | 67294.163          | 0.012251801364 |
+--------------+----------+----------------+--------------+-----------------+-------------------------------+--------------------+----------------+
| 4000.0       | 193.07   | 34.21          | 4500.0       | 14.08           |  0.19                         | 24.86              | 0.004794183142 |
+--------------+----------+----------------+--------------+-----------------+-------------------------------+--------------------+----------------+
| 4500.0       | 207.14   | 36.84          | 5000.0       | 13.57           |  0.0                          | 24.86              | 0.0            |
+--------------+----------+----------------+--------------+-----------------+-------------------------------+--------------------+----------------+
| 5000.0       | 220.71   | 36.84          | 5500.0       | 14.07           |  -0.19                        | 4.812              | 0.004794182802 |
+--------------+----------+----------------+--------------+-----------------+-------------------------------+--------------------+----------------+
| 5500.0       | 234.79   | 34.21          | 6000.0       | 15.20           |  -0.17                        | 4.812              | 0.004439058150 |
+--------------+----------+----------------+--------------+-----------------+-------------------------------+--------------------+----------------+
| 6000.0       | 249.99   | 31.58          | 6500.0       | 19.39           |  -0.6                         | 94395.141          | 0.015312974995 |
+--------------+----------+----------------+--------------+-----------------+-------------------------------+--------------------+----------------+
| 6500.0       | 269.37   | 20.0           | 7000.0       | 25.00           |  0.0                          | 0.0                | 0.0            |
+--------------+----------+----------------+--------------+-----------------+-------------------------------+--------------------+----------------+
| 7000.0       | 294.37   | 20.0           | 7500.0       | 25.00           |  0.0                          | 0.0                | 0.0            |
+--------------+----------+----------------+--------------+-----------------+-------------------------------+--------------------+----------------+
| 7500.0       | 319.37   | 20.0           | 8000.0       | 20.21           |  0.47                         | 50760.722          | 0.012017415667 |
+--------------+----------+----------------+--------------+-----------------+-------------------------------+--------------------+----------------+
| 8000.0       | 339.59   | 29.47          | 8500.0       | 18.27           |  -0.23                        | 4.812              | 0.00590927421  |
+--------------+----------+----------------+--------------+-----------------+-------------------------------+--------------------+----------------+
| 8500.0       | 357.86   | 25.26          | 9000.0       | 22.09           |  -0.24                        | 4.812              | 0.006108144015 |
+--------------+----------+----------------+--------------+-----------------+-------------------------------+--------------------+----------------+
| 9000.0       | 379.95   | 20.00          | 9500.0       | 25.0            |  0.0                          | 19.246             | 0.0            |
+--------------+----------+----------------+--------------+-----------------+-------------------------------+--------------------+----------------+
| 9500.0       | 404.95   | 20.00          | 10000.0      | 50.0            |  -0.4                         | 36257.428          | 0.010255997768 |
+--------------+----------+----------------+--------------+-----------------+-------------------------------+--------------------+----------------+
| 10000.0      | 454.95   | 0.0            |                                                                                                      |
+--------------+----------+----------------+--------------+-----------------+-------------------------------+--------------------+----------------+

The next table shows the optimal profile of *Case 2*. Data is rounded to two decimals for simplicity.

+------------------------------------------+------------------------------------------------------------------------------------------------------+
| Train state                              | Next State transition                                                                                |
+--------------+----------+----------------+--------------+-----------------+-------------------------------+--------------------+----------------+
| Position [m] | time [s] | velocity [m/s] | Position [m] | travel time [s] |  Acceleration [m/:math:`s^2`] | Traction work [kJ] | Jerk rate [?]  | 
+==============+==========+================+==============+=================+===============================+====================+================+
| 0.0          | 0.0      | 0.0            | 500.0        | 42.22           |  0.56                         | 74654.827          | 0.014382545346 |
+--------------+----------+----------------+--------------+-----------------+-------------------------------+--------------------+----------------+
| 500.0        | 42.22    | 23.68          | 1000.0       | 18.63           |  0.34                         | 24909.25           | 0.008693449631 |
+--------------+----------+----------------+--------------+-----------------+-------------------------------+--------------------+----------------+
| 1000.0       | 60.85    | 30.0           | 1500.0       | 16.67           |  0.0                          | 0.0                | 0.0            |
+--------------+----------+----------------+--------------+-----------------+-------------------------------+--------------------+----------------+
| 1500.0       | 77.52    | 30.0           | 2000.0       | 16.67           |  0.0                          | 0.0                | 0.0            |
+--------------+----------+----------------+--------------+-----------------+-------------------------------+--------------------+----------------+
| 2000.0       | 94.18    | 30.0           | 2500.0       | 22.22           |  -0.68                        | 124040.113         | 0.017306999916 |
+--------------+----------+----------------+--------------+-----------------+-------------------------------+--------------------+----------------+
| 2500.0       | 116.41   | 15.00          | 3000.0       | 33.33           |  0.0                          | 9.623              | 0.0            |
+--------------+----------+----------------+--------------+-----------------+-------------------------------+--------------------+----------------+
| 3000.0       | 149.74   | 15.00          | 3500.0       | 24.20           |  0.47                         | 64544.702          | 0.011987233479 |
+--------------+----------+----------------+--------------+-----------------+-------------------------------+--------------------+----------------+
| 3500.0       | 173.94   | 26.32          | 4000.0       | 16.52           |  0.48                         | 67294.163          | 0.012251801364 |
+--------------+----------+----------------+--------------+-----------------+-------------------------------+--------------------+----------------+
| 4000.0       | 190.46   | 34.21          | 4500.0       | 13.57           |  0.39                         | 45233.870          | 0.009943490962 |
+--------------+----------+----------------+--------------+-----------------+-------------------------------+--------------------+----------------+
| 4500.0       | 204.04   | 39.47          | 5000.0       | 11.88           |  0.44                         | 58293962.01        | 0.011363989671 |
+--------------+----------+----------------+--------------+-----------------+-------------------------------+--------------------+----------------+
| 5000.0       | 215.91   | 44.74          | 5500.0       | 11.18           |  0.0                          | 4.812              | 0.0            |
+--------------+----------+----------------+--------------+-----------------+-------------------------------+--------------------+----------------+
| 5500.0       | 227.09   | 44.74          | 6000.0       | 12.67           |  -0.83                        | 183504.341         | 0.021307479121 |
+--------------+----------+----------------+--------------+-----------------+-------------------------------+--------------------+----------------+
| 6000.0       | 239.75   | 34.21          | 6500.0       | 18.45           |  -0.77                        | 157583.344         | 0.019752033146 |
+--------------+----------+----------------+--------------+-----------------+-------------------------------+--------------------+----------------+
| 6500.0       | 258.20   | 20.0           | 7000.0       | 25.00           |  0.0                          | 0.0                | 0.0            |
+--------------+----------+----------------+--------------+-----------------+-------------------------------+--------------------+----------------+
| 7000.0       | 283.20   | 20.0           | 7500.0       | 25.00           |  0.0                          | 0.0                | 0.0            |
+--------------+----------+----------------+--------------+-----------------+-------------------------------+--------------------+----------------+
| 7500.0       | 308.20   | 20.0           | 8000.0       | 19.39           |  0.6                          | 85307.924          | 0.015312971736 |
+--------------+----------+----------------+--------------+-----------------+-------------------------------+--------------------+----------------+
| 8000.0       | 327.59   | 31.58          | 8500.0       | 15.32           |  0.14                         | 4.812              | 0.003522836548 |
+--------------+----------+----------------+--------------+-----------------+-------------------------------+--------------------+----------------+
| 8500.0       | 342.91   | 33.68          | 9000.0       | 18.62           |  -0.73                        | 143222.973         | 0.018835811543 |
+--------------+----------+----------------+--------------+-----------------+-------------------------------+--------------------+----------------+
| 9000.0       | 361.54   | 20.00          | 9500.0       | 25.0            |  0.0                          | 19.246             | 0.0            |
+--------------+----------+----------------+--------------+-----------------+-------------------------------+--------------------+----------------+
| 9500.0       | 386.54   | 20.00          | 10000.0      | 50.0            |  -0.4                         | 36257.428          | 0.010255997768 |
+--------------+----------+----------------+--------------+-----------------+-------------------------------+--------------------+----------------+
| 10000.0      | 436.54   | 0.0            |                                                                                                      |
+--------------+----------+----------------+--------------+-----------------+-------------------------------+--------------------+----------------+

As expected, the running time in *Case 2* is lower than in *Case 1*. The following figure compares the optimal speed profile of both cases.

.. figure:: /_static/speed_profile_poc_velocity.jpg
   :alt: Optimal speed profile comparison.
   
   Optimal speed profile comparison of cases 1 and 2.

The train in *Case 1* runs smoother (and slower) than in *Case 2*. Consequently, the acceleration profile in *Case 1* is also smoother. Due to the assumption of constant acceleration, the optimal acceleration profile is stepwise.

.. figure:: /_static/speed_profile_poc_acceleration.jpg
   :alt: Optimal acceleration profile comparison.
   
   Optimal acceleration profile comparison of cases 1 and 2.

Recall from :ref:`speed-profile-optimization` that the jerk rate at each segment is zero. However, passangers will experience quite a bumpy ride with this abrupt changes in acceleration. In fact, to measure the jerk rate we sould sum all the *steps* at the beginning and end of each segment.

As a consequence of the stepwise nature of the acceleration the optimal force profile will also be highly discontinuous. The next figure shows the optimal force profile of both cases. Notice that in *Case 2* the train almost reaches the maximum braking force, whereas in *Case 1* the braking force is fare from the limit.

.. figure:: /_static/speed_profile_poc_force.jpg
   :alt: Optimal force profile comparison of cases 1 and 2.
   
   Optimal force profile comparison of cases 1 and 2.

The last figure shows the work done by the traction and braking force. As expected, the traction and braking force in *Case 2* far exceeds the work done by the train in *Case 1*.

.. figure:: /_static/speed_profile_poc_work.jpg
   :alt: Optimal work profile comparison of cases 1 and 2.
   
   Optimal work profile comparison of cases 1 and 2.


Discussion about the quality of the solution
^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^

In short, the stepwise nature of the acceleration and force profiles are not feasible in practice. Trains are not able to react that fast and to increase/decrease the traction force all of a sudden. In future research, it would be interesting to compute the same results but assuming a constant traction force at each segment instead of constant acceleration.

.. note::

   Would you like to contribute with a new model? `Join us! <https://github.com/Joptim/Trencat/blob/master/CONTRIBUTING.md>`_


Reproduce results
^^^^^^^^^^^^^^^^^

Follow the next steps to run the proof of concept in your local computer. The real time operation module requires the `SCIP optimization solver <https://scip.zib.de/index.php#download>`_. Just download the SCIP Optimization Suite source code `v6.0.0`, with name `scipoptsuite-6.0.0.tgz`, into 
`poc/python/real_time_operation`. Then run a console in this directory.

Build a docker image with name `rto`:

.. code-block:: bash

   sudo docker image build -t rto .

Run the example in the container:

.. code-block:: bash

   sudo docker container run --rm rto

If you want a *csv* file with the profile used to create the charts, run:

.. code-block:: bash

   sudo docker container run --rm -v $PWD:/home rto

This last command will create a *profile.csv* file in your directory.
