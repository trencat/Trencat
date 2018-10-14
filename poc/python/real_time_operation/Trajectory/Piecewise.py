"""TODO"""

from collections import namedtuple

Piece = namedtuple('Piece', 'a b domain')
Piece.__doc__ = """TODO"""

class Piecewise:
    """TODO"""

    @staticmethod
    def check(pieces):
        """TODO"""
        # TODO: Implement function that checks if pieces are continuous
        return True

    def __init__(self, pieces):
        """TODO"""

        self._pieces = None
        self.pieces = pieces

    @property
    def pieces(self):
        """TODO"""
        return self._pieces

    @pieces.setter
    def pieces(self, value):
        if not Piecewise.check(value):
            return ValueError('Pieces are not continuous')
        self._pieces = value
        #TODO: Numpy piecewise here

    def __getitem__(self, item):
        return self.pieces[item]
