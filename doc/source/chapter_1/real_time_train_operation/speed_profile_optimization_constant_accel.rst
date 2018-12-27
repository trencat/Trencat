Computing physical magnitudes assuming constant acceleration
------------------------------------------------------------

Consider the situation in the next figure. The train is in position :math:`s_i` with velocity :math:`v_i` at time :math:`t_i` and it has to accelerate to achieve velocity :math:`v_j` at position :math:`s_j`. We want to know how much time the train needs to carry out this transition, how much `work <https://en.wikipedia.org/wiki/Work_(physics)>`_ it will do and what's the train average `power <https://en.wikipedia.org/wiki/Power_(physics)>`_ over this segment.

.. figure:: /_static/speed_profile_introduction_4.jpg
   :alt: Moving from one state to the next one.
   
   Moving from one state to the next one.

Segment time
''''''''''''

Let's assume that this transition is done with constant acceleration. The time needed for this transition is computed using kinemic formula :math:`\Delta s = \frac{1}{2}(v_j + v_i)\Delta t`:

.. math::

    \Delta t = \frac{2\Delta s}{v_j + v_i},\quad \text{i.e.} \quad t_j = t_i + \frac{2(s_j-s_i)}{v_j+v_i}.

.. note::

   The previous formula is not valid if both :math:`v_i` and :math:`v_j` equal zero. Since the train is stopped it implies that :math:`s_i =  s_j`, which contradicts the situation depticted (:math:`s_j > s_i`).

Segment acceleration
''''''''''''''''''''

Time train acceleration is straightforward to compute:

.. math::

    a = \frac{\Delta v}{\Delta t} = \frac{v_j-v_i}{t_j-t_i}.

Work done by the train in a segment
'''''''''''''''''''''''''''''''''''

The computation of the work done by the train is a bit tricky. A way to compute the work done by the train in this segment is:

.. math::

   W_{ij} = \int_{t_i}^{t_j}u(t)v(t)dt

Recall again Newton's equation :eq:`traindynamicseq`. Due to the space discretization, the term :math:`mg\sin(s)` is constant in the entire segment. In addition, the term :math:`l_t(s)` is also constant in the entire segment. Acceleration :math:`a` is also constant along the track by assumption. Therefore, :eq:`traindynamicseq` can be refactored as

.. math::

   m\rho a = u(t) - C_1 - C_2 v^2(t),

where :math:`C_1` and :math:`C_2` are constants. It is important to remark that :math:`C_2` is always non-negative. Notice also from this expression that the traction/braking force is not constant along the segment. The work done by the train in this segment is computed as the integral of traction force times velocity:

.. math::

   \begin{array}{rl}
   W_{ij} =& \int_{t_i}^{t_j}u(t)v(t)dt\\
     =& \int_{t_i}^{t_j} (m\rho a + C_1)v(t)dt + \int_{t_i}^{t_j}C_2v^3(t)dt\\
     \stackrel{v(t) = v_i + a(t-t_i)}=& (m\rho a^2 + aC_1)\int_{t_i}^{t_j}[v_i + a (t-t_i)] dt + a^3 C_2\int_{t_i}^{t_j}[v_i + a (t-t_i)]^3dt\\
     \stackrel{\tilde{t} = v_i + a(t-t_i)}=& (m\rho a^2 + aC_1)\int_{v_i}^{v_i + a (t_j-t_i)}a^{-1}\tilde{t}d\tilde{t} + a^3 C_2\int_{v_i}^{v_i +  a(t_j - t_i)}a^{-3}\tilde{t}^3d\tilde{t}\\
     =& (m\rho a + C_1)\frac{(v_i + a(t_j-t_i))^2 - v_i^2}{2} + C_2\frac{(v_i + a(t_j-t_i))^4 - v_i^4}{4}.
   \end{array}

Work done only by the traction force in a segment
.................................................

Until here everything is ok, but let's go a little bit further. Notice that the previous formula includes the work done by the traction force and the braking force. It is more accurate, though, to minimise only the work done by the traction force since the train consumes energy from the grid only when :math:`u(t) > 0` but not when :math:`u(t) < 0`. Work is then computed as

.. math::

   \overline{W_{ij}} = \int_{t_i}^{t_j}\max(u(t), 0)v(t)dt = \int_{t_i}^{t_j}u(t)v(t)dt - \int_{t\in\{t_i \leq t \leq t_j | u(t) < 0\}}u(t)v(t)dt.

The set :math:`\{t_i \leq t \leq t_j | u(t) < 0\}` is computed as follows:

.. math::

   \begin{array}{rl}
   u(t) < 0 \Leftrightarrow & m\rho a + C_1 + C_2v^2(t) < 0\\
   \stackrel{v(t)=v_i + a(t-t_i)}{\Leftrightarrow} & m\rho a + C_1 + C_2a[v_i + a(t-t_i)]^2 < 0\\
   \Leftrightarrow & m\rho a + C_1 + C_2a[v_i^2 + 2 v_ia(t-t_i) + a^2(t-t_i)^2] < 0\\
   \stackrel{\text{reorganise}}\Leftrightarrow & (m\rho a + C_1 + C_2av_i^2) + 2v_iC_2a^2(t-t_i) + C_2a^3(t-t_i)^2 < 0\\
   \end{array}

Let us refactor the last expression with :math:`\tilde{a} = m\rho a + C_1 + C_2av_i^2`, :math:`\tilde{b} = 2v_iC_2a^2` and :math:`\tilde{c} = C_2a^3`. The last inequality is equivalent to:

.. math::

   u(t) < 0 \Leftrightarrow \tilde{a} + \tilde{b}(t-t_i) + \tilde{c}(t-t_i)^2 < 0,

which is a second degree polinomial. The two times :math:`t` where the polynomial equals 0, together with the sign of the acceleration :math:`a` determines when :math:`u(t) < 0`. If acceleration is positive, then 

.. math::

    u(t) < 0 \Leftrightarrow t_i - \sqrt{\frac{\tilde{b}^2 - 4\tilde{a}\tilde{c}}{2\tilde{a}}}  < t < t_i + \sqrt{\frac{\tilde{b}^2 - 4\tilde{a}\tilde{c}}{2\tilde{a}}},

if the acceleration is negative, then

.. math::

   u(t) < 0 \Leftrightarrow t < t_i - \sqrt{\frac{\tilde{b}^2 - 4\tilde{a}\tilde{c}}{2\tilde{a}}}  \quad\text{or}\quad t > t_i + \sqrt{\frac{\tilde{b}^2 - 4\tilde{a}\tilde{c}}{2\tilde{a}}},


Such inequalities are valid only if the square root can be computed. For convenience, let us denote :math:`\underline{t} := t_i - \sqrt{\frac{\tilde{b}^2 - 4\tilde{a}\tilde{c}}{2\tilde{a}}}` and :math:`\overline{t} := t_i + \sqrt{\frac{\tilde{b}^2 - 4\tilde{a}\tilde{c}}{2\tilde{a}}}`. Finally, the work done only by the traction force is computed as 

.. math::

   \overline{W_{ij}} = 
   \left\{\begin{array}{rl}
      \int_{t_i}^{t_j}u(t)v(t)dt,& \text{if } a\neq 0, \frac{\tilde{b}^2 - 4\tilde{a}\tilde{c}}{2\tilde{a}} < 0,\\
      \int_{\min\{\overline{t}, t_j\}}^{t_j}u(t)v(t)dt,& \text{if } a\neq 0, \frac{\tilde{b}^2 - 4\tilde{a}\tilde{c}}{2\tilde{a}} > 0,\\
      C_1, & \text{if } a\neq 0, \frac{\tilde{b}^2 - 4\tilde{a}\tilde{c}}{2\tilde{a}} = 0, C_1 \geq 0,\\
      0, &\text{otherwise}.
   \end{array}\right.

Segment Average power
'''''''''''''''''''''

