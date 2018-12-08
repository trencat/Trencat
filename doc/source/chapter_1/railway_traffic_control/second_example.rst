.. _conflict-resolution-problem-second-example:

A second example
----------------

The :term:`CRP` formulation can implement more general situations than the mentioned in :ref:`conflict-resolution-problem` and :ref:`conflict-resolution-problem-model`. For instance, the following figure shows a real situation where two trains must circulate in different conditions (extracted from [DPP]_).

.. figure:: /_static/network_second_example.jpg
   :alt: Small network framework with two trains, track intersections and additional schedule restrictions.
   
   Small network framework with two trains, track intersections and additional schedule restrictions. 

This figure depicts the network situation at time :math:`t=0`. The network consists of 9 tracks, three conflicting segments and one station with two platforms. Train :math:`t_1` is a fast train that travels through segments 2, 4, 5, 7, 8 and 9 without stopping at platform 7. Train :math:`t_2` is a slow train running through segments 3, 4, 5, 6, 8 and 9 and stops at platform 6.

Since train :math:`t_1` is a fast train, it will enter a segment at high speed only if the signal aspect is green, i.e., only if both the entering segment and the next one are free of trains. Therefore, it is important to schedule the movements of boths trains to ensure that two empty blocks are empty for train :math:`t_1`.

The next figure shows the alternative graph for this network.

.. figure:: /_static/second_example_alternative_graph.jpg
   :alt: Alternative graph of the second example.
   
   Alternative graph of the second example.

The weights :math:`f_{t_i,p_2}` and :math:`f_{t_2,p_4}` correnspond to the time at which the two trains are expected to reach the end of their block sections. In addition, :math:`\delta_{dep}` represents the scheduled departure time of train :math:`t_2`. The weight from :math:`p_{wait}` to :math:`p_8` (not depicted) is the scheduled dwell time.

**Further constraints**

Consider now that train :math:`t_1` is a slow train that must stop at 7 and carries people that must change to train :math:`t_2`. Train :math:`t_2` is currently empty and will take some of the passengers from train :math:`t_1`. The following actions must happen:

   - Train :math:`t_1` must arrive first at platform 6. It must wait until train :math:`t_2` arrives, plus an extra time to allow passengers to commute.
   - Train :math:`t_2` must depart before train :math:`t_1`.

How would you model the alternative graph?
