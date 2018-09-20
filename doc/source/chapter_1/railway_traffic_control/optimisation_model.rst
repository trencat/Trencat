.. _conflict-resolution-problem-model:

Optimization model
------------------

This section aims to give a rigorous mathematical formulation of the :term:`CRP`\ . Before going through the model, it is strongly recommend reading :ref:`conflict-resolution-problem` to motivate the problem and introduce notation with a simple example.

.. role:: raw-html(raw)
   :format: html

Sets and Parameters
^^^^^^^^^^^^^^^^^^^

Consider: 
   - T = {t\ :sub:`1`\ , t\ :sub:`2`\ ,..., t\ :sub:`n`\ } the set of trains that are currently circulating or waiting in garages until the departure time.
   - P = {p\ :sub:`1`\ , p\ :sub:`2`\ ,..., p\ :sub:`m`\ } the set of operations that must be performed by trains in T with running times F = {f\ :raw-html:`<sub>p<sub>1</sub></sub>`\ , f\ :raw-html:`<sub>p<sub>2</sub></sub>`\ ,..., f\ :raw-html:`<sub>p<sub>m</sub></sub>`\ }.
   - Two additional *dummy* operations p\ :sub:`0` and p\ :sub:`x` to indicate the *start* and *stop* of all trains respectively. Each train t\ :sub:`i` will need a running time of f\ :raw-html:`<sub>t<sub>i</sub>,p<sub>0</sub>`\ .

Consider a train t\ :sub:`i` and two consecutive operations p\ :sub:`r` and p\ :sub:`s` that this train must perform. If p\ :sub:`s` is the next operation after p\ :sub:`r`\ , then we will be denote p\ :sub:`s`\ = :raw-html:`&sigma;<sub>i</sub>`\ (p\ :sub:`r`\ ). In addition, the next operation after p\ :sub:`s` is denoted as :raw-html:`&sigma;<sub>i</sub>(p<sub>s</sub>)` = :raw-html:`&sigma;<sup>2</sup><sub>i</sub>`\ (p\ :sub:`r`\ ), etc. Therefore, a train t\ :sub:`i` will perform operations P\ :sub:`i` := {p\ :sub:`0`, :raw-html:`&sigma;<sub>i</sub>(p<sub>0</sub>)`\ , :raw-html:`&sigma;<sup>2</sup><sub>i</sub>(p<sub>0</sub>)`\ , ..., :raw-html:`&sigma;<sup>k</sup><sub>i</sub>(p<sub>0</sub>)`\ , ..., p\ :sub:`x`\ } in the given order. Notice that all trains perform operations p\ :sub:`0` and p\ :sub:`x`\ .

Let *G=(N,E,A)*\  be a directed graph where nodes in *N* represent operations to be performed for each train, edges in *E* indicate the order of the operations and their running times and *A* contains pairs of alternative edges of conflicting operations. Such graph *G* is called the **alternative graph**.

