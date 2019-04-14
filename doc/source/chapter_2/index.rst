.. _designing-railway-simulator:

###########################
Designing railway simulator
###########################

This chapter gives an insight of the internal structure of the simulator. It introduces the role of all modules and their interactions with no programming language specific implementation details.

.. _simulator-internal-structure:

****************
Modules overview
****************

*TRENCAT* is a simulator platform to carry out custom research in Train Automation Technologies. The software is divided in modules carrying out unique tasks. Such modules can then be substituted by custom modules with personal research. Thanks to task modularity, custom modules can be developed with virtually any programming language. Modules are classified in a tight hierarchical structure according to their objective, as shown in the following figure.

.. figure:: /_static/simulator_hierarchical_structure.jpg
   :alt: Simulator hierarchical structure.
   
   Simulator hierarchical structure.

Modules in the *real time* layer work 24/7 with real time generated data. They are high demanding modules and may run in different computers. Modules in the *near real time* layer are run periodically (for example, every hour) or on demand and operate with almost real time data and historical data (for example, data from the last 10 hours).

Modules share data by communicating (love this lemma!). The next figure depicts briefly the communication flow between modules.
   
.. figure:: /_static/simulator_modules_communications.jpg
   :alt: Upward communication between modules.
   
   Communication between modules.
   
In the following sections de modules and communications are described in detail.

***************
Modules details
***************

.. _designing_railway_simulator_ato:

Automatic Train Operation
===========================

The objective of the :term:`ATO` module is to drive the train between two points in safety conditions. To perform this task, the following information is needed:

   - **Static information**:
      - **Train static information**: mass, length, maximum traction force, maximum braking force and train-specific coefficients.
      - **Railway infrastructure**: distance between origin and destinations, maximum velocity at each segment, segment slopes, segment bend radii, existence of tunnels (and their lengths) between origin and destination, specific track related coefficients, location of semaphores.
   - **Dynamic information**:
      - **Current driving state**: Train position, velocity and acceleration. Updated train mass taking into account the number of passengers.
      - **Currant railway state**: Status of the following semaphores. It may also use the position, velocity and acceleration of (former) nearby trains.
      - **Schedule**: Scheduled time between origin and destination.
      - **Alerts** broadcasted by the Control Centre.

With these information the :term:`ATO` computes a *driving curve plan* between the two points as explained in :ref:`real-time-train-operation`. Next, the :term:`ATO` establishes a high frequency closed loop feedback with the onboard :ref:`designing_railway_simulator_train` computer.

   1.  It reads/requests current driving information.
   2.  It responds with a setpoint to follow. Such setpoint indicates the traction force to be done by the engine or the braking force to be applied to the brakes.

During the trip, if the train deviates significantly from the *driving curve plan*, a new *plan* is quickly computed on the fly.

The :term:`ATO` should implement security measures (:term:`ATP`) such as:

   1.  If the train overruns a semaphore red signal, trigger the emergency brakes and stop the train completely.
   2.  Do not allow the train to exceed the segment speed limits. However, if it does exceed it, trigger the service brakes to reduce velocity.
   3.  If there is a train at close distance (up to a certain threshold depending on current velocity), trigger the emergency brakes and stop the train completely.
   4.  If the train deviates significantly from the *driving curve plan* and the computation *on the fly* of a new one is taking more time than expected, trigger the service brake and stop the train. Once the train has stopped, compute a new *driving curve* from the current position and continue the trip.

.. _designing_railway_simulator_train:

Train
=====

The train module simulates the train onboard computer.

   - It stores and provides all the data that :ref:`designing_railway_simulator_ato` needs to perform its calculations.
   - It receives the setpoints from :ref:`designing_railway_simulator_ato` and applies them, if possible, i.e. it sends the setpoints to the engine.
   - A builtin :term:`ATP` module takes care of safety and will take over control in the following situations:
    
      1.  If the train overruns a semaphore red signal, trigger the emergency brakes and stop the train completely.
      2.  If the train exceeds the segment speed limits, trigger the emergency brajes and stop the train completely.
      3.  If there is a train at close distance (up to a certain threshold depending on current velocity), trigger the emergency brakes and stop the train completely.
      4.  If no feedback has been received from the :ref:`designing_railway_simulator_ato` within a certain amount of milliseconds, trigger the emergency brakes and stop the train completely.
      5.  If the train cannot update dynamic infrastructure information within a certain amount of seconds, trigger the emergency brakes and stop the train completely.
   - During idle periods, it computes statistics and :term:`KPI`\s, compresses data, sends it to a database.
   - The train also sends heartbeat signals to the Control Centre to check if the connection is still alive. The Control Centre will broadcast an alert to all trains if any train looses connection with it.

.. _designing_railway_simulator_control_centre:

Control Centre
==============

The :ref:`designing_railway_simulator_control_centre` is the core of the simulator. It controls and orchestrates the entire ecosystem.

   - It receives real time information from all trains: train positions, velocities, accelerations, statuses (either if the train is stopped, running, etc) and any other data of interest (how much people each train is carrying, etc).
   - It controls railway infrastructure according to real time information by setting semaphore signals.
   - It detects train collisions. In case a collision occurs, the Control Centre sets railway semaphore signals to prevent more collisions. It also broadcasts alerts to other trains warning about the incident.
   - It monitors the schedule of all trains from the :ref:`designing_railway_simulator_rtc_rsp` and detects when a train is delayed or being delayed. When necessary, it balances the network by requesting new reschedules and rolling stock updates to the :ref:`designing_railway_simulator_rtc_rsp`.
   - It boots up and shuts down new trains remotely according to the :term:`RSP`.
   - It coordinates the status of each platform, for example, how many people are there in the platforms, how many of them step in and out of the train according to the :ref:`designing_railway_simulator_demand`, etc.
   - It performs the actions set by the user via :ref:`designing_railway_simulator_scada`.
   - It computes statistics and :term:`KPI`\s, compresses data, and sends it to a database.

.. _designing_railway_simulator_rtc_rsp:

Railway Traffic Control and Rolling Stock Planning
==================================================

The :term:`RTC` and :term:`RSP` are in charge of computing train timetables and rolling stock respectively for the next hours according to current and future estimated demand. To do so, they take near real time input data from the database for its calculations. The optimal reschedule and rolling stock are sent back to the Control Centre, which works to accomplish the new orders.


.. _designing_railway_simulator_demand:

Demand module
=============

The :ref:`designing_railway_simulator_demand` simulates the demand of people in a railway system and serves as a tool to test automatic infrastructure demand response. In practice, it introduces new people into train stations together with the root that each one must follow in real time. Sophisticated Demand Managers may require infrastructure information to decide people’s routing and current real time information to model advanced scenarios. For instance, in a metro railway system, when some trains halt due to another train break down people may decide reroute their trip, which may increase passengers demand in other lines and stations.

.. _designing_railway_simulator_scada:

SCADA
=====

The :term:`SCADA` provides the end user with a human interface to monitor, supervise and have full manual control of the entire ecosystem.
