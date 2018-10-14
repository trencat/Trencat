"""TODO"""

from collections import namedtuple

Segment = namedtuple('Segment', 'length max_speed slope bend_radius tunnel')

class Track:
    """TODO"""

    def __init__(self, length, max_speed, slope, bend_radius, tunnel):
        """TODO"""

        assert (len(length) == len(max_speed) == len(slope) == len(bend_radius) == len(tunnel))

        self.length = length
        self.max_speed = max_speed
        self.slope = slope
        self.bend_radius = bend_radius
        self.tunnel = tunnel

    def __getitem__(self, item):
        return Segment(length=self.length[item], max_speed=self.max_speed[item], slope=self.slope[item],
                       bend_radius=self.bend_radius[item], tunnel=self.tunnel[item])

    def __len__(self):
        return len(self.length)
