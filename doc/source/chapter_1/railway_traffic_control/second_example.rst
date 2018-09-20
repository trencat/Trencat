.. _conflict-resolution-problem-second-example:

A second example
----------------

.. role:: raw-html(raw)
   :format: html

The :term:`CRP` formulation can implement more general situations than the mentioned in :ref:`conflict-resolution-problem` and :ref:`conflict-resolution-problem-model`. For instance, the following figure shows a real situation where two trains must circulate in different conditions (extracted from [DPP]_).

.. figure:: /_static/network_second_example.jpg
   :alt: Small network framework with two trains, track intersections and additional schedule restrictions.
   
   Small network framework with two trains, track intersections and additional schedule restrictions. 

This figure depicts the network situation at time t=0. The network consists of 9 tracks, three conflicting segments and one station with two platforms. Train t\ :sub:`1` is a fast train that travels through segments 2, 4, 5, 7, 8 and 9 without stopping at platform 7. Train t\ :sub:`2` is a slow train running through segments 3, 4, 5, 6, 8 and 9 and stops at platform 6.

Since train t\ :sub:`1` is a fast train, it will enter a segment at high speed only if the signal aspect is green, i.e., only if both the entering segment and the next one are free of trains. Therefore, it is important to schedule the movements of boths trains to ensure that two empty blocks are empty for train t\ :sub:`1`.

The next figure shows the alternative graph for this network.

.. figure:: /_static/second_example_alternative_graph.jpg
   :alt: Alternative graph of the second example.
   
   Alternative graph of the second example.

The weights f\ :raw-html:`<sub>t<sub>1</sub>,p<sub>2</sub></sub>` and f\ :raw-html:`<sub>t<sub>2</sub>,p<sub>4</sub></sub>` correnspond to the time at which the two trains are expected to reach the end of their block sections. In addition, \ :raw-html:`&delta;<sub>dep</sub>` represents the scheduled departure time of train t\ :sub:`2`. The weight from p\ :sub:`wait` to p\ :sub:`8` (not depicted) is the scheduled dwell time.

**Further constraints**

Consider now that train t\ :sub:`1` is a slow train that must stop at 7 and carries people that must change to train t\ :sub:`2`. Train t\ :sub:`2` is currently empty and will take some of the passengers from train t\ :sub:`1`. The following actions must happen:

   - Train t\ :sub:`1` must arrive first at platform 6. It must wait until train t\ :sub:`2` arrives, plus an extra time to allow passengers to commute.
   - Train t\ :sub:`2` must depart before train t\ :sub:`1`.

How would you model the alternative graph?

Previous topic: :ref:`conflict-resolution-problem-model`.
   
Next: :ref:`optimal-rolling-stock-planning`.