Regarding the average power of the train, it is computed as:

.. math::

   P_{ij} = \frac{\Delta W_{ij}}{\Delta t},\quad \text{or} \quad \overline{P_{ij}} = \frac{\Delta \overline{W_{ij}}}{\Delta t}.

Maximum traction/braking force in a segment
'''''''''''''''''''''''''''''''''''''''''''

In section :ref:`speed-profile-optimization` it was shown that while building the graph we need to know if the traction/braking force required to move from one state to the next one exceeds the train's maximum traction/braking force.
As mentioned earlier, we can compute explicity the traction force in a segment with

.. math::

   u(t) = m\rho a + C_1 + C_2v^2(t)

Since :math:`C_2` is non-negative, :math:`u(t)` is a second degree polinomial with a positive quadratic coefficient, it is straightforward to see that

.. math::

   \min_{t\in[t_i, t_j]}\{u(t)\} = u(t_i)\quad \text{and}\quad \max_{t\in[t_i, t_j]}\{u(t)\} = u(t_j).


Jerk rate
'''''''''

The jerk rate is highly related with passenger's discomfort. The higher the rate, the higher the discomfort. The authors [WNBS]_ propose measuring the jerk rate in a segment :math:`[s_i, s_j]` as the sum of the change rate of the traction/braking force in absolute value ove the segment:

.. math::

   \begin{array}{rl}
   J_{ij} :=& \int_{t_i}^{t_j} \left|\frac{du(t)}{dt}\right| dt\\
   =& \int_{t_i}^{t_j} \left|\frac{d}{dt}(m\rho a + C_1 + C_2v^2(t))\right|dt\\
   =& \int_{t_i}^{t_j} \left|2C_2v(t)\frac{dv(t)}{dt}\right|dt\\
   \stackrel{\frac{dv(t)}{dt} = a}{=}& \int_{t_i}^{t_j} \left|2C_2v(t)a\right|dt\\
   \stackrel{C_2 \geq 0}{=}& 2C_2|a|\int_{t_i}^{t_j} \left|v(t)\right|dt\\
   \stackrel{v(t) \geq 0}{=}& 2C_2|a|\int_{t_i}^{t_j} v(t)dt\\
   =& 2C_2|a|\Delta s = 2C_2|a|(s_j - s_i)\\
   \end{array}

To compute the Jerk rate we recall that we assumed constant acceleration, we know the coefficient :math:`C_2` is always non-negative and that the integral of velocity over the period time :math:`[t_i, t_j]` is precisely the displacement of the train in this segment.



