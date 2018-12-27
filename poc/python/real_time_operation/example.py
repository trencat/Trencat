import math
import pandas as pd
from src import Train
from src import Profile


length = [500] * 20
max_speed = [50] * 2 + [30] * 3 + [15] + [50] * 7 + [20] * 2 + [40] * 3 + [20] * 2
min_speed = [0] * 20
slope = [-0.0070967741935483875] * 2 + \
        [-0.000967741935483871] * 3 + \
        [0.001935483870967742] + \
        [0.005] * 4 + \
        [0.000967741935483871] * 3 + \
        [-0.002903225806451613] * 2 + \
        [-0.0070967741935483875] + \
        [0.000967741935483871] * 2 + \
        [0.003870967741935484] + \
        [0.006774193548387097]
bend_radius = [math.inf] * 20
tunnel = [False] * 20

track = pd.DataFrame({'length': pd.Series(length, dtype=float),
                      'min_speed': pd.Series(min_speed, dtype=float),
                      'max_speed': pd.Series(max_speed, dtype=float),
                      'slope': pd.Series(slope, dtype=float),
                      'bend_radius': pd.Series(bend_radius, dtype=float),
                      'tunnel': pd.Series(tunnel, dtype=bool), })

train = Train(mass=5.07e5, massfactor=1.06, velocity=0.0, max_brake=4.475e5, max_traction=3e5,
              basic_resistance_params=(0.014 / 5.07e5, 2.564e-5 / 5.07e5))

profile = Profile(train=train, track=track)
success = profile.optimal_profile(timespan=(450, 455), start_velocity=0.0, end_velocity=0.0)

if success:
    print(profile.nodes, profile.segments)

    dataframe = profile.export(timestep=0.1)

    print(dataframe)

    # dataframe to csv
    # with open('profile.csv', 'w') as f:
    #     dataframe.to_csv(f)
