from queue import Queue
import networkx as nx
import pyscipopt as scip


def generate_profile(train, track, timespan, end_velocity, start_velocity=None, approximation=20):
    """Optimal trajectory profile for a single train.

        The generated profile consists of a (time, distance, velocity) curve that drives the train along the track
        within a timespan window. The optimal driving style (i.e., accelerating, cruising, coasting and braking patterns)
        minimises the electrical energy consumed from the grid.

        Args:
            train (:obj:`Train`): Train instance.
            track (:obj:`Track`): Track instance.
            timespan (:obj:`tuple` of float): Time window in seconds that the train must take to drive along the track.
            end_velocity (float): The velocity that the train must have at the end of the given track.
            start_velocity (float): Train velocity at the beginning of the track. If None, the current train
                velocity will be used.
            approximation (int): Approximation level of the algorithm. The higher this value, the more precise the
                profile will be and more time might be required to compute it.

        Returns:
            Returns a dictionary with the optimal profile, the estimated time of arrival (eta) [s] and the total
            electrical energy consumption from the grid [joules]. The optimal profile includes velocities [m/s],
            distances [m], times [s], traction (positive) and braking (negative) forces [newtons] and electrical
            energy consumptions from the grid [joules]. If the train cannot physically reach the destination
            (for example, due to a segment with high slope) it returns None.

        Note:
            The train is required to drive within a timespan window (A, B). If the train is not able to drive along the
            track before B seconds, the train try to travel as fast as possible in order to minimise the delay (B + delay).
            Analogously, if the train is not able to drive along the track spending at least A seconds, the train will
            travel as slow as possible in order to minimise the time difference (A - A_minimum).

        Note:
            Notice that the train might not be able to physically reach the destination for certain values of the slope.
    """

    # Initialise variables
    start_velocity, N = start_velocity if start_velocity is not None else train.velocity, approximation

    model = scip.Model('Trajectory')
    punctuality_factor = 1e5
    objective_scale_factor = 1e5
    s, S = model.addVar('s', vtype='C', lb=0), model.addVar('S', vtype='C', lb=0)
    objective, time_cons = punctuality_factor * (s + S), 0

    # Compute velocities map
    velocities = {}
    for segment in track:
        if segment.max_speed not in velocities:
            velocities[segment.max_speed] = [i * segment.max_speed/(N - 1) for i in range(N)]

    # Build network
    G = nx.DiGraph()
    G.add_node('origin', input=1, output=0, velocity=start_velocity, segment_idx=0, visited=False)
    G.add_node('destination', input=0, output=1, velocity=end_velocity, segment_idx=20, visited=False)
    pending_nodes = Queue()
    pending_nodes.put('origin')

    while not pending_nodes.empty():
        node = pending_nodes.get_nowait()
        if G.nodes[node]['visited']:
            continue

        segment_idx = G.nodes[node]['segment_idx']
        input_cons = sum(G.edges[(pre_node, node)]['var'] for pre_node in G.predecessors(node)) + G.nodes[node]['input']
        output_cons = G.nodes[node]['output']

        if segment_idx < len(track):
            max_speed = min(track.max_speed[segment_idx], track.max_speed[max(0, segment_idx - 1)])
            for vel_idx, velocity in enumerate(velocities[max_speed]):
                v, next_v = G.nodes[node]['velocity'], velocity
                next_v = velocity if segment_idx < len(track) - 1 else end_velocity
                ds = track.length[segment_idx]
                a = (next_v ** 2 - v ** 2) / (2 * ds)  # Acceleration
                u = train.mass * train.massfactor * a + train.resistance(track[segment_idx], v)  # Traction force

                if not (-train.max_brake <= u <= train.max_traction):
                    continue
                if next_v + v < 1e-6:
                    continue
                dt = 2 * ds / (next_v + v)  # Time difference
                if dt > timespan[0] and dt > timespan[1]:
                    continue

                next_node = f'{segment_idx}.{vel_idx}' if segment_idx < len(track) - 1 else 'destination'
                energy = max(0.0, u * ds / objective_scale_factor)
                if not G.has_node(next_node):
                    G.add_node(next_node, input=0, output=0, velocity=next_v, segment_idx=segment_idx + 1, visited=False)
                var = model.addVar(f'{node}-{next_node}', vtype='B')
                G.add_edge(node, next_node, t=dt, force=u, energy=energy, var=var)

                # Constraints
                time_cons += dt * var
                output_cons += var
                objective += energy * var

                pending_nodes.put(next_node)
                G.nodes[node]['visited'] = True
        else:
            G.nodes[node]['visited'] = True
        model.addCons(output_cons - input_cons == 0)
    model.addCons(timespan[0] - s <= time_cons)
    model.addCons(time_cons <= timespan[1] + S)

    model.setObjective(objective)
    model.setRealParam('limits/gap', 1e-3)
    model.hideOutput()  # Todo: Enable in debug mode.

    # Run network flow optimization
    model.optimize()

    if model.getStatus() != "optimal" and not model.getSols():
        return None

    plan = {'profile': [{'velocity': start_velocity, 'distance': 0, 'time': 0, 'force': 0, 'energy_consumption': 0}]}

    node, stop = 'origin', False
    while not stop:
        for next_node in G.successors(node):
            edge = G.edges[(node, next_node)]
            val = model.getVal(edge['var'])
            if val > 0.5:
                plan['profile'].append({'velocity': G.nodes[next_node]['velocity'],
                                        'distance': track.length[G.nodes[node]['segment_idx']],
                                        'time': edge['t'],
                                        'force': edge['force'],
                                        'energy_consumption': edge['energy'] * objective_scale_factor})
                node = next_node
                stop = True if node == 'destination' else False
                break
    plan['eta'] = sum(segment_profile['time'] for segment_profile in plan['profile'])
    plan['total_energy_consumption'] = sum(segment_profile['energy_consumption'] for segment_profile in plan['profile'])

    return plan
