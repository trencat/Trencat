from src import tracks, trains
from src import ClosedLoopControl

controller = ClosedLoopControl(tracks, trains, safety=10)
controller.compute_schedule(laps=3, offset=2)

# Files are dumped into 'model.txt'
