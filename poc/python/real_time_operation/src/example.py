import math
from Track import Track
from Train import Train
from profile import generate_profile

length = [500] * 20
max_speed = [50] * 2 + [30] * 3 + [15] + [50] * 7 + [20] * 2 + [40] * 3 + [20] * 2
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

train = Train(mass=5.07e5, massfactor=1.06, velocity=0, max_brake=4.475e5, max_traction=3e5,
              basic_resistance_params=(0.014 / 5.07e5, 2.564e-5 / 5.07e5))
track = Track(length=length, max_speed=max_speed, slope=slope, bend_radius=bend_radius, tunnel=tunnel)
trajectory = generate_profile(train=train, track=track, timespan=(450, 455), end_velocity=0)

print(trajectory)