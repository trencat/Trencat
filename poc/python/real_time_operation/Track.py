"""TODO"""

from collections import namedtuple

Segment = namedtuple('Segment', 'length max_speed slope bend_radius tunnel')


class Track:
    """Class gathering track information."""

    def __init__(self, length, max_speed, slope, bend_radius, tunnel):
        """Initialise track object.

            Args:
                length (:obj:`tuple` or :obj:`list` of float): Length of each track segment.
                max_speed (:obj:`tuple` or :obj:`list` of float): Maximum velocity allowed at each track segment.
                slope (:obj:`tuple` or :obj:`list` of float): Slope of each track segment.
                bend_radius (:obj:`tuple` or :obj:`list` of float): Bend radius of each track segment. The value of a
                non-curved segment must be `math.inf`.
                tunnel (:obj:`tunnel` or :obj:`list` of bool): Boolean tuple indicating if each segment track runs in a
                tunnel or not.
        """
        if len(length) == len(max_speed) == len(slope) == len(bend_radius) == len(tunnel):
            raise ValueError('All tuples must have the same length.')

        self.length = length
        self.max_speed = max_speed
        self.slope = slope
        self.bend_radius = bend_radius
        self.tunnel = tunnel

    def __getitem__(self, item):
        """Get a :obj:`Segment` with the data of a track segment."""
        return Segment(length=self.length[item], max_speed=self.max_speed[item], slope=self.slope[item],
                       bend_radius=self.bend_radius[item], tunnel=self.tunnel[item])

    def __len__(self):
        """Number of track segments."""
        return len(self.length)
