import math
import operator
import pyscipopt as opt
from contextlib import suppress
from functools import reduce, lru_cache
from .Piecewise import Piece, Piecewise
from .Matrix import Matrix


def summation(iterable):
    """Sum of all elements in an iterable object."""
    return reduce(operator.add, iterable)


def prod(iterable):
    """Product of all elements in an iterable object."""
    return reduce(operator.mul, iterable, 1)


class Simple:
    """TODO"""

    def __init__(self, train, track, timespan, end_velocity, smooth_factor, punctuality_factor, start_time=0,
                 start_velocity=None):
        """TODO"""
        self.train = train
        self.track = track
        self.timespan = timespan
        self.start_velocity = start_velocity if start_velocity is not None else self.train.velocity
        self.end_velocity = end_velocity
        self.smooth_factor = smooth_factor
        self.punctuality_factor = punctuality_factor
        self.start_time = start_time
        self.N = len(self.track.length)
        self._min_kinetic_energy_threshold = 0.1
        self._min_velocity_threshold = None
        self.Estart = None
        self.Eend = None
        self.eps = 1e-8
        self.piecewise_approximation = None

        # Optimisation parameters
        self.model = None
        self.variables = None
        self.constraints = None

    def prepare_parameters(self):
        """TODO"""
        self.N = len(self.track.length)
        self._min_velocity_threshold = math.sqrt(2 * self._min_kinetic_energy_threshold)
        if self.start_velocity <= self._min_velocity_threshold:
            self.Estart = self._min_kinetic_energy_threshold
        else:
            self.Estart = 0.5 * self.start_velocity**2
        if self.end_velocity <= self._min_velocity_threshold:
            self.Eend = self._min_kinetic_energy_threshold
        else:
            self.Eend = 0.5 * self.end_velocity**2

        self.piecewise_approximation = self.compute_piecewise_approximation(self.Estart, self.track.max_speed)

    # @staticmethod
    def compute_piecewise_approximation(self, Estart, max_speed):
        """TODO"""
        # TODO: Implement piecewise affine approximation computations based on max_speed
        # TODO Make this method static when implemented

        piecewise_approximation = {
            'init': Piecewise((Piece(a=-4.6463e-4, b=0.0734, domain=(self.Estart, 80.8)),
                               Piece(a=-4.6463e-4, b=0.0734, domain=(80.8, 200)),
                               Piece(a=-4.6463e-4, b=0.0734, domain=(200, 312.5)))),
            15: Piecewise((Piece(a=-5.0943e-4, b=0.0767, domain=(self.Estart, 71.2)),
                           Piece(a=-1.7393e-4, b=0.0528, domain=(71.2, 100)),
                           Piece(a=-1.7393e-4, b=0.0528, domain=(100, 122.5)))),
            20: Piecewise((Piece(a=-3.1153e-4, b=0.0665, domain=(self.Estart, 115)),
                           Piece(a=-6.7188e-5, b=0.0384, domain=(115, 150)),
                           Piece(a=-6.7188e-5, b=0.0384, domain=(150, 200)))),
            30: Piecewise((Piece(a=-9.4977e-5, b= 0.0443, domain=(self.Estart, 240)),
                           Piece(a=-2.3470e-5, b=0.0272, domain=(240, 300)),
                           Piece(a=-2.3470e-5, b=0.0272, domain=(300, 450)))),
            40: Piecewise((Piece(a=-4.4240e-5, b=0.0346, domain=(self.Estart, 415)),
                           Piece(a=-9.6462e-6, b=0.0202, domain=(415, 600)),
                           Piece(a=-9.6462e-6, b=0.0202, domain=(600, 800)))),
            50: Piecewise((Piece(a=-1.8122e-5, b=0.0251, domain=(self.Estart, 640)),
                           Piece(a=-6.2127e-6, b=0.0175, domain=(640, 900)),
                           Piece(a=-6.2127e-6, b=0.0175, domain=(900, 1250)))),
            'end': Piecewise((Piece(a=-1.4458e-4, b=0.0534, domain=(self.Estart, 229.9)),
                              Piece(a=-1.4514e-6, b=0.0235, domain=(229.9, 320)),
                              Piece(a=-1.4514e-6, b=0.0235, domain=(320, 450))))
        }

        return piecewise_approximation

    @lru_cache()
    def f(self, k):
        """TODO"""
        if k == 0 and self.start_velocity <= self._min_velocity_threshold:
            return self.piecewise_approximation['init']

        if k == self.N and self.end_velocity <= self._min_velocity_threshold:
            return self.piecewise_approximation['end']

        return self.piecewise_approximation[self.track.max_speed[k]]

    @lru_cache()
    def line_resistance(self, k):
        """TODO"""
        segment = self.track[k]
        return ((self.train.slope_resistance(segment) + self.train.curve_resistance(segment))/self.train.mass,
                self.train.tunnel_resistance(segment, velocity=1)/self.train.mass)

    @lru_cache()
    def zeta(self):
        """TODO"""
        return 1 / (self.train.mass * self.train.massfactor)

    @lru_cache()
    def eta(self, k):
        """TODO"""
        return -2 * (self.train.basic_resistance_params[1] + self.line_resistance(k)[1]) / self.train.massfactor

    @lru_cache()
    def gamma(self, k):
        """TODO"""
        return -(self.train.basic_resistance_params[0] + self.line_resistance(k)[0]) / self.train.massfactor

    @lru_cache()
    def a(self, k):
        """TODO"""
        return math.exp(self.eta(k) * self.track.length[k])

    @lru_cache()
    def b(self, k):
        """TODO"""
        return (self.a(k) - 1) * self.zeta() / self.eta(k)

    @lru_cache()
    def c(self, k):
        """TODO"""
        return (self.a(k) - 1) * self.gamma(k) / self.eta(k)

    @lru_cache()
    def A(self, k):
        return Matrix([[self.a(k), 0],
                       [self.track.length[k] * (self.f(k)[2].a + self.a(k) * self.f(k+1)[2].a), 1]])

    @lru_cache()
    def B(self, k):
        return Matrix([[self.b(k)],
                       [self.track.length[k] * self.f(k+1)[2].a * self.b(k)]])

    @lru_cache()
    def C1(self, k):
        fk = self.f(k)
        return self.track.length[k] * Matrix([[0, 0, 0],
                                              [-fk[2].b, fk[1].b - fk[2].b,
                                               fk[0].b - fk[1].b + fk[2].b]])

    @lru_cache()
    def C2(self, k):
        fkk = self.f(k + 1)
        return self.track.length[k] * Matrix([[0, 0, 0],
                                              [-fkk[2].b, fkk[1].b - fkk[2].b, fkk[0].b - fkk[1].b + fkk[2].b]])

    @lru_cache()
    def D1(self, k):
        fk = self.f(k)
        return self.track.length[k] * Matrix([[0, 0, 0],
                                              [-fk[2].a, fk[1].a - fk[2].a, fk[0].a - fk[1].a + fk[2].a]])

    @lru_cache()
    def D2(self, k):
        fkk = self.f(k + 1)
        return self.track.length[k] * Matrix([[0, 0, 0],
                                              [-fkk[2].a, fkk[1].a - fkk[2].a, fkk[0].a - fkk[1].a + fkk[2].a]])

    @lru_cache()
    def e(self, k):
        fk, fkk = self.f(k), self.f(k + 1)
        return Matrix([[self.c(k)], [self.track.length[k] * (fkk[2].a * self.c(k) + fk[2].b + fkk[2].b)]])

    @lru_cache()
    def R1(self, k):
        """Matrix multiplying delta variables."""
        fk = self.f(k)
        # Constraint 1
        R1 = Matrix([[-1, 1, 0], [-1, 0, 1], [1, 1, -1]])

        # Constraint 2.1 & 2.2
        R1 |= (-1) * Matrix.diag([0] * 3)
        R1 |= Matrix.diag([fk[0].domain[0] - fk[i].domain[1] for i in range(3)])

        # Constraint 2.3 & 2.4
        R1 |= (-1) * Matrix.diag([fk[0].domain[0] - fk[i].domain[1] for i in range(3)])
        R1 |= Matrix.diag([0] * 3)

        # Constraint 3.1 & 3.2
        R1 |= Matrix.diag([0] * 3)
        R1 |= Matrix.diag([fk[0].domain[0] - fk[i].domain[1] - self.eps for i in range(2)] + [0])

        return R1

    @lru_cache()
    def R3(self):
        """Matrix multiplying z variables."""

        # Constraint 1
        R3 = Matrix.zeros(3, 3)

        # Constraint 2.1 & 2.2
        R3 |= Matrix.id(3)
        R3 |= (-1) * Matrix.id(3)

        # Constraint 2.3 & 2.4
        R3 |= (Matrix.id(3) | (-1) * Matrix.id(3))

        # Constraint 3.1 & 3.2
        R3 |= (Matrix.zeros(3, 3) | Matrix.zeros(3, 3))

        return R3

    @lru_cache()
    def R5(self):
        """Matrix multiplying u variables."""

        # Constraint 1
        R5 = Matrix.zeros(3, 1)

        # Constraint 2.1 & 2.2
        R5 |= (Matrix.zeros(3, 1) | Matrix.zeros(3, 1))

        # Constraint 2.3 & 2.4
        R5 |= (Matrix.zeros(3, 1) | Matrix.zeros(3, 1))

        # Constraint 3.1 & 3.2
        R5 |= (Matrix.zeros(3, 1) | Matrix.zeros(3, 1))

        return R5

    @lru_cache()
    def R6(self):
        """Matrix multiplying X."""

        # Constraint 1
        R6 = Matrix.zeros(3, 2)

        # Constraint 2.1 & 2.2
        R6 |= (Matrix.zeros(3, 2) | Matrix.zeros(3, 2))

        # Constraint 2.3 & 2.4
        R6 |= Matrix([[-1, 0], [-1, 0], [-1, 0]])
        R6 |= Matrix([[1, 0], [1, 0], [1, 0]])

        # Constraint 3.1 & 3.2
        R6 |= Matrix([[1, 0], [1, 0], [0, 0]])
        R6 |= Matrix([[-1, 0], [-1, 0], [0, 0]])

        return R6

    @lru_cache()
    def R7(self, k):
        """Independent term."""
        fk = self.f(k)
        # Constraint 1
        R7 = Matrix([[0], [0], [1]])

        # Constraint 2.1 & 2.2
        R7 |= (Matrix.zeros(3, 1) | Matrix.zeros(3, 1))

        # Constraint 2.3 & 2.4
        R7 |= (-1) * Matrix([[fk[0].domain[0] - fk[i].domain[1] + fk[i].domain[1]] for i in range(3)])
        R7 |= Matrix([[0 + fk[i].domain[1]] for i in range(3)])

        # Constraint 3.1 & 3.2
        R7 |= Matrix([[fk[0].domain[1] - 0], [fk[1].domain[1] - 0], [0]])
        R7 |= Matrix([[fk[0].domain[1] - self.eps], [-fk[1].domain[1] - self.eps], [0]])

        return R7

    # @lru_cache()
    def X(self, k, solution=None):
        """Model dynamics."""

        if k == 0:
            return Matrix([[self.Estart], [self.start_time]])

        if solution is not None:
            d, z, traction = solution['d'], solution['z'], solution['traction']
        else:
            d, z, traction = self.vars['d'], self.vars['z'], self.vars['traction']

        term1 = prod(self.A(j) for j in range(k)) * Matrix([[self.Estart], [self.start_time]])
        term2 = summation(prod(self.A(j) for j in range(i + 1, k)) * self.B(i) * traction[i][0] for i in range(k))
        term3 = prod(self.A(j) for j in range(1, k)) * self.C1(0) * d[0]
        if k == 1:
            term4, term7 = Matrix.zeros(2, 1), Matrix.zeros(2, 1)
        else:
            term4 = summation(prod(self.A(j) for j in range(i + 1, k)) * (self.A(i) * self.C2(i - 1) + self.C1(i)) * d[i] for i in range(1, k))
            term7 = summation(prod(self.A(j) for j in range(i + 1, k)) * (self.A(i) * self.D2(i - 1) + self.D1(i)) * z[i] for i in range(1, k))
        term5 = self.C2(k-1) * d[k]
        term6 = prod(self.A(j) for j in range(1, k)) * self.D1(0) * z[0]

        term8 = self.D2(k-1) * z[k]
        term9 = summation(prod(self.A(j) for j in range(i + 1, k)) * self.e(i) for i in range(k))

        return term1 + term2 + term3 + term4 + term5 + term6 + term7 + term8 + term9

    def build_model(self):
        self.model = model = opt.Model("optimal_trajectory_single_train")

        # Variables
        d = [Matrix([[self.model.addVar(f'd_{k}_{i}', vtype='B')] for i in range(3)]) for k in range(self.N + 1)]
        # TODO: Remove -model.infinity() lower bound in newer versions of pyscipopt
        z = [Matrix([[self.model.addVar(f'z_{k}_{i}', vtype='C', lb=-model.infinity(), ub=None)] for i in range(3)]) for k in range(self.N + 1)]
        traction = Matrix([[model.addVar(f'u_{k}', vtype='C', lb=-self.train.max_brake, ub=self.train.max_traction)] for k in range(self.N)])
        w = Matrix([[self.model.addVar(f'w_{k}', vtype='C', lb=0)] for k in range(self.N - 1)])
        delay = model.addVar('delay', vtype='C', lb=0, ub=None)
        self.vars = {'d': d, 'z': z, 'traction': traction, 'w': w, 'delay': delay}

        # Objective function
        model.setObjective(
            sum(traction[k][0] * self.track.length[k] for k in range(self.N)) +
            sum(self.smooth_factor * w[k][0] for k in range(self.N-1)) +
            self.punctuality_factor * delay)

        # Constraints
        self.constraints = cons = []
        for k in range(self.N):
            if k < self.N-1:
                model.addCons(w[k][0] >= traction[k + 1][0] - traction[k][0])
                model.addCons(w[k][0] >= traction[k][0] - traction[k + 1][0])
                cons.append(w[k][0] >= traction[k + 1][0] - traction[k][0])
                cons.append(w[k][0] >= traction[k][0] - traction[k + 1][0])

            block = self.R1(k) * d[k] + self.R3() * z[k] + self.R5() * traction[k][0] + self.R6() * self.X(k) + (-1) * self.R7(k)
            for i in range(block.rows):
                with suppress(ValueError):
                    model.addCons(block[i][0] <= 0)
                    cons.append(block[i][0] <= 0)

        # Ending conditions
        block = self.A(self.N-1) * self.X(self.N-1) + self.B(self.N-1) * traction[self.N-1][0] + self.C1(self.N-1) * d[self.N-1] + self.D1(self.N-1) * z[self.N-1] + self.e(self.N-1)
        model.addCons(self.Eend == block[0][0])
        model.addCons(self.start_time + self.timespan <= block[1][0])
        model.addCons(block[1][0] == self.start_time + self.timespan + delay)

    def solve(self):
        self.prepare_parameters()
        self.build_model()
        self.print_model()  # TODO: Only in debug mode
        self.model.writeStatistics(filename="debug.log")  # TODO: Only in debug mode
        self.model.optimize()

        # TODO Return a nicer object
        d, z, w, traction, delay = self.vars['d'], self.vars['z'], self.vars['w'], self.vars['traction'], self.vars['delay']
        sol = {'d': {}, 'z': {}, }
        for k in range(self.N + 1):
            sol['d'][k] = Matrix([[self.model.getVal(d[k][i][0])] for i in range(3)])
            sol['z'][k] = Matrix([[self.model.getVal(z[k][i][0])] for i in range(3)])

        sol['traction'] = Matrix([[self.model.getVal(traction[k][0])] for k in range(self.N)])
        sol['w'] = Matrix([[self.model.getVal(w[k][0])] for k in range(self.N-1)])
        sol['delay'] = self.model.getVal(delay)
        sol['kinetic_energy'] = [self.X(k, sol)[0][0] for k in range(self.N + 1)]
        sol['time'] = [self.X(k, sol)[1][0] for k in range(self.N + 1)]
        sol['velocity'] = [math.sqrt(2 * E) for E in sol['kinetic_energy']]

        return sol

    def print_model(self):
        with open('debug_cons.log', 'w') as f:
            for cons in self.constraints:
                f.write(str(cons.expr) + " <= " + str(cons.rhs) + '\n')
