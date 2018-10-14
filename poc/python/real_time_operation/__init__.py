GRAVITY_ACCELERATION = 9.80665  # In meters/(second**2)

from .Track import Track
from .Train import Train
from .Trajectory import Simple

__all__ = ['Track', 'Train', 'Simple', 'GRAVITY_ACCELERATION']
