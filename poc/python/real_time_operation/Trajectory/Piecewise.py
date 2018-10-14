"""TODO"""

from collections import namedtuple

Piece = namedtuple('Piece', 'a b domain')
Piece.__doc__ = """Linear function of a :obj:`Piecewise` function.

    A Piece is a linear function of the form:
        f(x) = a * x + b, with x in the interval [domain[0], domain[1]]

    Args:
        a (float): Linear coefficient.
        b (float): Independent term.
        domain (:obj:`tuple` or :obj:`list` of float): Domain of the linear function.
        
    For example, Piece(a=0.5, b=1, domain=(3, 7))."""


class Piecewise:
    """Wrapper class around :ref:`numpy.piecewise` implementing piecewise affine functions.

    An N-piecewise affine function of one variable x in the domain [alpha, beta] is a function of the form
        f(x) = a_i * x + b_i for x in the interval [z_i, z_(i+1)], with
        alpha = z_0 < z_1 < ... < z_(N-1) < z_N = beta.

    This class give access to the parameters a_i, b_i, z_i and z_(i+1).
    """

    @staticmethod
    def check(pieces):
        """Check if a piecewise function is continuous."""
        # TODO: Implement function that checks if pieces are continuous
        return True

    def __init__(self, pieces):
        """Initialise Piecewise object.

            Args:
                pieces (:obj:`tuple` or :obj:`list` of :obj:`Piece`). Pieces must be ordered form lower to higher domain
                    values.
        """
        self.pieces = pieces

    @property
    def pieces(self):
        """Gets or sets the pieces of the piecewise function"""
        return self._pieces

    @pieces.setter
    def pieces(self, value):
        if not Piecewise.check(value):
            raise ValueError('Pieces are not continuous')
        self._pieces = value
        # TODO: numpy.piecewise instance here

    def __getitem__(self, item):
        return self.pieces[item]

    def __call__(self, x):
        """Evaluate piecewise function at a given point."""
        raise NotImplementedError("Piecewise function evaluation is not supported yet.")
        # TODO: Call numpy.piecewise
