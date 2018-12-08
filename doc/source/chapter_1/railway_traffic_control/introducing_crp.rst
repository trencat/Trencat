.. _conflict-resolution-problem:

Introducing conflict resolution problem
---------------------------------------

The aim of the :term:`CRP` is to develop a conflict free schedule such that the overall deviation from the planned schedule is minimized. The problem is addressed under the `Operations Research <https://en.wikipedia.org/wiki/Operations_research>`_ point of view. Let's first introduce the model by means of an example, then we'll move to a rigorous mathematical formulation and finally, we'll elaborate a more complex example.

The next figure shows a simple example involving two trains. Train t\ :sub:`1` must depart in 5 minutes from *Station A* and heads to *Station B* while train t\ :sub:`2` must depart now from *Station C* to *Station D*. The running time of each track segment is shown below. As it can be seen, both trains require the *track segment 3* in the same period of time.

.. figure:: /_static/simple_network_introduction.jpg
   :alt: Small railway network example.
   
   Small network framework with two trains four stations and five track segments.

The next step is to build an *alternative graph* which models this situaltion. An **alternative graph** is a directed graph :math:`G=(N,E,A)` where nodes in :math:`N` are the operations that this train must perform and edges in :math:`E` indicate the order of such operations. We will talk about set :math:`A` later. For instance, next figure shows an (incomplete) alternative graph for this setup.

.. figure:: /_static/incomplete_alternative_graph_introduction.jpg
   :alt: Incomplete alternative graph.
   
   Incomplete alternative graph representing train operations and their running times.

We use the following notation:
   
   - Operations :math:`p_r`\ : Track segment :math:`r` is occupied by a train. In other words, there is a train running through segment :math:`r` with a running time of :math:`f_{p_r}`\ .
   - Operation :math:`p_B`\ : A train enters *Station B* with a running time of :math:`f_{p_B}` minutes.
   - Operation :math:`p_D`\ : A train enters *Station D* with a running time of :math:`f_{p_D}` minutes.

Additionally two *dummy* nodes are added which indicate the start time of the first operation of the alternative graph to be performed (node :math:`p_0`\) and end time when all operations have finished (node :math:`p_x`\ ). Let's take a close look to the alternative graph for train :math:`t_1`\ :

   - At time :math:`t=0`, train :math:`t_1` must wait *5* minutes before departure.
   - Then, at any time :math:`t\geq 5` the train must run through segment *1* in *12* minutes.
   - Then, at any time :math:`t\geq 18` the train must run through segment *3* in *12* minutes.
   - Then, at any time :math:`t\geq 30` the train must run through segment *4* in *8* minutes.
   - Then, at any time :math:`t\geq 38` the train must enter the station in *30* seconds.
   - Then, its trip is finised.

The same reasoning applies to train :math:`t_2`\ . Thus, if we define **optimisation variables** :math:`s_{t_1, p_1}` as the starting time of operation :math:`p_r`  performed by train :math:`t_i` and :math:`s_{p_0}` and :math:`s_{p_x}` as the starting times of *dummy* operations :math:`p_0` and :math:`p_x` respectively, it is straightforward to see that:

=============================================== ==============================================
Constraints for train :math:`t_1`               Constraints for train :math:`t_2`
=============================================== ==============================================
:math:`s_{t_1,p_1} \geq s_{p_0} + f_{t_1,p_0}`  :math:`s_{t_2,p_2} \geq s_{p_0} + f_{t_2,p_0}`
:math:`s_{t_1,p_3} \geq s_{t_1, p_1} + f_{p_1}` :math:`s_{t_2,p_3} \geq s_{t_2,p_2} + f_{p_2}`
:math:`s_{t_1,p_4} \geq s_{t_1, p_3} + f_{p_3}` :math:`s_{t_2,p_5} \geq s_{t_2,p_3} + f_{p_3}`
:math:`s_{t_1,p_B} \geq s_{t_1, p_4} + f_{p_4}` :math:`s_{t_2,p_D} \geq s_{t_2,p_5} + f_{p_5}`
:math:`s_{p_x} \geq s_{t_1, p_B} + f_{p_B}`     :math:`s_{p_x} \geq s_{t_2,p_D} + f_{p_D}`
=============================================== ==============================================

Notice, however, that there is a conflict between two trains since both demand the same operation :math:`p_3`. We must impose that either :math:`t_1` precedes :math:`t_2` or vice versa. For example, we may impose that :math:`t_2` will run through segment :math:`p_3` after :math:`t_1` plus an extra safety time :math:`f_{p_4, p_3}`, or vice versa. This behaviour is accomplished with the following constraints.

   - :math:`s_{t_2,p_3} \geq s_{t_1,p_4} + f_{p_4,p_3}` **or** :math:`s_{t_1,p_3} \geq s_{t_2,p_5} + f_{p_5,p_3}`, **but not both**.

Last disjunction constraints can be linearised by introducing a large positive number :math:`M` and a binary variable :math:`x_{t_1p_3,t_2p_3}` that determines if train :math:`t_1` must perform :math:`p_3` before :math:`t_2` performs :math:`p_3` (value *1*) or the opposite (value *0*).

   - :math:`s_{t_2,p_3}\geq s_{t_1,p_4} + f_{p_4,p_3} - M(1 - x_{t_1p_3,t_2p_3})`
   - :math:`s_{t_1,p_3}\geq s_{t_2,p_5} + f_{p_5,p_3} - M(1 - x_{t_1p_3,t_2p_3})`

This situation is represented in the alternative graph by adding a couple of edges as shown next.

.. figure:: /_static/complete_alternative_graph_introduction.jpg
   :alt: Complete alternative graph
   
   Complete alternative graph with a pair of alternative edges to solve the conflict.

This pair of edges of conflicting operations that we added in the figure are added to set :math:`A`. Thus, set :math:`A` is defined as *the set containing pairs of edges of conflicting operations*.

Imagine that instead of forcing one of the trains to stop in the middle of a track we prefer them to wait in the station. In this case, the alternative graph would be the following.

.. figure:: /_static/complete2_alternative_graph_introduction.jpg
   :alt: Complete alternative graph 2.
   
   Complete alternative graph with different consequences when solving the conflict.

Finally, our goal is to minimize the global makespan of the two trains. Then, the objective function *minimizes* :math:`s_{p_x} - s_{p_0}`\ .

There are some important remarks. First, notice, that this objective function does not consider any penalty term on the binary decision variable :math:`x_{t_1p_3,t_2p_3}`\ . If one of the two trains had higher priority, then the objective function should include a penalty term for this variable. Second, the proposed objective function may be too simple in more complex situations. For example, if there is any conflict in the network that causes train delays, it would be interesting to minimise schedule deviation times instead of the global makespan. Third, notice that until now the running times have been considered parameters. Instead, we could consider them as variables imposing that they should be greater than minimum track segment threshold.

.. note:: We deliberately want to keep the model simple in the first version of the project. Therefore, all these nice extra features will not be added immediately. The idea is that the optimisation model will grow in future releases.