More specifically,

   - N := P :raw-html:`&cup;` {p\ :sub:`0`\ } :raw-html:`&cup;` {p\ :sub:`x`\ }.
   - E := :raw-html:`&cup;<sub>t<sub>i</sub>&straightepsilon;T</sub>` E\ :raw-html:`<sub>t<sub>i</sub></sub>`, where E\ :raw-html:`<sub>t<sub>i</sub></sub>` := :raw-html:`&cup;<sub>k=0,...,|P<sub>i</sub>|-1<sub>` { (:raw-html:`&sigma;<sup>k</sup><sub>i</sub>(p<sub>0</sub>)`, :raw-html:`&sigma;<sup>k+1</sup><sub>i</sub>(p<sub>0</sub>)`) }
   - A := { ( (p\ :sub:`r`, :raw-html:`&sigma;<sub>i</sub>(p<sub>r</sub>)`),(p\ :sub:`s`, :raw-html:`&sigma;<sub>j</sub>(p<sub>s</sub>))` ) | for all trains t\ :sub:`i`, t\ :sub:`j` :raw-html:`&straightepsilon;` T and for all operations p\ :sub:`r` :raw-html:`&straightepsilon;` P\ :sub:`i` and p\ :sub:`s` :raw-html:`&straightepsilon;` P\ :sub:`j` such that p\ :sub:`r` and p\ :sub:`s` are conflicting operations}.

Variables
^^^^^^^^^

   - s\ :raw-html:`<sub>t<sub>i</sub>,p<sub>r</sub></sub>` = Start time of operation p\ :sub:`r` :raw-html:`&straightepsilon;` P\ :sub:`i`. Defined for each train t\ :sub:`i` :raw-html:`&straightepsilon;` T and each operation p\ :sub:`j` :raw-html:`&straightepsilon;` P\ :sub:`i`.
   - s\ :raw-html:`<sub>p<sub>0</sub></sub>`\ = Start time of operation p\ :sub:`0`\ .
   - s\ :raw-html:`<sub>p<sub>n</sub></sub>`\ = Start time of operation p\ :sub:`n`\ .
   - x\ :raw-html:`<sub>t<sub>i</sub>p<sub>r</sub>,t<sub>j</sub>p<sub>s</sub></sub>` = Binary variable that takes the value 1 if train t\ :sub:`i` performs p\ :sub:`r`  before train t\ :sub:`j` performs p\ :sub:`s` and takes the value *0* otherwise. Defined for all pairs ( (p\ :sub:`r`, :raw-html:`&sigma;<sub>i</sub>(p<sub>r</sub>)`),(p\ :sub:`s`, :raw-html:`&sigma;<sub>j</sub>(p<sub>s</sub>))` ) :raw-html:`&straightepsilon;` A.

Objective function
^^^^^^^^^^^^^^^^^^
The objective function minimises the makespan of all train operations:

*minimise* s\ :raw-html:`<sub>p<sub>n</sub></sub>` - s\ :raw-html:`<sub>p<sub>0</sub></sub>`

Constraints
^^^^^^^^^^^
Subject to the following constraints.

   - s\ :raw-html:`<sub>&sigma;<sub>i</sub>(p<sub>0</sub>)</sub> &ge;` s\ :raw-html:`<sub>p<sub>0</sub></sub>` + f\ :raw-html:`<sub>t<sub>i</sub>,p<sub>0</sub></sub>`\ . Defined for all train t\ :sub:`i` :raw-html:`&straightepsilon;` T.
   - s\ :raw-html:`<sub>t<sub>i</sub>,&sigma;<sub>i</sub>(p<sub>r</sub>)</sub> &ge;` s\ :raw-html:`<sub>t<sub>i</sub>,p<sub>r</sub></sub>` + f\ :raw-html:`<sub>p<sub>r</sub></sub>`\ , for all train t\ :sub:`i` :raw-html:`&straightepsilon;` T and for all (p\ :sub:`r`, :raw-html:`&sigma;<sub>i</sub>(p<sub>r</sub>)`\ ) :raw-html:`&straightepsilon;` E\ :raw-html:`<sub>t<sub>i</sub></sub>`\ .
   - s\ :raw-html:`<sub>t<sub>j</sub>,p<sub>s</sub> &ge;` s\ :raw-html:`<sub>t<sub>i</sub>,&sigma;<sub>i</sub>(p<sub>r</sub>)</sub>` + f\ :raw-html:`<sub>&sigma;<sub>i</sub>(p<sub>r</sub>), p<sub>s</sub></sub>` **or** s\ :raw-html:`<sub>t<sub>i</sub>,p<sub>r</sub> &ge;` s\ :raw-html:`<sub>t<sub>j</sub>,&sigma;<sub>j</sub>(p<sub>s</sub>)</sub>` + f\ :raw-html:`<sub>&sigma;<sub>j</sub>(p<sub>s</sub>), p<sub>r</sub></sub>` for all ( (p\ :sub:`r`, :raw-html:`&sigma;<sub>i</sub>(p<sub>r</sub>)`),(p\ :sub:`s`, :raw-html:`&sigma;<sub>j</sub>(p<sub>s</sub>))` ) :raw-html:`&straightepsilon;` A.

The disjoint constraint can be linearsided by introducing binary variables and a large value *M* as follows:

s\ :raw-html:`<sub>t<sub>j</sub>,p<sub>s</sub> &ge;` s\ :raw-html:`<sub>t<sub>i</sub>,&sigma;<sub>i</sub>(p<sub>r</sub>)</sub>` + f\ :raw-html:`<sub>&sigma;<sub>i</sub>(p<sub>r</sub>), p<sub>s</sub></sub>` - M (1 - x\ :raw-html:`<sub>t<sub>i</sub>p<sub>r</sub>,t<sub>j</sub>p<sub>s</sub></sub>`)

s\ :raw-html:`<sub>t<sub>i</sub>,p<sub>r</sub> &ge;` s\ :raw-html:`<sub>t<sub>j</sub>,&sigma;<sub>j</sub>(p<sub>s</sub>)</sub>` + f\ :raw-html:`<sub>&sigma;<sub>j</sub>(p<sub>s</sub>), p<sub>r</sub></sub>` - M x\ :raw-html:`<sub>t<sub>i</sub>p<sub>r</sub>,t<sub>j</sub>p<sub>s</sub></sub>`

Previous topic: :ref:`conflict-resolution-problem`.   

Next topic: :ref:`conflict-resolution-problem-second-example`.