import math

GRAVITY_ACCELERATION = 9.80665  # In meters/(second**2)


class Train:
    """Train class implementing the model described in documentation.

    The motion of the train can be described by Newton's law:
        m * rho * (dv/dt) = u(t) - Rb(v) - Rline(x, v), where

    - m: is the mass of the train [Kg]
    - rho: is the mass factor of the train
    - x = x(t): position of the train at time t [m].
    - v = v(t): velocity of the train at time t [m/s].
    - dv/dt: Derivative of the velocity with respect to time variable, i.e, train acceleration [m/s**2).
    - u(t): Traction traction force (if positive) or braking force (if negative) [Newtons].
    - Rb(v): Basic resistance of the train at a given speed. Rb(v) = m * (a1 + a2 * v**2), where a1 and a2 coefficients
             depend on the train characteristics and the wind speed.
    - Rline(x, v): Line resistance of a train at a given position and speed. It is the addition of the slope resistance
                   at a given point, the curve resistance at a given point and the tunnel resistance at a given point
                   and velocity.
       - The slope resistance is computes as mass * gravity_acceleration * sin(slope(x)).
       - The curve resistance is computed as:
          * mass * 6.3 / (r(x) - 55) if the bend radius r(x) at position x is >= 300 meters,
          * mass * 4.91 / (r(x) - 30) if the bend radius r(x) at position x is < 300 meters.
       - The tunnel resistance, under the assumption that there is a limiting gradient (i.e., there is a maximum
         gradient that can be climbed without the help of a second power unit) is computed as:
            1.296 * 1e-9 * tunnel_length * mass * gravity_acceleration * v**2.
    """

    def __init__(self, mass, massfactor, velocity, max_brake, max_traction, basic_resistance_params):
        """Initialise train object:

            Args:
                mass (float): Mass in [kg].
                massfactor (float): Mass factor.
                velocity (float): Current train velocity in [m/s].
                max_brake (float): Maximum service brake force in [Newtons]. Usually, the service brake is 0.75 the
                emergency brake.
                max_traction (float or :obj:`Piecewise`): Maximum traction force in [Newtons].
                basic_resistance_params (:obj:`tuple` of float): Parameters a1 and a2 used to compute basic resistance.
        """
        self.mass = mass
        self.massfactor = massfactor
        self.velocity = velocity
        self.max_brake = abs(max_brake)
        self.max_traction = max_traction
        self.basic_resistance_params = basic_resistance_params

    def basic_resistance(self, velocity=None):
        """Compute basic resistance of the train at a given velocity.

            If velocity isn't provided, current train velocity will be used.
        """
        if velocity is None:
            velocity = self.velocity
        a, b = self.basic_resistance_params
        return self.mass * (a + b * velocity**2)

    def slope_resistance(self, segment):
        """Compute slope resistance of the train at a given track segment."""
        return self.mass * GRAVITY_ACCELERATION * math.sin(segment.slope)

    def curve_resistance(self, segment):
        """Compute the curve resistance of the train at a given track segment."""
        if segment.bend_radius <= 30:
            return math.inf
        if segment.bend_radius < 300:
            return self.mass * 4.91 / (segment.bend_radius - 30)
        else:
            return self.mass * 6.3 / (segment.bend_radius - 55)

    def tunnel_resistance(self, segment, velocity=None):
        """Compute the tunnel resistance of the train at given track segment and velocity.

            If velocity isn't provided, current train velocity will be used.
        """
        if not segment.tunnel:
            return 0.0
        if velocity is None:
            velocity = self.velocity
        return 1.296 * 1e-9 * segment.length * self.mass * GRAVITY_ACCELERATION * velocity**2

    def line_resistance(self, segment, velocity=None):
        """Compute the line resistance at given track segment and velocity.

            If velocity isn't provided, current train velocity will be used.
        """
        if velocity is None:
            velocity = self.velocity
        return self.slope_resistance(segment) + \
               (self.curve_resistance(segment) if velocity != 0 else 0) + \
               self.tunnel_resistance(segment, velocity)

    def resistance(self, segment, velocity=None):
        """Compute the train resistance at given track segment and velocity.

            If velocity isn't provided, current train velocity will be used.
        """
        if velocity is None:
            velocity = self.velocity
        return self.basic_resistance(velocity) + self.line_resistance(segment, velocity)
