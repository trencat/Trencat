import math
import unittest
import pandas as pd
import hypothesis
from hypothesis import given, example, settings, strategies as st
from src import Train, Profile


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

    @given(profile_input())
    @example([0, 0, (450, 455)])
    @settings(timeout=hypothesis.unlimited, deadline=None)
    def test_feasibility(self, input_list):
        """Test case based on Y. Wang, M. Zhang, J. Ma and X. Zhou, “Survey on Driverless Train Operation for Urban Rail
         Transit Systems,” Urban Rail Transit, vol. 2, no. 3-4, pp. 106--113, 01 Dec 2016.
        """
        train_velocity, train_end_velocity, timespan = input_list
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

        train = Train(mass=5.07e5, massfactor=1.06, velocity=train_velocity, max_brake=4.475e5, max_traction=3e5,
                      basic_resistance_params=(0.014 / 5.07e5, 2.564e-5 / 5.07e5))

        profile = Profile(train=train, track=track)

        success = profile.optimal_profile(timespan=timespan, end_velocity=train_end_velocity)

        if success:
            # TODO: Check that solution is valid
            pass

        del train_velocity, train_end_velocity, timespan, max_speed, min_speed, slope, bend_radius, tunnel, track,\
            train, profile, success


