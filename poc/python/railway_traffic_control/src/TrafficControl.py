from itertools import dropwhile, takewhile, tee
import pyscipopt as scip


class ClosedLoopControl:
    """Railway Traffic Control implementation for closed loop tracks.

        Closed loop tracks do not have neither start nor end.
        Trains circulate cyclically and can not overtake preceding trains.
        This is the case of metro/underground/tube railways.
    """

    def __init__(self, tracks, trains, safety):
        """Initialise Traffic Control

            Args:
                tracks (dict): Dict with track ids and travel time.
                trains (dict): Dict with train ids and their current positions.
                safety (float): Safety time between conflicting operations.
        """

        self.tracks = tracks
        self.trains = trains
        self.safety = safety  # Safety time between conflicting operations
        self.model = None
        self.vars, self.cons = {}, {}

    def add_var(self, name, offset, **kwargs):
        """Adds and returns a variable to the model.

            If the variable already exists in the model, nothing happens and
            such variable is returned.

            Args:
                name (str): Variable name.
                offset (float): Variable lower bound.

            kwargs allows to store extra info of the variable.
            """

        if name not in self.vars:
            self.vars[name] = {'name': name, 
                               'var': self.model.addVar(name, vtype='C',
                                                        lb=offset),
                               **kwargs}

        return self.vars[name]

    def add_cons(self, expression, name='Cons'):
        """Adds a constraint to the model."""
 
        if name not in self.cons:
            self.cons[name] = []

        self.cons[name].append(expression)
        self.model.addCons(expression)


    def get_operations(self, from_track, laps):
        """Get a list of forthcoming operations, starting from the given one.

            The list starts and ends at the same point.

            Args:
                from_track (str): Id of the starting track.
                laps (int): Number of laps to perform.
        """

        # Shortcut
        tracks = self.tracks

        lap = list(dropwhile(lambda x, track=from_track: x != track, tracks.keys())) +\
              list(takewhile(lambda x, track=from_track: x != track, tracks.keys()))

        return ['start'] + lap * laps + ['end']

    def build_model(self, laps, offset):
        """Build optimization problem."""

        # Operations to be performed for each train
        operations = {train_id: self.get_operations(train['position'], laps=laps)
                      for train_id, train in self.trains.items()}
        conflicts = {}

        counter, stop = 0, False
        while not stop:
            for train, ops in operations.items():
                track, next_track = ops[counter], ops[counter+1]
                op_name = f'{train}_{track}_{counter}'
                next_op_name = f'{train}_{next_track}_{counter+1}'
                time = self.tracks[track]['time'] if track != 'start' else 0

                if track == 'start':
                    op_name = track
                if next_track == 'end':
                    next_op_name = next_track

                # Variable
                op = self.add_var(op_name, offset, time=time)
                next_op = self.add_var(next_op_name, offset, time=time)

                # Constraint
                self.add_cons( next_op['var'] >= op['var'] + time)

                # Conflicts constraints
                if track in conflicts:
                    latest = conflicts[track]['next_op']

                    self.add_cons(op['var'] >= latest['var'] + self.safety, name='Conf')

                # Store next operation for conflict resolution
                if track != 'start':
                    conflicts[track] = {'op': op,
                        'next_op': next_op if next_track != 'end' else op}

            counter += 1
            stop = True if counter == len(ops) - 1 else False

        # Fix dummy start variable
        self.model.chgVarUb(self.vars['start']['var'], offset)

    def export_model(self, filename):
        """Export model in LP format into a file."""

        self.model.writeProblem(filename)

    def get_results(self, filename):
        """Export results into a file in readable format."""

        with open(filename, 'w') as file:
            # Print solution status
            file.write(f'Solution: {self.model.getStatus()}\n')
            file.write(f'Objective function: {self.model.getObjVal()}\n')

            for var_name, var in self.vars.items():
                file.write(f'{var_name}: {self.model.getVal(var["var"])}\n')

    def compute_schedule(self, laps, offset=0):
        """Compute train schedule.

            Results are dumped into a file named 'model.txt'.

            Args:
                laps (int): Number of laps that trains must run.
                offset (float): Offset.
        """

        # Initialise model
        self.model = scip.Model('Schedule')
        self.build_model(laps, offset)

        # Create objective function
        start, end = self.vars['start']['var'], self.vars['end']['var']
        objective = self.model.setObjective(end - start, sense='minimize')

        # Export file
        # self.export_model('model.txt')

        # Solve
        # self.model.hideOutput()
        self.model.optimize()

        # Obtain results
        self.get_results('solution.txt')

