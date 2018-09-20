.. _designing-railway-simulator:

###########################
Designing railway simulator
###########################

This chapter gives an insight of the internal structure of the simulator. It introduces the role of all modules and their interactions with no programming language specific implementation details.

.. _simulator-internal-structure:

****************
Modules overview
****************

*TrenCAT* is a simulator platform prepared for any student who wants to put into practice his/her own research in Train Automation Technologies. This assumption have an impact on software design. As a consequence, the software is divided in modules, each one with a specific task to carry out. Such modules are classified in a tight hierarchical structure according to their contribution, as shown in the following figure.

.. figure:: /_static/simulator_hierarchical_structure.jpg
   :alt: Simulator hierarchical structure.
   
   Simulator hierarchical structure.

Modules in the *real time* layer work 24/7 with real time generated data. Modules in the *near real time* layer are run periodically (for example, every hour) or on demand and operate with almost real time and historical data (for example, data from the last X hours until now).

High demanding modules, namely modules in the *real time* layer, run in different computers, whereas not-so-high-demanding modules, namely those in the *near real time* layer, may either run in the same computer or different computers. This design facilitates user custom module integration, as seen in section <X>.

Modules share data by communing (love this lemma!). Communication constantly flows *upwards* and *downwards* as depicted next.
   
.. figure:: /_static/simulator_modules_upwards_communications.jpg
   :alt: Upward communication between modules.
   
   Upward communication between modules.
   
.. figure:: /_static/simulator_modules_downwards_communications.jpg
   :alt: Downward communication between modules.
   
   Downward communication between modules.

There are many things going on here. Let's introduce them step by step.

***************
Modules details
***************

   
Automatic Train Supervision modules and Train 
=============================================

A Train module simulates a real train. This Train, however, will never move unless someone drives it. The :term:`ATS` module is in charge of driving a Train between two points within safety limits. The next figure depicts the interaction between the :term:`ATS` and the Train.
   
.. figure:: /_static/ato_to_train.jpg
   :alt: Communication between :term:`ATO` and Train.
   
   Communication between :term:`ATO` and Train.

The :term:`ATS` establishes a high frequency closed loop feedback with the Train: it requests current driving information and responds with a decision such as *accelerate* or *break*.

A Train stores and makes accessible all the necessary data that is needed to drive it:

   - **Train static information:** mass, number of units that it consists of, length, maximum acceleration/brake rates, etc.
   - **Dynamic information**:
  
      - **Train operation:** self position, velocity and acceleration.
      - **Infrastructure:** scheduled travel time between two stations, distance to the next station, maximum velocity at each segment, distance to the next semaphore and its current signal, track slope, curve bend radius, distance to the next tunnel (and its length), etc.	 
      - **Ecosystem**: Any data of interest such as other trains positions, velocities and accelerations.
	  
When the train is in motion, the :term:`ATS` requests detailed data to the Train. With all this information, it must decide whether the train must accelerate, coast, break or trigger the emergency break. Once decided, the :term:`ATS` sends a setpoint to the Train, which will perform the requested action if its feasible. The Train will then update the internal data, which will be requested again by the :term:`ATS` to decide next setpoint. And so on.

Additionally, the Train integrates :term:`ATP` (Automatic Train Protection), a safety layer which will take over control in some specific bad scenarios such as breaking safety boundaries (red signal overrun, maximum speed exceeding, etc), loosing connection with :term:`ATS` or :term:`ATS` shutdown requested from higher priority modules.

Finally, the Train got safely to a station and stopped. Now the :term:`ATS` waits until Train confirmation. Meanwhile, the Train will communicate with the Train Manager to update the number of passengers that get on and off the train, and send some collected statistics to the Train Manager.

In future releases the train may implement other modules which may be accessible via requests to the Train process.

.. note::
	The :term:`ATS` module will run in its own process and users are free to implement their own :term:`ATS` in any programming language as long as it sticks with the protocol. Users are encouraged to use *TrenCAT* as a platform to test their :term:`ATS` implementations and see how such implementations react to many different simulated scenarios. By default, *TrenCAT* implements an :term:`ATS` based on chapter :ref:`speed-profile-optimization`.


	
Train and The Train Manager
===========================

