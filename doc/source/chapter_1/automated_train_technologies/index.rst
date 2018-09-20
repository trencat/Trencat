.. _introduction-att-technologies:

*********************************************
Automated Train Techonolgies: An introduction
*********************************************

Definitions and motivations
===========================

Traditionally, trains were driven manually by train drivers. As soon as wailway transport grew in popularity, authorities realised the importance of train coordination since more and more trains were circulating and sharing the same tracks. At the beginning signalling was based in track side light signals and even human body language, which was inherently unefficient. In addition, train circulation was heavily dependent on driver knowledge, skills and expertise, which was also really unefficient. Sadly, lots of accidents took place due to poor communication technology and a null safety infrastructure. Check [C]_ for a history of railway railway signalling.  

Nowadays, train public transport has achieved an impressive level of safety, punctuality, speed and energy efficiency due to automation. At the time of writing, the `L’Union Internationale des Transports Publics <http://www.uitp.org/>`_ (International Association of Public Transport, :term:`UITP`) defines automation in metro systems as follows [UITP-PK]_:

	"In metro systems, automation refers to the process by which responsibility for operation management of the trains is transferred from the driver to the train control system."

In other words, trains are able to be driven by computer software without the presence of a train driver. In addition, :term:`UITP` defines various *degrees of automation* (or *Grades of Automation*, :term:`GoA`) according to which basic functions of train operations are the responsibility of staff, and which are the responsibility of the system itself [UITP-PK]_.

- **GoA0**: It corresponds to on-sight train operation, similar to tram running in street traffic.
- **GoA1**: It corresponds to manual train operation where a train driver controls starting and stopping, operation of doors and handling of emergencies or sudden diversions.
- **GoA2**: It corresponds to *semi-automatic train operation* (:term:`STO`) where starting and stopping is automated, but a driver operates the doors, drives the train if needed and handles emergencies.
- **GoA3**: It corresponds to *driverless train operation* (:term:`DTO`) where starting and stopping are automated but a train attendant operates the doors and drives the train in case of emergencies.
- **GoA4**: It corresponds to *unattended train operation* (:term:`UTO`) where starting and stopping, operation of doors and handling of emergencies are fully automated without any on train staff.

Train automation brings many benefits.

- **Security improvement**: Train automation implement new active and passive security measures in all facilities and trains. For instance, stations are provided with *platform screen doors* (:term:`PSDs`) that separate physically the platform and the track and protect passengers from accidentally falling down into the track area. In addition, automated systems reduce significantly human risk factors (due to driver’s psychological condition, extreme weather that may affect the sights of drivers, etc). Well-designed :term:`UTO` systems have proven to be more reliable than conventional metros [UITP-PK]_.

- **Better quality of service**: All automated trains work synchronously to offer better quality of service. As a consequence, the same rolling stock is able to transport more passengers in less time with full security. This translates to shorter waiting times and reliability and punctuality enhancement. In addition, increasing demand in peak hours or big events can be quickly met by injecting more trains. In fact, automated metro becomes affordable for smaller cities: when trains run more frequently, the system does not need to be oversized to cope with peak demand.

- **Energy efficiency**: Energy-efficient train operation is critical as the rising of energy prices and environmental concerns. In particular, the train operations account for about 80% in the whole energy consumption of metro systems. What’s more, optimised train driving strategies could reduce as much as 20% energy consumption [YTYXHG]_. Not only metro becomes more environmentally friendly, but also reduces operating expenses (:term:`OPEX`).

Structure
=========

This section is mainly based on [WNBS]_ and [YTYXHG]_\ .

A railway system essentially consists of three elements: **infrastructure** (line tracks, stations, signalling equipments, etc.), a **rolling stock** of trains circulating on tracks and the **operation rules** taking care of safety and operation efficiency. Additionally, railway systems can be classified in two types: **interurban** and **urban** systems. In interurban systems, trains share a limited resource of tracks and line overlaps and trains usually overtake and meet each other. In urban systems lines are not so scarce since tracks are separated from each other and each direction of the line has a dedicated infrastructure. This project focuses on underground railway infrastructure, which is a particular case of urban systems. However, all concepts introduced here are still valid for interurban systems.

In railway systems, the operation of trains follow a clear hierarchical framework with five levels: **scheduling**, **real-time (re)scheduling**, **remote traffic control**,  **interlocking and signalling** and **train & infrastructure control**. 

.. figure:: /_static/hierarchical_structure.jpg
   :alt: Hierarchical structure of the railway system.
   
   Hierarchical structure of the railway system (obtained from [WNBS]_\ ).

**Scheduling**

First, the railway transportation system is formulated on the basis of an extensive planning stage, which consists in deciding how many convoys are running at each period of time and planning a timetable. This planning stage is carried out a long time before the real-time operations taking into account demand estimation. Next, the railway managers need to assign the available resources, including the rolling stock and crew duties to the trips in this timetable.

**Real-time (re)scheduling**

During real-time operations, convoys may not adhere to the planning due to many external factors, such as failures, delays, interruptions, issues in the track infrastructure or a significant increase/decrease of passengers demand among many others. Hence, during real-time operation, the planning (this is, the rolling stock, timetable, etc) is usually rescheduled several times with real-time data collected.

**Interlocking and signalling**

The scheduled (or rescheduled) planning is communicated to local traffic centers, who set routes and track speed limits through interlocking systems and signaling systems (semaphores, needle exchange devices, etc). Side track devices (such as beacons) communicate with trains and provide them with updated schedules and track speed limits. Conversely, trains provide side track devices with real-time onboard operating data to feed the Local Traffic Center and the Traffic Management Center.

**Train Control**

With the given scheduled timetable, line infrastructure data and internal on-board computer data, the on-board computer generates a speed profile from the current station to the next one. This is, it decides the acceleration, cruising, coasting and braking periods to be carried out until next station.

The following figure summarises how all levels integrate in the railway system.

.. figure:: /_static/railway_traffic_control_train_operation_relation.jpg
   :alt: Relation between railway traffic control and train operation
   
   Railway traffic control and train operation integration (obtained from [YTYXHG]_\ ).

.. As seen in :ref:`benefits-automation`, there are three main goals of train automation: to provide passengers with a better quality of service, to reduce financial costs by increasing energy efficiency, which also makes this means of transport more environmentally friendly and to ensure maximum security to passengers. Every decision taken by automated train algorithms must take these three goals into account.

.. Next, the scheduled (or rescheduled) planning is communicated to each in-service train, whose goal is to conduct a safe, scheduled and efficient travel that meets with the objectives in :ref:`benefits-automation`. With the given scheduled timetable, line infrastructure data and internal on-board computer data, the on-board computer generates a speed profile from the current station to the next one. This is, it decides the acceleration, cruising, coasting and braking periods to be carried out until next station. The whole train journey is supervised by :term:`ATP` and :term:`ATS` systems, which will break the train if security conditions are not being fulfilled (for example, by exceeding speed limit, overrunning a red sempathore, etc).

Previous topic: :ref:`optimal-layout-design`.

Next topic: :ref:`railway-traffic-control`.