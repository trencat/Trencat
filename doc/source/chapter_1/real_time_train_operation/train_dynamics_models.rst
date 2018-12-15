.. _train-dynamic-models:

Train dynamic models
--------------------

Mainly based on [YTYXHG]_ and [WNBS]_.

There are two main dynamics models depending on the type of train: **single-point control models** and **multi-point control models**. Single-point models offer good results in urban rail transit systems, where running resistances can be neglected with respect to the traction/braking forces. These models are too simple for heavy-haul trains, which consist of a large number of long vehicles carrying a lot of cargo. Multi-point models are more suitable for these types of trains.

Single-point control models
^^^^^^^^^^^^^^^^^^^^^^^^^^^

In single-point control models, a train that consists of more than one vehicle is simplified as a single-point, which means that its dynamics can be modelled with a Newton's equation. Considering the time as the independent variable, the model states as follows:

.. math::
   :label: traindynamicseq

   m\rho\frac{dv(t)}{dt} = u(t) - R_b(t) - R_l(s(t), v(t))

.. math::

   \frac{ds(t)}{dt} = v(t)

where
   - :math:`m`: Mass of the train.
   - :math:`\rho`: Rotating mass factor.
   - :math:`v(t)`: Velocity of the train at time :math:`t`.
   - :math:`s(t)`: Position of the train at time :math:`t`.
   - :math:`u(t)`: Control variable representing the traction force (if positive) or braking force (if negative).
   - :math:`R_b(v(t))`: Basic resistance including roll resintance and air resistance.
   - :math:`R_l(s(t),v(t))`: Resitance causes by track grade, curves and tunnels.

.. note::

   For convenience we will denote :math:`s:=s(t)` and :math:`v:=v(t)`.

First, the traction/braking force are bounded by :math:`u_{min}(v) \leq u(t) \leq u_{max}`(v)`. According to [WNBS]_\ , usually :math:`u_{min}` and :math:`u_{max}` are considered as constant parameters although in reality they are a function of velocity due to many factors such as the maximum adhesion and characteristics of the power equipment. The expression of :math:`u_{max}` can be modelled as a quadratic piecewise function of the form

.. math::

   u_{max}(v) = c_{1,j} + c_{2,j}v + c_{3,j}v^2 \text{, for } v\in [v_j, v_{j+1}]

or 

.. math::

   u_{max}(v) = c_{h,j}v^{-1} \text{, for } v\in [v_j, v_{j+1}]

for :math:`j = 1, 2, ..., M-1`, and for constants :math:`c_{1,j}`, :math:`c_{2,j}`, :math:`c_{3,j}`, :math:`c_{h,j}`, :math:`v_j`, :math:`v_{j+1}` determined by the characteristics of the train.

.. figure:: /_static/maximum_traction_force_ATO.jpg
   :alt: Maximum traction force depending on train speed.
   
   Maximum traction force :math:`u_{max}` depending on train speed (obtained from [WNBS]_\ ).

Regarding :math:`u_{min}`, according to [WNBS]_\ , under normal circumstances the maximum braking force is just *0.75* times the train's full braking force, which is reserved for emergency stop. This reduction also contributes in passengers riding comfort.

Second, the basic resitance :math:`R_b(v)` is expressed as follows:

.. math::

   R_b(v) = m(a_1 + a_2v^2),

where again, coefficients :math:`a_1` and :math:`a_2` depend on the characteristics of the train and wind speed.

Third, :math:`R_l(s,v)` is a function of both space and velocity, expressed as follows

.. math::

   R_l(s,v) = mg\sin(\beta(s)) + f_c(r(s)) + f_t(l_t(s), v),

The first term corresponds to the acceleration/deacceleration in non-zero slope segments. The mass of the train is :math:`m`, the gravitational acceleration is :math:`g` and the slope at position :math:`s` is :math:`beta(s)`. The second term refers to the deacceleration caused by the increase of the friction of the wheels in curves, being :math:`r(s)` the radius of the curve at position :math:`s`. The third term models the deacceleration caused by the air friction inside tunnels, being :math:`l_t(s)` the length of the tunnel at position :math:`s`.

Terms :math:`f_c(\cdot)` and :math:`f_t(\cdot)` are based on empirical formulas. For instance, [WNBS]_ uses the following expression:

.. math::

   f_c(r(s)) = \frac{6.3m}{r(s) - 55} \quad \text{if } r(s)\geq 300\text{ meters}

.. math::

   f_c(r(s)) = \frac{4.91m}{r(s) - 55} \quad \text{if } r(s)\le 300\text{ meters}

Finally, a train is affected by air resistance inside a tunnel. The value of such resistance depends on the form of the tunnel as well as the smoothnes of its walls and the train walls. If there is a limiting gradient in the tunnel (which is defined as the maximum railway gradient that can be climed without the help of a second power unit), then the following expression applies [WNBS]_:

.. math::

   f_t(l_t(s), v) = 1.296\cdot10^{-9}l_t(s)mgv^2.

On the contrary, if there is no limiting gradient, the following expression applies:

.. math::

   f_t(l_t(s), v) = 1.3\cdot10^{-7}l_t(s)mg.

Finally, :math:`f_t(l_t(s), v)` equals zero outside tunnels.


Multi-point control models
^^^^^^^^^^^^^^^^^^^^^^^^^^

In multi-point models, all vehicle's positions, velocities and accelerations are taken into account. In addition, couplers between vehicles are modelled as linear springs. Basically, multi-point control models are a natural extension of single-point models, which means that they contain the same Newton's equations, but now with additional terms regarding train interactions during traction and braking periods. The following figure shows the key concepts of these types of models.

.. figure:: /_static/multi-point_control_model_ATO.jpg
   :alt: Illustration of the forces involved in multi-point control models.
   
   Illustration of the forces involved in multi-point control models (obtained from [YTYXHG]_\ ).

In this project we are focusing on single-point models. Therefore, Newton's equations are not introduced in this section. Check [YTYXHG]_ for further references.