The Train Manager is a key process in the *TrenCAT* infrastructure as it controls and orchestrates the entire ecosystem.

   - It receives information in real time of each train: position, velocity, acceleration, statuses (stopped, running, how much people each train is carrying, etc).
   - It monitors the status and events of each station, for example, how many people are there in the platforms.
   - It monitors the schedule of all trains from the :term:`RTC` and detects when a train is delayed or being delayed. When necessary, it balances the network by requesting new reschedules and rolling stock updates to the :term:`RTC` and :term:`RSP` and propagates the response to all trains.
   - It determines how many passengers step in and out of a train at each station according to the decisions taken by the Demand Module.
   - It controls railway infrastructure according to real time information by setting semaphore signals.
   - It detects train collisions. In case a collision occurs, the Train Manager modifies railway semaphore signals to prevent more collisions. It also notifies other trains about the incident.

Finally, the Train Manager stores periodically historical data in a database for later use and future analysis.
   
Next figures depict the communications that take place between each Train and the Train Manager.

.. figure:: /_static/train_to_train_manager.jpg
   :alt: Communication from the Train to the Train Manager.
   
   Communication from the Train to the Train Manager.

.. figure:: /_static/train_manager_to_train.jpg
   :alt: Communication from the Train Manager to the Train.
   
   Communication from the Train Manager to the Train.
   

Demand Module and Train Manager
===============================

The Demand Manager simulates the demand of people in a railway system. In practice, it introduces new people into train stations together with the root that each one must follow in real time. Sophisticated Demand Managers may require infrastructure information to decide people's routing and current real time information to model advanced scenarios. For instance, in a metro railway system, when some trains halt due to another train break down people may decide reroute their trip, which may increase passengers demand in other lines and stations. The next figure depicts the communications between the Demand Manager and the Train Manager.

.. figure:: /_static/demand_manager_to_train.jpg
   :alt: Communication between the Demand Manager and the Train Manager.
   
   Communication between the Demand Manager and the Train Manager.

.. note::
	The Demand Manager will run in its own process and users are free to implement their own manager in any programming language as long as it sticks with the protocol. Users are encouraged to use *TrenCAT* as a platform to test their Demand Manager implementations and see how such implementations react to many different simulated scenarios. Currently *TrenCAT* has not started designing this manager yet.

   
Railway Traffic Control, Rolling Stock Planning and :term:`SCADA`
=================================================================

The :term:`RTC` and :term:`RSP` are in charge of computing train timetables and rolling stock respectively for the next hours according to current and future estimated demand. To do so, they take near real time input data from the database for its calculations. The optimal reschedule and rolling stock are sent back to the Train Manager, which works to accomplish the new orders. In addition, the Train Manager stores the new orders in the historical database for later use and future analysis.

In parallel, the :term:`SCADA` is responsible for displaying real time data (from the Train Manager) and historical data (from the database) in a nice user interface. The :term:`SCADA` allows to...
   - Monitor the entire infrastructure, showing information reported by trains, the train manager, the demand manager, the :term:`RTC` and :term:`RSP`.
   - Control manually the entire infrastructure, allowing to set and monitor any kind of testing scenarios.

Analogously to the :term:`RTC` and :term:`RSP` orders, user actions taken via the :term:`SCADA` interface are communicated to the Train Manager, which works to accomplish the new orders. Additionally, the Train Manager stores the new orders in the historical database for later use and future analysis.

The following two figures briefly depict the communications that take place between the Operaing Control Center, the :term:`SCADA`, the Train Manager and the database.

.. figure:: /_static/train_manager_to_near_real_time.jpg
   :alt: Communication from the Train Manager to all processes in the near real time layer.
   
   Communication from the Train Manager to all processes in the near real time layer.
   
.. figure:: /_static/near_real_time_to_train_manager.jpg
   :alt: Communication from the near real time layer to the Train Manager.
   
   Communication from the near real time layer to the Train Manager.

.. note::
	The :term:`RTC` and :term:`RSP` modules will run in their own processes (either periodically and on-demand) and users are free to implement their modules in any programming language as long as they sticks with the protocol. Users are encouraged to use *TrenCAT* as a platform to test their implementations and see how they react to many different simulated scenarios. By default, *TrenCAT* implements a :term:`RTC` based on chapter :ref:`railway-traffic-control` and a :term:`RSP` based on chapter :ref:`optimal-rolling-stock-planning`.


Previous topic: :ref:`introduction-railway-infrastructure-design-theory`.