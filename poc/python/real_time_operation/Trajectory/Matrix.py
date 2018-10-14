"""TODO"""

# TODO: Try to use python libraries instead of custom library

class Matrix:
    """Basic matrix support to be used in optimisation models."""

    @staticmethod
    def check_dimensions(value):
        """Checks if the input list has correct dimensions."""
        # TODO: Find better implementation
        check_cols = set([len(row) for row in value])  # TODO: Implement only with generator, instead of list
        if len(check_cols) != 1 or check_cols.pop() == 0:
            raise ValueError('Trying to define a matrix with inconsistent dimensions.')
        return True

    def __init__(self, value):
        """TODO"""
        self.matrix = value


    @property
    def matrix(self):
        """TODO"""
        return self._matrix

    @matrix.setter
    def matrix(self, value):
        """TODO"""
        if not Matrix.check_dimensions(value):
            raise ValueError("Matrix dimensions don't match")
        self._matrix = value

    @property
    def rows(self):
        """TODO"""
        return len(self._matrix)

    @property
    def cols(self):
        """TODO"""
        return len(self._matrix[0])

    def __and__(self, other):
        """Join two matrices horizontally"""
        if self.rows == other.rows:
            return Matrix([self[i] + other[i] for i in range(self.rows)])
        else:
            raise ValueError('Trying to join matrices of different dimensions.')

    def __or__(self, other):
        """Join two matrices vertically"""
        if self.cols == other.cols:
            return Matrix([self[i] for i in range(self.rows)] + [other[i] for i in range(other.rows)])
        else:
            raise ValueError('Trying to join matrices of different dimensions.')

    def __ior__(self, other):
        """Join two matrices vertically"""
        return self.__or__(other)

    # Operators
    def __getitem__(self, key):
        """Get a row at the specified position"""
        return self.matrix[key]

    def __add__(self, other):
        """Addition operation."""
        if isinstance(other, Matrix):
            return Matrix.__add__matrix__matrix(self, other)

    def __radd__(self, other):
        return self.__add__(other)

    def __mul__(self, other):
        """Multiplication operation."""
        if isinstance(other, Matrix):
            return Matrix.__mul__matrix__matrix(self, other)
        else:
            return Matrix.__mul__matrix__scalar(self, other)

    def __rmul__(self, other):
        """Matrix multiplication"""
        return self.__mul__(other)

    @staticmethod
    def __add__matrix__matrix(m1, m2):
        """Addition of two matrices."""
        if m1.rows != m2.rows or m1.cols != m2.cols:
            raise ValueError('Trying to add matrices of different dimensions.')

        range_rows, range_cols = range(m1.rows), range(m1.cols)
        return Matrix([[m1[i][j] + m2[i][j] for j in range_cols] for i in range_rows])

    @staticmethod
    def __mul__matrix__scalar(m1, num):
        """Multiply matrix by a number."""

        range_rows, range_cols = range(m1.rows), range(m1.cols)
        return Matrix([[num * m1[i][j] for j in range_cols] for i in range_rows])

    @staticmethod
    def __mul__matrix__matrix(m1, m2):
        """Addition of two matrices."""
        if m1.cols != m2.rows:
            raise ValueError('Trying to multiply matrices of different dimensions.')

        m1_rows, m1_cols, m2_cols = range(m1.rows), range(m1.cols), range(m2.cols)
        return Matrix([[sum([m1[i][k] * m2[k][j] for k in m1_cols]) for j in m2_cols] for i in m1_rows])

    # TODO: zeros, id, diag instead of Zeros, Id, Diag

    # Utilities
    @staticmethod
    def Zeros(rows, cols):
        """Define a zero matrix with the specified dimensions"""
        return Matrix([[0 for j in range(cols)] for i in range(rows)])

    @staticmethod
    def Id(size):
        """Define the identity matrix with the specified dimensions"""
        return Matrix([[1 if i == j else 0 for j in range(size)] for i in range(size)])

    @staticmethod
    def Diag(iterable):
        """Create a diagonal matrix with the given iterable"""
        return Matrix([[iterable[i] if i == j else 0 for j in range(len(iterable))] for i in range(len(iterable))])

    def __repr__(self):
        """String representation of a matrix"""
        return repr(self.matrix)
