.. _conflict-resolution-problem:

Introducing conflict resolution problem
---------------------------------------

.. role:: raw-html(raw)
   :format: html

The aim of the :term:`CRP` is to develop a conflict free schedule such that the overall deviation from the planned schedule is minimized. The problem is addressed under the `Operations Research <https://en.wikipedia.org/wiki/Operations_research>`_ point of view. Let's first introduce the model by means of an example, then we'll move to a rigorous mathematical formulation and finally, we'll elaborate a more complex example.

The next figure shows a simple example involving two trains. Train t\ :sub:`1` must depart in 5 minutes from *Station A* and heads to *Station B* while train t\ :sub:`2` must depart now from *Station C* to *Station D*. The running time of each track segment is shown below. As it can be seen, both trains require the *track segment 3* in the same period of time.

.. figure:: /_static/simple_network_introduction.jpg
   :alt: Small railway network example.
   
   Small network framework with two trains four stations and five track segments.

The next step is to build an *alternative graph* which models this situaltion. An **alternative graph** is a directed graph *G=(N,E,A)* where nodes in *N* are the operations that this train must perform and edges in *E* indicate the order of such operations. We will talk about *A* later. For instance, next figure shows an (incomplete) alternative graph for this setup.

.. figure:: /_static/incomplete_alternative_graph_introduction.jpg
   :alt: Incomplete alternative graph.
   
   Incomplete alternative graph representing train operations and their running times.

We use the following notation:
   
   - Operations p\ :sub:`r`\ : Track segment *r* is occupied by a train. In other words, there is a train running through segment *r* with a running time of f\ :sub:`r`\ .
   - Operation p\ :sub:`B`\ : A train enters *Station B* with a running time of f\ :sub:`B` minutes.
   - Operation p\ :sub:`D`\ : A train enters *Station D* with a running time of f\ :sub:`D` minutes.

Additionally two *dummy* nodes are added which indicate the start time of the first operation of the alternative graph to be performed (node p\ :sub:`0`\) and end time when all operations have finished (node p\ :sub:`x`\ ). Let's take a close look to the alternative graph for train t\ :sub:`1`\ :

   - At time *t=0*, train t\ :sub:`1` must wait *5* minutes before departure.
   - Then, at any time *t* :raw-html:`&ge;`\ *5* the train must run through segment *1* in *12* minutes.
   - Then, at any time *t* :raw-html:`&ge;`\ *18* the train must run through segment *3* in *12* minutes.
   - Then, at any time *t* :raw-html:`&ge;`\ *30* the train must run through segment *4* in *8* minutes.
   - Then, at any time *t* :raw-html:`&ge;`\ *38* the train must enter the station in *30* seconds.
   - Then, its trip is finised.

The same reasoning applies for train t\ :sub:`2`\ . Thus, if we define **optimisation variables** s\ :raw-html:`<sub>t<sub>i</sub>,p<sub>r</sub></sub>` as the starting time of operation p\ :sub:`r`  performed by train t\ :sub:`i` and s\ :raw-html:`<sub>p<sub>0</sub></sub>` and s\ :raw-html:`<sub>p<sub>x</sub></sub>` as the starting times of *dummy* operations p\ :sub:`0` and p\ :sub:`x` respectively, it is straightforward to see that:

============================================================================================================================================================ ============================================================================================================================================================
Constraints for train t\ :sub:`1`                                                                                                                            Constraints for train t\ :sub:`2`                                            
============================================================================================================================================================ ============================================================================================================================================================
s\ :raw-html:`<sub>t<sub>1</sub>,p<sub>1</sub></sub> &ge;` s\ :raw-html:`<sub>p<sub>0</sub></sub>` +  f\ :raw-html:`<sub>t<sub>1</sub>,p<sub>0</sub></sub>`  s\ :raw-html:`<sub>t<sub>2</sub>,p<sub>2</sub></sub> &ge;` s\ :raw-html:`<sub>p<sub>0</sub></sub>` +  f\ :raw-html:`<sub>t<sub>2</sub>,p<sub>0</sub></sub>` 
s\ :raw-html:`<sub>t<sub>1</sub>,p<sub>3</sub></sub> &ge;` s\ :raw-html:`<sub>t<sub>1</sub>,p<sub>1</sub></sub>` +  f\ :raw-html:`<sub>p<sub>1</sub></sub>`  s\ :raw-html:`<sub>t<sub>2</sub>,p<sub>3</sub></sub> &ge;` s\ :raw-html:`<sub>t<sub>2</sub>,p<sub>2</sub></sub>` +  f\ :raw-html:`<sub>p<sub>2</sub></sub>` 
s\ :raw-html:`<sub>t<sub>1</sub>,p<sub>4</sub></sub> &ge;` s\ :raw-html:`<sub>t<sub>1</sub>,p<sub>3</sub></sub>` +  f\ :raw-html:`<sub>p<sub>3</sub></sub>`  s\ :raw-html:`<sub>t<sub>2</sub>,p<sub>5</sub></sub> &ge;` s\ :raw-html:`<sub>t<sub>2</sub>,p<sub>3</sub></sub>` +  f\ :raw-html:`<sub>p<sub>3</sub></sub>` 
s\ :raw-html:`<sub>t<sub>1</sub>,p<sub>B</sub></sub> &ge;` s\ :raw-html:`<sub>t<sub>1</sub>,p<sub>4</sub></sub>` +  f\ :raw-html:`<sub>p<sub>4</sub></sub>`  s\ :raw-html:`<sub>t<sub>2</sub>,p<sub>D</sub></sub> &ge;` s\ :raw-html:`<sub>t<sub>2</sub>,p<sub>5</sub></sub>` +  f\ :raw-html:`<sub>p<sub>5</sub></sub>` 
s\ :raw-html:`<sub>p<sub>x</sub></sub> &ge;` s\ :raw-html:`<sub>t<sub>1</sub>,p<sub>B</sub></sub>` +  f\ :raw-html:`<sub>p<sub>B</sub></sub>`                s\ :raw-html:`<sub>p<sub>x</sub></sub> &ge;` s\ :raw-html:`<sub>t<sub>2</sub>,p<sub>D</sub></sub>` +  f\ :raw-html:`<sub>p<sub>D</sub></sub>`
============================================================================================================================================================ ============================================================================================================================================================

