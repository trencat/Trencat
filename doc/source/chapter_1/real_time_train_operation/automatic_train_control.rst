.. role:: raw-html(raw)
   :format: html

.. _automatic-train-control:

Automatic Train Control
-----------------------

Mainly based on [YTYXHG]_ and [WNBS]_.


Train dynamic models
^^^^^^^^^^^^^^^^^^^^

There are two main train operation models depending on the type of train: **single-point control models** and **multi-point control models**. Single-point models offer good results in urban rail transit systems, where running resistances can be neglected with respect to the traction/braking forces. These models are too simple for heavy-haul trains, which consist of a large number of long vehicles carrying a lot of cargo. Multi-point models are more suitable for these types of trains.

Single-point control models
^^^^^^^^^^^^^^^^^^^^^^^^^^^

In single-point control models, a train that consists of more than one vehicle is simplified as a single-point, which means that its dynamics can be modelled with a Newton's equation. Considering the time as the independent variable, the model states as follows:

m\ :raw-html:`&rho;v&#775;`\ (t) = u(t) - R\ :sub:`b`\ (t) - R\ :sub:`l`\ (s,v)

:raw-html:`s&#775;`\ (t) = v(t)

where
   - m := Mass of the train.
   - :raw-html:`&rho;` := Rotating mass factor.
   - v(t) := Velocity of the train at time t.
   - s(t) := Position of the train at time t.
   - u(t) := Control variable representing the traction force (if positive) or braking force (if negative).
   - R\ :sub:`b`\ (v) := Basic resistance includeing roll resintance and air resistance.
   - R\ :sub:`l`\ (s,v) := Resitance cause by track grade, curves and tunnels.

The model has several constraints that are introduced next.

First, the traction/braking force are bounded by u\ :sub:`min`\ (v) :raw-html:`&le;` u(t) :raw-html:`&le;` u\ :sub:`max`\ (v). According to [WNBS]_\ , usually u\ :sub:`min` and u\ :sub:`max` are considered as constant parameters (for instance, the European Train Control System and the German train  control system) although in reality they are a function of velocity due to many factors such as the maximum adhesion and characteristics of the power equipment.  The expression of u\ :sub:`max` can be modelled as a is a piecewise function that can take the the value

u\ :sub:`max`\ (v) = c\ :sub:`1,j` + c\ :sub:`2,j`\ v + c\ :sub:`3,j`\ v\ :sup:`2`, for v :raw-html:`&straightepsilon;` [v\ :sub:`j`, v\ :sub:`j+1`]

or 

u\ :sub:`max`\ (v) = c\ :sub:`h,j`\ v\ :sup:`-1` for v :raw-html:`&straightepsilon;` [v\ :sub:`j`, v\ :sub:`j+1`]

for j = 1, 2, ..., M-1, and for constants c\ :sub:`1,j`\ , c\ :sub:`2,j`\ , c\ :sub:`3,j`\ , c\ :sub:`h,j`\ , v\ :sub:`j`\ , v\ :sub:`j+1` determined by the characteristics of the train. According to [WNBS]_\ , under normal circumstances the maximum fraking force is just 0.75 times the train's full braking force, which is reserved for emergency stop. This reduction also contributes in passengers riding comfort.

.. figure:: /_static/maximum_traction_force_ATO.jpg
   :alt: Maximum traction force depending on train speed.
   
   Maximum traction force depending on train speed (obtained from [WNBS]_\ ).

Second, the basic resitance R\ :sub:`b`\ (v) is expressed as follows:

R\ :sub:`b`\ (v) = m(a\ :sub:`1` + a\ :sub:`2`\ v\ :sup:`2`),

where again, coefficients a\ :sub:`1` and a\ :sub:`2` depend on the characteristics of the train and wind speed.

Third, R\ :sub:`l`\ (s,v) is a function of both space and velocity, expressed as follows

R\ :sub:`l`\ (s,v) = m·g·sin(\ :raw-html:`&beta;`\ (s)) + f\ :sub:`c`\ (r(s)) + f\ :sub:`t`\ (l\ :sub:`t`\ (s), v),

The first term corresponds to the acceleration/deacceleration in non-zero slope segments. The mass of the train is m, the gravitational acceleration is g and the slop is :raw-html:`beta`\ (s). The second term refers to the deacceleration caused by the increase of the friction of the wheels in curves, being r(s) the radius of the curve at position s. The third term models the deacceleration caused by the air friction inside tunnels, being l\ :sub:`t`\ (s) the length of the tunnel.

Terms f\ :sub:`c`\ (·) and f\ :sub:`t`\ (·) are based on empyrical formulas. For instance, [WNBS]_ uses the following expression:

f\ :sub:`c`\ (r(s)) = 6.3·m·(r(s) - 55)\ :sup:`-1` if r(s) :raw-html:`&ge;` 300m

f\ :sub:`c`\ (r(s)) = 4.91·m·(r(s) - 30)\ :sup:`-1` if r(s) < 300m.

Finally, a train is affected by air resistance inside a tunnel. The value of such resistance depends on the form of the tunnel as well as the smoothnes of its walls and the train walls. If there is a limiting gradient in the tunnel (a limiting gradient is defined as the maximum railway gradient that can be climed without the help of a second power unit), then the following expression applies:

f\ :sub:`t`\ (l\ :sub:`t`\ (s), v) = 1.296·10\ :sup:`-9`\ ·l\ :sub:`t`\ (s)·m·g·v\ :sup:`2`\ .

On the contrary, if there is no limiting gradient, the following expression applies:

f\ :sub:`t`\ (l\ :sub:`t`\ (s), v) = 1.3·10\ :sup:`-7`\ ·l\ :sub:`t`\ (s)·m·g.

Finally, f\ :sub:`t`\ (l\ :sub:`t`\ (s), v) equals 0 outside tunnels.


Multi-point control models
^^^^^^^^^^^^^^^^^^^^^^^^^^

In multi-point models, all vehicle's positions, velocities and accelerations are taken into account. In addition, couplers between vehicles are modelled as linear springs. Basically, multi-point control models are a natural extension of single-point models, which means that they contain the same Newton's equations, but now with additional terms regarding train interactions during traction and braking periods. The following figure shows the key concepts of these types of models.

.. figure:: /_static/multi-point_control_model_ATO.jpg
   :alt: Illustration of the forces involved in multi-point control models.
   
   Illustration of the forces involved in multi-point control models (obtained from [YTYXHG]_\ ).

In this project we are focusing on single-point models. Therefore, Newton's equations are not introduced in this section. Check [YTYXHG]_ for further references.
However, if multi-point models were implemented, documentation would include enough documentation and references.

Previous topic: :ref:`real-time-train-operation`.

Next topic: :ref:`speed-profile-optimization`.
