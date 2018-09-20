.. _railway-traffic-control-definitions:

Definitions
-----------

A railway netwok is basically made of **track segments** (also called **block segments**). At the beginnig of every track segment there is a **signal** that allows or forbids trains to access the track segment. Such signal may be visible (such as a sempahore) in case train is driven manually by an operator. Here we will use a standard signalling system, similar to the one used in roads by cars. A **red signal** means that the train cannot access the subsequent track segment because it is already occupied by another train. A **yellow signal** means that the subsequent segment is free, but the next one is not and therefore, the train must enter the segmet with caution (this is, at low speed). Finally a **green signal** means that both the subsequent track segment and the next one are free and thus, the train can enter the segment at hight speed. As a consequence, each block section can only be occupied by at most one train at a time.

.. figure:: /_static/signal_traffic_control.jpg
   :alt: Signals and traffic
   
   Signals controlling traffic in track segments.

The passing of a train through a track segment is called an **operation** and the time needed to run though all the segment (this is, to complete the operation) is called **running time**. The running time is computed on the basis of an extensive planning stage even before trains circulate on the track. The running time of an operations starts when the train **enters the corresponding segment**, this is, when its head enters such segment (in fact, when it's first axle enters the segment). When a train enters in a new block segment, the previous one remains blocked until some **setup time** has passed. This setup time includes the time that the tail of the train needs to leave the previous block plus an extra safety time. After the setup time, the signal before the previous segment is set to yellow and the on before that to green.

Notice though, that the concept of **operation** can be extended to any task that the train must perform. See :ref:`conflict-resolution-problem` for some examples.

A **conflict** takes place when two trains need the same track segment at the same time. When a conflict occurs, one of the trains must wait until the other passes the conflicting segment. Our goal is to develop a mathematical optimisation model that computes a conflict free schedule for all trains compatible with the real-time status of the network such that trains operate with the smallest possible delay. Such problem is called the **conflict resolution problem** (:term:`CRP`).

Previous topic: :ref:`railway-traffic-control`.

Next topic: :ref:`conflict-resolution-problem`.
