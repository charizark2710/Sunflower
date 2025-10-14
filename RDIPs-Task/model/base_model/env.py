from transformers import RobertaTokenizer, RobertaModel
from utils.common import DEVICE

from gymnasium.envs.registration import register
import gymnasium as gym
from gymnasium import spaces
import numpy as np
import psutil
import time

CPU_RATE = 3.5e9

register(
    id='task-eval-v0',
    entry_point='v0_task_eval_env:TaskEvalEnv',
)

SAVE_DIR = "./saved_model"
LR = 1e-4

class TaskEvalEnv(gym.Env):
    def __init__(self):
        super(TaskEvalEnv, self).__init__()
        psutil.cpu_percent(interval=5)
        self.memory_total = psutil.virtual_memory().total
        self.action_space = spaces.Discrete(2)  # Skip or Execute
        self.observation_space = spaces.Box(
            low=np.array([0.0, 0.0, 0.0]),
            high=np.array([1.0, 100.0, 100.0]),
            shape=(3,),
            dtype=np.float32
        )

        # placeholders for context updated externally
        self.result = {}
        self.task = {}

    def update_context(self, result=None, task=None):
        """
        Called externally before env.step(action).
        Stores the latest result and task context.
        """
        if task is not None:
            self.task = task
        if result is not None:
            self.result = result

    def step(self, action):
        """
        Gym-compatible step — uses stored self.result and self.task.
        """
        certainty, ic_est, cycle_est = self.task
        result = self.result
        reward = 0.0
        executed = int(action) == 1

        actual_ic = ic_est
        actual_cycle = cycle_est

        percent_deduction = 1.0
        if executed:
            mem_usage = psutil.virtual_memory().used
            cpu_usage = psutil.cpu_percent(interval=0)
            mem_percent = mem_usage / self.memory_total
            cpu_percent = cpu_usage / 100.0

            # apply highest thresholds first
            if mem_percent > 0.9:
                percent_deduction *= 2.0
            elif mem_percent > 0.7:
                percent_deduction *= 1.5
            elif mem_percent > 0.5:
                percent_deduction *= 1.2

            if cpu_percent > 0.9:
                percent_deduction *= 2.0
            elif cpu_percent > 0.7:
                percent_deduction *= 1.5
            elif cpu_percent > 0.5:
                percent_deduction *= 1.2

        # --- Reward logic ---
        if (0.5 <= certainty <= 0.7) and executed:
            if result is None:
                reward = -5.0
            else:
                actual_ic = result.get("actual_ic", ic_est)
                actual_cycle = result.get("actual_cycle", cycle_est)
                log_ic = np.log10(max(actual_ic, 1e-12))
                log_cycle = np.log10(max(actual_cycle, 1e-12))
                log_ic_est = np.log10(max(ic_est, 1e-12))
                log_cycle_est = np.log10(max(cycle_est, 1e-12))
                ic_diff = abs(log_ic - log_ic_est)
                cycle_diff = abs(log_cycle - log_cycle_est)
                reward = 1.0 if (ic_diff < 0.3 and cycle_diff < 0.3) else -1.0

        elif certainty > 0.7:
            if executed:
                if result is None:
                    reward = -5.0
                else:
                    reward = 2.0
                    actual_ic = result.get("actual_ic", ic_est)
                    actual_cycle = result.get("actual_cycle", cycle_est)
            else:
                reward = 0.5

        elif certainty < 0.5:
            if result is None:
                reward = -7.0
            else:
                actual_ic = result.get("actual_ic", ic_est)
                actual_cycle = result.get("actual_cycle", cycle_est)
                # Add accuracy check like in medium certainty case
                log_ic = np.log10(max(actual_ic, 1e-12))
                log_cycle = np.log10(max(actual_cycle, 1e-12))
                log_ic_est = np.log10(max(ic_est, 1e-12))
                log_cycle_est = np.log10(max(cycle_est, 1e-12))
                ic_diff = abs(log_ic - log_ic_est)
                cycle_diff = abs(log_cycle - log_cycle_est)
                reward = 0.5 if (ic_diff < 0.3 and cycle_diff < 0.3) else -0.5

        done = True
        obs = self._get_state()
        return obs, reward / percent_deduction, done

    def _get_state(self):
        """Return [mem_usage_ratio, cpu_usage%]."""
        cpu_usage = psutil.cpu_percent(interval=0)
        mem_usage = psutil.virtual_memory().used / self.memory_total
        return np.array([mem_usage, cpu_usage], dtype=np.float32)

    def reset(self, *, seed=None, options=None):
        super().reset(seed=seed)
        self.result = {}
        self.task = {}
        time.sleep(1)
        return self._get_state()
