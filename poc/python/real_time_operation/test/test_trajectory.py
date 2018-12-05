import math
import unittest
from hypothesis import given, example, settings, strategies as st
from src import Train, Track, generate_profile
#from real_time_operation import Track, Train, generate_profile


@st.composite
def profile_input(draw):
    """Generate examples for hypothesis module."""
    # Random data
    train_velocity = draw(st.floats(min_value=0, max_value=50), label='train_velocity')
    train_end_velocity = draw(st.floats(min_value=0, max_value=20), label='train_end_velocity')
    timespan = draw(st.tuples(st.floats(min_value=300, max_value=1000),
                              st.floats(min_value=300, max_value=1000)).map(sorted).filter(lambda x: x[0] < x[1]),
                    label='timespan')

    return train_velocity, train_end_velocity, timespan


class TestProfile(unittest.TestCase):
    # TODO: Write more tests.

    @given(profile_input())
    @example([0, 0, (450, 455)])
    @settings(deadline=None)
    def test_feasibility(self, input_list):
        """Test case based on Y. Wang, M. Zhang, J. Ma and X. Zhou, “Survey on Driverless Train Operation for Urban Rail
         Transit Systems,” Urban Rail Transit, vol. 2, no. 3-4, pp. 106--113, 01 Dec 2016.
        """
        train_velocity, train_end_velocity, timespan = input_list
        self.length = [500] * 20
        self.max_speed = [50] * 2 + [30] * 3 + [15] + [50] * 7 + [20] * 2 + [40] * 3 + [20] * 2
        self.slope = [-0.0070967741935483875] * 2 + \
                [-0.000967741935483871] * 3 + \
                [0.001935483870967742] + \
                [0.005] * 4 + \
                [0.000967741935483871] * 3 + \
                [-0.002903225806451613] * 2 + \
                [-0.0070967741935483875] + \
                [0.000967741935483871] * 2 + \
                [0.003870967741935484] + \
                [0.006774193548387097]
        self.bend_radius = [math.inf] * 20
        self.tunnel = [False] * 20

        train = Train(mass=5.07e5, massfactor=1.06, velocity=train_velocity, max_brake=4.475e5, max_traction=3e5,
                      basic_resistance_params=(0.014 / 5.07e5, 2.564e-5 / 5.07e5))
        track = Track(length=self.length, max_speed=self.max_speed, slope=self.slope, bend_radius=self.bend_radius,
                      tunnel=self.tunnel)
        trajectory = generate_profile(train=train, track=track, timespan=timespan, end_velocity=train_end_velocity)

        self.assertIsNotNone(trajectory)
