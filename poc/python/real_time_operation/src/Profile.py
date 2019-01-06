import networkx as nx
import numpy as np
import pandas as pd
import pyscipopt as scip
from math import sqrt
from itertools import accumulate
from queue import Queue


# Zero value threshold
threshold = 1e-6


class Profile:

    def __init__(self, train, track):
        """Initialise Profile.

            Args:
                train (:obj:`Train`): Train instance.
                track (:obj:`pd.Dataframe`): Pandas dataframe with track information. The DataFrame rows represent each
                    track segment and they are expected to be ordered as they come up in the track.
        """

        self.train = train
        self.track = track
        self.nodes = []
        self.segments = []

    def _work(self, C1, C2, a, v, dt):
        """Compute work done by the traction/braking force over a period of time.

        Parameters have the same meaning as in function :func:`traction_work`."""

        if dt < threshold:
            return 0

        return 0.5 * (self.train.mass * self.train.massfactor * a + C1) * ((v + a * dt) ** 2 - v ** 2) + \
               0.25 * C2 * ((v + a * dt) ** 4 - v ** 4)

    def _traction_work(self, C1, C2, a, v, dt):
        """Computes the work done only by the traction force in a segment.

            Args:
                C1 (float): Parameter C1 (see documentation).
                C2 (float): Parameter C2 (see documentation).
                a (float): Train acceleration.
                v (float): Train velocity at time .
                dt (float): Timespan in seconds
        """

        A = self.train.mass * self.train.massfactor * a + C1 + C2 * a * (v ** 2)
        if abs(A) < threshold:
            return 0

        B = 2 * v * C2 * (a ** 2)
        C = C2 * (a ** 3)
        switch = (B ** 2 - 4 * A * C) / (2 * A)
        if abs(switch) <= threshold:
            if C1 >= 0:
                return C1
        elif switch < 0:
            return self._work(C1, C2, a, v, dt)
        elif switch > 0:
            return self._work(C1, C2, a, v, max(0.0, dt - sqrt(switch)))

        return 0

    def optimal_profile(self, timespan, start_velocity=None, end_velocity=0.0, approximation=20):
        """Compute optimal trajectory profile assuming constant acceleration (see documentation).

            The generated profile consists of a (time, distance, velocity) curve that drives the train along the track
            within a timespan window. The optimal driving style (i.e., accelerating, cruising, coasting and braking patterns)
            minimises the work done by the traction force.

            Args:
                timespan (:obj:`tuple` of float): Time window in seconds that the train must take to drive along the track.
                start_velocity (float): Train velocity at the beginning of the track. If None, the current train
                    velocity will be used.
                end_velocity (float): The velocity that the train must have at the end of the given track. Default value
                    is zero.
                approximation (int): Approximation level of the algorithm. The higher this value, the more precise the
                    profile will be and more time might be required to compute it.

            Returns:
                True success and False otherwise. Optimal profile is stored in self.nodes and self.segments.

            Note:
                The train is required to drive within a timespan window (A, B). If the train is not able to drive along the
                track before B seconds, the train try to travel as fast as possible in order to minimise the delay (B + delay).
                Analogously, if the train is not able to drive along the track spending at least A seconds, the train will
                travel as slow as possible in order to minimise the time difference (A - A_minimum).

            Note:
                Notice that the train might not be able to physically reach the destination for certain values of slope,
                bend radius, etc.
        """

        # Initialise variables
        model = scip.Model('Trajectory')
        punctuality_factor, objective_scale_factor = 1e5, 1e5
        s, S = model.addVar('s', vtype='C', lb=0), model.addVar('S', vtype='C', lb=0)
        objective, time_cons = punctuality_factor * (s + S), 0

        # Build network
        G = nx.DiGraph()
        G.add_node('origin', input=1, output=0,
                   velocity=start_velocity if start_velocity is not None else self.train.velocity,
                   segment_idx=0, visited=False)
        G.add_node('destination', input=0, output=1, velocity=end_velocity, segment_idx=len(self.track), visited=False)
        pending_nodes = Queue()
        pending_nodes.put('origin')

        while not pending_nodes.empty():
            node = pending_nodes.get_nowait()
            if G.nodes[node]['visited']:
                continue

            segment_idx = G.nodes[node]['segment_idx']
            input_cons = sum(G.edges[(pre_node, node)]['var'] for pre_node in G.predecessors(node)) + G.nodes[node][
                'input']
            output_cons = G.nodes[node]['output']

            if segment_idx < len(self.track):
                min_speed = self.track.min_speed.iloc[segment_idx] if segment_idx < len(self.track) - 1 else end_velocity
                if segment_idx == len(self.track) - 1:
                    max_speed = end_velocity
                elif self.track.max_speed.iloc[segment_idx + 1] > self.track.max_speed.iloc[segment_idx]:
                    max_speed = self.track.max_speed.iloc[segment_idx]
                else:
                    max_speed = self.track.max_speed.iloc[segment_idx + 1]

                v, ds = G.nodes[node]['velocity'], self.track.length.iloc[segment_idx]
                C1 = self.train.resistance(segment=self.track.iloc[segment_idx], velocity=0)
                C2 = self.train.resistance(segment=self.track.iloc[segment_idx], velocity=1) - C1

                for next_v_idx, next_v in enumerate(np.linspace(start=min_speed, stop=max_speed, endpoint=True,
                                                                dtype=float,
                                                                num=approximation if abs(max_speed - min_speed) > threshold else 1)):
                    # Time
                    if next_v + v < 1e-6:
                        continue
                    dt = 2 * ds / (next_v + v)  # Time difference

                    # Acceleration
                    a = (next_v - v) / dt  # Acceleration

                    # Max traction/brake force
                    max_traction = self.train.mass * self.train.massfactor * a + C1 + C2 * (v ** 2)
                    max_brake = self.train.mass * self.train.massfactor * a + C1 + C2 * (next_v ** 2)
                    if max_traction > self.train.max_traction or abs(max_brake) > self.train.max_brake:
                        continue

                    # Traction parameters
                    force = (self.train.mass * self.train.massfactor * a + C1, C2)

                    # Work and jerk rate
                    segment_work = self._traction_work(C1=C1, C2=C2, a=a, v=v, dt=dt) / objective_scale_factor
                    jerk = 2 * C2 * abs(a) * ds

                    # Graph update
                    next_node = f'{segment_idx + 1}.{next_v_idx}' if segment_idx < len(self.track) - 1 else 'destination'
                    if not G.has_node(next_node):
                        G.add_node(next_node, input=0, output=0, velocity=next_v, segment_idx=segment_idx + 1,
                                   visited=False)
                    var = model.addVar(f'{node}-{next_node}', vtype='B')
                    G.add_edge(node, next_node, var=var, timespan=dt, work=segment_work, jerk=jerk, force=force,
                               acceleration=a)
                    pending_nodes.put(next_node)

                    # Constraints
                    time_cons += dt * var
                    output_cons += var
                    objective += segment_work * var

            G.nodes[node]['visited'] = True

            model.addCons(output_cons - input_cons == 0)
        model.addCons(timespan[0] - s <= time_cons)
        model.addCons(time_cons <= timespan[1] + S)

        model.setObjective(objective)
        model.setRealParam('limits/gap', 1e-3)
        model.hideOutput()  # Todo: Enable in debug mode.

        # Run network flow optimization
        model.optimize()

        # model.printStatistics()

        # Generate output
        if model.getStatus() != "optimal" and not model.getSols():
            return False

        self.nodes, self.segments = [], []

        edges = filter(lambda x: model.getVal(G.edges[x]['var']) > 0.5, G.edges)
        edges = sorted(edges, key=lambda x: G.nodes[x[0]]['segment_idx'])

        time = 0
        self.nodes.append({'velocity': start_velocity, 'time': time, })
        for (node, next_node) in edges:
            edge = G.edges[(node, next_node)]
            time += edge['timespan']
            self.nodes.append({'velocity': G.nodes[next_node]['velocity'], 'time': time,})
            self.segments.append({'timespan': edge['timespan'],
                                  'work': edge['work'] * objective_scale_factor,
                                  'jerk': edge['jerk'],
                                  'force': edge['force'],
                                  'acceleration': edge['acceleration']})

        self.nodes[-1]['time'] = self.nodes[-2]['time'] + self.segments[-1]['timespan']

        return True

    def export(self, timestep=0.1):
        """Export a profile to a DataFrame.

            Args:
                timestep (float): Timestep in seconds.

            Returns:
                Dataframe with train motion profile.
        """

        df = pd.DataFrame(columns=['time', 'segment_no', 'position', 'velocity', 'acceleration', 'force',
                                   'power', 'work', 'work_traction', 'jerk', 'slope', 'max_speed', 'bend_radius',
                                   'tunnel'], dtype=float)
        df.segment_no.astype(int)
        df.tunnel.astype(bool)

        breakpoints = [0] + list(accumulate(self.track.length))

        for index, (n, s, bp) in enumerate(zip(self.nodes, self.segments, breakpoints)):
            time, timespan, time_end = n['time'], s['timespan'], n['time'] + s['timespan']
            time_now, time_then = time, min(time + timestep, time_end)
            dt = time_then - time_now
            velocity_now, velocity_then = n['velocity'], n['velocity'] + s['acceleration'] * dt
            position_now, position_then = bp, bp + 0.5 * (n['velocity'] + velocity_then) * dt
            ds = position_then - position_now

            while time_then <= time_end:
                force = s['force'][0] + s['force'][1] * (velocity_then ** 2)
                df = df.append({'time': time_then,
                                'segment_no': index,
                                'position': position_then,
                                'velocity': velocity_then,
                                'acceleration': s['acceleration'],
                                'force': force,
                                'power': force * velocity_then,
                                'work': force * ds,
                                'work_traction': max(0.0, force) * ds,
                                'jerk': 2 * s['force'][1] * abs(s['acceleration']) * ds,
                                'slope': self.track.slope[index],
                                'max_speed': self.track.max_speed[index],
                                'bend_radius': self.track.bend_radius[index],
                                'tunnel': self.track.tunnel[index]}, ignore_index=True)

                if abs(time_end - time_then) <= threshold:
                    break

                time_now, time_then = time_then, min(time_then + timestep, time_end)
                dt = time_then - time_now
                velocity_now, velocity_then = velocity_then, velocity_then + s['acceleration'] * dt
                position_now, position_then = position_then, position_then + 0.5 * (velocity_now + velocity_then) * dt
                ds = position_then - position_now

        return df

    def total_time(self):
        return self.nodes[-1]['time']

    def total_work(self):
        return sum(segment['work'] for segment in self.segments)

    def total_jerk(self):
        return sum(segment['jerk'] for segment in self.segments)
