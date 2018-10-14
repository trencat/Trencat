"""TODO"""

import math
from real_time_operation import GRAVITY_ACCELERATION


class Train:
    """TODO"""

    def __init__(self, mass, massfactor, velocity, max_brake, max_traction, basic_resistance_params):
        """TODO"""
        self.mass = mass
        self.massfactor = massfactor
        self.velocity = velocity
        self.max_brake = max_brake
        self.max_traction = max_traction  # Can be either float or piecewise affine
        self.basic_resistance_params = basic_resistance_params

    def basic_resistance(self, velocity=None):
        """TODO"""
        if velocity is None:
            velocity = self.velocity
        a, b = self.basic_resistance_params
        return self.mass * (a + b * velocity**2)

    def slope_resistance(self, segment):
        """TODO"""
        return self.mass * GRAVITY_ACCELERATION * math.sin(segment.slope)

    def curve_resistance(self, segment):
        """TODO"""
        if segment.bend_radius <= 30:
            return math.inf
        if segment.bend_radius < 300:
            return self.mass * 4.91 / (segment.bend_radius - 30)
        else:
            return self.mass * 6.3 / (segment.bend_radius - 55)

    def tunnel_resistance(self, segment, velocity=None):
        """TODO"""
        if not segment.tunnel:
            return 0.0
        if velocity is None:
            velocity = self.velocity
        return 1.296 * 1e-9 * segment.length * self.mass * GRAVITY_ACCELERATION * velocity**2

    def line_resistance(self, segment, velocity=None):
        """TODO"""
        if velocity is None:
            velocity = self.velocity
        return self.slope_resistance(segment) + self.curve_resistance(segment) + self.tunnel_resistance(segment, velocity)

    def resistance(self, segment, velocity=None):
        """TODO"""
        if velocity is None:
            velocity = self.velocity
        return self.basic_resistance(velocity) + self.line_resistance(segment, velocity)
