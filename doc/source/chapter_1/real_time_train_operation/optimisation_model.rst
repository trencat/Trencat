.. _speed-profile-optimization-model:

Optimization model
------------------

This section aims to give a rigorous mathematical formulation of the :term:`SPO`\ . Before going through the model, it is strongly recommended to read :ref:`speed-profile-optimization` and :ref:`speed-profile-optimization-computations-constant-acceleration` to motivate the problem and introduce the model notation with a simple example.

The problem is formulated as an optimal network flow problem with few additional constraints.

Sets and Parameters
^^^^^^^^^^^^^^^^^^^

The following figure depicts the model notation that will be used.

.. figure:: /_static/speed_profile_optimization_notation.jpg
   :alt: Figure supporting model notation.
   
   Model notation depicted.

Track parameters:
   - :math:`L \in \mathbb{R}`: Length of the track [m].  
   - :math:`S=\{s_0, s_1, \dots, s_n\}`: Set of positions [m] that define the segments of the track, i.e.
   
      .. math::
         [0, L] = \bigcup_{i=0}^{n-1} [s_i, s_{i+1}) \cup \{L\},

      where :math:`s_0 = 0` and :math:`s_n=L`. It is assumed that 

      .. math::
         s_0 \lneq s_1 \lneq \dots \lneq s_n.

   - :math:`\underline{v_i}, \overline{v_i}`: Minimum and maximum velocities [m/s] of segment :math:`[s_i, s_{i+1})`, for :math:`i\in \{0,\dots, n-1\}`.

   - :math:`V_i`: Finite discrete set of velocities [m/s] that the train can achieve at position :math:`s_i`. :math:`V_0 = \{v_{start}\}`, :math:`V_n = \{v_{end}\}` and :math:`V_i=\{v_i^0, v_i^1, \dots, v_i^{m_i}\}` for :math:`i\in\{1, 2, \dots, n-1\}`. More precisely,

      .. math::
         \underline{v_i} = v_i^0 \lneq v_i^1 \lneq \dots \lneq v_i^{m_i}, \text{ where}

      .. math::
         v_i^{m_i} = \left\{\begin{array}{rl}
         \overline{v_i} & \text{if } \overline{v_{i+1}} \gneq \overline{v_i},\\
         \overline{v_i} & \text{otherwise}.\\
         \end{array}\right.

   - :math:`\beta_i`: Slope [radians] of segment :math:`[s_i, s_{i+1})`, for :math:`i\in \{0,\dots, n-1\}`.
   - :math:`r_i`: Bend radius [m] of segment :math:`[s_i, s_{i+1})`, for :math:`i\in \{0,\dots, n-1\}`.
   - :math:`\text{tunnel}_i \in \{0, 1\}`: Binary parameter that equals :math:`1` if the segment is built inside a tunnel and :math:`0` otherwise.

Train parameters (check :ref:`train-dynamic-models` for further information):
   - :math:`m \in \mathbb{R}`: Train mass [kg].
   - :math:`\rho \in \mathbb{R}`: Rotating mass factor.
   - :math:`u_{max} \in \mathbb{R}`: Maximum traction force [N].
   - :math:`u_{min} \in \mathbb{R}`: Maximum braking force in absolute value [N]:
   - :math:`a_1, a_2 \in \mathbb{R}`: Parameters of the basic resistance :math:`R_b(v)`.

Profile parameters:
   - :math:`T_{min}, T_{max} \in \mathbb{R}`: Minimum and maximum travel time allowed [s].
   - :math:`v_{start}, v_{end}`: Start and end velocities [m/s] of the train at the beginning and end of the track respectively.


**Network flow parameters:**

Consider a graph :math:`G(V, E)`, where :math:`V` is the set of vertices and :math:`E` is the set of edges. The set of vertices is represented by all states of the train, i.e. all possible positions and velocities of the train.

   .. math::
      V := \bigcup_{i=0}^n V_i,

:math:`E` is the set of all transitions between states. Each link represents the transition of the train between two states. More precisely,

   .. math::
      E := \bigcup_{i=0}^{n-1} E_i,\quad E_i := V_i \times V_{i+1},\quad E_n := \emptyset

where :math:`\times` means the cartesian product of the two sets. Let us denote an element of this set as :math:`e_i^{jk} := (v_i^j, v_{i+1}^k)`.

   - :math:`b_i^j`: Input/output flow of node :math:`v_i^j`, for :math:`i\in\{0,\dots, n\},\ j\in \{0,\dots, m_i\}`. More precisely, :math:`b_0^0 = 1`, :math:`b_n^0 = -1` and :math:`b_i^j = 0` otherwise.
   - :math:`c_i^{jk}`: Cost of the transition from :math:`v_i^j` to :math:`v_{i+1}^k`, for each edge :math:`i\in \{0, \dots, n-1\}` and for each edge :math:`e_i^{jk}\in E_i`. We can consider this cost to be the travel time between the two states [s], the (traction) work that the train must do to reach the next state [J], the jerk rate in this segment, etc. In this case, we will consider the traction work [J] as computed in :ref:`speed-profile-optimization-computations-constant-acceleration`. Such calculation involves all train parameters, slope, bend radius and tunnel parameters among others.
   - :math:`t_i^{jk}`: Time span [s] of the transition from :math:`v_i^j` to :math:`v_{i+1}^k`, for each :math:`i\in \{0, \dots, n-1\}` and for each edge :math:`e_i^{jk}\in E_i`. Refer to :ref:`speed-profile-optimization-computations-constant-acceleration` to see how it is computed.

Variables
^^^^^^^^^

   - :math:`x_i^{jk}`: Binary variable that equals one if the train reaches state :math:`v_{i+1}^k` from state :math:`v_i^k`, and 0 otherwise, for each :math:`i\in \{0, \dots, n-1\}` and for each edge :math:`e_i^{jk}\in E_i`.

Objective function
^^^^^^^^^^^^^^^^^^

The objective function minimises the total work done by the traction force.

   .. math::
      \min_x\sum_{i=0}^{n-1} \sum_{e_i^{jk}\in E_i}c_i^{jk}x_i^{jk}

Constraints
^^^^^^^^^^^

   - **Node equilibrium** (network flow problem): The flow entering a node must be equal to the flow leaving this node

      .. math::
         \sum_{e_{i}^{jk}\in E_i} x_{i}^{jk} - \sum_{e_{i-1}^{hj}\in E_{i-1}} x_{i-1}^{hj} = b_i^j, \quad \forall i\in\{1,2,\dots,n-1\}\text{ and }\forall j\text{ such that }v_i^j\in V_i.

      For convenience, consider :math:`E_{-1} = E_n := \emptyset`.

   - **Punctuality constraints:** The travel time is lower and upper bounded.

      .. math::
         T_{min}\leq\sum_{i=0}^{n-1} \sum_{e_i^{jk}\in E_i}t_i^{jk}x_i^{jk} \leq T_{max}

   - **Feasibility constraint:** If the train cannot drive from :math:`v_i^j` to :math:`v_{i+1}^k` due to physical limitations (for example, the traction force needed exceeds the maximum traction force), then the corresponding variable is fixed to zero.

      .. math::
         x_i^{jk} = 0 \text{ if the train cannot physically reach } v_{i+1}^k \text{ from } v_i^k,\ \forall i\in\{0,1, \dots, n-1\}, \forall e_i^{jk}\in E_i.

      .. note::
         For better performance, it is recomended to not create the variable and the edge while building the graph and the optimization model.

   - **Velocity constraint**: The train velocity cannot exceed


In all, we obtain a Mixed Integer Linear Programming Problem.

Avoiding infeasibility
^^^^^^^^^^^^^^^^^^^^^^

If the travel time window :math:`[T_{min}, T_{max}]` is too early, the train will not be physically able to arrive on time and the problem becomes infeasible. In this case, we desire the train to drive as fast as allowed (and as fast as the engine can get) to the next station.
Analogously, if the time window is too far and the train cannot drive that slow, the problem becomes also feasible. In this case, we want the train to drive as slow as possible.

We accomplish both behaviours by introducing an extra variable

   - :math:`\delta\in \mathbb{R}`: Non negative time deviation.

First we modify the punctuality constraint to give some flexibility to the train schedule

   .. math::
      T_{min} - \delta \leq\sum_{e_i^{jk}\in E_i}t_i^{jk}x_i^{jk} \leq T_{max} + \delta,

Finally, :math:`\delta` comes with a large penalty :math:`M` in the objective function

   .. math::
      \min_x\sum_{i=0}^{n-1} \sum_{e_i^{jk}\in E_i}c_i^{jk}x_i^{jk} + M\delta.
