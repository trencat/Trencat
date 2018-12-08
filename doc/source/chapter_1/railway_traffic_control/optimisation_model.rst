.. _conflict-resolution-problem-model:

Optimization model
------------------

This section aims to give a rigorous mathematical formulation of the :term:`CRP`\ . Before going through the model, it is strongly recommend reading :ref:`conflict-resolution-problem` to motivate the problem and introduce notation with a simple example.

Sets and Parameters
^^^^^^^^^^^^^^^^^^^

Consider: 
   - :math:`T=\{t_1,t_2,\dots,t_n\}` the set of trains that are currently circulating or waiting in garages until the departure time.
   - :math:`P=\{p_1,p_2,\dots,p_m\}` the set of operations that must be performed by trains in :math:`T` with running times :math:`F = \{f_{p_1}, f_{p_2},\dots,f_{p_n}\}`.
   - Two additional *dummy* operations :math:`p_0` and :math:`p_x` to indicate the *start* and *stop* of all trains respectively. Each train :math:`t_i` has associated a waiting time :math:`f_{t_i,p_0}`\ .

Consider a train :math:`t_i` and two consecutive operations :math:`p_r` and :math:`p_s` that this train must perform. If :math:`p_s` is the next operation after :math:`p_r`\ , then we will be denote :math:`p_s:=\sigma_i(p_r)`. The next operation after :math:`p_s` is denoted as :math:`\sigma_i(p_s)=\sigma_i^2(p_r)`, and so on. Therefore, a train :math:`t_i` will perform operations :math:`P_i:=\{p_0,\sigma_i(p_0), \sigma_i^2(p_0),\dots, \sigma_i^k(p_0), ..., p_x\}` in the given order. Notice that all trains perform operations :math:`p_0` and :math:`p_x`.

Let :math:`G=(N,E,A)` be a directed graph where nodes in :math:`N` represent operations to be performed for each train, edges in :math:`E` indicate the order of the operations and their running times and :math:`A` contains pairs of alternative edges of conflicting operations. Such graph :math:`G` is called the **alternative graph**.

More specifically,

   - :math:`N := P\cup\{p_0\}\cup\{p_x\}`.
   - :math:`E := \cup_{t_i\in T} E_{t_i}`, where :math:`E_{t_i} := \cup_{k=0,\dots, \lvert P_i\lvert - 1} \{\sigma_i^k(p_0), \sigma_i^{k+1}(p_0)\}`.
   - :math:`A := \{(p_r,\sigma_i(p_r)),(p_s,\sigma_j(p_s)) | \text{for all trains $t_i, t_j\in T$ and for all operations $p_r\in P_i$ and $p_s\in P_j$ such that $p_r$ and $p_s$ are conflicting operations}\}`

Variables
^^^^^^^^^

   - :math:`s_{t_i, p_r}`: Start time of operation :math:`p_r\in P_i`. Defined for each train :math:`t_i\in T` and each operation :math:`p_j\in P_i`.
   - :math:`s_{p_0}`: Start time of operation :math:`p_0`.
   - :math:`s_{p_x}`: Start time of operation :math:`p_x`.
   - :math:`x_{t_ip_r, t_jp_s}`: Binary variable that takes the value *1* if train :math:`t_i` performs :math:`p_r`  before train :math:`t_j` performs :math:`p_s` and takes the value *0* otherwise. Defined for all pairs :math:`((p_r,\sigma_i(p_r)), (p_s,\sigma_j(p_s)))\in A`.

Objective function
^^^^^^^^^^^^^^^^^^
The objective function minimises the makespan of all train operations:

.. math::

   \min_{s, x} s_{p_n} - s_{p_0}

Constraints
^^^^^^^^^^^
Subject to the following constraints.

   - :math:`s_{\sigma_i(p_0)}\geq s_{p_0} + f_{t_i,p_0}`, for all train :math:`t_i\in T`.
   - :math:`s_{t_i,\sigma_i(p_r)} \geq s_{t_i,p_r} + f_{p_r}`, for all train :math:`t_i\in T` and for :math:`(p_r, \sigma_i(p_r))\in E_{t_i}`.
   - :math:`s_{t_i,\sigma_i(p_r)} \geq s_{t_i,p_r} + f_{p_r}`, for all train :math:`t_i\in T` and for all :math:`(p_r,\sigma_i(p_r))\in E_{t_i}`.
   - :math:`s_{t_j,p_s}\geq s_{t_i,\sigma_i(p_r)} + f_{\sigma_i(p_r),p_s}` **or** :math:`s_{t_i,p_r}\geq s_{t_j,\sigma_j(p_s)} + f_{\sigma_j(p_s),p_r}` for all :math:`((p_r,\sigma_i(p_r)), (p_s,\sigma_j(p_s)))\in A`.

Finally, the disjoint constraint can be linearsided by introducing binary variables and a large value :math:`M` as follows:

.. math::

   s_{t_j,p_s} \geq s_{t_i,\sigma_i(p_r)} + f_{\sigma_i(p_r), p_s} - M(1 - x_{t_ip_r, t_jp_s})

.. math::

   s_{t_i,p_r} \geq s_{t_j,\sigma_j(p_s)} + f_{\sigma_j(p_s), p_r} - M(1 - x_{t_ip_r, t_jp_s})

In all, we obtain a Mixed Integer Linear Programming Problem.