Notice, however, that there is a conflict between two trains since both demand the same operation p\ :sub:`3`. If trains circulate according to schedule, the train t\ :sub:`2` is the first one to enter segment *3* and train t\ :sub:`1` will have to wait three minutes until the former one leaves the segment, plus an extra safety time. We must impose that either t\ :sub:`1` precedes t\ :sub:`2` or the opposite. This is accomplished with the following constraints.

   - s\ :raw-html:`<sub>t<sub>2</sub>,p<sub>3</sub></sub>` :raw-html:`&ge;` s\ :raw-html:`<sub>t<sub>1</sub>,p<sub>4</sub></sub>` + f\ :raw-html:`<sub>p<sub>4</sub>,p<sub>3</sub></sub>` **or** s\ :raw-html:`<sub>t<sub>1</sub>,p<sub>3</sub></sub>` :raw-html:`&ge;` s\ :raw-html:`<sub>t<sub>2</sub>,p<sub>5</sub></sub>` + f\ :raw-html:`<sub>p<sub>5</sub>,p<sub>3</sub></sub>`\ , **but not both**.

Last disjunction constraints can be linearised by introducing a large positive number *M* and a binary variable x\ :raw-html:`<sub>t<sub>1</sub>p<sub>3</sub>,t<sub>2</sub>p<sub>3</sub></sub>` that determines if train t\ :sub:`1` must perform p\ :sub:`3` before t\ :sub:`2` performs p\ :sub:`3` (value *1*) or the opposite (value *0*).

   - s\ :raw-html:`<sub>t<sub>2</sub>,p<sub>3</sub></sub>` :raw-html:`&ge;` s\ :raw-html:`<sub>t<sub>1</sub>,p<sub>4</sub></sub>` + f\ :raw-html:`<sub>p<sub>4</sub>,p<sub>3</sub></sub>` - M (1 - x\ :raw-html:`<sub>t<sub>1</sub>p<sub>3</sub>,t<sub>2</sub>p<sub>3</sub></sub>`)
   - s\ :raw-html:`<sub>t<sub>1</sub>,p<sub>3</sub></sub>` :raw-html:`&ge;` s\ :raw-html:`<sub>t<sub>2</sub>,p<sub>5</sub></sub>` + f\ :raw-html:`<sub>p<sub>5</sub>,p<sub>3</sub></sub>` - M x\ :raw-html:`<sub>t<sub>1</sub>p<sub>3</sub>,t<sub>2</sub>p<sub>3</sub></sub>`

This situation is represented in the alternative graph by adding a couple of edges as shown next.

.. figure:: /_static/complete_alternative_graph_introduction.jpg
   :alt: Complete alternative graph
   
   Complete alternative graph with a pair of alternative edges to solve the conflict.

This pair of edges of conflicting operations that we added in the figure are added to set *A*. Thus, set *A* is defined as *the set containing pairs of edges of conflicting operations*. Notice that this pair of edges must be added even if the two trains do not coincide in time in track segment *3*.

Imagine that instead of forcing one of the trains to stop in the middle of the track we prefer them to wait in the station. In this case, the alternative graph would be the following.

.. figure:: /_static/complete2_alternative_graph_introduction.jpg
   :alt: Complete alternative graph 2.
   
   Complete alternative graph with different consequences when solving the conflict.

Finally, our goal is to minimize the global makespan of the two trains. Then, the objective function *minimizes* s\ :raw-html:`<sub>p<sub>x</sub></sub>` - s\ :raw-html:`<sub>p<sub>0</sub></sub>`\ .

There are some important remarks. First, notice, that this objective function does not consider any penalty term on the binary decision variable x\ :raw-html:`<sub>t<sub>1</sub>p<sub>3</sub>,t<sub>2</sub>p<sub>3</sub></sub>`\ . If one of the two trains had higher priority, then the objective function should include a penalty term with x\ :raw-html:`<sub>t<sub>1</sub>p<sub>3</sub>,t<sub>2</sub>p<sub>3</sub></sub>`\ . Second, this objective function may be too simple in more complex situations. For example, if there is any conflict in the network that causes train delays, it would be interesting to minimise schedule deviation times instead of the global makespan. Third, notice that until now the running times have been considered parameters. Instead, we could consider them as variables imposing that they should be greater than minimum track segment threshold.

.. note:: We deliberately want to keep the model simple in the first version of the project. Therefore, all these nice extra features will not be added immediately. The idea is that the optimisation model will grow in future releases.

Previous topic: :ref:`railway-traffic-control-definitions`.

Next topic: :ref:`conflict-resolution-problem-model`.
